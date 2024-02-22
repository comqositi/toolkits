package utils

import (
	"database/sql"
	"errors"
	"time"

	"github.com/jinzhu/copier"
)

var (
	timeToStringConverter = copier.TypeConverter{
		SrcType: time.Time{},
		DstType: copier.String,
		Fn: func(src interface{}) (interface{}, error) {
			s, ok := src.(time.Time)
			if !ok {
				return nil, errors.New("src type not matching")
			}
			return s.Format("2006-01-02 15:04:05"), nil
		},
	}

	stringToSqlNullStringConverter = copier.TypeConverter{
		SrcType: copier.String,
		DstType: sql.NullString{},
		Fn: func(src interface{}) (interface{}, error) {
			s, ok := src.(string)
			if !ok {
				return nil, errors.New("src type not matching")
			}

			if s == "" {
				return sql.NullString{
					Valid: false,
				}, nil
			}

			return sql.NullString{
				String: s,
				Valid:  true,
			}, nil
		},
	}

	SqlNullTimeToStringConverter = copier.TypeConverter{
		SrcType: sql.NullTime{},
		DstType: copier.String,
		Fn: func(src interface{}) (interface{}, error) {
			s, ok := src.(sql.NullTime)
			if !ok {
				return nil, errors.New("src type not matching")
			}
			if !s.Valid {
				return "", nil
			}
			return s.Time.Format("2006-01-02 15:04:05"), nil
		},
	}
)

// Copy things
func Copy(toValue interface{}, fromValue interface{}) (err error) {
	return copier.CopyWithOption(toValue, fromValue, copier.Option{
		Converters: []copier.TypeConverter{
			timeToStringConverter,
			stringToSqlNullStringConverter,
			SqlNullTimeToStringConverter,
		},
	})
}
