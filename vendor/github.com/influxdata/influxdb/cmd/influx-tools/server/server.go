package server

import (
	"time"

	"github.com/influxdata/influxdb/services/meta"
	"github.com/influxdata/influxdb/tsdb"
	"go.uber.org/zap"
)

type Interface interface {
	Open(path string) error
	Close()
	MetaClient() MetaClient
	TSDBConfig() tsdb.Config
	Logger() *zap.Logger
}

type MetaClient interface {
	Database(name string) *meta.DatabaseInfo
	RetentionPolicy(database, name string) (rpi *meta.RetentionPolicyInfo, err error)
	ShardGroupsByTimeRange(database, policy string, min, max time.Time) (a []meta.ShardGroupInfo, err error)
	CreateRetentionPolicy(database string, spec *meta.RetentionPolicySpec, makeDefault bool) (*meta.RetentionPolicyInfo, error)
	UpdateRetentionPolicy(database, name string, rpu *meta.RetentionPolicyUpdate, makeDefault bool) error
	CreateDatabase(name string) (*meta.DatabaseInfo, error)
	CreateDatabaseWithRetentionPolicy(name string, spec *meta.RetentionPolicySpec) (*meta.DatabaseInfo, error)
	DeleteShardGroup(database, policy string, id uint64) error
	CreateShardGroup(database, policy string, timestamp time.Time) (*meta.ShardGroupInfo, error)
}
