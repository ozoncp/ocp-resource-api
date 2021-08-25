package repo

import (
	"context"
	"github.com/ozoncp/ocp-resource-api/internal/models"
)

type Repo interface {
	AddEntity(ctx context.Context, entity models.Resource) (*models.Resource, error)
	AddEntities(ctx context.Context, entities []models.Resource) error
	ListEntities(ctx context.Context, limit uint64, offset uint64) (*[]models.Resource, error)
	DescribeEntity(ctx context.Context, entityId uint64) (*models.Resource, error)
	RemoveEntity(ctx context.Context, entityId uint64) error
}
