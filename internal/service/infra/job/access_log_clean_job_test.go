package job

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

type fakeAccessLogCleanService struct {
	called      bool
	exceedDay   int
	deleteLimit int
	err         error
}

func (f *fakeAccessLogCleanService) CleanAccessLog(_ context.Context, exceedDay int, deleteLimit int) (int, error) {
	f.called = true
	f.exceedDay = exceedDay
	f.deleteLimit = deleteLimit
	return 5, f.err
}

func TestAccessLogCleanJob_GetHandlerName(t *testing.T) {
	svc := &fakeAccessLogCleanService{}
	job := &AccessLogCleanJob{apiAccessLogService: svc}
	require.Equal(t, "accessLogCleanJob", job.GetHandlerName())
}

func TestAccessLogCleanJob_Execute(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &fakeAccessLogCleanService{}
		job := &AccessLogCleanJob{apiAccessLogService: svc}

		err := job.Execute(context.Background(), "")
		require.NoError(t, err)
		require.True(t, svc.called)
		require.Equal(t, accessLogCleanRetainDay, svc.exceedDay)
		require.Equal(t, accessLogCleanDeleteLimit, svc.deleteLimit)
	})

	t.Run("service error", func(t *testing.T) {
		svc := &fakeAccessLogCleanService{err: errors.New("boom")}
		job := &AccessLogCleanJob{apiAccessLogService: svc}

		err := job.Execute(context.Background(), "")
		require.Error(t, err)
		require.True(t, svc.called)
	})
}
