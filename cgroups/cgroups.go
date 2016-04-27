package cgroups

import (
	"fmt"
	"github.com/intelsdi-x/snap-plugin-utilities/ns"
	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const (
	// namespace vendor prefix
	NS_VENDOR = "intel"
	// namespace os prefix
	NS_CLASS = "linux"
	// namespace plugin name
	NS_PLUGIN = "cgroups"
	// version of plugin
	VERSION = 1
)

type CgManager interface {
}

type cgroups struct {
	mountPoints    map[string]string
	cgroupManagers map[string]*Cgroup
	statsRefMap    map[string]string
	cgstat         *cgroupstat
}

func NewCgroups() (*cgroups, error) {
	cgstat := NewCgroupstat()
	mountPoints, err := cgstat.discoverMountPoints()
	if err != nil {
		return nil, err
	}
	c := &cgroups{
		mountPoints:    mountPoints,
		cgroupManagers: make(map[string]*Cgroup),
		statsRefMap:    make(map[string]string),
		cgstat:         cgstat}

	err = c.refreshCgroups()
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *cgroups) CollectMetrics(metricTypes []plugin.PluginMetricType) ([]plugin.PluginMetricType, error) {
	c.refreshCgroups()
	c.fetchStats()
	c.buildStatsRefMap()
	hostname := "localhost"
	if hostname_, err := os.Hostname(); err == nil {
		hostname = hostname_
	}
	queryMap := map[string]*plugin.PluginMetricType{}
	for _, metricQuery := range metricTypes {
		regexStr := "^" + strings.Replace(strings.Join(metricQuery.Namespace()[3:], "/"), "*", "[^/]*", -1) + "$"
		queryMap[regexStr] = &metricQuery
		//FIXME:removeit \/
		fmt.Printf("  rx: %s \n", regexStr)
	}
	metrics := []plugin.PluginMetricType{}
	for queryExpr, metricQuery := range queryMap {
		matchFound := false
		for nsEntry, path := range c.statsRefMap {
			if match, _ := regexp.MatchString(queryExpr, nsEntry); !match {
				continue
			}
			manager := c.cgroupManagers[path]
			nsSplit := strings.Split(nsEntry, "/")[strings.Count(path, "/")+1:]
			data := ns.GetValueByNamespace(manager.stats, nsSplit)
			mt := plugin.PluginMetricType{
				Data_:      data,
				Timestamp_: time.Now(),
				Namespace_: append([]string{NS_VENDOR, NS_CLASS, NS_PLUGIN}, strings.Split(nsEntry, "/")...),
				Version_:   VERSION,
				Source_:    hostname,
				Tags_:      copyTags(metricQuery.Tags_),
				Labels_:    metricQuery.Labels_[:]}
			metrics = append(metrics, mt)
			matchFound = matchFound || true
		}
		//TODO: found nothing for the query?
		//TODO: remove metric if match found?
	}
	return metrics, nil
}

func copyTags(tags map[string]string) (res map[string]string) {
	res = map[string]string{}
	for k, v := range tags {
		res[k] = v
	}
	return
}

func (c *cgroups) GetMetricTypes(plugin.PluginConfigType) ([]plugin.PluginMetricType, error) {
	mts := []plugin.PluginMetricType{}
	if err := c.fetchStats(); err != nil {
		return mts, err
	}
	for path, manager := range c.cgroupManagers {
		namespace := []string{}
		ns.FromCompositeObject(manager.stats, "", &namespace)
		for _, nsentry := range namespace {
			fullNs := filepath.Join(NS_VENDOR, NS_CLASS, NS_PLUGIN, path, nsentry)
			mts = append(mts, plugin.PluginMetricType{
				Namespace_: strings.Split(fullNs, "/")})
		}
	}
	return mts, nil
}

func (c *cgroups) buildStatsRefMap() {
	for path, manager := range c.cgroupManagers {
		namespace := []string{}
		ns.FromCompositeObject(
			manager.stats,
			"",
			&namespace,
			ns.InspectEmptyContainers(ns.AlwaysFalse),
			ns.InspectNilPointers(ns.AlwaysFalse))
		for _, nsEntry := range namespace {
			c.statsRefMap[filepath.Join(path, nsEntry)] = path
			//FIXME:removeit \/
			fmt.Printf("  sr: %s\n", nsEntry)
		}
	}
}

func (c *cgroups) GetConfigPolicy() (*cpolicy.ConfigPolicy, error) {
	return cpolicy.New(), nil
}

func (c *cgroups) refreshCgroups() error {
	markManagersDirty(c.cgroupManagers)
	c.cgstat.discoverCgroupsFromFs(c.mountPoints, c.cgroupManagers)
	discardDirtyManagers(c.cgroupManagers)
	return nil
}

func (c *cgroups) fetchStats() error {
	for _, manager := range c.cgroupManagers {
		if stats, err := manager.GetStats(); err != nil {
			return err
		} else {
			manager.stats = stats
		}
	}
	return nil
}
