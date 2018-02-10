package coderun

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type RailsSuite struct {
	suite.Suite
	Resource       *Resource
	CRDockerMock   *CRDockerMock
	ExecMock       *ExecMock
	RunEnvironment RunEnvironment
}

func (suite *RailsSuite) SetupTest() {
	suite.Resource = RailsResource()
	suite.CRDockerMock = &CRDockerMock{}
	suite.ExecMock = &ExecMock{}
	suite.RunEnvironment = RunEnvironment{CRDocker: suite.CRDockerMock, Exec: suite.ExecMock.Exec}
}

func (suite *RailsSuite) TestRegisterOnCmd() {
	assert.True(suite.T(), suite.Resource.RegisterOnCmd("rails server"))
}

func (suite *RailsSuite) TestDoesntRegisterOnWrongCmd() {
	assert.False(suite.T(), suite.Resource.RegisterOnCmd("asdf"))
}

func (suite *RailsSuite) TestSetup() {
	d := suite.CRDockerMock
	d.On("Pull", "ruby:2.3")
	suite.Resource.Setup(suite.RunEnvironment)
	d.AssertExpectations(suite.T())
}

func (suite *RailsSuite) TestRun() {
	d := suite.CRDockerMock
	d.On("Run", mock.AnythingOfType(fmt.Sprintf("%T", dockerRunConfig{})))
	suite.Resource.Run(suite.RunEnvironment)
	d.AssertExpectations(suite.T())
}

func TestRailsSuite(t *testing.T) {
	suite.Run(t, new(RailsSuite))
}
