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
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/intelsdi-x/snap-plugin-utilities/ns"
	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"
	"github.com/intelsdi-x/snap/core"
	libcgroups "github.com/opencontainers/runc/libcontainer/cgroups"
	libcgroupsfs "github.com/opencontainers/runc/libcontainer/cgroups/fs"
)

const (
	// namespace vendor prefix
	NS_VENDOR = "intel"
	// namespace os prefix
	NS_CLASS = "linux"
	// namespace plugin name
	NS_PLUGIN = "cgroups"
	// version of plugin
	VERSION = 2
)

type mountpoint struct {
	subsystem   string
	subsystemNs string
	namespace   string
	metrics     []string
	path        string
	stats       *libcgroups.Stats
}

type cgroups struct {
	mountPoints []*mountpoint
}

func NewCgroups(test bool) (*cgroups, error) {
	c := &cgroups{}

	// Do not get mountpoints automatically when testing
	if !test {
		c.getMountPoints()
	}

	return c, nil
}

func (c *cgroups) CollectMetrics(metricTypes []plugin.MetricType) ([]plugin.MetricType, error) {
	for i, metricType := range metricTypes {
		for _, mountPoint := range c.mountPoints {

			// Mount point part of metric namespace match mount point namespace
			if strings.HasPrefix(metricType.Namespace().String(), "/"+strings.Join(mountPoint.getFullNamespace(""), "/")) {

				// Get cgroup metrics for matching mount point
				c.getMountPointMetrics(mountPoint)

				for _, metric := range mountPoint.metrics {

					// Metric namespace match mount point cgroup metric
					if metricType.Namespace().String() == "/"+strings.Join(mountPoint.getFullNamespace(metric), "/") {
						data := ns.GetValueByNamespace(mountPoint.stats, strings.Split(strings.Join([]string{mountPoint.subsystemNs, metric}, "/"), "/"))

						metricTypes[i].Data_ = data
						metricTypes[i].Timestamp_ = time.Now()
						metricTypes[i].Version_ = VERSION
					}
				}
			}
		}
	}

	return metricTypes, nil
}

func (c *cgroups) GetMetricTypes(plugin.ConfigType) ([]plugin.MetricType, error) {
	metrics := []plugin.MetricType{}

	for _, mountPoint := range c.mountPoints {
		c.getMountPointMetrics(mountPoint)

		for _, metric := range mountPoint.metrics {
			mt := plugin.MetricType{
				Namespace_: core.NewNamespace(mountPoint.getFullNamespace(metric)...),
				Version_:   VERSION,
			}

			metrics = append(metrics, mt)
		}
	}
	return metrics, nil
}

func (c *cgroups) GetConfigPolicy() (*cpolicy.ConfigPolicy, error) {
	return cpolicy.New(), nil
}

func (c *cgroups) newMountPoint(subsystem string, namespace string, path string) (*mountpoint, error) {
	return &mountpoint{subsystem: subsystem, subsystemNs: subsystem + "_stats", namespace: namespace, path: path}, nil
}

func (c *cgroups) addMountPoint(mountPoint *mountpoint) {
	c.mountPoints = append(c.mountPoints, mountPoint)
}

func (c *cgroups) getMountPoints() error {
	// Get all subsystems
	subsysAll, err := libcgroups.GetAllSubsystems()

	if err != nil {
		return err
	}

	// Get root mountpoints for subsystems
	for _, subsys := range subsysAll {
		mountPointPath, err := libcgroups.FindCgroupMountpoint(subsys)

		if err != nil {
			return err
		}

		// Get inner mountpoints
		filepath.Walk(mountPointPath, func(path string, f os.FileInfo, err error) error {
			if f.IsDir() {
				// Generate mountpoint namespace
				nsSplit := strings.Split(strings.TrimPrefix(strings.TrimPrefix(path, mountPointPath), "/"), "/")
				for i, part := range nsSplit {
					nsSplit[i] = ns.ReplaceNotAllowedCharsInNamespacePart(part)
				}
				namespace := strings.Join(nsSplit, "/")

				// Create mountpoint
				mountPoint, err := c.newMountPoint(subsys, namespace, path)
				if err == nil {
					c.addMountPoint(mountPoint)
				}
			}
			return nil
		})
	}
	return nil
}

func (c *cgroups) getMountPointMetrics(mountPoint *mountpoint) error {
	// Get stats from mountpoint
	manager := libcgroupsfs.Manager{Paths: make(map[string]string)}
	manager.Paths[mountPoint.subsystem] = mountPoint.path
	stats, err := manager.GetStats()

	if err != nil {
		return err
	}

	// Find all cgroup stat metrics for mountpoint
	metrics := []string{}
	switch mountPoint.subsystem {
	case "blkio":
		ns.FromCompositeObject(stats.BlkioStats, "", &metrics, ns.InspectEmptyContainers(ns.AlwaysFalse))
	case "cpu":
		ns.FromCompositeObject(stats.CpuStats, "", &metrics, ns.InspectEmptyContainers(ns.AlwaysFalse))
	case "hugetlb":
		ns.FromCompositeObject(stats.HugetlbStats, "", &metrics, ns.InspectEmptyContainers(ns.AlwaysFalse))
	case "memory":
		ns.FromCompositeObject(stats.MemoryStats, "", &metrics, ns.InspectEmptyContainers(ns.AlwaysFalse))
	case "pids":
		ns.FromCompositeObject(stats.PidsStats, "", &metrics, ns.InspectEmptyContainers(ns.AlwaysFalse))
	}

	mountPoint.stats = stats
	mountPoint.metrics = metrics

	return nil
}

func (m *mountpoint) getFullNamespace(metric string) []string {
	result := []string{}

	// Skip empty mountpoint namespace
	if m.namespace == "" {
		result = strings.Split(strings.Join([]string{NS_VENDOR, NS_CLASS, NS_PLUGIN, m.subsystemNs, metric}, "/"), "/")
	} else {
		result = strings.Split(strings.Join([]string{NS_VENDOR, NS_CLASS, NS_PLUGIN, m.subsystemNs, m.namespace, metric}, "/"), "/")
	}

	return result
}
