package coderun

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type RubySuite struct {
	suite.Suite
	Resource       *Resource
	CRDockerMock   *CRDockerMock
	ExecMock       *ExecMock
	RunEnvironment RunEnvironment
}

func (suite *RubySuite) SetupTest() {
	suite.Resource = RubyResource()
	suite.CRDockerMock = &CRDockerMock{}
	suite.ExecMock = &ExecMock{}
	suite.RunEnvironment = RunEnvironment{CRDocker: suite.CRDockerMock, Exec: suite.ExecMock.Exec}
}

func (suite *RubySuite) TestRegisterOnCmd() {
	assert.True(suite.T(), suite.Resource.RegisterOnCmd("ruby"))
}

func (suite *RubySuite) TestDoesntRegisterOnWrongCmd() {
	assert.False(suite.T(), suite.Resource.RegisterOnCmd("asdf"))
}

func (suite *RubySuite) TestSetup() {
	d := suite.CRDockerMock
	d.On("Pull", "ruby:2.3")
	suite.Resource.Setup(suite.RunEnvironment)
	d.AssertExpectations(suite.T())
}

func (suite *RubySuite) TestRun() {
	d := suite.CRDockerMock
	d.On("Run", mock.AnythingOfType(fmt.Sprintf("%T", dockerRunConfig{})))
	suite.Resource.Run(suite.RunEnvironment)
	d.AssertExpectations(suite.T())
}

func TestRubySuite(t *testing.T) {
	suite.Run(t, new(RubySuite))
}
