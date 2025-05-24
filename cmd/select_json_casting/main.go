package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type UsersProps struct {
	ID          int64
	Value       sql.NullInt64
	Additionals []byte
}

type Additionals struct {
	Nome  string `json:"nome"`
	Idade int64  `json:"idade"`
}

type Users struct {
	ID          int64
	Value       int64
	Additionals Additionals
}

func (up *UsersProps) toDomain() (Users, error) {
	var u Users
	u.ID = up.ID

	if up.Value.Valid {
		u.Value = up.Value.Int64
	} else {
		u.Value = 0
	}

	if len(up.Additionals) > 0 {
		err := json.Unmarshal(up.Additionals, &u.Additionals)
		if err != nil {
			return Users{}, fmt.Errorf("additionals unmarshal error: %w", err)
		}
	}

	return u, nil
}

func main() {
	dsn := "root:12345678@tcp(127.0.0.1:3306)/test_database?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("db connection error: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("ping db error: %v", err)
	}

	rows, err := db.Query("SELECT id, value, additionals FROM users")
	if err != nil {
		log.Fatalf("select error: %v", err)
	}
	defer rows.Close()

	var domainUsers []Users

	for rows.Next() {
		var up UsersProps

		err := rows.Scan(&up.ID, &up.Value, &up.Additionals)
		if err != nil {
			log.Printf("scan error %v", err)
			continue
		}

		domainUser, err := up.toDomain()
		if err != nil {
			log.Printf("toDomain error - id %d: %v", up.ID, err)
			continue
		}

		domainUsers = append(domainUsers, domainUser)
	}

	if err = rows.Err(); err != nil {
		log.Fatalf("iteration error: %v", err)
	}

	for _, user := range domainUsers {
		fmt.Printf("ID: %d, Value: %d, Additionals: %+v\n", user.ID, user.Value, user.Additionals)
	}
}
