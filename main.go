// +build linux

/*
http://www.apache.org/licenses/LICENSE-2.0.txt


Copyright 2016 Intel Corporation

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at`

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"github.com/intelsdi-x/pulse-collector-libcontainer/cgroups"
	"github.com/intelsdi-x/snap/control/plugin"
	"os"
	"strings"
)

func main() {

	if cgplugin, err := cgroups.NewCgroups(); err != nil {
		fmt.Printf("Error: %s \n", err)
		return
	} else {
		plugin.Start(plugin.NewPluginMeta(
			cgroups.NS_VENDOR,
			cgroups.VERSION,
			plugin.CollectorPluginType,
			[]string{},
			[]string{plugin.SnapGOBContentType},
			plugin.ConcurrencyCount(1)),
			cgplugin,
			os.Args[1])
		mts := []plugin.PluginMetricType{}
		mts = append(mts, plugin.PluginMetricType{
			Namespace_: strings.Split("intel/linux/cgroups/user_slice/*/cpu_stats/cpu_usage/total_usage", "/"),
			Tags_:      map[string]string{}})
	}

}
