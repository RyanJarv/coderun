package coderun

import "github.com/stretchr/testify/mock"

func NewRegistryMock() *RegistryMock {
	return &RegistryMock{}
}

type RegistryMock struct {
	mock.Mock
}

func (r *RegistryMock) AddAt(l int, s *StepCallback)                  { r.Called(l, s) }
func (r *RegistryMock) AddBefore(search *StepSearch, s *StepCallback) { r.Called(search, s) }
func (r *RegistryMock) AddAfter(search *StepSearch, s *StepCallback)  { r.Called(search, s) }
func (r *RegistryMock) Run()                                          { r.Called() }
