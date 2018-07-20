package gobreak

import (
	"testing"

	"github.com/sony/gobreaker"
	"github.com/stretchr/testify/assert"
)

func TestErrorToEvent(t *testing.T) {
	assert.Equal(t, "too-many-requests", ErrorToEvent(gobreaker.ErrTooManyRequests))
	assert.Equal(t, "circuit-open", ErrorToEvent(gobreaker.ErrOpenState))
	assert.Equal(t, "success", ErrorToEvent(nil))
}
