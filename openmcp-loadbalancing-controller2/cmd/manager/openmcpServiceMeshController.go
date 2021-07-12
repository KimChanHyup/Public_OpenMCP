package main

import (
	"log"
	"openmcp/openmcp/openmcp-loadbalancing-controller2/pkg/DestinationRule"
	"openmcp/openmcp/openmcp-loadbalancing-controller2/pkg/OpenMCPVirtualService"

	"admiralty.io/multicluster-controller/pkg/cluster"
	"admiralty.io/multicluster-controller/pkg/manager"

	"openmcp/openmcp/omcplog"
	"openmcp/openmcp/util/clusterManager"
	"openmcp/openmcp/util/controller/logLevel"
	"openmcp/openmcp/util/controller/reshape"
)

func serviceMeshController() {
	logLevel.KetiLogInit()

	for {
		omcplog.V(2).Info("Start OpenMCPServiceMeshController")

		cm := clusterManager.NewClusterManager()

		host_ctx := "openmcp"
		namespace := "openmcp"

		host_cfg := cm.Host_config

		live := cluster.New(host_ctx, host_cfg, cluster.Options{})

		ghosts := []*cluster.Cluster{}

		for _, ghost_cluster := range cm.Cluster_list.Items {
			ghost_ctx := ghost_cluster.Name
			ghost_cfg := cm.Cluster_configs[ghost_ctx]

			ghost := cluster.New(ghost_ctx, ghost_cfg, cluster.Options{})
			ghosts = append(ghosts, ghost)
		}

		vs_cont, _ := OpenMCPVirtualService.NewController(live, ghosts, namespace, cm)
		dr_cont, _ := DestinationRule.NewController(live, ghosts, namespace, cm)
		reshape_cont, _ := reshape.NewController(live, ghosts, namespace, cm)
		loglevel_cont, _ := logLevel.NewController(live, ghosts, namespace)

		m := manager.New()
		m.AddController(vs_cont)
		m.AddController(dr_cont)
		m.AddController(reshape_cont)
		m.AddController(loglevel_cont)

		stop := reshape.SetupSignalHandler()

		if err := m.Start(stop); err != nil {
			log.Fatal(err)
		}
	}
}
