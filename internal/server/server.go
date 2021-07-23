package server

import (
	etcd "github.com/go-kratos/etcd/registry"
	"github.com/go-kratos/kratos-layout/internal/conf"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/google/wire"
	clientv3 "go.etcd.io/etcd/client/v3"
	google_grpc "google.golang.org/grpc"
	"time"
)

// ProviderSet is server providers.
var ProviderSet = wire.NewSet(NewHTTPServer, NewGRPCServer, NewRegistrar)

// NewRegistrar 服务注册
func NewRegistrar(conf *conf.Registry) registry.Registrar {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   conf.Etcd.Addr,
		DialTimeout: 3 * time.Second,
		DialOptions: []google_grpc.DialOption{google_grpc.WithBlock()},
	})
	if err != nil {
		panic(err)
	}
	r := etcd.New(client)
	return r
}
