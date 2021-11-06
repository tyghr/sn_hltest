package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"

	_ "github.com/go-sql-driver/mysql"
)

var (
	db      *sqlx.DB
	mysqlLB = "haproxy_mysql_lb:3306"
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

	http.HandleFunc("/mysql", handleMysqlSearchUser())

	err = http.ListenAndServe("0.0.0.0:8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func connectMysql() error {
	dbUrl := "testuser:testpass@tcp(" + mysqlLB + ")/sntest?parseTime=true"
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
