package repo

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/ozoncp/ocp-resource-api/internal/models"
)

type Repo interface {
	AddEntity(ctx context.Context, entity models.Resource) (*models.Resource, error)
	AddEntities(ctx context.Context, entities []models.Resource, span opentracing.Span) error
	ListEntities(ctx context.Context, limit uint64, offset uint64) (*[]models.Resource, error)
	DescribeEntity(ctx context.Context, entityId uint64) (*models.Resource, error)
	RemoveEntity(ctx context.Context, entityId uint64) error
	UpdateEntity(ctx context.Context, entityId uint64, userId uint64, resourceType uint64, status uint64) (*models.Resource, error)
}
