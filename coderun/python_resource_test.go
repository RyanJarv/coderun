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
}

func (suite *PythonSuite) SetupTest() {
	suite.Resource = PythonResource()
	suite.CRDockerMock = &CRDockerMock{}
	suite.ExecMock = &ExecMock{}
	suite.RunEnvironment = RunEnvironment{CRDocker: suite.CRDockerMock, Exec: suite.ExecMock.Exec}
}

func (suite *PythonSuite) TestRegisterOnCmd() {
	assert.True(suite.T(), suite.Resource.RegisterOnCmd("python"))
}

func (suite *PythonSuite) TestDoesntRegisterOnWrongCmd() {
	assert.False(suite.T(), suite.Resource.RegisterOnCmd("asdf"))
}

func (suite *PythonSuite) TestSetup() {
	d := suite.CRDockerMock
	d.On("Pull", "python:3")
	suite.Resource.Setup(suite.RunEnvironment)
	d.AssertExpectations(suite.T())
}

func (suite *PythonSuite) TestRun() {
	d := suite.CRDockerMock
	d.On("Run", mock.AnythingOfType(fmt.Sprintf("%T", dockerRunConfig{})))
	suite.Resource.Run(suite.RunEnvironment)
	d.AssertExpectations(suite.T())
}

func TestPythonSuite(t *testing.T) {
	suite.Run(t, new(PythonSuite))
}
