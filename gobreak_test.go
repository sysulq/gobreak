package gobreak

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDo(t *testing.T) {
	ctx := context.TODO()
	err := Do(ctx, "test do", func(context.Context) error {
		return nil
	}, nil)
	assert.Nil(t, err)

	err = Do(ctx, "test do", func(context.Context) error {
		return errors.New("failed")
	}, nil)
	assert.Equal(t, errors.New("failed"), err)

	err = Do(ctx, "test do", func(context.Context) error {
		return errors.New("failed")
	}, func(context.Context, error) error {
		return nil
	})
	assert.Nil(t, err)

	err = Do(ctx, "test do", func(context.Context) error {
		return errors.New("failed")
	}, func(context.Context, error) error {
		return errors.New("fallback")
	})
	assert.Equal(t, errors.New("fallback"), err)

	err = Do(ctx, "test do", func(context.Context) error {
		return errors.New("failed")
	}, nil)
	assert.Equal(t, errors.New("circuit breaker is open"), err)
}

func TestDoDelay(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.TODO(), 1*time.Second)
	defer cancel()
	err := Do(ctx, "delay", func(context.Context) error {
		time.Sleep(2 * time.Second)
		return nil
	}, nil)
	assert.Equal(t, context.DeadlineExceeded, err)
}

func TestGoCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.TODO())
	err := Go(ctx, "go cancel", func(context.Context) error {
		time.Sleep(1 * time.Second)
		return nil
	}, nil)
	cancel()
	assert.Equal(t, context.Canceled, <-err)
}

func TestDoPanic(t *testing.T) {
	ctx, _ := context.WithTimeout(context.TODO(), 1*time.Second)
	err := Do(ctx, "TestDoPanic", func(context.Context) error {
		panic("panic")
		return nil
	}, nil)

	assert.Equal(t, errors.New("command panics"), err)
}

func TestGo(t *testing.T) {
	ctx := context.TODO()
	err := Go(ctx, "test go", func(context.Context) error {
		return nil
	}, nil)
	assert.Nil(t, <-err)

	err = Go(ctx, "test go", func(context.Context) error {
		return errors.New("failed")
	}, nil)
	assert.Equal(t, errors.New("failed"), <-err)

	err = Go(ctx, "test go", func(context.Context) error {
		return errors.New("failed")
	}, func(context.Context, error) error {
		return nil
	})
	assert.Nil(t, <-err)

	err = Go(ctx, "test go", func(context.Context) error {
		return errors.New("failed")
	}, func(context.Context, error) error {
		return errors.New("fallback")
	})
	assert.Equal(t, errors.New("fallback"), <-err)

	err = Go(ctx, "test go", func(context.Context) error {
		return errors.New("failed")
	}, nil)
	assert.Equal(t, errors.New("circuit breaker is open"), <-err)

	err = Go(ctx, "test go", func(context.Context) error {
		return errors.New("failed")
	}, nil)
	assert.Equal(t, errors.New("circuit breaker is open"), <-err)
}
func TestGoNormal(t *testing.T) {
	ctx, _ := context.WithTimeout(context.TODO(), 2*time.Second)
	err := Go(ctx, "normal", func(context.Context) error {
		time.Sleep(1 * time.Second)
		return nil
	}, nil)

	assert.Equal(t, nil, <-err)
}

func TestGoDoubleCheckLock(t *testing.T) {
	ctx, _ := context.WithTimeout(context.TODO(), 2*time.Second)
	err := Go(ctx, "normal", func(context.Context) error {
		time.Sleep(1 * time.Second)
		return nil
	}, nil)

	assert.Equal(t, nil, <-err)
}

func TestGoDelay(t *testing.T) {
	ctx, _ := context.WithTimeout(context.TODO(), 1*time.Second)
	err := Go(ctx, "delay", func(context.Context) error {
		time.Sleep(2 * time.Second)
		return nil
	}, nil)

	assert.Equal(t, context.DeadlineExceeded, <-err)
}

func TestGoDelayFallback(t *testing.T) {
	ctx, _ := context.WithTimeout(context.TODO(), 1*time.Second)
	err := Go(ctx, "delay and fallback", func(context.Context) error {
		time.Sleep(2 * time.Second)
		return nil
	}, func(context.Context, error) error {
		return nil
	})

	assert.Equal(t, nil, <-err)
}

func TestGoPanic(t *testing.T) {
	ctx, _ := context.WithTimeout(context.TODO(), 1*time.Second)
	err := Go(ctx, "panic", func(context.Context) error {
		panic("panic")
		return nil
	}, nil)

	assert.Equal(t, errors.New("command panics"), <-err)
}

func TestGoPanicFallBack(t *testing.T) {
	ctx, _ := context.WithTimeout(context.TODO(), 1*time.Second)
	err := Go(ctx, "panic and fallback", func(context.Context) error {
		panic("panic")
		return nil
	}, func(context.Context, error) error {
		return errors.New("fallback")
	})

	assert.Equal(t, errors.New("fallback"), <-err)
}

func BenchmarkNormal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		func() error {
			return nil
		}()
	}
}

func BenchmarkDo(b *testing.B) {
	ctx, _ := context.WithTimeout(context.TODO(), 1*time.Second)
	for i := 0; i < b.N; i++ {
		Do(ctx, "test", func(context.Context) error {
			return nil
		}, nil)
	}
}

func BenchmarkGo(b *testing.B) {
	ctx, _ := context.WithTimeout(context.TODO(), 1*time.Second)
	for i := 0; i < b.N; i++ {
		<-Go(ctx, "test", func(context.Context) error {
			return nil
		}, nil)
	}
}

func BenchmarkDoFail(b *testing.B) {
	ctx, _ := context.WithTimeout(context.TODO(), 1*time.Second)
	for i := 0; i < b.N; i++ {
		Do(ctx, "test", func(context.Context) error {
			return errors.New("fail")
		}, nil)
	}
}

func BenchmarkGoFail(b *testing.B) {
	ctx, _ := context.WithTimeout(context.TODO(), 1*time.Second)
	for i := 0; i < 1000; i++ {
		Go(ctx, "test", func(context.Context) error {
			return errors.New("fail")
		}, nil)
	}
	time.Sleep(1 * time.Second)
	b.Fail()
}
