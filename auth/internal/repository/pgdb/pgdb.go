// Package pgdb provides functionalities to manage the PostgreSQL database
// that stores authentication tokens.
package main

import (
	"auth/pkg/model"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

// Manage the database that stores the token

type Repository interface {
	Get(ctx context.Context) (*model.Token, error)
	Put(ctx context.Context, token *model.Token) error
}

type DBParams struct {
	dbName   string
	host     string
	user     string
	password string
}

type TokenStoreDB struct {
	db *sql.DB
	sync.RWMutex
}

func (t *TokenStoreDB) Get(ctx context.Context) (*model.Token, error) {
	t.RLock()
	defer t.RUnlock()

	var token model.Token

	selectQuery := `SELECT * FROM tokens ORDER BY lastupdate DESC LIMIT 1`
	err := t.db.QueryRowContext(ctx, selectQuery).Scan(&token.ID, &token.Key, &token.Lastupdate) // more performatic for single row query as does not create *Rows object
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (t *TokenStoreDB) Put(ctx context.Context, token *model.Token) error {
	t.Lock()
	defer t.Unlock()

	token.Lastupdate = time.Now().Format(time.RFC3339)

	insertQuery := `INSERT INTO tokens (key, lastupdate) VALUES ($1, $2)`
	_, err := t.db.ExecContext(ctx, insertQuery, token.Key, token.Lastupdate)
	if err != nil {
		return err
	}
	return nil
}

func (t *TokenStoreDB) createTable() error {
	createQuery := `CREATE TABLE IF NOT EXISTS tokens (
						id SERIAL PRIMARY KEY, 
						key TEXT, 
						lastupdate TIMESTAMP)`

	_, err := t.db.Exec(createQuery)
	if err != nil {
		return err
	}
	return nil
}

func (t *TokenStoreDB) openConnection(config DBParams) (Repository, error) {
	var err error
	// create connection string and attempt to connect
	connectionString := fmt.Sprintf("dbname=%s host=%s user=%s password=%s sslmode=disable", config.dbName, config.host, config.user, config.password)
	t.db, err = sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to open db connection: %w", err)
	}

	// establish connection
	err = t.db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to establish db connection: %w", err)
	}

	// create table if it does not exist
	err = t.createTable()
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	db := &TokenStoreDB{db: t.db}
	return db, nil
}

func NewTokenStoreDB(config DBParams) (Repository, error) {
	t := &TokenStoreDB{}
	return t.openConnection(config)
}

func main() {
	// for testing!
	token := &model.Token{Key: "34adade1-6ac4-4a5a-a394-2c47177a9311.95c5eb2f-e8a8-4f48-8bf2-fa2882f6c607.3dcda8a1-a6ef-4964-adcc-d0a5e1b8eebb"}
	/*
		export TSTORE_DB_NAME="postgres"
		export TSTORE_DB_HOST="localhost"
		export TSTORE_DB_USER="postgres"
		export TSTORE_DB_PASSWORD="postgres"
	*/
	store, err := NewTokenStoreDB(DBParams{dbName: os.Getenv("TSTORE_DB_NAME"),
		host:     os.Getenv("TSTORE_DB_HOST"),
		user:     os.Getenv("TSTORE_DB_USER"),
		password: os.Getenv("TSTORE_DB_PASSWORD")})

	if err != nil {
		panic(err)
	}

	err = store.Put(context.Background(), token)
	if err != nil {
		panic(err)
	}

	token, err = store.Get(context.Background())
	if err != nil {
		panic(err)
	}

	bytes, err := json.MarshalIndent(token, "", " ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bytes))
}
