package services

import (
	"context"
	"fmt"
	"github.com/gocql/gocql"
	"strings"

	. "go-service/internal/models"
)

type CqlUserService struct {
	Cluster *gocql.ClusterConfig
}

func NewUserService(db *gocql.ClusterConfig) *CqlUserService {
	return &CqlUserService{Cluster: db}
}


func (m *CqlUserService) All(ctx context.Context) (*[]User, error) {
	session, err := m.Cluster.CreateSession()
	if err != nil{
		return nil, err
	}
	query := "select id, username, email, phone, date_of_birth from users"
	rows := session.Query(query).Iter()
	var result []User
	var user User
	for rows.Scan(&user.Id, &user.Username, &user.Phone, &user.Email, &user.DateOfBirth) {
		result = append(result, user)
	}
	return &result, nil
}

func (m *CqlUserService) Load(ctx context.Context, id string) (*User, error) {
	session, err := m.Cluster.CreateSession()
	if err != nil{
		return nil, err
	}
	var user User
	query := "select id, username, email, phone, date_of_birth from users where id = ?"
	err = session.Query(query, id).Scan(&user.Id, &user.Username, &user.Email, &user.Phone, &user.DateOfBirth)
	if err != nil {
		errMsg := err.Error()
		if strings.Compare(fmt.Sprintf(errMsg), "0 row(s) returned") == 0 {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &user, nil
}

func (m *CqlUserService) Insert(ctx context.Context, user *User) (int64, error) {
	session, err := m.Cluster.CreateSession()
	if err != nil{
		return 0, err
	}
	query := "insert into users (id, username, email, phone, date_of_birth) values (?, ?, ?, ?, ?)"
	err = session.Query(query, user.Id, user.Username, user.Email, user.Phone, user.DateOfBirth).Exec()
	if err != nil {
		return -1, nil
	}
	return 1, nil
}

func (m *CqlUserService) Update(ctx context.Context, user *User) (int64, error) {
	session, err := m.Cluster.CreateSession()
	if err != nil{
		return 0, err
	}
	query := "update users set username = ?, email = ?, phone = ?, date_of_birth = ? where id = ?"
	err = session.Query(query, user.Username, user.Email, user.Phone, user.DateOfBirth, user.Id).Exec()
	if err != nil {
		return -1, err
	}
	return 1, nil
}

func (m *CqlUserService) Delete(ctx context.Context, id string) (int64, error) {
	session, err := m.Cluster.CreateSession()
	if err != nil{
		return 0, err
	}
	query := "delete from users where id = ?"
	er1 := session.Query(query, id).Exec()
	if er1 != nil {
		return -1, er1
	}
	return 1, nil
}
