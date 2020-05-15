// Copyright 2020 The Kubernetes Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package decode

import (
	"fmt"
	"math"
	"time"

	"k8s.io/klog"

	"cluster-metric-collector/pkg/stats"
	"cluster-metric-collector/pkg/storage"
	"k8s.io/apimachinery/pkg/api/resource"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
)

func DecodeBatch(summary *stats.Summary) (*storage.MetricsBatch, error) {
	fmt.Println("Func DecodeBatch Called")
	res := &storage.MetricsBatch{
		IP:   summary.IP,
		Node: storage.NodeMetricsPoint{},
		Pods: make([]storage.PodMetricsPoint, len(summary.Pods)),
	}

	var errs []error
	errs = append(errs, decodeNodeStats(&summary.Node, &res.Node)...)
	if len(errs) != 0 {
		// if we had errors providing node metrics, discard the data point
		// so that we don't incorrectly report metric values as zero.
	}

	num := 0
	for _, pod := range summary.Pods {
		podErrs := decodePodStats(&pod, &res.Pods[num])
		errs = append(errs, podErrs...)
		if len(podErrs) != 0 {
			// NB: we explicitly want to discard pods with partial results, since
			// the horizontal pod autoscaler takes special action when a pod is missing
			// metrics (and zero CPU or memory does not count as "missing metrics")

			// we don't care if we reuse slots in the result array,
			// because they get completely overwritten in decodePodStats
			continue
		}
		num++
	}
	res.Pods = res.Pods[:num]
	return res, utilerrors.NewAggregate(errs)
}

func decodeNodeStats(nodeStats *stats.NodeStats, target *storage.NodeMetricsPoint) []error {
	fmt.Println("Func decodeNodeStats Called")
	timestamp, err := getScrapeTimeNode(nodeStats.CPU, nodeStats.Memory, nodeStats.Network, nodeStats.Fs)
	if err != nil {
		// if we can't get a timestamp, assume bad data in general
		return []error{fmt.Errorf("unable to get valid timestamp for metric point for node %q, discarding data: %v", nodeStats.NodeName, err)}
	}
	*target = storage.NodeMetricsPoint{
		Name: nodeStats.NodeName,
		MetricsPoint: storage.MetricsPoint{
			Timestamp: timestamp,
		},
	}
	var errs []error
	if err := decodeCPU(&target.CpuUsage, nodeStats.CPU); err != nil {
		errs = append(errs, fmt.Errorf("unable to get CPU for node %q, discarding data: %v", nodeStats.NodeName, err))
	}
	if err := decodeMemory(&target.MemoryUsage, nodeStats.Memory); err != nil {
		errs = append(errs, fmt.Errorf("unable to get memory for node %q, discarding data: %v", nodeStats.NodeName, err))
	}
	if err := decodeNetworkRX(&target.NetworkRxUsage, nodeStats.Network); err != nil {
		errs = append(errs, fmt.Errorf("unable to get Network for node %q, discarding data: %v", nodeStats.NodeName, err))
	}
	if err := decodeNetworkTX(&target.NetworkTxUsage, nodeStats.Network); err != nil {
		errs = append(errs, fmt.Errorf("unable to get Network for node %q, discarding data: %v", nodeStats.NodeName, err))
	}
	if err := decodeFs(&target.FsUsage, nodeStats.Fs); err != nil {
		errs = append(errs, fmt.Errorf("unable to get FS for node %q, discarding data: %v", nodeStats.NodeName, err))
	}
	return errs
}

func decodePodStats(podStats *stats.PodStats, target *storage.PodMetricsPoint) []error {
	fmt.Println("Func decodePodStats Called")

	timestamp, err := getScrapeTimePod(podStats.CPU, podStats.Memory)
	if err != nil {
		// if we can't get a timestamp, assume bad data in general
		return []error{fmt.Errorf("unable to get valid timestamp for metric point for pod %q, discarding data: %v", podStats.PodRef.Name, err)}
	}

	// completely overwrite data in the target
	*target = storage.PodMetricsPoint{
		Name:      podStats.PodRef.Name,
		Namespace: podStats.PodRef.Namespace,
		MetricsPoint: storage.MetricsPoint{
			Timestamp: timestamp,
		},
		//Containers: make([]storage.ContainerMetricsPoint, len(podStats.Containers)),
	}

	var errs []error
	if err := decodeCPU(&target.CpuUsage, podStats.CPU); err != nil {
		errs = append(errs, fmt.Errorf("unable to get CPU for pod %q, discarding data: %v", podStats.PodRef.Name, err))
	}
	if err := decodeMemory(&target.MemoryUsage, podStats.Memory); err != nil {
		errs = append(errs, fmt.Errorf("unable to get memory for pod %q, discarding data: %v", podStats.PodRef.Name, err))
	}
	if err := decodeNetworkRX(&target.NetworkRxUsage, podStats.Network); err != nil {
		errs = append(errs, fmt.Errorf("unable to get network RX for pod %q, discarding data: %v", podStats.PodRef.Name, err))
	}
	if err := decodeNetworkTX(&target.NetworkTxUsage, podStats.Network); err != nil {
		errs = append(errs, fmt.Errorf("unable to get network TX for pod %q, discarding data: %v", podStats.PodRef.Name, err))
	}
	if err := decodeFs(&target.FsUsage, podStats.EphemeralStorage); err != nil {
		errs = append(errs, fmt.Errorf("unable to get Fs for pod %q, discarding data: %v", podStats.PodRef.Name, err))
	}

	return errs
}

