package coderun

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMatchCommandOrExt(t *testing.T) {
	assert := assert.New(t)
	assert.True(MatchCommandOrExt([]string{"test.sh"}, "bash", ".sh"))
	assert.True(MatchCommandOrExt([]string{"test.sh", "-r", "fsic"}, "bash", ".sh"))
	assert.True(MatchCommandOrExt([]string{"bash", "test"}, "bash", ".sh"))
	assert.True(MatchCommandOrExt([]string{"/usr/bin/bash"}, "bash", ".sh"))

	assert.False(MatchCommandOrExt([]string{"qwerpoiu"}, "bash", ".sh"))
	assert.False(MatchCommandOrExt([]string{"/usr/bin/notbash"}, "bash", ".sh"))
	assert.False(MatchCommandOrExt([]string{"test.sh.not"}, "bash", ".sh"))
}

type ExecMock struct {
	mock.Mock
}

func (m *ExecMock) Exec(image ...string) string {
	args := m.Called(image)
	return args.String(0)
}
