package main

import (
	olistmediator "auth/internal/controller"
	httphandler "auth/internal/handler/http"
	"auth/internal/repository"
	file "auth/internal/repository/file"
	"auth/internal/repository/pgdb"
	"net/http"
	"os"
)

// Returns the token that client can use to access Olist ERP API, refreshes the token automatically and allow to store a new token
func main() {
	fileRepo, err := file.NewTokenStoreFile()
	if err != nil {
		panic(err)
	}
	dbRepo, err := pgdb.NewTokenStoreDB(pgdb.DBParams{
		DbName:   os.Getenv("TSTORE_DB_NAME"),
		Host:     os.Getenv("TSTORE_DB_HOST"),
		User:     os.Getenv("TSTORE_DB_USER"),
		Password: os.Getenv("TSTORE_DB_PASSWORD"),
	})

	compositeRepo := &repository.CompositeTokenRepository{
		Primary:   dbRepo,
		Secondary: fileRepo,
	}

	ctrl := olistmediator.New(compositeRepo)
	h := httphandler.New(ctrl)
	http.Handle("/auth", http.HandlerFunc(h.GetToken))
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}

}
