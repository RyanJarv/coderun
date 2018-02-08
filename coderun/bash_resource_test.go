package coderun

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBashRegisterOnCmd(t *testing.T) {
	bash := BashProvider{}
	assert.True(t, bash.RegisterOnCmd("bash"))
}

type MockRunEnvironment struct {
	mock.Mock
}

func (m *MockRunEnvironment) CRDocker(image string) CRDocker {
	args := m.Called(image)
	return args.Get(0).(CRDocker)
}

func TestBashSetup(t *testing.T) {
	testEnv := new(MockRunEnvironment)

	testEnv.On("Pull", "bash")
	BashProvider{}.Setup(testEnv)
	testEnv.AssertExpectations(t)
}
