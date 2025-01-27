// Package pgdb provides functionalities to manage the PostgreSQL database
// that stores authentication tokens.
package pgdb

import (
	encrypt "auth/internal/security"
	"auth/pkg/model"
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

// Manage the database that stores the token

type Repository interface {
	Get(ctx context.Context) (*model.Token, error)
	Put(ctx context.Context) error
}

type DBParams struct {
	DbName   string
	Host     string
	User     string
	Password string
}

type TokenStoreDB struct {
	db   *sql.DB
	data *model.Token
	sync.RWMutex
}

func (t *TokenStoreDB) Get(ctx context.Context) (*model.Token, error) {
	t.RLock()
	defer t.RUnlock()

	selectQuery := `SELECT * FROM tokens ORDER BY lastupdate DESC LIMIT 1`
	err := t.db.QueryRowContext(ctx, selectQuery).Scan(&t.data.ID, &t.data.Key, &t.data.Lastupdate) // more performatic for single row query as does not create *Rows object
	if err != nil {
		return nil, err
	}

	byte_token := []byte(t.data.Key)
	decrypted_token, err := encrypt.DecryptAES(byte_token)
	if err != nil {
		return nil, err
	}
	t.data.Key = string(decrypted_token)
	return t.data, nil
}

func (t *TokenStoreDB) Put(ctx context.Context) error {
	t.Lock()
	defer t.Unlock()
	t.data.Lastupdate = time.Now().Format(time.RFC3339)
	fmt.Println("reached here")
	insertQuery := `INSERT INTO tokens (key, lastupdate) VALUES ($1, $2)`
	_, err := t.db.ExecContext(ctx, insertQuery, t.data.Key, t.data.Lastupdate)
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
	connectionString := fmt.Sprintf("dbname=%s host=%s user=%s password=%s sslmode=disable", config.DbName, config.Host, config.User, config.Password)
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

	return t, nil
}

func NewTokenStoreDB(config DBParams) (Repository, error) {
	t := &TokenStoreDB{data: &model.Token{}}
	return t.openConnection(config)
}

/*
func main() {
	// for testing!

	store, err := NewTokenStoreDB(DBParams{dbName: os.Getenv("TSTORE_DB_NAME"),
		host:     os.Getenv("TSTORE_DB_HOST"),
		user:     os.Getenv("TSTORE_DB_USER"),
		password: os.Getenv("TSTORE_DB_PASSWORD")})

	if err != nil {
		panic(err)
	}

	err = store.Put(context.Background())
	if err != nil {
		panic(err)
	}

	token, err := store.Get(context.Background())
	if err != nil {
		panic(err)
	}

	bytes, err := json.MarshalIndent(token, "", " ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bytes))
}
*/
