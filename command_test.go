package gobreak

import (
	"testing"

	"github.com/sony/gobreaker"
	"github.com/stretchr/testify/assert"
)

func TestErrorToEvent(t *testing.T) {
	assert.Equal(t, "too-many-requests", errorToEvent(gobreaker.ErrTooManyRequests))
	assert.Equal(t, "circuit-open", errorToEvent(gobreaker.ErrOpenState))
	assert.Equal(t, "success", errorToEvent(nil))
}
