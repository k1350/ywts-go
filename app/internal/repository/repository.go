package repository

import (
	"database/sql"
	"encoding/json"
	"log"
	"strconv"

	_ "github.com/go-sql-driver/mysql"

	"app/internal/configs"
)

type MyDB struct {
	db *sql.DB
}

var Db MyDB

func (m *MyDB) Connect() error {
	sqlConnString := configs.GetConnString("ywts")
	var err error
	m.db, err = sql.Open("mysql", sqlConnString)
	return err
}

func (m *MyDB) Close() {
	if m.db != nil {
		m.db.Close()
	}
}

func (m *MyDB) Count(query string, args ...interface{}) (int, error) {
	var i string
	if err := m.db.QueryRow(query, args...).Scan(&i); err != nil {
		return -1, err
	}
	count, _ := strconv.Atoi(i)
	return count, nil
}

func (m *MyDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	result, err := m.db.Exec(query, args...)
	return result, err
}

func (m *MyDB) Fetch(query string, args ...interface{}) (string, error) {
	rows, err := m.db.Query(query, args...)
	defer rows.Close()

	tableData := make([]map[string]interface{}, 0)
	columns, err := rows.Columns()
	if err != nil {
		return "", err
	}

	count := len(columns)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)

	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		entry := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v
		}
		tableData = append(tableData, entry)
	}

	err = rows.Err()
	if err != nil {
		return "", err
	}

	jresults, err := json.Marshal(tableData)
	if err != nil {
		return "", err
	}
	return string(jresults), nil
}

func (m *MyDB) Transaction(txFunc func(*sql.Tx) error) error {
	tx, err := m.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			log.Println("recover")
			tx.Rollback()
			panic(p)
		} else if err != nil {
			log.Println("rollback")
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()
	err = txFunc(tx)
	return err
}
