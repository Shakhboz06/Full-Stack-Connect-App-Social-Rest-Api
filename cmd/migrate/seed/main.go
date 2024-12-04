package main

import (
	"log"
	"go-project/internal/db"
	"go-project/internal/env"
	"go-project/internal/store"
)

func main() {
	addr := env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/app?sslmode=disable")
	conn, err := db.New(addr, 3, 3, "15m")

	if err != nil{
		log.Fatal(err)
	}

	defer conn.Close()
	
	store := store.NewPostgresStorage(conn)
	db.Seed(store, conn)
}
