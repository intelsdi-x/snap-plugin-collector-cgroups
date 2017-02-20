// +build linux

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
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
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

	subsystemID      = 3
	mountPointPathID = 4
)

type mountpoint struct {
	subsystem string
	namespace string
	metrics   []string
	path      string
	stats     *libcgroups.Stats
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
	var data interface{}
	var metrics []plugin.MetricType

	for i, metricType := range metricTypes {
		matchAny := false

		for _, mountPoint := range c.mountPoints {

			// Cgroup subsystem match with namespace subsystem
			if metricType.Namespace()[subsystemID].Value == mountPoint.subsystem {

				// Cgroup mountpoint match with namespace mountpoint
				if matchSlice(strings.Split(metricType.Namespace()[mountPointPathID].Value, "/"), strings.Split(mountPoint.namespace, "/")) {

					mountPoint.getMetrics()

					for _, mpMetric := range mountPoint.metrics {
						// Get the metric part from namespace
						nsMetric := strings.Join(metricType.Namespace().Strings()[mountPointPathID+1:], "/")

						// Remove percpu_usage value suffix for matching
						if strings.HasPrefix(nsMetric, "cpu_usage/percpu_usage/") {
							nsMetric = strings.TrimSuffix(nsMetric, "/value")
						}

						// Mountpoint metric match with namespace metric part
						if matchSlice(strings.Split(nsMetric, "/"), strings.Split(mpMetric, "/")) {

							// Asterisk exists only in empty containers on static metrics list
							if strings.Contains(mpMetric, "*") {
								data = nil
							} else {
								// CpuStats structure can have only one JSON field name and it is "cpu_stats" even for "cpuacct_stats", so replace is needed
								subsystemEntry := strings.Replace(mountPoint.subsystem, "cpuacct_stats", "cpu_stats", -1)

								// Get metric data
								data = ns.GetValueByNamespace(mountPoint.stats, strings.Split(strings.Join([]string{subsystemEntry, mpMetric}, "/"), "/"))
							}

							metricTypes[i].Namespace_ = mountPoint.getFullNamespace(mountPoint.namespace, mpMetric)
							metricTypes[i].Data_ = data
							metricTypes[i].Timestamp_ = time.Now()
							metricTypes[i].Version_ = VERSION

							assignMetricMeta(&metricTypes[i], allMetrics)

							metrics = append(metrics, metricTypes[i])
							matchAny = true
						}
					}
				}
			}
		}

		if !matchAny {
			metrics = append(metrics, metricType)
		}
	}

	return metrics, nil
}

func (c *cgroups) GetMetricTypes(plugin.ConfigType) ([]plugin.MetricType, error) {
	metrics := []plugin.MetricType{}
	exists := make(map[string]bool)

	for _, mountPoint := range c.mountPoints {
		mountPoint.getMetrics()

		for _, metric := range mountPoint.metrics {
			// Create namespace with dynamic mount path for dynamic metric name
			metricNamespace := mountPoint.getFullNamespace("*", dynamicToAsterisk(metric))

			// Ignore duplicates
			if !exists[metricNamespace.String()] {
				mt := plugin.MetricType{
					Namespace_: metricNamespace,
					Version_:   VERSION,
				}
				metrics = append(metrics, mt)

				exists[metricNamespace.String()] = true
			}
		}
	}
	return metrics, nil
}

func (c *cgroups) GetConfigPolicy() (*cpolicy.ConfigPolicy, error) {
	return cpolicy.New(), nil
}

func (c *cgroups) newMountPoint(subsystem string, namespace string, path string) (*mountpoint, error) {
	if namespace == "" {
		namespace = "root"
	}

	return &mountpoint{subsystem: subsystem, namespace: namespace, path: path}, nil
}

func (c *cgroups) addMountPoint(mountPoint *mountpoint) {
	c.mountPoints = append(c.mountPoints, mountPoint)
}

