// ++build unit

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
	"errors"
	"strings"
	"testing"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/core"
	. "github.com/smartystreets/goconvey/convey"
)

var mountPointMocks = []mountpoint{
	{subsystem: "blkio", subsystemNs: "blkio_stats", namespace: "", path: "/sys/fs/cgroup/blkio"},
	{subsystem: "cpu", subsystemNs: "cpu_stats", namespace: "", path: "/sys/fs/cgroup/cpu"},
	{subsystem: "hugetlb", subsystemNs: "hugetlb_stats", namespace: "", path: "/sys/fs/cgroup/hugetlb"},
	{subsystem: "memory", subsystemNs: "memory_stats", namespace: "", path: "/sys/fs/cgroup/memory"},
	{subsystem: "pids", subsystemNs: "pids_stats", namespace: "system_slice", path: "/sys/fs/cgroup/pids/system.slice"},
}

var nameSpaceMocks = []plugin.MetricType{
	{Namespace_: core.NewNamespace(strings.Split("intel/linux/cgroups/blkio_stats/io_service_bytes_recursive/0/major", "/")...)},
	{Namespace_: core.NewNamespace(strings.Split("intel/linux/cgroups/cpu_stats/cpu_usage/total_usage", "/")...)},
	{Namespace_: core.NewNamespace(strings.Split("intel/linux/cgroups/hugetlb_stats/1GB/max_usage", "/")...)},
	{Namespace_: core.NewNamespace(strings.Split("intel/linux/cgroups/hugetlb_stats/1GB/usage", "/")...)},
	{Namespace_: core.NewNamespace(strings.Split("intel/linux/cgroups/memory_stats/usage/usage", "/")...)},
	{Namespace_: core.NewNamespace(strings.Split("intel/linux/cgroups/pids_stats/system_slice/current", "/")...)},
	{Namespace_: core.NewNamespace(strings.Split("intel/linux/cgroups/pids_stats/system_slice/limit", "/")...)},
}

func TestNewCgroups(t *testing.T) {
	cgplugin := new(cgroups)
	err := errors.New("")

	Convey("NewCgroups should not panic", t, func() {

		So(func() { cgplugin, err = NewCgroups(true) }, ShouldNotPanic)

		Convey("NewCgroups should not return error", func() {
			So(err, ShouldBeNil)
		})

		Convey("NewCgroups should not return nil", func() {
			So(cgplugin, ShouldNotBeNil)
		})
	})
}

func TestCollectMetrics(t *testing.T) {
	cgplugin, _ := NewCgroups(true)

	for _, mpMock := range mountPointMocks {
		mp, _ := cgplugin.newMountPoint(mpMock.subsystem, mpMock.namespace, mpMock.path)
		cgplugin.addMountPoint(mp)
	}

	metrics := []plugin.MetricType{}
	err := errors.New("")

	Convey("CollectMetrics should not panic", t, func() {

		So(func() { metrics, err = cgplugin.CollectMetrics(nameSpaceMocks) }, ShouldNotPanic)

		Convey("CollectMetrics should not return error", func() {
			So(err, ShouldBeNil)
		})

		Convey("CollectMetrics should not return nil", func() {
			So(metrics, ShouldNotBeNil)
		})

		Convey("CollectMetrics should not return empty value", func() {
			So(metrics, ShouldNotBeEmpty)
		})

		Convey("CollectMetrics should return 7 items", func() {
			So(len(metrics), ShouldEqual, 7)
		})
	})
}

func TestGetMetricTypes(t *testing.T) {
	cgplugin, _ := NewCgroups(true)
	cfg := plugin.NewPluginConfigType()

	for _, mpMock := range mountPointMocks {
		mp, _ := cgplugin.newMountPoint(mpMock.subsystem, mpMock.namespace, mpMock.path)
		cgplugin.addMountPoint(mp)
	}

	metrics := []plugin.MetricType{}
	err := errors.New("")

	Convey("GetMetricTypes should not panic", t, func() {

		So(func() { metrics, err = cgplugin.GetMetricTypes(cfg) }, ShouldNotPanic)

		Convey("GetMetricTypes should not return error", func() {
			So(err, ShouldBeNil)
		})

		Convey("GetMetricTypes should not return nil", func() {
			So(metrics, ShouldNotBeNil)
		})

		Convey("GetMetricTypes should not return empty value", func() {
			So(metrics, ShouldNotBeEmpty)
		})

		Convey("GetMetricTypes should return more than 50 items", func() {
			So(len(metrics), ShouldBeGreaterThan, 50)
		})

	})
}

func TestGetMountPoints(t *testing.T) {
	cgplugin, _ := NewCgroups(true)

	for _, mpMock := range mountPointMocks {
		mp, _ := cgplugin.newMountPoint(mpMock.subsystem, mpMock.namespace, mpMock.path)
		cgplugin.addMountPoint(mp)
	}

	err := errors.New("")

	Convey("GetMountPoints should not panic", t, func() {

		So(func() { err = cgplugin.getMountPoints() }, ShouldNotPanic)

		Convey("GetMountPoints should not return error", func() {
			So(err, ShouldBeNil)
		})

		Convey("GetMountPoints should return more than 50 items", func() {
			So(len(cgplugin.mountPoints), ShouldBeGreaterThan, 50)
		})

	})
}
