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
	"encoding/json"
	. "github.com/smartystreets/goconvey/convey"
	//"github.com/stretchr/testify/mock"
	"os"
	"path/filepath"
	"testing"
	"time"
	"github.com/stretchr/testify/mock"
)

type JsonWalker struct {
	fs interface{}
}

type DummyFileInfo struct {
	name  string
	isDir bool
	sys   interface{}
}

func (i *DummyFileInfo) Name() string {
	return i.name
}

func (i *DummyFileInfo) IsDir() bool {
	return i.isDir
}

func (i *DummyFileInfo) Size() int64 {
	return 0
}

func (i *DummyFileInfo) Mode() os.FileMode {
	if i.isDir {
		return os.ModeDir
	}
	return 0
}

func (i *DummyFileInfo) ModTime() time.Time {
	return time.Time{}
}

func (i *DummyFileInfo) Sys() interface{} {
	return i.sys
}

func NewJsonWalker(jsonSource string) *JsonWalker {
	walker := new(JsonWalker)
	var root interface{}
	json.Unmarshal([]byte(jsonSource), &(root))
	walker.fs = root
	return walker
}

func (w *JsonWalker) Walk(path string, walkFunc filepath.WalkFunc) error {
	node := seek(w.fs, path)
	walk(node, path, walkFunc)
	return nil
}

func seek(root interface{}, seekPath string) interface{} {
	var result interface{}
	walk(root, "/", func(path string, info os.FileInfo, _ error) error {
		if result != nil {
			return filepath.SkipDir
		} else if path == seekPath {
			result = info.Sys()
			return filepath.SkipDir
		}
		return nil
	})
	return result
}

func basename(path string) string {
	base := filepath.Base(path)
	if base == "." {
		return "/"
	}
	return base
}

func walk(node interface{}, path string, walkFunc filepath.WalkFunc) error {
	var err error
	if dirNode, isDir := node.(map[string]interface{}); isDir {
		err = walkFunc(path, &DummyFileInfo{basename(path), true, dirNode}, nil)
		if err == filepath.SkipDir {
			return nil
		}
		for k, subNode := range dirNode {
			err = walk(subNode, filepath.Join(path, k), walkFunc)
			if err == filepath.SkipDir {
				return nil
			}
		}
	} else {
		err = walkFunc(path, &DummyFileInfo{basename(path), false, node}, nil)
		if err == filepath.SkipDir {
			return err
		}
	}
	return nil
}

var cgFs = `{
"cgfs": {
	"cpu": {
		"cg_c": {
		},
		"docker": {
			"dk_c": {
			},
			"dk_mc": {
			}
		}
	},
	"memory": {
		"cg_m": {
		},
		"docker": {
			"dk_m": {
			},
			"dk_mc": {
			}
		}
	}
}
}
`

func TestDiscoverCgroupsFromFs(t *testing.T) {
	Convey("Given a custom cgroups fs", t, func() {
		fsWalker := NewJsonWalker(cgFs)
		Convey("and a cgroupstat object initialized", func() {
			cgstat := NewCgroupstat()
			cgstat.fsWalker = fsWalker
			mockGw := &MockCgroupGateway{}
			mockGw.On("GetAllSubsystems").Return([]string{"cpu", "memory"}, nil)
			mockGw.On("FindCgroupMountpoint", "cpu").Return("/cgfs/cpu", nil)
			mockGw.On("FindCgroupMountpoint", "memory").Return("/cgfs/memory", nil)
			mockGw.On("NewCgManager", mock.AnythingOfType("string")).Return(new(MockCgManager))
			cgstat.cgroupGw = mockGw
			var mountPoints map[string]string
			Convey("when  discoverMountPoints is called", func() {
				var err error
				mountPoints, err = cgstat.discoverMountPoints()
				Convey("Then error should not be reported", func() {
					So(err, ShouldBeNil)
				})
				Convey("Then mountpoints for preconfigured subsystems should be reported", func() {
					So(mountPoints, ShouldContainKey, "cpu")
					So(mountPoints, ShouldContainKey, "memory")
					So(mountPoints["cpu"], ShouldEqual, "/cgfs/cpu")
					So(mountPoints["memory"], ShouldEqual, "/cgfs/memory")
				})
				Convey("and  when discoverCgroupsFromFs is called", func() {
					managers := map[string]*Cgroup{}
					cgstat.discoverCgroupsFromFs(mountPoints, managers)
					keys := []string{}
					for k, _ := range managers {
						keys = append(keys, k)
					}
					Convey("Then discovered cgroups should merge those in  cpu and in  memory subsystem", func() {
						So(managers, ShouldContainKey, "/")
						So(managers, ShouldContainKey, "cg_c")
						So(managers, ShouldContainKey, "cg_m")
						So(managers, ShouldContainKey, "docker")
						So(managers, ShouldContainKey, "docker/dk_m")
						So(managers, ShouldContainKey, "docker/dk_c")
						So(managers, ShouldContainKey, "docker/dk_mc")
					})
				})
			})

			mockGw.AssertExpectations(t)
		})
	})
}
