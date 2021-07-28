package repo

import "github.com/ozoncp/ocp-resource-api/internal/models"

type Repo interface {
	AddEntities(entities []models.Resource) error
	ListEntities(limit, offset uint64) ([]models.Resource, error)
	DescribeEntity(entityId uint64) (*models.Resource, error)
}
