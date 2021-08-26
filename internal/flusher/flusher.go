package flusher

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/ozoncp/ocp-resource-api/internal/models"
	"github.com/ozoncp/ocp-resource-api/internal/repo"
	"github.com/ozoncp/ocp-resource-api/internal/utils"
)

type Flusher interface {
	Flush(ctx context.Context, resources []models.Resource, span opentracing.Span) []models.Resource
}

type flusher struct {
	chunkSize    uint
	resourceRepo repo.Repo
	tracer       *opentracing.Tracer
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

func (f *flusher) Flush(ctx context.Context, resources []models.Resource, parentSpan opentracing.Span) []models.Resource {
	var span opentracing.Span
	if parentSpan != nil {
		span = opentracing.GlobalTracer().StartSpan("Flush", opentracing.ChildOf(parentSpan.Context()))
		defer span.Finish()
	}

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
		errAddEntities := f.resourceRepo.AddEntities(ctx, chunk, span)
		if errAddEntities != nil {
			return resources[int(f.chunkSize)*i:]
		}
	}

	return make([]models.Resource, 0)
}
