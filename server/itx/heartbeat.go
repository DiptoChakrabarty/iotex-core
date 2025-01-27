// Copyright (c) 2019 IoTeX Foundation
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package itx

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime"
	"strconv"
	"strings"

	"github.com/iotexproject/go-fsm"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"

	"github.com/iotexproject/iotex-core/consensus"
	"github.com/iotexproject/iotex-core/consensus/scheme"
	"github.com/iotexproject/iotex-core/consensus/scheme/rolldpos"
	"github.com/iotexproject/iotex-core/dispatcher"
	"github.com/iotexproject/iotex-core/p2p"
	"github.com/iotexproject/iotex-core/pkg/log"
	"github.com/iotexproject/iotex-core/pkg/version"
	statedb "github.com/iotexproject/iotex-core/state"
)

// TODO: HeartbeatHandler opens encapsulation of a few structs to inspect the internal status, we need to find a better
// approach to do so in the future

var heartbeatMtc = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "iotex_heartbeat_status",
		Help: "Node heartbeat status.",
	},
	[]string{"status_type", "source"},
)

var versionMtc = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "iotex_version_status",
		Help: "Node software version status.",
	},
	[]string{"type", "value"},
)

func init() {
	prometheus.MustRegister(heartbeatMtc)
	prometheus.MustRegister(versionMtc)
}

// HeartbeatHandler is the handler to periodically log the system key metrics
type HeartbeatHandler struct {
	s *Server
	l *zap.Logger
}

// NewHeartbeatHandler instantiates a HeartbeatHandler instance
func NewHeartbeatHandler(s *Server, cfg p2p.Network) *HeartbeatHandler {
	return &HeartbeatHandler{
		s: s,
		l: log.L().With(zap.String("networkAddr", fmt.Sprintf("%s:%d", cfg.Host, cfg.Port))),
	}
}

