package coderun

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type GoSuite struct {
	suite.Suite
	Resource       *Resource
	CRDockerMock   *CRDockerMock
	RunEnvironment *RunEnvironment
	ProviderEnv    dockerProviderEnv
}

func (suite *GoSuite) SetupTest() {
	suite.Resource = GoResource()
	suite.CRDockerMock = &CRDockerMock{}
	suite.RunEnvironment = &RunEnvironment{CRDocker: suite.CRDockerMock}
	suite.ProviderEnv = dockerProviderEnv{}
}

func (suite *GoSuite) TestRegister() {
	suite.RunEnvironment.Cmd = []string{"go"}
	assert.True(suite.T(), suite.Resource.Register(suite.RunEnvironment, suite.ProviderEnv))
}

func (suite *GoSuite) TestDoesntRegisterOnWrongCmd() {
	suite.RunEnvironment.Cmd = []string{"asdf"}
	assert.False(suite.T(), suite.Resource.Register(suite.RunEnvironment, suite.ProviderEnv))
}

func (suite *GoSuite) TestSetup() {
	d := suite.CRDockerMock
	d.On("Pull", mock.AnythingOfType("string"))
	d.On("Run", mock.Anything)
	suite.Resource.Setup(suite.RunEnvironment, suite.ProviderEnv)
	d.AssertExpectations(suite.T())
}

func (suite *GoSuite) TestRun() {
	m := suite.CRDockerMock
	m.On("Run", mock.AnythingOfType(fmt.Sprintf("%T", dockerRunConfig{})))
	suite.Resource.Run(suite.RunEnvironment, suite.ProviderEnv)
	m.AssertExpectations(suite.T())
}

func TestGoSuite(t *testing.T) {
	suite.Run(t, new(GoSuite))
}
