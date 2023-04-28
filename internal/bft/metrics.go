package bft

import metrics "github.com/SmartBFT-Go/consensus/pkg/api"

const (
	nameBlackListNodeID = "blackid"
	nameReasonFailAdd   = "reason"

	reasonRequestMaxBytes      = "MAX_BYTES"
	reasonSemaphoreAcquireFail = "SEMAPHORE_ACQUIRE_FAIL"
)

var countOfRequestPoolOpts = metrics.GaugeOpts{
	Namespace:    "consensus",
	Subsystem:    "bft",
	Name:         "pool_count_of_elements",
	Help:         "Number of elements in the consensus request pool.",
	LabelNames:   []string{},
	StatsdFormat: "%{#fqname}",
}

var countOfFailAddRequestToPoolOpts = metrics.CounterOpts{
	Namespace:    "consensus",
	Subsystem:    "bft",
	Name:         "pool_count_of_fail_add_request",
	Help:         "Number of requests pool insertion failure.",
	LabelNames:   []string{nameReasonFailAdd},
	StatsdFormat: "%{#fqname}.%{" + nameReasonFailAdd + "}",
}

// ForwardTimeout
var countOfLeaderForwardRequestOpts = metrics.CounterOpts{
	Namespace:    "consensus",
	Subsystem:    "bft",
	Name:         "pool_count_leader_forward_request",
	Help:         "Number of requests forwarded to the leader.",
	LabelNames:   []string{},
	StatsdFormat: "%{#fqname}",
}

var countTimeoutTwoStepOpts = metrics.CounterOpts{
	Namespace:    "consensus",
	Subsystem:    "bft",
	Name:         "pool_count_timeout_two_step",
	Help:         "Number of times requests reached second timeout.",
	LabelNames:   []string{},
	StatsdFormat: "%{#fqname}",
}

var countOfDeleteRequestPoolOpts = metrics.CounterOpts{
	Namespace:    "consensus",
	Subsystem:    "bft",
	Name:         "pool_count_of_delete_request",
	Help:         "Number of elements removed from the request pool.",
	LabelNames:   []string{},
	StatsdFormat: "%{#fqname}",
}

var countOfRequestPoolAllOpts = metrics.CounterOpts{
	Namespace:    "consensus",
	Subsystem:    "bft",
	Name:         "pool_count_of_elements_all",
	Help:         "Total amount of elements in the request pool.",
	LabelNames:   []string{},
	StatsdFormat: "%{#fqname}",
}

var latencyOfRequestPoolOpts = metrics.HistogramOpts{
	Namespace:    "consensus",
	Subsystem:    "bft",
	Name:         "pool_latency_of_elements",
	Help:         "The average request processing time, time request resides in the pool.",
	Buckets:      []float64{0.005, 0.01, 0.015, 0.05, 0.1, 1, 10},
	LabelNames:   []string{},
	StatsdFormat: "%{#fqname}",
}

// MetricsRequestPool encapsulates request pool metrics
type MetricsRequestPool struct {
	CountOfRequestPool          metrics.Gauge
	CountOfFailAddRequestToPool metrics.Counter
	CountOfLeaderForwardRequest metrics.Counter
	CountTimeoutTwoStep         metrics.Counter
	CountOfDeleteRequestPool    metrics.Counter
	CountOfRequestPoolAll       metrics.Counter
	LatencyOfRequestPool        metrics.Histogram

	labels []string
}

// NewMetricsRequestPool create new request pool metrics
func NewMetricsRequestPool(p *metrics.CustomerProvider) *MetricsRequestPool {
	countOfRequestPoolOptsTmp := p.NewGaugeOpts(countOfRequestPoolOpts)
	countOfFailAddRequestToPoolOptsTmp := p.NewCounterOpts(countOfFailAddRequestToPoolOpts)
	countOfLeaderForwardRequestOptsTmp := p.NewCounterOpts(countOfLeaderForwardRequestOpts)
	countTimeoutTwoStepOptsTmp := p.NewCounterOpts(countTimeoutTwoStepOpts)
	countOfDeleteRequestPoolOptsTmp := p.NewCounterOpts(countOfDeleteRequestPoolOpts)
	countOfRequestPoolAllOptsTmp := p.NewCounterOpts(countOfRequestPoolAllOpts)
	latencyOfRequestPoolOptsTmp := p.NewHistogramOpts(latencyOfRequestPoolOpts)
	return &MetricsRequestPool{
		CountOfRequestPool:          p.NewGauge(countOfRequestPoolOptsTmp).With(p.LabelsForWith()...),
		CountOfFailAddRequestToPool: p.NewCounter(countOfFailAddRequestToPoolOptsTmp),
		CountOfLeaderForwardRequest: p.NewCounter(countOfLeaderForwardRequestOptsTmp).With(p.LabelsForWith()...),
		CountTimeoutTwoStep:         p.NewCounter(countTimeoutTwoStepOptsTmp).With(p.LabelsForWith()...),
		CountOfDeleteRequestPool:    p.NewCounter(countOfDeleteRequestPoolOptsTmp).With(p.LabelsForWith()...),
		CountOfRequestPoolAll:       p.NewCounter(countOfRequestPoolAllOptsTmp).With(p.LabelsForWith()...),
		LatencyOfRequestPool:        p.NewHistogram(latencyOfRequestPoolOptsTmp).With(p.LabelsForWith()...),
		labels:                      p.LabelsForWith(),
	}
}

