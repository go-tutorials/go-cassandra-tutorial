package cql

import (
	"errors"
	s "github.com/core-go/sql"
	"github.com/gocql/gocql"
	"reflect"
)

func BuildParam(i int) string {
	return "?"
}

func Patch(cluster *gocql.ClusterConfig, table string, model map[string]interface{}, modelType reflect.Type) (int64, error) {
	session, err := cluster.CreateSession()
	if err != nil{
		return 0, err
	}
	idcolumNames, idJsonName := s.FindNames(modelType)
	columNames := s.FindJsonName(modelType)
	query, value := s.BuildPatch(table, model, columNames, idJsonName, idcolumNames, BuildParam)
	if query == "" {
		return 0, errors.New("fail to build query")
	}
	err = session.Query(query, value...).Exec()
	if err != nil {
		return -1, err
	}
	return 1, nil
}
