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
	RunEnvironment RunEnvironment
}

func (suite *GoSuite) SetupTest() {
	suite.Resource = GoResource()
	suite.CRDockerMock = &CRDockerMock{}
	suite.RunEnvironment = RunEnvironment{CRDocker: suite.CRDockerMock}
}

func (suite *GoSuite) TestRegisterOnCmd() {
	assert.True(suite.T(), suite.Resource.RegisterOnCmd("go"))
}

func (suite *GoSuite) TestDoesntRegisterOnWrongCmd() {
	assert.False(suite.T(), suite.Resource.RegisterOnCmd("asdf"))
}

func (suite *GoSuite) TestSetup() {
	d := suite.CRDockerMock
	d.On("Pull", mock.AnythingOfType("string"))
	d.On("Run", mock.Anything)
	suite.Resource.Setup(suite.RunEnvironment)
	d.AssertExpectations(suite.T())
}

func (suite *GoSuite) TestRun() {
	m := suite.CRDockerMock
	m.On("Run", mock.AnythingOfType(fmt.Sprintf("%T", dockerRunConfig{})))
	suite.Resource.Run(suite.RunEnvironment)
	m.AssertExpectations(suite.T())
}

func TestGoSuite(t *testing.T) {
	suite.Run(t, new(GoSuite))
}
