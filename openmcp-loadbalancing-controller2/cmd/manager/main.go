/*
Copyright 2018 The Multicluster-Controller Authors.

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

package main

import (
	//"sync"
	//"time"

	"openmcp/openmcp/openmcp-loadbalancing-controller2/pkg/DestinationRuleWeight"
	"sync"
	"time"

	"github.com/oschwald/geoip2-golang"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

//var GeoDB, GeoErr = geoip2.Open("/root/GeoLite2-City.mmdb")
var GeoDB, GeoErr = geoip2.Open("/root/dbip-city-lite-2021-07.mmdb")

func main() {
	var wg sync.WaitGroup
	wg.Add(3)

	go reverseProxy()
	go serviceMeshController()
	//go OpenMCPVirtualService.SyncWeight()
	time.Sleep(time.Second * 2)
	go DestinationRuleWeight.AnalyticWeight()

	wg.Wait()

}
