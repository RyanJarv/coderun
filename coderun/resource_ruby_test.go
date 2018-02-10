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
	ProviderEnv    IProviderEnv
}

func (suite *RubySuite) SetupTest() {
	suite.Resource = RubyResource()
	suite.CRDockerMock = &CRDockerMock{}
	suite.ExecMock = &ExecMock{}
	suite.RunEnvironment = RunEnvironment{}
	suite.ProviderEnv = dockerProviderEnv{CRDocker: suite.CRDockerMock, Exec: suite.ExecMock.Exec}
}

func (suite *RubySuite) TestRegister() {
	suite.RunEnvironment.Cmd = []string{"ruby"}
	assert.True(suite.T(), suite.Resource.Register(suite.RunEnvironment, suite.ProviderEnv))
}

func (suite *RubySuite) TestDoesntRegisterOnWrongCmd() {
	suite.RunEnvironment.Cmd = []string{"asdf"}
	assert.False(suite.T(), suite.Resource.Register(suite.RunEnvironment, suite.ProviderEnv))
}

func (suite *RubySuite) TestSetup() {
	d := suite.CRDockerMock
	d.On("Pull", "ruby:2.3")
	suite.Resource.Setup(suite.RunEnvironment, suite.ProviderEnv)
	d.AssertExpectations(suite.T())
}

func (suite *RubySuite) TestRun() {
	d := suite.CRDockerMock
	d.On("Run", mock.AnythingOfType(fmt.Sprintf("%T", dockerRunConfig{})))
	suite.Resource.Run(suite.RunEnvironment, suite.ProviderEnv)
	d.AssertExpectations(suite.T())
}

func TestRubySuite(t *testing.T) {
	suite.Run(t, new(RubySuite))
}
