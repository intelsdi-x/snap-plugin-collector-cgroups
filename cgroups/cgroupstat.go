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
	"github.com/intelsdi-x/snap-plugin-utilities/ns"
	lcgroups "github.com/opencontainers/runc/libcontainer/cgroups"
	"github.com/opencontainers/runc/libcontainer/cgroups/fs"
	"github.com/opencontainers/runc/libcontainer/configs"
	"os"
	"path/filepath"
	"strings"
)

type Cgroup struct {
	CgManager
	dirty bool
	stats *lcgroups.Stats
}

type cgroupstat struct {
	fsWalker FsWalker
	cgroupGw CgroupGateway
}

type FsWalker interface {
	Walk(path string, walkFunc filepath.WalkFunc) error
}

type CgManager interface {
	GetStats() (*lcgroups.Stats, error)
	Paths() map[string]string
	SetPaths(paths map[string]string)
	AddPath(name, path string)
}

type OsCgManager struct {
	*fs.Manager
}

func (m *OsCgManager) GetStats() (*lcgroups.Stats, error) {
	return m.Manager.GetStats()
}

func (m *OsCgManager) Paths() map[string]string {
	return m.Manager.Paths
}

func (m *OsCgManager) SetPaths(paths map[string]string) {
	m.Manager.Paths = make(map[string]string)
	for k, v := range paths {
		m.Manager.Paths[k] = v
	}
}

func (m *OsCgManager) AddPath(name, path string) {
	m.Manager.Paths[name] = path
}

type CgroupGateway interface {
	GetAllSubsystems() ([]string, error)
	FindCgroupMountpoint(string) (string, error)
	NewCgManager(name string) CgManager
}

type OsFsWalker struct {
}

func (*OsFsWalker) Walk(path string, walkFunc filepath.WalkFunc) error {
	return filepath.Walk(path, walkFunc)
}

type OsCgroupGateway struct {
}

func (*OsCgroupGateway) GetAllSubsystems() ([]string, error) {
	return lcgroups.GetAllSubsystems()
}

func (*OsCgroupGateway) FindCgroupMountpoint(system string) (string, error) {
	return lcgroups.FindCgroupMountpoint(system)
}

func (*OsCgroupGateway) NewCgManager(name string) CgManager {
	return newCgManager(name)
}

func newCgManager(name string) *OsCgManager {
	res := &OsCgManager{
		Manager: &fs.Manager{
			Cgroups: &configs.Cgroup{Name: name},
			Paths:   make(map[string]string)}}
	return res
}

func NewCgroupstat() *cgroupstat {
	return &cgroupstat{
		cgroupGw: new(OsCgroupGateway),
		fsWalker: new(OsFsWalker)}
}

func (s *cgroupstat) discoverMountPoints() (mountPoints map[string]string, err error) {
	var subsystems []string
	if subsystems, err = s.cgroupGw.GetAllSubsystems(); err != nil {
		return nil, err
	}
	type ListedSubsystem struct {
		name  string
		paths []string
	}
	mountPoints = make(map[string]string)
	for _, subsystem := range subsystems {
		point, err := s.cgroupGw.FindCgroupMountpoint(subsystem)
		if err != nil {
			return nil, err
		}
		mountPoints[subsystem] = point
	}
	return
}

func relative(parent, dest string) (string, error) {
	rel, err := filepath.Rel(parent, dest)
	if err != nil {
		return rel, err
	}
	if rel == "." {
		return "/", nil
	}
	return rel, nil
}

func (s *cgroupstat) scanCgroups(mountPoint string) (paths []string) {
	s.fsWalker.Walk(mountPoint, func(path string, info os.FileInfo, walkError error) (err error) {
		if walkError != nil {
			return nil
		}
		if info.IsDir() {
			relPath, err := relative(mountPoint, path)
			if err != nil {
				return err
			}
			pathSplit := strings.Split(relPath, "/")
			for i, part := range pathSplit {
				pathSplit[i] = ns.ReplaceNotAllowedCharsInNamespacePart(part)
			}
			relPath = strings.Join(pathSplit, "/")
			paths = append(paths, relPath)
		}
		return nil
	})
	return
}

func (s *cgroupstat) discoverCgroupsFromFs(mountPoints map[string]string, managers map[string]*Cgroup) {
	//TODO: scan subsystems concurrently
	cgroupPaths := make(map[string][]string)
	for subsystem, mountPoint := range mountPoints {
		cgroupPaths[subsystem] = s.scanCgroups(mountPoint)
	}
	// combine cgroup paths from different subsystems into result map
	for subsystem, mountPoint := range mountPoints {
		for _, path := range cgroupPaths[subsystem] {
			if manager, haveMgr := managers[path]; !haveMgr {
				//newMgr := &(Cgroup{
				//	Manager: &fs.Manager{
				//		Cgroups: &configs.Cgroup{Name: path},
				//		Paths:   make(map[string]string)}})
				newMgr := &(Cgroup{
					CgManager: s.cgroupGw.NewCgManager(path)})
				managers[path] = newMgr
			} else {
				manager.dirty = false
			}
			manager := managers[path]
			manager.AddPath(subsystem, filepath.Join(mountPoint, path))
		}
	}
}

func markManagersDirty(managers map[string]*Cgroup) {
	for _, manager := range managers {
		manager.dirty = true
		manager.SetPaths(map[string]string{})
	}
}

func discardDirtyManagers(managers map[string]*Cgroup) {
	for key, manager := range managers {
		if manager.dirty {
			delete(managers, key)
		}
	}
}
