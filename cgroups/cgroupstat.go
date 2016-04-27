package cgroups

import (
	lcgroups "github.com/opencontainers/runc/libcontainer/cgroups"
	"github.com/opencontainers/runc/libcontainer/cgroups/fs"
	"github.com/opencontainers/runc/libcontainer/configs"
	"os"
	"path/filepath"
)

type Cgroup struct {
	*fs.Manager
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

type CgroupGateway interface {
	GetAllSubsystems() ([]string, error)
	FindCgroupMountpoint(string) (string, error)
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
			if err == nil {
				paths = append(paths, relPath)
			}
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
				newMgr := &(Cgroup{
					Manager: &fs.Manager{
						Cgroups: &configs.Cgroup{Name: path},
						Paths:   make(map[string]string)}})
				managers[path] = newMgr
			} else {
				manager.dirty = false
			}
			manager := managers[path]
			manager.Paths[subsystem] = filepath.Join(mountPoint, path)
		}
	}
}

func markManagersDirty(managers map[string]*Cgroup) {
	for _, manager := range managers {
		manager.dirty = true
		manager.Paths = make(map[string]string)
	}
}

func discardDirtyManagers(managers map[string]*Cgroup) {
	for key, manager := range managers {
		if manager.dirty {
			delete(managers, key)
		}
	}
}
