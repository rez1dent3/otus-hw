package hw05parallelexecution

import (
	"errors"
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("error at n=0 or m=0", func(t *testing.T) {
		data := []struct {
			n, m int
			err  error
		}{
			{n: 0, m: 0, err: ErrErrorsLimitExceeded},
			{n: 0, m: 1, err: ErrErrorsLimitExceeded},
			{n: 1, m: 0, err: ErrErrorsLimitExceeded},
			{n: 1, m: 1, err: nil},
		}

		for _, datum := range data {
			err := Run([]Task{}, datum.n, datum.m)
			require.Truef(t, errors.Is(err, datum.err), "actual err - %v", err)
		}
	})

	t.Run("compare n & taskCount", func(t *testing.T) {
		data := []struct {
			workers, tasks int
		}{
			{workers: 100, tasks: 50},
			{workers: 50, tasks: 100},
			{workers: 100, tasks: 100},
		}

		for _, datum := range data {
			tasks := make([]Task, 0, datum.tasks)

			var runTasksCount int32
			for i := 0; i < datum.tasks; i++ {
				tasks = append(tasks, func() error {
					atomic.AddInt32(&runTasksCount, 1)
					return nil
				})
			}

			err := Run(tasks, datum.workers, 1)

			require.Nil(t, err)
			require.LessOrEqual(t, runTasksCount, int32(datum.tasks))
		}
	})

	t.Run("eventually. tasks without errors", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		for i := 0; i < tasksCount; i++ {
			tasks = append(tasks, func() error {
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 1

		require.Eventually(
			t,
			func() bool {
				require.Nil(t, Run(tasks, workersCount, maxErrorsCount))
				return true
			},
			time.Duration(tasksCount)*time.Millisecond/time.Duration(workersCount),
			time.Millisecond,
		)

		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
	})

	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 10
		maxErrorsCount := 23
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
	})

	t.Run("tasks without errors", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep

			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 1

		start := time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.NoError(t, err)

		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})
}
