package priorities

import (
	"math"
	ketiresource "openmcp/openmcp/openmcp-scheduler/pkg/resourceinfo"

	v1 "k8s.io/api/core/v1"
)

type QosPriority struct {
	prescoring   map[string]int64
	betweenScore int64
}

const (
	minScore int64 = 0
	midScore int64 = (minScore + maxScore) / 2
	maxScore int64 = 10
)

func (pl *QosPriority) Name() string {
	return "QosPriority"
}
func (pl *QosPriority) PreScore(pod *ketiresource.Pod, clusterInfo *ketiresource.Cluster, check bool) int64 {
	var clusterScore int64

	for _, node := range clusterInfo.Nodes {
		var nodeScore int64
		for _, pod := range node.Pods {

			// get PodQOSClass from v1.Pod
			qos := pod.Pod.Status.QOSClass

			switch qos {
			case v1.PodQOSGuaranteed:
				nodeScore += minScore
			case v1.PodQOSBurstable:
				nodeScore += midScore
			case v1.PodQOSBestEffort:
				nodeScore += maxScore
			}
		}
		node.NodeScore = nodeScore
		clusterScore += nodeScore
	}
	if !check {
		if len(pl.prescoring) == 0 {
			pl.prescoring = make(map[string]int64)
		}
		pl.prescoring[clusterInfo.ClusterName] = clusterScore
	} else {
		pl.betweenScore = pl.prescoring[clusterInfo.ClusterName] - int64(clusterScore)
		pl.betweenScore = int64(math.Abs(float64(pl.betweenScore)))

	}
	return clusterScore
}

func (pl *QosPriority) Score(pod *ketiresource.Pod, clusterInfo *ketiresource.Cluster, replicas int32, clustername string) int64 {
	if clustername == clusterInfo.ClusterName {
		score := pl.prescoring[clusterInfo.ClusterName] - pl.betweenScore
		return score
	}
	score := pl.prescoring[clusterInfo.ClusterName]
	return score

}