func decodeCPU(target *resource.Quantity, cpuStats *stats.CPUStats) error {
	if cpuStats == nil || cpuStats.UsageNanoCores == nil {
		return fmt.Errorf("missing cpu usage metric")
	}

	*target = *uint64Quantity(*cpuStats.UsageNanoCores, -9)
	return nil
}

func decodeMemory(target *resource.Quantity, memStats *stats.MemoryStats) error {
	if memStats == nil || memStats.WorkingSetBytes == nil {
		return fmt.Errorf("missing memory usage metric")
	}

	*target = *uint64Quantity(*memStats.WorkingSetBytes, 0)
	target.Format = resource.BinarySI

	return nil
}
func decodeNetworkRX(target *resource.Quantity, netStats *stats.NetworkStats) error {
	if netStats == nil || netStats.Interfaces[0].RxBytes == nil {
		return fmt.Errorf("missing network RX usage metric")
	}
	var RX_Usage uint64 = 0
	for _, Interface := range netStats.Interfaces {
		RX_Usage = RX_Usage + *Interface.RxBytes
	}

	*target = *uint64Quantity(RX_Usage, 0)
	target.Format = resource.BinarySI

	return nil
}
func decodeNetworkTX(target *resource.Quantity, netStats *stats.NetworkStats) error {
	if netStats == nil || netStats.Interfaces[0].TxBytes == nil {
		return fmt.Errorf("missing network RX usage metric")
	}
	var TX_Usage uint64 = 0
	for _, Interface := range netStats.Interfaces {
		TX_Usage = TX_Usage + *Interface.TxBytes
	}

	*target = *uint64Quantity(TX_Usage, 0)
	target.Format = resource.BinarySI

	return nil
}
func decodeFs(target *resource.Quantity, FsStats *stats.FsStats) error {
	if FsStats == nil || FsStats.UsedBytes == nil {
		return fmt.Errorf("missing memory usage metric")
	}

	*target = *uint64Quantity(*FsStats.UsedBytes, 0)
	target.Format = resource.BinarySI

	return nil
}
func getScrapeTimePod(cpu *stats.CPUStats, memory *stats.MemoryStats) (time.Time, error) {
	fmt.Println("Func getScrapeTime Called")
	// Ensure we get the earlier timestamp so that we can tell if a given data
	// point was tainted by pod initialization.

	var earliest *time.Time
	if cpu != nil && !cpu.Time.IsZero() && (earliest == nil || earliest.After(cpu.Time.Time)) {
		earliest = &cpu.Time.Time
	}

	if memory != nil && !memory.Time.IsZero() && (earliest == nil || earliest.After(memory.Time.Time)) {
		earliest = &memory.Time.Time
	}

	if earliest == nil {
		return time.Time{}, fmt.Errorf("no non-zero timestamp on either CPU or memory")
	}

	return *earliest, nil
}
func getScrapeTimeNode(cpu *stats.CPUStats, memory *stats.MemoryStats, network *stats.NetworkStats, fs *stats.FsStats) (time.Time, error) {
	fmt.Println("Func getScrapeTime Called")
	// Ensure we get the earlier timestamp so that we can tell if a given data
	// point was tainted by pod initialization.

	var earliest *time.Time
	if cpu != nil && !cpu.Time.IsZero() && (earliest == nil || earliest.After(cpu.Time.Time)) {
		earliest = &cpu.Time.Time
	}

	if memory != nil && !memory.Time.IsZero() && (earliest == nil || earliest.After(memory.Time.Time)) {
		earliest = &memory.Time.Time
	}

	if network != nil && !network.Time.IsZero() && (earliest == nil || earliest.After(network.Time.Time)) {
		earliest = &network.Time.Time
	}

	if fs != nil && !fs.Time.IsZero() && (earliest == nil || earliest.After(fs.Time.Time)) {
		earliest = &fs.Time.Time
	}

	if earliest == nil {
		return time.Time{}, fmt.Errorf("no non-zero timestamp on either CPU or memory")
	}

	return *earliest, nil
}

// uint64Quantity converts a uint64 into a Quantity, which only has constructors
// that work with int64 (except for parse, which requires costly round-trips to string).
// We lose precision until we fit in an int64 if greater than the max int64 value.
func uint64Quantity(val uint64, scale resource.Scale) *resource.Quantity {
	// easy path -- we can safely fit val into an int64
	if val <= math.MaxInt64 {
		return resource.NewScaledQuantity(int64(val), scale)
	}

	klog.V(1).Infof("unexpectedly large resource value %v, loosing precision to fit in scaled resource.Quantity", val)

	// otherwise, lose an decimal order-of-magnitude precision,
	// so we can fit into a scaled quantity
	return resource.NewScaledQuantity(int64(val/10), resource.Scale(1)+scale)
}
