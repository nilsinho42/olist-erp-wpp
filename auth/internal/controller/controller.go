package controller

import (
	"auth/pkg/model"
	"context"
)

type tokenRepository interface {
	Get(ctx context.Context) (*model.Token, error)
	Put(ctx context.Context) error
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

func (c *Controller) Put(ctx context.Context) error {
	err := c.repo.Put(ctx)
	if err != nil {
		return err
	}
	return nil
}
