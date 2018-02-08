package coderun

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBashRegisterOnCmd(t *testing.T) {
	bash := BashResource()
	assert.True(t, bash.RegisterOnCmd("bash"))
}

func TestBashSetup(t *testing.T) {
	m := &CRDockerMock{}
	m.On("Pull", "bash")
	BashResource().Setup(RunEnvironment{CRDocker: m})
	m.AssertExpectations(t)
}

func TestBashRun(t *testing.T) {
	m := &CRDockerMock{}
	m.On("Run", mock.MatchedBy(func(c dockerRunConfig) bool {
		return c.Image == "ubuntu"
	}))
	BashResource().Run(RunEnvironment{CRDocker: m})
	m.AssertExpectations(t)
}
