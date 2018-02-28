package coderun

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type BundlerSuite struct {
	suite.Suite
	Resource       *Resource
	CRDockerMock   *CRDockerMock
	ExecMock       *ExecMock
	RunEnvironment *RunEnvironment
	ProviderEnv    IProviderEnv
}

func (suite *BundlerSuite) SetupTest() {
	suite.Resource = BundlerResource()
	suite.CRDockerMock = &CRDockerMock{}
	suite.ExecMock = &ExecMock{}
	suite.RunEnvironment = &RunEnvironment{
		CRDocker: suite.CRDockerMock,
		Exec:     suite.ExecMock.Exec,
	}
	suite.ProviderEnv = dockerProviderEnv{}
	Logger = SetupLogger("error")
}

func (suite *BundlerSuite) TestRegister() {
	suite.RunEnvironment.Cmd = []string{"ruby"}
	assert.True(suite.T(), suite.Resource.Register(suite.RunEnvironment, suite.ProviderEnv))
}

func (suite *BundlerSuite) TestDoesntRegisterOnWrongCmd() {
	suite.RunEnvironment.Cmd = []string{"asdf"}
	assert.False(suite.T(), suite.Resource.Register(suite.RunEnvironment, suite.ProviderEnv))
}

func (suite *BundlerSuite) TestSetup() {
	d := suite.CRDockerMock
	d.On("Pull", "ruby:2.3")
	suite.Resource.Setup(suite.RunEnvironment, suite.ProviderEnv)
	d.AssertExpectations(suite.T())
}

func TestBundlerSuite(t *testing.T) {
	suite.Run(t, new(BundlerSuite))
}
