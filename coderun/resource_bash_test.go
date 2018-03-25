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
	Provider     IProvider
	Resource     *BashResource
	CRDockerMock *CRDockerMock
	RegistryMock *RegistryMock
}

func (suite *BashSuite) SetupTest() {
	suite.CRDockerMock = &CRDockerMock{}
	suite.RegistryMock = NewRegistryMock()
	env := &RunEnvironment{registry: suite.RegistryMock}
	suite.Resource = &BashResource{bash: suite.CRDockerMock, env: env}
}

func (suite *BashSuite) TestRegister() {
	suite.RegistryMock.On("AddAt", mock.AnythingOfType("int"), mock.Anything)
	suite.Resource.env.SetCmd([]string{"bash", "bash.sh"})
	assert.True(suite.T(), suite.Resource.Register(suite.Resource.env, suite.Provider))
}

func (suite *BashSuite) TestDoesntRegisterOnWrongCmd() {
	suite.Resource.env.SetCmd([]string{"asdf"})
	assert.False(suite.T(), suite.Resource.Register(suite.Resource.env, suite.Provider))
}

func (suite *BashSuite) TestSetup() {
	suite.CRDockerMock.On("Pull", mock.AnythingOfType("string"))
	suite.Resource.Setup(&StepCallback{}, &StepCallback{})
	suite.CRDockerMock.AssertExpectations(suite.T())
}

func (suite *BashSuite) TestRun() {
	suite.CRDockerMock.On("Run", mock.AnythingOfType(fmt.Sprintf("%T", dockerRunConfig{})))
	suite.Resource.Run(&StepCallback{}, &StepCallback{})
	suite.CRDockerMock.AssertExpectations(suite.T())
}

func TestBashSuite(t *testing.T) {
	suite.Run(t, new(BashSuite))
}
