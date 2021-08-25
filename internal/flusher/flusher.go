package flusher

import (
	"github.com/ozoncp/ocp-resource-api/internal/models"
	"github.com/ozoncp/ocp-resource-api/internal/repo"
	"github.com/ozoncp/ocp-resource-api/internal/utils"
)

type Flusher interface {
	Flush(resources []models.Resource) []models.Resource
}

type flusher struct {
	chunkSize    uint
	resourceRepo repo.Repo
}

func NewFlusher(
	chunkSize uint,
	resourceRepo repo.Repo,
) Flusher {
	return &flusher{
		chunkSize:    chunkSize,
		resourceRepo: resourceRepo,
	}
}

func (f *flusher) Flush(resources []models.Resource) []models.Resource {
	var err error
	if f.resourceRepo == nil {
		return resources
	}
	if len(resources) == 0 {
		return resources
	}
	chunks, err := utils.SplitToBulksResource(resources, f.chunkSize)
	if err != nil {
		return resources
	}

	for i := range chunks {
		chunk := chunks[i]
		errAddEntities := f.resourceRepo.AddEntities(nil, chunk)
		if errAddEntities != nil {
			return resources[int(f.chunkSize)*i:]
		}
	}

	return make([]models.Resource, 0)
}
