package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func connect() (bool, error) {
	dblocal, err := sql.Open("mysql", GetEnvKey("DB_CONF"))
	log.Println("[INFO] Connect with MySQL server...")
	if err != nil {
		return false, err
	}
	dblocal.SetConnMaxLifetime(time.Minute * 3)
	dblocal.SetMaxOpenConns(10)
	dblocal.SetMaxIdleConns(10)

	db = dblocal
	log.Println("[INFO] Connected!")
	return true, nil
}

type GetProps struct {
	Where          string
	Select_columns string
	Limit          int
	Orderby        string
	Groupby        string
}

func Get[T interface{}](table string, options *GetProps) ([]T, error) {
	var rows *sql.Rows
	var sql_text string

	sql_text += "SELECT"

	if db == nil {
		_, err := connect()
		if err != nil {
			log.Println("[ERROR] Error on connect mysql bank!")
			return nil, err
		}
	}

	if options != nil {
		if options.Select_columns != "" {
			sql_text += " " + options.Select_columns
		} else {
			sql_text += " *"

		}

		sql_text += " FROM " + table

		if options.Where != "" {
			sql_text += " WHERE " + options.Where
		}
		if options.Limit != 0 {
			sql_text += " LIMIT " + fmt.Sprint(options.Limit)
		}
		if options.Groupby != "" {
			sql_text += " GROUP BY " + fmt.Sprint(options.Groupby)
		}
		if options.Orderby != "" {
			sql_text += " ORDER BY " + fmt.Sprint(options.Orderby)
		}
	} else {
		sql_text += " * FROM " + table
	}

	log.Println(sql_text)

	rows, err := db.Query(sql_text)

	if err != nil {
		log.Println("[ERROR] Error on Query!" + err.Error())
		return nil, err
	}

	columns, err := rows.Columns()
	if err != nil {
		log.Println("[ERROR] Error on Select Columns!")
		return nil, err
	}

	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	var response []T
	for rows.Next() {

		err = rows.Scan(scanArgs...)
		if err != nil {
			log.Println("[ERROR] Error on Get Row!")
			return nil, err
		}

		// Now do something with the data.
		// Here we just print each column as a string.
		var value any
		item := make(map[string]any)
		for i, col := range values {
			// Here we can check if the value is nil (NULL value)
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			item[columns[i]] = value

		}

		jsonStr, err := json.Marshal(item)
		if err != nil {
			log.Println(err)
		}

		// Convert json string to struct
		var table_data T
		if err := json.Unmarshal(jsonStr, &table_data); err != nil {
			log.Println(err)
		}
		response = append(response, table_data)
	}
	// log.Println(response)

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return response, nil
}

func GetAll[T interface{}](table string) ([]T, error) {
	return Get[T](table, nil)
}

func GetFirst[T interface{}](table string, options *GetProps) (*T, error) {
	response, err := Get[T](table, options)

	if len(response) > 0 {
		return &response[0], err
	}

	return nil, err

}

func Set[T interface{}](table string, value T) (int64, error) {
	if db == nil {
		_, err := connect()
		if err != nil {
			log.Println("[ERROR] Error on connect mysql bank!")
			return 0, err
		}
	}
	strjson, err := json.Marshal(value)
	if err != nil {
		log.Println(err)
		return 0, err

	}

	var mapper map[string]interface{}
	json.Unmarshal([]byte(strjson), &mapper)

	values := make([]string, 0)
	keys := make([]string, 0)
	interrogations := make([]string, 0)

	for k, v := range mapper {
		if v != "" && v != nil {
			values = append(values, "'"+fmt.Sprint(v)+"'")
			interrogations = append(interrogations, "?")
			keys = append(keys, k)
		}
	}

	sql_text := "INSERT INTO `" + table + "` (" + strings.Join(keys, ", ") + ") VALUES (" + strings.Join(values, ", ") + ")"
	// log.Println(sql_text)

	prep, err := db.Prepare(sql_text)

	if err != nil {
		log.Println("Erro ao Preparar Insert => ", err.Error())
		return 0, err
	}

	response, err := prep.Exec()

	if err != nil {
		log.Println("Erro ao Executar Insert => ", err.Error())
		return 0, err
	}

	id, err := response.LastInsertId()
	if err != nil {
		log.Println("Erro ao Contar Insert => ", err.Error())
		return 0, err
	}
	// log.Println(id)
	return id, nil
}

func Update[T interface{}](table string, value T, where string) (int64, error) {
	if db == nil {
		_, err := connect()
		if err != nil {
			log.Println("[ERROR] Error on connect mysql bank!")
			return 0, err
		}
	}
	strjson, err := json.Marshal(value)
	if err != nil {
		log.Println(err)
		return 0, err

	}

	var mapper map[string]interface{}
	json.Unmarshal([]byte(strjson), &mapper)

	values := make([]string, 0)

	for k, v := range mapper {
		if v != "" && v != nil {
			values = append(values, k+" = '"+fmt.Sprint(v)+"'")
		}
	}

	sql_text := "UPDATE `" + table + "` SET " + strings.Join(values, ", ") + " WHERE " + where
	log.Println(sql_text)

	prep, err := db.Prepare(sql_text)

	if err != nil {
		log.Println("Erro ao Preparar Update => ", err.Error())
		return 0, err
	}

	response, err := prep.Exec()

	if err != nil {
		log.Println("Erro ao Executar Update => ", err.Error())
		return 0, err
	}

	id, err := response.RowsAffected()
	if err != nil {
		log.Println("Erro ao Contar Update => ", err.Error())
		return 0, err
	}
	// log.Println(id)
	return id, nil
}

func Delete[T interface{}](table string, where string) (int64, error) {
	if db == nil {
		_, err := connect()
		if err != nil {
			log.Println("[ERROR] Error on connect mysql bank!")
			return 0, err
		}
	}

	sql_text := "DELETE FROM `" + table + "` WHERE " + where
	log.Println(sql_text)

	prep, err := db.Prepare(sql_text)

	if err != nil {
		log.Println("Erro ao Preparar Update => ", err.Error())
		return 0, err
	}

	response, err := prep.Exec()

	if err != nil {
		log.Println("Erro ao Executar Update => ", err.Error())
		return 0, err
	}

	id, err := response.RowsAffected()
	if err != nil {
		log.Println("Erro ao Contar Update => ", err.Error())
		return 0, err
	}
	// log.Println(id)
	return id, nil
}

func RunSQL(sql string) (bool, error) {
	if db == nil {
		_, err := connect()
		if err != nil {
			log.Println("[ERROR] Error on connect mysql bank!")
			return false, err
		}
	}

	prep, err := db.Prepare(sql)

	if err != nil {
		log.Println("Erro ao Preparar Insert => ", err.Error())
		return false, err
	}

	_, err = prep.Exec()

	if err != nil {
		log.Println("Erro ao Executar SQL => ", err.Error())
		return false, err
	}

	return true, nil
}
