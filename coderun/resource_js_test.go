package coderun

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type JsSuite struct {
	suite.Suite
	Resource       *Resource
	CRDockerMock   *CRDockerMock
	ExecMock       *ExecMock
	RunEnvironment RunEnvironment
	ProviderEnv    IProviderEnv
}

func (suite *JsSuite) SetupTest() {
	suite.Resource = JsResource()
	suite.CRDockerMock = &CRDockerMock{}
	suite.ExecMock = &ExecMock{}
	suite.RunEnvironment = RunEnvironment{}
	suite.ProviderEnv = dockerProviderEnv{CRDocker: suite.CRDockerMock, Exec: suite.ExecMock.Exec}
}

func (suite *JsSuite) TestRegister() {
	suite.RunEnvironment.Cmd = []string{"node"}
	assert.True(suite.T(), suite.Resource.Register(suite.RunEnvironment, suite.ProviderEnv))
}

func (suite *JsSuite) TestDoesntRegisterOnWrongCmd() {
	suite.RunEnvironment.Cmd = []string{"asdf"}
	assert.False(suite.T(), suite.Resource.Register(suite.RunEnvironment, suite.ProviderEnv))
}

func (suite *JsSuite) TestSetup() {
	e := suite.ExecMock
	e.On("Exec", mock.Anything).Return("cmd")
	suite.Resource.Setup(suite.RunEnvironment, suite.ProviderEnv)
	e.AssertExpectations(suite.T())
}

func (suite *JsSuite) TestRun() {
	m := suite.CRDockerMock
	m.On("Run", mock.AnythingOfType(fmt.Sprintf("%T", dockerRunConfig{})))
	suite.Resource.Run(suite.RunEnvironment, suite.ProviderEnv)
	m.AssertExpectations(suite.T())
}

func TestJsSuite(t *testing.T) {
	suite.Run(t, new(JsSuite))
}