// Log executes the logging logic
func (h *HeartbeatHandler) Log() {
	// Network metrics
	p2pAgent := h.s.P2PAgent()

	// Dispatcher metrics
	dp, ok := h.s.Dispatcher().(*dispatcher.IotxDispatcher)
	if !ok {
		h.l.Error("dispatcher is not the instance of IotxDispatcher")
		return
	}
	numDPEvts := dp.EventQueueSize()
	totalDPEventNumber := 0
	events := []string{}
	for event, num := range numDPEvts {
		totalDPEventNumber += num
		events = append(events, event+":"+strconv.Itoa(num))
	}
	dpEvtsAudit, err := json.Marshal(dp.EventAudit())
	if err != nil {
		h.l.Error("error when serializing the dispatcher event audit map.", zap.Error(err))
		return
	}

	ctx := context.Background()
	peers, err := p2pAgent.Neighbors(ctx)
	if err != nil {
		h.l.Debug("error when get neighbors.", zap.Error(err))
		peers = nil
	}
	numPeers := len(peers)
	h.l.Debug("Node status.",
		zap.Int("numPeers", numPeers),
		zap.String("pendingDispatcherEvents", "{"+strings.Join(events, ", ")+"}"),
		zap.String("pendingDispatcherEventsAudit", string(dpEvtsAudit)))

	heartbeatMtc.WithLabelValues("numPeers", "node").Set(float64(numPeers))
	heartbeatMtc.WithLabelValues("pendingDispatcherEvents", "node").Set(float64(totalDPEventNumber))
	// chain service
	for _, c := range h.s.chainservices {
		// Consensus metrics
		cs, ok := c.Consensus().(*consensus.IotxConsensus)
		if !ok {
			h.l.Info("consensus is not the instance of IotxConsensus.")
			return
		}
		rolldpos, ok := cs.Scheme().(*rolldpos.RollDPoS)
		numPendingEvts := 0
		consensusEpoch := uint64(0)
		consensusHeight := uint64(0)
		height := c.Blockchain().TipHeight()

		var consensusMetrics scheme.ConsensusMetrics
		var state fsm.State
		if ok {
			numPendingEvts = rolldpos.NumPendingEvts()
			state = rolldpos.CurrentState()

			// RollDpos Consensus Metrics
			consensusMetrics, err = rolldpos.Metrics()
			if err != nil {
				if height > 0 || errors.Cause(err) != statedb.ErrStateNotExist {
					h.l.Error("failed to read consensus metrics", zap.Error(err))
					return
				}
			}
			consensusEpoch = consensusMetrics.LatestEpoch
			consensusHeight = consensusMetrics.LatestHeight
		} else {
			h.l.Debug("scheme is not the instance of RollDPoS")
		}

		// Block metrics
		actPoolSize := c.ActionPool().GetSize()
		actPoolCapacity := c.ActionPool().GetCapacity()
		targetHeight := c.BlockSync().TargetHeight()

		h.l.Debug("chain service status",
			zap.Int("rolldposEvents", numPendingEvts),
			zap.String("fsmState", string(state)),
			zap.Uint64("blockchainHeight", height),
			zap.Uint64("actpoolSize", actPoolSize),
			zap.Uint64("actpoolCapacity", actPoolCapacity),
			zap.Uint32("chainID", c.ChainID()),
			zap.Uint64("targetHeight", targetHeight),
			zap.Uint64("concensusEpoch", consensusEpoch),
			zap.Uint64("consensusHeight", consensusHeight),
		)

		chainIDStr := strconv.FormatUint(uint64(c.ChainID()), 10)
		heartbeatMtc.WithLabelValues("consensusEpoch", chainIDStr).Set(float64(consensusHeight))
		heartbeatMtc.WithLabelValues("consensusRound", chainIDStr).Set(float64(consensusEpoch))
		heartbeatMtc.WithLabelValues("pendingRolldposEvents", chainIDStr).Set(float64(numPendingEvts))
		heartbeatMtc.WithLabelValues("blockchainHeight", chainIDStr).Set(float64(height))
		heartbeatMtc.WithLabelValues("actpoolSize", chainIDStr).Set(float64(actPoolSize))
		heartbeatMtc.WithLabelValues("actpoolGasInPool", chainIDStr).Set(float64(c.ActionPool().GetGasSize()))
		heartbeatMtc.WithLabelValues("actpoolCapacity", chainIDStr).Set(float64(actPoolCapacity))
		heartbeatMtc.WithLabelValues("targetHeight", chainIDStr).Set(float64(targetHeight))
		heartbeatMtc.WithLabelValues("packageVersion", version.PackageVersion).Set(1)
		heartbeatMtc.WithLabelValues("packageCommitID", version.PackageCommitID).Set(1)
		heartbeatMtc.WithLabelValues("goVersion", version.GoVersion).Set(1)
	}

	// Mem metrics
	memMetrics()
}

func memMetrics() {
	bToMb := func(b uint64) uint64 {
		return b / 1024 / 1024
	}
	var memStat runtime.MemStats
	runtime.ReadMemStats(&memStat)
	heartbeatMtc.WithLabelValues("allocatedHeapObjects", "node").Set(float64(bToMb(memStat.Alloc)))
	heartbeatMtc.WithLabelValues("totalAllocatedHeapObjects", "node").Set(float64(bToMb(memStat.TotalAlloc)))
	heartbeatMtc.WithLabelValues("stackInUse", "node").Set(float64(bToMb(memStat.StackInuse)))
	heartbeatMtc.WithLabelValues("stackFromOS", "node").Set(float64(bToMb(memStat.StackSys)))
	heartbeatMtc.WithLabelValues("totalFromOS", "node").Set(float64(bToMb(memStat.Sys)))
	heartbeatMtc.WithLabelValues("heapInUse", "node").Set(float64(bToMb(memStat.HeapInuse)))
	heartbeatMtc.WithLabelValues("heapFromOS", "node").Set(float64(bToMb(memStat.HeapSys)))
	heartbeatMtc.WithLabelValues("heapIdle", "node").Set(float64(bToMb(memStat.HeapIdle)))
	heartbeatMtc.WithLabelValues("heapReleased", "node").Set(float64(bToMb(memStat.HeapReleased)))
	heartbeatMtc.WithLabelValues("numberOfGC", "node").Set(float64(memStat.NumGC))
	heartbeatMtc.WithLabelValues("numberOfRoutines", "node").Set(float64(runtime.NumGoroutine()))
}
