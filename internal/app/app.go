package app

import (
	"context"
	"github.com/core-go/health"
	"github.com/gocql/gocql"
	"time"

	"go-service/internal/handler"
	"go-service/internal/service"
	"go-service/pkg/cql"
)

const (
	Keyspace = `masterdata`

	CreateKeyspace = `create keyspace if not exists masterdata with replication = {'class':'SimpleStrategy', 'replication_factor':1}`

	CreateTable = `
					create table if not exists users (
					id varchar,
					username varchar,
					email varchar,
					phone varchar,
					date_of_birth date,
					primary key (id)
	)`
)

type ApplicationContext struct {
	Health *health.Handler
	User   *handler.UserHandler
}

func NewApp(ctx context.Context, config Config) (*ApplicationContext, error) {
	// connect to the cluster
	cluster := gocql.NewCluster(config.Cql.PublicIp)
	cluster.Consistency = gocql.Quorum
	cluster.ProtoVersion = 4
	cluster.Timeout = time.Second * 1000
	cluster.ConnectTimeout = time.Second * 1000
	cluster.Authenticator = gocql.PasswordAuthenticator{Username: config.Cql.UserName, Password: config.Cql.Password}
	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}
	//defer session.Close()

	// create keyspaces
	err = session.Query(CreateKeyspace).Exec()
	if err != nil {
		return nil, err
	}

	//switch keyspaces
	session.Close()
	cluster.Keyspace = Keyspace
	session, err = cluster.CreateSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()

	// create table
	err = session.Query(CreateTable).Exec()
	if err != nil {
		return nil, err
	}

	userService := service.NewUserService(cluster)
	userHandler := handler.NewUserHandler(userService)

	cqlChecker := cql.NewHealthChecker(cluster)
	healthHandler := health.NewHandler(cqlChecker)

	return &ApplicationContext{
		Health: healthHandler,
		User:   userHandler,
	}, nil
}
