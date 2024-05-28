package service

import (
	"context"

	"github.com/nandanugg/HaloSusterTestCasesPSW3B2/entity"
)

func (c *NipService) AddNurseUsedAccount(ctx context.Context, usr *entity.UsedUser) error {
	return c.repo.PostUsedNurseAccount(ctx, usr)
}

func (c *NipService) GetNurseUsedAccount(ctx context.Context) (*entity.UsedUser, error) {
	return c.repo.GetUsedNurseAccount(ctx)
}
