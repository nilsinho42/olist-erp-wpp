package repository

import (
	"auth/pkg/model"
	"context"
	"fmt"
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

func NewCompositeRepository(primary, secondary TokenRepository) *CompositeTokenRepository {
	return &CompositeTokenRepository{
		Primary:   primary,
		Secondary: secondary,
	}
}

func (c *CompositeTokenRepository) Get(ctx context.Context) (*model.Token, error) {
	token, err := c.Primary.Get(ctx)
	if err != nil {
		log.Printf("Primary repository failed: %v, falling back to secondary", err)
		return c.Secondary.Get(ctx)
	}
	return token, nil
}

func (c *CompositeTokenRepository) Put(ctx context.Context) error {
	key, ok := ctx.Value(model.ContextKey).(string)
	if !ok {
		return fmt.Errorf("key not found in context")
	}
	tokenPrimary, err := c.Primary.Get(ctx)
	if err != nil {
		log.Printf("Primary repository failed: %v", err)
		return err
	}
	tokenPrimary.Key = key
	err = c.Primary.Put(ctx)
	if err != nil {
		log.Printf("Primary repository failed: %v", err)
		return err
	}

	tokenSecondary, err := c.Secondary.Get(ctx)
	if err != nil {
		log.Printf("Secondary repository failed: %v", err)
		return err
	}
	tokenSecondary.Key = key
	err = c.Secondary.Put(ctx)
	if err != nil {
		log.Printf("Secondary repository failed: %v", err)
	}

	return err
}
