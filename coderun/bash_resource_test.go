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
	RunEnvironment RunEnvironment
}

func (suite *BashSuite) SetupTest() {
	suite.Resource = BashResource()
	suite.CRDockerMock = &CRDockerMock{}
	suite.RunEnvironment = RunEnvironment{CRDocker: suite.CRDockerMock}
}

func (suite *BashSuite) TestRegisterOnCmd() {
	assert.True(suite.T(), suite.Resource.RegisterOnCmd("bash"))
}

func (suite *BashSuite) TestDoesntRegisterOnWrongCmd() {
	assert.False(suite.T(), suite.Resource.RegisterOnCmd("asdf"))
}

func (suite *BashSuite) TestSetup() {
	d := suite.CRDockerMock
	d.On("Pull", mock.AnythingOfType("string"))
	suite.Resource.Setup(suite.RunEnvironment)
	d.AssertExpectations(suite.T())
}

func (suite *BashSuite) TestRun() {
	m := suite.CRDockerMock
	m.On("Run", mock.AnythingOfType(fmt.Sprintf("%T", dockerRunConfig{})))
	suite.Resource.Run(suite.RunEnvironment)
	m.AssertExpectations(suite.T())
}

func TestBashSuite(t *testing.T) {
	suite.Run(t, new(BashSuite))
}
