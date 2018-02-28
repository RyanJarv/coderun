package coderun

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type BashSuite struct {
	suite.Suite
	Resource       *Resource
	CRDockerMock   *CRDockerMock
	RunEnvironment *RunEnvironment
	ProviderEnv    IProviderEnv
}

func (suite *BashSuite) SetupTest() {
	suite.Resource = BashResource()
	suite.CRDockerMock = &CRDockerMock{}
	suite.RunEnvironment = &RunEnvironment{CRDocker: suite.CRDockerMock}
	suite.ProviderEnv = dockerProviderEnv{}
}

func (suite *BashSuite) TestRegister() {
	suite.RunEnvironment.Cmd = []string{"bash"}
	assert.True(suite.T(), suite.Resource.Register(suite.RunEnvironment, suite.ProviderEnv))
}

func (suite *BashSuite) TestDoesntRegisterOnWrongCmd() {
	suite.RunEnvironment.Cmd = []string{"asdf"}
	assert.False(suite.T(), suite.Resource.Register(suite.RunEnvironment, suite.ProviderEnv))
}

func (suite *BashSuite) TestSetup() {
	d := suite.CRDockerMock
	d.On("Pull", mock.AnythingOfType("string"))
	suite.Resource.Setup(suite.RunEnvironment, suite.ProviderEnv)
	d.AssertExpectations(suite.T())
}

func (suite *BashSuite) TestRun() {
	m := suite.CRDockerMock
	m.On("Run", mock.AnythingOfType(fmt.Sprintf("%T", dockerRunConfig{})))
	suite.Resource.Run(suite.RunEnvironment, suite.ProviderEnv)
	m.AssertExpectations(suite.T())
}

func TestBashSuite(t *testing.T) {
	suite.Run(t, new(BashSuite))
}
