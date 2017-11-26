package gobreak

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDo(t *testing.T) {
	err := Do("test", func() error {
		return nil
	}, nil)
	assert.Nil(t, err)

	err = Do("test", func() error {
		return errors.New("failed")
	}, nil)
	assert.Equal(t, errors.New("failed"), err)

	err = Do("test", func() error {
		return errors.New("failed")
	}, func(error) error {
		return nil
	})
	assert.Nil(t, err)

	err = Do("test", func() error {
		return errors.New("failed")
	}, func(error) error {
		return errors.New("fallback")
	})
	assert.Equal(t, errors.New("fallback"), err)

	err = Do("test", func() error {
		return errors.New("failed")
	}, nil)
	assert.Equal(t, errors.New("circuit breaker is open"), err)

	err = Do("test", func() error {
		return errors.New("failed")
	}, nil)
	assert.Equal(t, errors.New("circuit breaker is open"), err)
}

func TestDoDelay(t *testing.T) {
	err := Do("delay", func() error {
		time.Sleep(1 * time.Second)
		return nil
	}, nil)
	assert.Nil(t, err)
}

func BenchmarkNormal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		func() error {
			return nil
		}()
	}
}

func BenchmarkDo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Do("test", func() error {
			return nil
		}, nil)
	}
}

func BenchmarkDoFail(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Do("test", func() error {
			return errors.New("fail")
		}, nil)
	}
}
