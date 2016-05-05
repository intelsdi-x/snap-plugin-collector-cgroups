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
	"github.com/stretchr/testify/mock"
	lcgroups "github.com/opencontainers/runc/libcontainer/cgroups"
)

type MockCgroupGateway struct {
	mock.Mock
}

func (m *MockCgroupGateway) GetAllSubsystems() ([]string, error) {
	ret := m.Called()
	return ret.Get(0).([]string), ret.Error(1)
}

func (m *MockCgroupGateway) FindCgroupMountpoint(subsystem string) (string, error) {
	ret := m.Called(subsystem)
	return ret.String(0), ret.Error(1)
}

func (m *MockCgroupGateway) NewCgManager(name string) CgManager {
	ret := m.Called(name)
	return ret.Get(0).(CgManager)
}

type MockCgManager struct {
	mock.Mock
}

func (m *MockCgManager) GetStats() (*lcgroups.Stats, error) {
	ret := m.Called()
	return ret.Get(0).(*lcgroups.Stats), ret.Error(1)
}

func (m *MockCgManager) Paths() map[string]string {
	ret := m.Called()
	return ret.Get(0).(map[string]string)
}

func (m *MockCgManager) SetPaths(paths map[string]string) {
	m.Called(paths)
	return
}

func (m *MockCgManager) AddPath(name, path string) {
	m.Called(name, path)
	return
}