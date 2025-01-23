package repository

import (
	"auth/pkg/model"
	"context"
	"log"
)

type TokenRepository interface {
	Get(ctx context.Context) (*model.Token, error)
	Put(ctx context.Context) error
}
type CompositeTokenRepository struct {
	Primary   TokenRepository
	Secondary TokenRepository
}

func (c *CompositeTokenRepository) Get(ctx context.Context) (*model.Token, error) {
	token, err := c.Primary.Get(ctx)
	if err != nil {
		log.Printf("Primary repository failed: %v, falling back to secondary", err)
		return c.Secondary.Get(ctx)
	}
	return token, nil
}

func (c *CompositeTokenRepository) Put(ctx context.Context, token *model.Token) error {
	err := c.Primary.Put(ctx)
	if err != nil {
		log.Printf("Primary repository failed: %v, falling back to secondary", err)
		err = c.Secondary.Put(ctx)
		return err
	}

	return err
}
