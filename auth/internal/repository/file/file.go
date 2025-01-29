package file

// Manage the file that stores the token

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	encrypt "auth/internal/security"

	"auth/pkg/model"
)

type Repository interface {
	Get(ctx context.Context) (*model.Token, error)
	Put(ctx context.Context) error
}

type TokenStoreFile struct {
	sync.RWMutex
	data *model.Token
}

func (t *TokenStoreFile) Get(_ context.Context) (*model.Token, error) {
	t.RLock()
	defer t.RUnlock()

	const filename = ".token_store.json"

	dir := os.Getenv("TOKEN_STORE_DIR")
	if dir == "" {
		return t.data, fmt.Errorf("TOKEN_STORE_DIR not set")
	}
	filepath := strings.Join([]string{dir, filename}, string(os.PathSeparator))
	file, err := os.Open(filepath)
	if err != nil {
		return t.data, err
	}

	defer file.Close()
	err = json.NewDecoder(file).Decode(&t.data)
	if err != nil {
		return t.data, err
	}

	byte_token := []byte(t.data.Key)
	decrypted_token, err := encrypt.DecryptAES(byte_token)
	if err != nil {
		return nil, err
	}

	t.data.Key = string(decrypted_token)
	return t.data, err
}

func (t *TokenStoreFile) Put(_ context.Context) error {
	t.Lock()
	defer t.Unlock()

	// decrypt> 3) encrypt the token
	encrypted_token, err := encrypt.EncryptAES([]byte(t.data.Key))
	if err != nil {
		return err
	}

	// decrypt> 4) update the token
	t.data.Key = string(encrypted_token)

	// decrypt> 5) update the last updated time
	t.data.ID++
	t.data.Lastupdate = time.Now().Format(time.RFC3339)

	// decrypt> 6) save to the json file
	const filename = ".token_store.json"
	bytes, err := json.MarshalIndent(t.data, "", "   ")
	if err != nil {
		return err
	}
	os.WriteFile(filename, bytes, 0644)

	return nil
}

func NewTokenStoreFile() (Repository, error) {
	d := &TokenStoreFile{data: &model.Token{}}

	return d, nil
}

// for testing!
/*
func main() {
	store, err := NewTokenStoreFile()
	if err != nil {
		panic(err)
	}

	// The TOKEN.KEY is encrypted (AES) using a key (hex-encoded) and encoded (BASE64)
	store.Put(context.Background())

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
