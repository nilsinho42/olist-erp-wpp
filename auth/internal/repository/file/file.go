package main

// Manage the file that stores the token

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	encrypt "auth/internal/security"

	"auth/pkg/model"
)

type Repository struct {
	sync.RWMutex
	data *model.Token
}

func (r *Repository) Get(_ context.Context) (*model.Token, error) {
	r.RLock()
	defer r.RUnlock()

	const filename = ".token_store.json"

	filepath, err := os.Getwd()
	if err != nil {
		return r.data, nil
	}

	filepath = strings.Join([]string{filepath, filename}, string(os.PathSeparator))
	file, err := os.Open(filepath)
	if err != nil {
		return r.data, err
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		return r.data, err
	}

	err = json.Unmarshal(data, &r.data)
	if err != nil {
		return nil, err
	}

	byte_token := []byte(r.data.Key)
	decrypted_token, err := encrypt.DecryptAES(byte_token)
	if err != nil {
		return nil, err
	}

	r.data.Key = string(decrypted_token)
	return r.data, err
}

func (r *Repository) Put(_ context.Context) error {
	r.Lock()
	defer r.Unlock()
	// decrypt> 3) encrypt the token
	encrypted_token, err := encrypt.EncryptAES([]byte(r.data.Key))
	if err != nil {
		return err
	}

	// decrypt> 4) update the token
	r.data.Key = string(encrypted_token)

	// decrypt> 5) update the last updated time
	r.data.Lastupdate = time.Now().Format(time.RFC3339)

	// decrypt> 6) save to the json file
	const filename = ".token_store.json"
	bytes, err := json.MarshalIndent(r.data, "", "   ")
	if err != nil {
		return err
	}
	os.WriteFile(filename, bytes, 0644)

	return nil
}

// for testing!
func main() {
	r := Repository{}
	r.data = &model.Token{}
	r.data.Key = "34adade1-6ac4-4a5a-a394-2c47177a9311.95c5eb2f-e8a8-4f48-8bf2-fa2882f6c607.3dcda8a1-a6ef-4964-adcc-d0a5e1b8eebb"

	// The TOKEN.KEY is encrypted (AES) using a key (hex-encoded) and encoded (BASE64)
	r.Put(context.Background())

	token, err := r.Get(context.Background())
	if err != nil {
		panic(err)
	}

	bytes, err := json.MarshalIndent(token, "", " ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bytes))
}
