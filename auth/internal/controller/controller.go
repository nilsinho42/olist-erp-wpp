package olistmediator

import (
	"auth/pkg/model"
	"context"
)

type tokenRepository interface {
	Get(ctx context.Context) (*model.Token, error)
	Put(ctx context.Context, token *model.Token) error
}

type Controller struct {
	repo tokenRepository
}

func New(repo tokenRepository) *Controller {
	return &Controller{repo}
}

func (c *Controller) Get(ctx context.Context) (*model.Token, error) {
	res, err := c.repo.Get(ctx)
	if err != nil {
		return nil, err
	}
	return res, err
}

func (c *Controller) Put(ctx context.Context, token *model.Token) error {
	err := c.repo.Put(ctx, token)
	if err != nil {
		return err
	}
	return nil
}
