package main

import (
	"auth/internal/controller"
	httphandler "auth/internal/handler/http"
	"auth/internal/refresher"
	"auth/internal/repository"
	file "auth/internal/repository/file"
	"auth/internal/repository/pgdb"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

// Returns the token that client can use to access Olist ERP API, refreshes the token automatically and allow to store a new token
func main() {
	// Refresh token every 1 minute
	go refresher.RefreshToken()

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
	if err != nil {
		panic(err)
	}
	compositeRepo := repository.NewCompositeRepository(fileRepo, dbRepo)

	ctrl := controller.New(compositeRepo)

	h := httphandler.New(ctrl)
	r := mux.NewRouter()
	r.HandleFunc("/auth", h.GetToken).Methods("GET")
	r.HandleFunc("/auth", h.PutToken).Methods("PUT")
	if err := http.ListenAndServe(":8081", r); err != nil {
		panic(err)
	}

}
