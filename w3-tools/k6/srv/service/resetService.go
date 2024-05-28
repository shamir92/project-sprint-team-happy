package service

import (
	"context"
)

func (c *NipService) Reset(ctx context.Context) error {
	c.itUsedIndexNIP = 0
	c.nurseUsedIndexNIP = 0
	return c.repo.Reset(ctx)
}
