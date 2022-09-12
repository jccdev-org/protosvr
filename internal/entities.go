package internal

import (
	"context"
	"errors"
	"fmt"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"regexp"
	"strings"
)

type EntityDb struct {
	Id          int64
	Name        string
	CreatedAt   pgtype.Timestamptz
	CreatedById int64
	UpdatedAt   pgtype.Timestamptz
	UpdatedById int64
	Deleted     bool
}

type EntityField struct {
	Name  string
	Type  string
	Value interface{}
}

type GetAllOptions struct {
	Take int64
	Skip int64
}

var ErrInvalidEntityName = errors.New("invalid entity name")

// GetCanonicalEntityName
// entity names are case-insensitive and only allow special chars [-,_]
// due to postgres limitation, - gets is treated the same as _
func GetCanonicalEntityName(value string) (string, error) {
	var res string
	nameReg, err := regexp.Compile("^[a-zA-Z0-9_-]*$")

	if err != nil {
		return res, WrapError(err)
	}

	if !nameReg.MatchString(value) {
		return res, WrapError(ErrInvalidEntityName)
	}

	res = value
	res = strings.ReplaceAll(res, "-", "_")
	res = strings.ToLower(res)

	return res, nil
}

func GetEntityMetaList() ([]*EntityDb, error) {
	var values []*EntityDb
	err := pgxscan.Select(context.Background(), Db, &values, "select * from entities")
	if err != nil {
		return values, WrapError(err)
	}
	return values, nil
}

func GetSingleEntity(entityName string, id int64) (map[string]interface{}, error) {
	res := make(map[string]interface{})
	query := fmt.Sprintf("select * from %s where id=$1", pgx.Identifier{entityName}.Sanitize())
	rows, err := Db.Query(context.Background(), query, id)
	if err != nil {
		return res, WrapError(err)
	}
	defer rows.Close()

	hasRow := rows.Next()
	if !hasRow {
		return nil, nil
	}

	res, err = ConvertDbRowToMap(rows)
	if err != nil {
		return nil, WrapError(err)
	}

	return res, nil
}

func GetAllEntities(entityName string, options *GetAllOptions) ([]map[string]interface{}, error) {
	res := make([]map[string]interface{}, 0)
	query := fmt.Sprintf("select * from %s", pgx.Identifier{entityName}.Sanitize())

	queryArgs := make([]interface{}, 0)

	if options != nil && options.Skip > 0 {
		queryArgs = append(queryArgs, options.Skip)
		query = query + fmt.Sprintf(" offset $%d", len(queryArgs))
	}

	if options != nil && options.Take > 0 {
		queryArgs = append(queryArgs, options.Take)
		query = query + fmt.Sprintf(" limit $%d", len(queryArgs))
	}

	rows, err := Db.Query(context.Background(), query, queryArgs...)
	defer rows.Close()

	if err != nil {
		return nil, WrapError(err)
	}

	// Iterate through the result set
	for rows.Next() {
		rowRes, err := ConvertDbRowToMap(rows)
		if err != nil {
			return nil, WrapError(err)
		}
		res = append(res, rowRes)
	}

	return res, nil
}

//func CreateEntity(entityName string, values map[string]interface{}) (map[string]interface{}, error) {
//	res := make(map[string]interface{})
//	query := fmt.Sprintf("select * from %s where id=$1", pgx.Identifier{entityName}.Sanitize())
//	rows, err := Db.Query(context.Background(), query, id)
//	if err != nil {
//		return res, WrapError(err)
//	}
//	defer rows.Close()
//
//	hasRow := rows.Next()
//	if !hasRow {
//		return nil, nil
//	}
//
//	res, err = ConvertDbRowToMap(rows)
//	if err != nil {
//		return nil, WrapError(err)
//	}
//
//	return res, nil
//}

func GetEntityMeta() {

}

func CreateEntityMeta() {

}
