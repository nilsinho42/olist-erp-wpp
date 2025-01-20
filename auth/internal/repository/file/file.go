package file

// Manage the file that stores the token

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"strings"
	"sync"

	encrypt "auth/internal/security"

	"github.com/nilsinho42/OlistERPMediator/auth/pkg/model"
)

type Repository struct {
	sync.RWMutex
	data map[string]*model.Token
}

func (r *Repository) Get(_ context.Context) (*model.Token, error) {
	r.RLock()
	defer r.RUnlock()

	const filename = ".token_store.json"
	var token *model.Token

	filepath, err := os.Getwd()
	if err != nil {
		return token, nil
	}

	filepath = strings.Join([]string{filepath, filename}, string(os.PathSeparator))
	file, err := os.Open(filepath)
	if err != nil {
		return token, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return token, err
	}

	err = json.Unmarshal(data, &token)

	decrypted_token, err := encrypt.DecryptAES([]byte("key"), token.Key)
	if err != nil {
		return nil, err
	}

	return decrypted_token, err
}
