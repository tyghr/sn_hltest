package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"

	_ "github.com/go-sql-driver/mysql"

	"github.com/tarantool/go-tarantool"
)

var (
	db   *sqlx.DB
	conn *tarantool.Connection
)

type User struct {
	UserName string `db:"username" json:"username"`
	SurName  string `db:"surname" json:"surname"`
	City     string `db:"city" json:"city"`
}

func main() {
	err := connectMysql()
	if err != nil {
		log.Fatal(err)
	}
	err = connectTarantool()
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/mysql", handleMysqlSearchUser())
	http.HandleFunc("/tarantool", handleTarantoolSearchUser())

	err = http.ListenAndServe("0.0.0.0:8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func connectMysql() error {
	dbUrl := "testuser:testpass@tcp(127.0.0.1:3306)/sntest?parseTime=true"
	var err error
	db, err = sqlx.Connect("mysql", dbUrl)
	if err != nil {
		return err
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	return nil
}

func mysqlSearchUser() ([]User, error) {
	users := []User{}
	if err := db.Ping(); err != nil {
		return users, err
	}

	sqlQuery := `SELECT username,surname,city FROM users
	WHERE username LIKE 'tu%' AND surname LIKE 'Ka%'
	ORDER BY id ASC
	`

	rows, err := db.Queryx(sqlQuery)
	if err != nil {
		return users, err
	}
	defer rows.Close()

	for rows.Next() {
		u := User{}
		err = rows.StructScan(&u)
		if err != nil {
			return users, err
		}
		users = append(users, u)
	}

	return users, rows.Err()
}

func handleMysqlSearchUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := mysqlSearchUser()
		if err != nil {
			log.Fatal(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(users)
	}
}

func connectTarantool() error {
	server := "127.0.0.1:3302"
	opts := tarantool.Opts{
		Timeout:       60 * time.Second,
		Reconnect:     1 * time.Second,
		MaxReconnects: 3,
		User:          "guest",
		Pass:          "",
	}
	var err error
	conn, err = tarantool.Connect(server, opts)
	if err != nil {
		return err
	}

	resp, err := conn.Ping()
	if err != nil {
		return err
	}
	if resp.Code != tarantool.OkCode {
		return fmt.Errorf("ping error")
	}
	return nil
}

func tarantoolSearchUser() (interface{}, error) {
	resp, err := conn.Call("find_users_by_username_and_surname", []interface{}{"tu", "ka"})
	if err != nil {
		return nil, err
	}
	if resp.Code != tarantool.OkCode {
		return nil, fmt.Errorf("ping error")
	}

	return resp.Data, nil
}

func handleTarantoolSearchUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := tarantoolSearchUser()
		if err != nil {
			log.Fatal(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(users)
	}
}
