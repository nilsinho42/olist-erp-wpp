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
	tokenPrimary, err := c.Primary.Get(ctx)
	if err != nil {
		log.Printf("Primary repository failed: %v", err)
		return err
	}

	tokenSecondary, err := c.Secondary.Get(ctx)
	if err != nil {
		log.Printf("Secondary repository failed: %v", err)
		return err
	}

	if tokenPrimary.Key != "" {
		err = c.Primary.Put(ctx)
		if err != nil {
			log.Printf("Primary repository failed: %v", err)
			return err
		}
		tokenSecondary.Key = tokenPrimary.Key
		err = c.Secondary.Put(ctx)
		if err != nil {
			log.Printf("Secondary repository failed: %v", err)
			return err
		}
	} else if tokenSecondary.Key != "" {
		err = c.Secondary.Put(ctx)
		if err != nil {
			log.Printf("Secondary repository failed: %v", err)
			return err
		}
		tokenPrimary.Key = tokenSecondary.Key
		err = c.Primary.Put(ctx)
		if err != nil {
			log.Printf("Secondary repository failed: %v", err)
			return err
		}
	}

	return err
}