// getMountPoints discovers all cgroup mountpoints recursively
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
				mountPoint, err := c.newMountPoint(subsys+"_stats", namespace, path)
				if err == nil {
					c.addMountPoint(mountPoint)
				}
			}
			return nil
		})
	}
	return nil
}

// getMetrics gets all metrics for desired mountpoint from OS
func (mp *mountpoint) getMetrics() error {
	// Get stats from mountpoint
	manager := libcgroupsfs.Manager{Paths: make(map[string]string)}
	manager.Paths[strings.TrimSuffix(mp.subsystem, "_stats")] = mp.path
	stats, err := manager.GetStats()

	if err != nil {
		return err
	}

	// Find all cgroup stat metrics for mountpoint
	metricList := []string{}

	switch mp.subsystem {
	case "blkio_stats":
		ns.FromCompositeObject(stats.BlkioStats, "", &metricList, ns.InspectEmptyContainers(ns.AlwaysTrue))
	case "cpu_stats", "cpuacct_stats":
		ns.FromCompositeObject(stats.CpuStats, "", &metricList, ns.InspectEmptyContainers(ns.AlwaysTrue))
	case "hugetlb_stats":
		ns.FromCompositeObject(stats.HugetlbStats, "", &metricList, ns.InspectEmptyContainers(ns.AlwaysTrue))
	case "memory_stats":
		ns.FromCompositeObject(stats.MemoryStats, "", &metricList, ns.InspectEmptyContainers(ns.AlwaysTrue))
	case "pids_stats":
		ns.FromCompositeObject(stats.PidsStats, "", &metricList, ns.InspectEmptyContainers(ns.AlwaysTrue))
	}

	mp.stats = stats
	mp.metrics = metricList

	return nil
}

// getFullNamespace constructs full snap namespace for desired mountpoint with provided mountpoint path and metric name
func (mp *mountpoint) getFullNamespace(mountpoint string, metric string) (result core.Namespace) {
	if mountpoint == "*" {
		result = core.NewNamespace(NS_VENDOR, NS_CLASS, NS_PLUGIN, mp.subsystem).AddDynamicElement("mountpoint", "Part of a mountpoint path")
	} else {
		result = core.NewNamespace(NS_VENDOR, NS_CLASS, NS_PLUGIN, mp.subsystem).AddStaticElements(strings.Split(mountpoint, "/")...)
	}

	if metric != "" {
		metricSplit := strings.Split(metric, "/")
		for i, metricEntry := range metricSplit {
			if metricEntry == "*" {
				result = result.AddDynamicElement("metric_id_"+strconv.Itoa(i), "Metric ID")
			} else {
				result = result.AddStaticElement(metricEntry)
			}
		}

		// Fix "ends with asterisk is not allowed" error for percpu_usage
		if strings.HasPrefix(metric, "cpu_usage/percpu_usage/") {
			result = result.AddStaticElement("value")
		}
	}

	return result
}

// dynamicToAsterisk replaces all dynamic namespace entries (like digit-only ID or capacity in hugetlb) to asterisks
func dynamicToAsterisk(ns string) string {
	digitEntry := regexp.MustCompile(`(^|\/)(\d+(B|KB|MB|GB)?)($|\/)`)
	return digitEntry.ReplaceAllString(ns, `$1*$4`)
}

// matchSlice matches 2 slices with asterisk support
func matchSlice(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] && a[i] != "*" && b[i] != "*" {
			return false
		}
	}

	return true
}

// assignMetricMeta assigs metadata to metric using predefined metadata slice
func assignMetricMeta(mt *plugin.MetricType, allMetrics []metric) {
	for _, metricMeta := range allMetrics {
		fmt.Println("Match ", metricMeta.ns, "with", mt.Namespace().String())
		if matchSlice(strings.Split(metricMeta.ns, "/"), strings.Split(mt.Namespace().String(), "/")) {
			mt.Description_ = metricMeta.description
			mt.Unit_ = metricMeta.unit
			break
		}
	}
}
