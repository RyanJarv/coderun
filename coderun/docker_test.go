package coderun

import "github.com/stretchr/testify/mock"

type CRDockerMock struct {
	ICRDocker
	mock.Mock
}

func (m *CRDockerMock) Pull(image string) {
	m.Called(image)
}

func (m *CRDockerMock) Run(c dockerRunConfig) {
	m.Called(c)
}
