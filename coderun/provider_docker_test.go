package coderun

import (
	"time"

	"github.com/stretchr/testify/mock"
)

type CRDockerMock struct {
	ICRDocker
	mock.Mock
}

func (m *CRDockerMock) Pull(image string) {
	m.Called(image)
}

func (m *CRDockerMock) Run(c DockerRunConfig) {
	m.Called(c)
}

func (m *CRDockerMock) Stop(name string) {
	m.Called(name)
}

func (m *CRDockerMock) Teardown(timeout time.Duration) {
	m.Called(timeout)
}

func (m *CRDockerMock) GetImageName() string {
	return "mock"
}

func (m *CRDockerMock) setImageName(name string) string {
	return "mock"
}

func (m *CRDockerMock) GetOrBuildImage(source string, cmds ...[]string) string {
	return "mock"
}

func (m *CRDockerMock) NewImageName() string {
	return "mock"
}
