/*
Copyright 2022 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"k8s.io/client-go/util/workqueue"
)

// This file sets the workqueue DefaultMetricsFactory to produce
// prometheus metrics for workqueue.

// Metrics subsystem and keys used by the workqueue.
const (
	WorkQueueSubsystem         = "workqueue"
	DepthKey                   = "depth"
	AddsKey                    = "adds_total"
	QueueLatencyKey            = "queue_duration_seconds"
	WorkDurationKey            = "work_duration_seconds"
	UnfinishedWorkKey          = "unfinished_work_seconds"
	LongestRunningProcessorKey = "longest_running_processor_seconds"
	RetriesKey                 = "retries_total"
	nameLabel                  = "name"
)

var (
	depth = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Subsystem: WorkQueueSubsystem,
		Name:      DepthKey,
		Help:      "Current depth of workqueue",
	}, []string{nameLabel})

	adds = prometheus.NewCounterVec(prometheus.CounterOpts{
		Subsystem: WorkQueueSubsystem,
		Name:      AddsKey,
		Help:      "Total number of adds handled by workqueue",
	}, []string{nameLabel})

	latency = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Subsystem: WorkQueueSubsystem,
		Name:      QueueLatencyKey,
		Help:      "How long in seconds an item stays in workqueue before being requested.",
		// custom buckets - [1s, 2s, 4s, 8s, 16s, 32s, 64s, 128s, 256s(~4min), 512s(~8min), 1024s(~17min), 2048 (~34min), 4096(~68min), 8192(~136min), 16384(~272min)  +Inf]
		Buckets: prometheus.ExponentialBuckets(1, 2, 15),
	}, []string{nameLabel})

	workDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Subsystem: WorkQueueSubsystem,
		Name:      WorkDurationKey,
		Help:      "How long in seconds processing an item from workqueue takes.",
		// custom buckets - [1s, 2s, 4s, 8s, 16s, 32s, 64s, 128s, 256s(~4min), 512s(~8min), 1024s(~17min), 2048 (~34min), 4096(~68min), +Inf]
		Buckets: prometheus.ExponentialBuckets(1, 2, 13),
	}, []string{nameLabel})

	unfinished = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Subsystem: WorkQueueSubsystem,
		Name:      UnfinishedWorkKey,
		Help: "How many seconds of work has done that " +
			"is in progress and hasn't been observed by work_duration. Large " +
			"values indicate stuck threads. One can deduce the number of stuck " +
			"threads by observing the rate at which this increases.",
	}, []string{nameLabel})

	longestRunningProcessor = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Subsystem: WorkQueueSubsystem,
		Name:      LongestRunningProcessorKey,
		Help: "How many seconds has the longest running " +
			"processor for workqueue been running.",
	}, []string{nameLabel})

	retries = prometheus.NewCounterVec(prometheus.CounterOpts{
		Subsystem: WorkQueueSubsystem,
		Name:      RetriesKey,
		Help:      "Total number of retries handled by workqueue",
	}, []string{nameLabel})
)

// workqueueMetricsProvider implements workqueue.MetricsProvider interface
type workqueueMetricsProvider struct {
}

func init() {
	prometheus.MustRegister(depth, adds, latency, workDuration, unfinished, longestRunningProcessor, retries)
	workqueue.SetProvider(workqueueMetricsProvider{})
}

func (workqueueMetricsProvider) NewDepthMetric(name string) workqueue.GaugeMetric {
	return depth.WithLabelValues(name)
}

func (workqueueMetricsProvider) NewAddsMetric(name string) workqueue.CounterMetric {
	return adds.WithLabelValues(name)
}

func (workqueueMetricsProvider) NewLatencyMetric(name string) workqueue.HistogramMetric {
	return latency.WithLabelValues(name)
}

func (workqueueMetricsProvider) NewWorkDurationMetric(name string) workqueue.HistogramMetric {
	return workDuration.WithLabelValues(name)
}

func (workqueueMetricsProvider) NewUnfinishedWorkSecondsMetric(name string) workqueue.SettableGaugeMetric {
	return unfinished.WithLabelValues(name)
}

func (workqueueMetricsProvider) NewLongestRunningProcessorSecondsMetric(name string) workqueue.SettableGaugeMetric {
	return longestRunningProcessor.WithLabelValues(name)
}

func (workqueueMetricsProvider) NewRetriesMetric(name string) workqueue.CounterMetric {
	return retries.WithLabelValues(name)
}
