package gobreak

import (
	"testing"

	"github.com/sony/gobreaker"
	"github.com/stretchr/testify/assert"
)

func TestErrorToEvent(t *testing.T) {
	assert.Equal(t, "too-many-requests", errorToEvent(gobreaker.ErrTooManyRequests))
}
