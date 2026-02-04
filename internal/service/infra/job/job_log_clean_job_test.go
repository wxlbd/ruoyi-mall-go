package job

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

type fakeJobLogCleanService struct {
	called      bool
	exceedDay   int
	deleteLimit int
	err         error
}

func (f *fakeJobLogCleanService) CleanJobLog(_ context.Context, exceedDay int, deleteLimit int) (int, error) {
	f.called = true
	f.exceedDay = exceedDay
	f.deleteLimit = deleteLimit
	return 5, f.err
}

func TestJobLogCleanJob_GetHandlerName(t *testing.T) {
	svc := &fakeJobLogCleanService{}
	job := &JobLogCleanJob{jobLogService: svc}
	require.Equal(t, "jobLogCleanJob", job.GetHandlerName())
}

func TestJobLogCleanJob_Execute(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &fakeJobLogCleanService{}
		job := &JobLogCleanJob{jobLogService: svc}

		err := job.Execute(context.Background(), "")
		require.NoError(t, err)
		require.True(t, svc.called)
		require.Equal(t, jobCleanRetainDay, svc.exceedDay)
		require.Equal(t, jobCleanDeleteLim, svc.deleteLimit)
	})

	t.Run("service error", func(t *testing.T) {
		svc := &fakeJobLogCleanService{err: errors.New("boom")}
		job := &JobLogCleanJob{jobLogService: svc}

		err := job.Execute(context.Background(), "")
		require.Error(t, err)
		require.True(t, svc.called)
	})
}