func (m *MetricsRequestPool) LabelsForWith(labelValues ...string) []string {
	result := make([]string, 0, len(m.labels)+len(labelValues))
	result = append(result, labelValues...)
	result = append(result, m.labels...)
	return result
}

var countBlackListOpts = metrics.GaugeOpts{
	Namespace:    "consensus",
	Subsystem:    "bft",
	Name:         "blacklist_count",
	Help:         "Count of nodes in blacklist on this channel.",
	LabelNames:   []string{},
	StatsdFormat: "%{#fqname}",
}

var nodesInBlackListOpts = metrics.GaugeOpts{
	Namespace:    "consensus",
	Subsystem:    "bft",
	Name:         "node_id_in_blacklist",
	Help:         "Node ID in blacklist on this channel.",
	LabelNames:   []string{nameBlackListNodeID},
	StatsdFormat: "%{#fqname}.%{" + nameBlackListNodeID + "}",
}

// MetricsBlacklist encapsulates blacklist metrics
type MetricsBlacklist struct {
	CountBlackList   metrics.Gauge
	NodesInBlackList metrics.Gauge

	labels []string
}

// NewMetricsBlacklist create new blacklist metrics
func NewMetricsBlacklist(p *metrics.CustomerProvider) *MetricsBlacklist {
	countBlackListOptsTmp := p.NewGaugeOpts(countBlackListOpts)
	nodesInBlackListOptsTmp := p.NewGaugeOpts(nodesInBlackListOpts)
	return &MetricsBlacklist{
		CountBlackList:   p.NewGauge(countBlackListOptsTmp).With(p.LabelsForWith()...),
		NodesInBlackList: p.NewGauge(nodesInBlackListOptsTmp),
		labels:           p.LabelsForWith(),
	}
}

func (m *MetricsBlacklist) LabelsForWith(labelValues ...string) []string {
	result := make([]string, 0, len(m.labels)+len(labelValues))
	result = append(result, labelValues...)
	result = append(result, m.labels...)
	return result
}

var consensusReconfigOpts = metrics.CounterOpts{
	Namespace:    "consensus",
	Subsystem:    "bft",
	Name:         "consensus_reconfig",
	Help:         "Number of reconfiguration requests.",
	LabelNames:   []string{},
	StatsdFormat: "%{#fqname}",
}

var latencySyncOpts = metrics.HistogramOpts{
	Namespace:    "consensus",
	Subsystem:    "bft",
	Name:         "consensus_latency_sync",
	Help:         "An average time it takes to sync node.",
	Buckets:      []float64{0.005, 0.01, 0.015, 0.05, 0.1, 1, 10},
	LabelNames:   []string{},
	StatsdFormat: "%{#fqname}",
}

// MetricsConsensus encapsulates consensus metrics
type MetricsConsensus struct {
	CountConsensusReconfig metrics.Counter
	LatencySync            metrics.Histogram
}

// NewMetricsConsensus create new consensus metrics
func NewMetricsConsensus(p *metrics.CustomerProvider) *MetricsConsensus {
	consensusReconfigOptsTmp := p.NewCounterOpts(consensusReconfigOpts)
	latencySyncOptsTmp := p.NewHistogramOpts(latencySyncOpts)
	return &MetricsConsensus{
		CountConsensusReconfig: p.NewCounter(consensusReconfigOptsTmp).With(p.LabelsForWith()...),
		LatencySync:            p.NewHistogram(latencySyncOptsTmp).With(p.LabelsForWith()...),
	}
}