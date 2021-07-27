package data

import (
	"context"
	"github.com/go-kratos/kratos-layout/internal/conf"
	"github.com/go-kratos/kratos-layout/internal/data/ent"
	"github.com/go-kratos/kratos/v2/log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewGreeterRepo)

// Data .
type Data struct {
	// TODO warpped database client
	db *ent.Client
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	// 连接数据库
	client, err := ent.Open(
		c.Database.Driver,
		c.Database.Source,
		ent.AlternateSchema(ent.SchemaConfig{}),
	)
	if err != nil {
		log.NewHelper(logger).Fatalf("failed opening connection to sqlite: %v", err)
	}
	defer func() {
		_ = client.Close()
	}()
	if err := client.Schema.Create(context.Background()); err != nil {
		log.NewHelper(logger).Fatalf("failed creating schema resources: %v", err)
	}
	return &Data{
		db: client,
	}, cleanup, nil
}
