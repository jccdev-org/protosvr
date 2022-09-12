package internal

import "github.com/jackc/pgx/v4"

func ConvertDbRowToMap(rows pgx.Rows) (map[string]interface{}, error) {
	res := make(map[string]interface{})
	fieldDescs := rows.FieldDescriptions()
	fieldLen := len(fieldDescs)
	fieldValues := make([]interface{}, fieldLen)

	for i := 0; i < fieldLen; i++ {
		var fieldValue interface{}
		fieldValues[i] = &fieldValue
	}

	err := rows.Scan(fieldValues...)
	if err != nil {
		return nil, WrapError(err)
	}

	for i, fieldDesc := range fieldDescs {
		fieldName := string(fieldDesc.Name)
		res[fieldName] = fieldValues[i]
	}

	return res, nil
}
