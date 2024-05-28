package service

import (
	"context"

	"github.com/nandanugg/HaloSusterTestCasesPSW3B2/entity"
)

func (c *NipService) AddItUsedAccount(ctx context.Context, usr *entity.UsedUser) error {
	return c.repo.PostUsedITAccount(ctx, usr)
}

func (c *NipService) GetItUsedAccount(ctx context.Context) (*entity.UsedUser, error) {
	return c.repo.GetUsedITAccount(ctx)
}
