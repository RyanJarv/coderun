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
}

func (suite *JsSuite) SetupTest() {
	suite.Resource = JsResource()
	suite.CRDockerMock = &CRDockerMock{}
	suite.ExecMock = &ExecMock{}
	suite.RunEnvironment = RunEnvironment{CRDocker: suite.CRDockerMock, Exec: suite.ExecMock.Exec}
}

func (suite *JsSuite) TestRegisterOnCmd() {
	assert.True(suite.T(), suite.Resource.RegisterOnCmd("node"))
}

func (suite *JsSuite) TestDoesntRegisterOnWrongCmd() {
	assert.False(suite.T(), suite.Resource.RegisterOnCmd("asdf"))
}

func (suite *JsSuite) TestSetup() {
	e := suite.ExecMock
	e.On("Exec", mock.Anything).Return("cmd")
	suite.Resource.Setup(suite.RunEnvironment)
	e.AssertExpectations(suite.T())
}

func (suite *JsSuite) TestRun() {
	m := suite.CRDockerMock
	m.On("Run", mock.AnythingOfType(fmt.Sprintf("%T", dockerRunConfig{})))
	suite.Resource.Run(suite.RunEnvironment)
	m.AssertExpectations(suite.T())
}

func TestJsSuite(t *testing.T) {
	suite.Run(t, new(JsSuite))
}
