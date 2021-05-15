package app

import (
	"context"
	"github.com/core-go/health"
	"github.com/gocql/gocql"
	"go-service/cql"
	"go-service/internal/handlers"
	"go-service/internal/services"
	"time"
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
	HealthHandler *health.HealthHandler
	UserHandler   *handlers.UserHandler
}

func NewApp(context context.Context, root Root) (*ApplicationContext, error) {
	// connect to the cluster
	cluster := gocql.NewCluster(root.Cql.PublicIp)
	cluster.Consistency = gocql.Quorum
	cluster.ProtoVersion = 4
	cluster.Timeout = time.Second * 1000
	cluster.ConnectTimeout = time.Second * 1000
	cluster.Authenticator = gocql.PasswordAuthenticator{Username: root.Cql.UserName, Password: root.Cql.Password}
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

	userService := services.NewUserService(cluster)
	userHandler := handlers.NewUserHandler(userService)

	cqlChecker := cql.NewHealthChecker(cluster)
	healthHandler := health.NewHealthHandler(cqlChecker)

	return &ApplicationContext{
		HealthHandler: healthHandler,
		UserHandler:   userHandler,
	}, nil
}
