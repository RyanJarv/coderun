package coderun

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type PythonSuite struct {
	suite.Suite
	Resource       *Resource
	CRDockerMock   *CRDockerMock
	ExecMock       *ExecMock
	RunEnvironment RunEnvironment
	ProviderEnv    IProviderEnv
}

func (suite *PythonSuite) SetupTest() {
	suite.Resource = PythonResource()
	suite.CRDockerMock = &CRDockerMock{}
	suite.ExecMock = &ExecMock{}
	suite.RunEnvironment = RunEnvironment{}
	suite.ProviderEnv = dockerProviderEnv{CRDocker: suite.CRDockerMock}
}

func (suite *PythonSuite) TestRegister() {
	suite.RunEnvironment.Cmd = []string{"python"}
	assert.True(suite.T(), suite.Resource.Register(suite.RunEnvironment, suite.ProviderEnv))
}

func (suite *PythonSuite) TestDoesntRegisterOnWrongCmd() {
	suite.RunEnvironment.Cmd = []string{"asdf"}
	assert.False(suite.T(), suite.Resource.Register(suite.RunEnvironment, suite.ProviderEnv))
}

func (suite *PythonSuite) TestSetup() {
	d := suite.CRDockerMock
	d.On("Pull", "python:3")
	suite.Resource.Setup(suite.RunEnvironment, suite.ProviderEnv)
	d.AssertExpectations(suite.T())
}

func (suite *PythonSuite) TestRun() {
	d := suite.CRDockerMock
	d.On("Run", mock.AnythingOfType(fmt.Sprintf("%T", dockerRunConfig{})))
	suite.Resource.Run(suite.RunEnvironment, suite.ProviderEnv)
	d.AssertExpectations(suite.T())
}

func TestPythonSuite(t *testing.T) {
	suite.Run(t, new(PythonSuite))
}
