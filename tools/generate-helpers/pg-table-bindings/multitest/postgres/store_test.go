// Code generated by pg-bindings generator. DO NOT EDIT.

//go:build sql_integration

package postgres

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stackrox/rox/pkg/features"
	"github.com/stackrox/rox/pkg/fixtures"
	"github.com/stackrox/rox/pkg/postgres/pgtest"
	"github.com/stackrox/rox/pkg/testutils/envisolator"
	"github.com/stretchr/testify/suite"
)

type MultikeyStoreSuite struct {
	suite.Suite
	envIsolator *envisolator.EnvIsolator
}

func TestMultikeyStore(t *testing.T) {
	suite.Run(t, new(MultikeyStoreSuite))
}

func (s *MultikeyStoreSuite) SetupTest() {
	s.envIsolator = envisolator.NewEnvIsolator(s.T())
	s.envIsolator.Setenv(features.PostgresDatastore.EnvVar(), "true")

	if !features.PostgresDatastore.Enabled() {
		s.T().Skip("Skip postgres store tests")
		s.T().SkipNow()
	}
}

func (s *MultikeyStoreSuite) TearDownTest() {
	s.envIsolator.RestoreAll()
}

func (s *MultikeyStoreSuite) TestStore() {
	source := pgtest.GetConnectionString(s.T())
	config, err := pgxpool.ParseConfig(source)
	if err != nil {
		panic(err)
	}
	pool, err := pgxpool.ConnectConfig(context.Background(), config)
	s.NoError(err)
	defer pool.Close()

	Destroy(pool)
	store := New(pool)

	testMultiKeyStruct := fixtures.GetTestMultiKeyStruct()
	foundTestMultiKeyStruct, exists, err := store.Get(testMultiKeyStruct.GetId())
	s.NoError(err)
	s.False(exists)
	s.Nil(foundTestMultiKeyStruct)

	s.NoError(store.Upsert(testMultiKeyStruct))
	foundTestMultiKeyStruct, exists, err = store.Get(testMultiKeyStruct.GetId())
	s.NoError(err)
	s.True(exists)
	s.Equal(testMultiKeyStruct, foundTestMultiKeyStruct)

	testMultiKeyStructCount, err := store.Count()
	s.NoError(err)
	s.Equal(testMultiKeyStructCount, 1)

	testMultiKeyStructExists, err := store.Exists(testMultiKeyStruct.GetId())
	s.NoError(err)
	s.True(testMultiKeyStructExists)
	s.NoError(store.Upsert(testMultiKeyStruct))

	foundTestMultiKeyStruct, exists, err = store.Get(testMultiKeyStruct.GetId())
	s.NoError(err)
	s.True(exists)
	s.Equal(testMultiKeyStruct, foundTestMultiKeyStruct)

	s.NoError(store.Delete(testMultiKeyStruct.GetId()))
	foundTestMultiKeyStruct, exists, err = store.Get(testMultiKeyStruct.GetId())
	s.NoError(err)
	s.False(exists)
	s.Nil(foundTestMultiKeyStruct)
}
