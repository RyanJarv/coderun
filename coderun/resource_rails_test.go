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
	ProviderEnv    IProviderEnv
}

func (suite *RailsSuite) SetupTest() {
	suite.Resource = RailsResource()
	suite.CRDockerMock = &CRDockerMock{}
	suite.ExecMock = &ExecMock{}
	suite.RunEnvironment = RunEnvironment{}
	suite.ProviderEnv = dockerProviderEnv{CRDocker: suite.CRDockerMock, Exec: suite.ExecMock.Exec}
}

func (suite *RailsSuite) TestRegister() {
	suite.RunEnvironment.Cmd = []string{"rails", "server"}
	assert.True(suite.T(), suite.Resource.Register(suite.RunEnvironment, suite.ProviderEnv.(dockerProviderEnv)))
}

func (suite *RailsSuite) TestDoesntRegisterOnWrongCmd() {
	suite.RunEnvironment.Cmd = []string{"asdf"}
	assert.False(suite.T(), suite.Resource.Register(suite.RunEnvironment, suite.ProviderEnv.(dockerProviderEnv)))
}

func (suite *RailsSuite) TestSetup() {
	d := suite.CRDockerMock
	d.On("Pull", "ruby:2.3")
	suite.Resource.Setup(suite.RunEnvironment, suite.ProviderEnv.(dockerProviderEnv))
	d.AssertExpectations(suite.T())
}

func (suite *RailsSuite) TestRun() {
	d := suite.CRDockerMock
	d.On("Run", mock.AnythingOfType(fmt.Sprintf("%T", dockerRunConfig{})))
	suite.Resource.Run(suite.RunEnvironment, suite.ProviderEnv)
	d.AssertExpectations(suite.T())
}

func TestRailsSuite(t *testing.T) {
	suite.Run(t, new(RailsSuite))
}
