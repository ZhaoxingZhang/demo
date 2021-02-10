package test

import (
	"errors"
	"fmt"
	"github.com/cenk/backoff"
	"testing"
	"time"
)

type getInterval func() error
func TestBackoff(t *testing.T) {
	i := 0
	policy := backoff.NewExponentialBackOff()
	policy.InitialInterval = 100 * time.Millisecond
	policy.MaxElapsedTime = 1000 * time.Millisecond
	fn := func() error{
		fmt.Printf("print time %v\n", i)
		i++
		return errors.New("err")
	}
	backoff.Retry(fn, policy)
}

func TestTimer(t *testing.T) {
	sleepFn := func(dura int64) {
		time.Sleep(time.Duration(dura))
	}
	for i:=0;i<10;i++ {
		func(i int) {
			defer printTimer(i, time.Now())
			sleepFn(int64(i)*1e8)
		}(i)
	}
}
func printTimer(method int, start time.Time) {
	fmt.Printf("[slow method] method %v cost %v\n",method, time.Since(start))
}