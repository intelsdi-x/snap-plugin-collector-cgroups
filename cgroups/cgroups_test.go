// +build unit

/*
http://www.apache.org/licenses/LICENSE-2.0.txt


Copyright 2016 Intel Corporation

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

package cgroups

import (
	//"encoding/json"
	//. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
	//"os"
	//"path/filepath"
	"testing"
	//"time"
	//lcgroups "github.com/opencontainers/runc/libcontainer/cgroups"
	//"github.com/intelsdi-x/snap-plugin-utilities/ns"
	"fmt"
	//"strings"
)

type MockCgstat struct {
	mock.Mock
	cgroupstat
}

func handleThat(s *cgroupstat) {
	fmt.Println(s.discoverMountPoints())
}

func TestCgroups_CollectMetrics(t *testing.T) {
	fmt.Println("Hello")
	s := NewCgroupstat()
	handleThat(s)
	handleThat(s)
	//stats := lcgroups.NewStats()
	//namespace := []string {}
	//ns.FromCompositeObject(stats, "", &namespace)
	//{
	//	jsonBytes, _ := json.MarshalIndent(namespace, "  ", "- ")
	//	fmt.Printf("%s \n", string(jsonBytes))
	//}
	//{
	//	jsonBytes, _ := json.MarshalIndent(stats, "  ", "- ")
	//	fmt.Printf("2/ %s \n", string(jsonBytes))
	//}
	//val := ns.GetValueByNamespace(stats, strings.Split("cpu_stats/cpu_usage", "/"))
	//var cpuUsage *lcgroups.CpuUsage
	//cpuUsage = val.(*lcgroups.CpuUsage)
	//cpuUsage.TotalUsage = 3218
	//{
	//	jsonBytes, _ := json.MarshalIndent(stats, "  ", "- ")
	//	fmt.Printf("3/ %s \n", string(jsonBytes))
	//}
}
