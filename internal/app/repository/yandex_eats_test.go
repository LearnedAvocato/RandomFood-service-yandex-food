package repository

import (
	"context"
	"testing"
	"yandex-food/internal/pkg/datastruct"

	"github.com/stretchr/testify/require"
)

var testLog = &datastruct.RequestLog{
	Latitude:  55.696233,
	Longitude: 37.570431,
}

func TestCreateRepository(t *testing.T) {
	ctx := context.Background()

	_, err := NewRepository(ctx)
	require.Nil(t, err)
}

func TestLogRequest(t *testing.T) {
	var err error
	ctx := context.Background()

	repo, err := NewRepository(ctx)
	require.Nil(t, err)

	err = repo.LogRequest(ctx, testLog)
	require.Nil(t, err)
}
