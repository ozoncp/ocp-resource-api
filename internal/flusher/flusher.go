package flusher

import (
	"errors"
	"fmt"
	"github.com/ozoncp/ocp-resource-api/internal/models"
	"github.com/ozoncp/ocp-resource-api/internal/repo"
	"github.com/ozoncp/ocp-resource-api/internal/utils"
)

var ErrRepoIsNil = errors.New("repo should not be nil")

type Flusher interface {
	Flush(resources []models.Resource) ([]models.Resource, error)
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

func (f *flusher) Flush(resources []models.Resource) ([]models.Resource, error) {
	var err error
	if f.resourceRepo == nil {
		return resources, ErrRepoIsNil
	}

	chunks, err := utils.SplitToBulksResource(resources, f.chunkSize)
	if err != nil {
		return resources, fmt.Errorf("error during flush: %w", err)
	}

	for i, _ := range chunks {
		chunk := chunks[i]
		errAddEntities := f.resourceRepo.AddEntities(chunk)
		if errAddEntities != nil {
			return resources[int(f.chunkSize)*i:], errAddEntities
		}
	}

	return make([]models.Resource, 0), nil
}
