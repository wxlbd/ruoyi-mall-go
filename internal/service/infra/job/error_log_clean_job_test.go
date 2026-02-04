package job

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

type fakeErrorLogCleanService struct {
	called      bool
	exceedDay   int
	deleteLimit int
	err         error
}

func (f *fakeErrorLogCleanService) CleanErrorLog(_ context.Context, exceedDay int, deleteLimit int) (int, error) {
	f.called = true
	f.exceedDay = exceedDay
	f.deleteLimit = deleteLimit
	return 5, f.err
}

func TestErrorLogCleanJob_GetHandlerName(t *testing.T) {
	svc := &fakeErrorLogCleanService{}
	job := &ErrorLogCleanJob{apiErrorLogService: svc}
	require.Equal(t, "errorLogCleanJob", job.GetHandlerName())
}

func TestErrorLogCleanJob_Execute(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &fakeErrorLogCleanService{}
		job := &ErrorLogCleanJob{apiErrorLogService: svc}

		err := job.Execute(context.Background(), "")
		require.NoError(t, err)
		require.True(t, svc.called)
		require.Equal(t, errorLogCleanRetainDay, svc.exceedDay)
		require.Equal(t, errorLogCleanDeleteLimit, svc.deleteLimit)
	})

	t.Run("service error", func(t *testing.T) {
		svc := &fakeErrorLogCleanService{err: errors.New("boom")}
		job := &ErrorLogCleanJob{apiErrorLogService: svc}

		err := job.Execute(context.Background(), "")
		require.Error(t, err)
		require.True(t, svc.called)
	})
}
