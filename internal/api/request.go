package api

import (
	"context"
	"fmt"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/ozoncp/ocp-resource-api/internal/flusher"
	"github.com/ozoncp/ocp-resource-api/internal/metrics"
	"github.com/ozoncp/ocp-resource-api/internal/models"
	"github.com/ozoncp/ocp-resource-api/internal/producer"
	"github.com/ozoncp/ocp-resource-api/internal/repo"
	desc "github.com/ozoncp/ocp-resource-api/pkg/ocp-resource-api"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const chunkSize = 10

//TODO add a tests as improvements
type api struct {
	desc.UnimplementedOcpResourceApiServer
	repo        repo.Repo
	flusher     flusher.Flusher
	msgProducer producer.Producer
}

func (a *api) CreateResourceV1(ctx context.Context, req *desc.CreateResourceRequestV1) (*desc.ResourceV1, error) {
	a.notifyProducer(ctx, "add")
	createdResource, err := a.repo.AddEntity(ctx, models.NewNotPersistedResource(req.GetUserId(), req.GetType(), req.GetStatus()))
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("request err: %v", err))
	}
	rsp := mapResourceToResourceV1(createdResource)
	metrics.IncReqCounter("add")
	return rsp, nil
}

func (a *api) DescribeResourceV1(ctx context.Context, req *desc.DescribeResourceRequestV1) (*desc.ResourceV1, error) {
	a.notifyProducer(ctx, "get")
	resource, err := a.repo.DescribeEntity(ctx, req.GetId())
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("request err: %v", err))
	}
	rsp := mapResourceToResourceV1(resource)
	return rsp, nil
}

func (a *api) ListResourcesV1(ctx context.Context, req *desc.ListResourcesRequestV1) (*desc.ListResourcesResponseV1, error) {
	a.notifyProducer(ctx, "list")
	resources, err := a.repo.ListEntities(ctx, req.GetLimit(), req.GetOffset())
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("request err: %v", err))
	}
	resourcesV1Slice := make([]*desc.ResourceV1, 0, len(resources))
	for i := range resources {
		resourcesV1Slice = append(resourcesV1Slice, mapResourceToResourceV1(&resources[i]))
	}
	rsp := desc.ListResourcesResponseV1{Resources: resourcesV1Slice}
	return &rsp, nil
}

func (a *api) RemoveResourceV1(ctx context.Context, req *desc.RemoveResourceRequestV1) (*desc.RemoveResourceResponseV1, error) {
	a.notifyProducer(ctx, "remove")
	err := a.repo.RemoveEntity(ctx, req.GetId())
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("request err: %v", err))
	}
	rsp := desc.RemoveResourceResponseV1{}
	metrics.IncReqCounter("remove")
	return &rsp, nil
}

func (a *api) MultiCreateResourcesV1(ctx context.Context, req *desc.MultiCreateResourceRequestV1) (*desc.MultiCreateResourceResponseV1, error) {
	a.notifyProducer(ctx, "multi add")
	tracer := opentracing.GlobalTracer()
	span := tracer.StartSpan("MultiCreateResourcesV1")
	defer span.Finish()
	resources := make([]models.Resource, 0, cap(req.GetResources()))
	for i, res := range req.GetResources() {
		resources[i] = models.NewNotPersistedResource(res.UserId, res.Type, res.Status)
	}
	ctx = opentracing.ContextWithSpan(ctx, span)
	notSaved := a.flusher.Flush(ctx, resources)
	if len(notSaved) != 0 {
		return nil, status.Error(codes.Internal, fmt.Sprint("flush err"))
	}
	rsp := desc.MultiCreateResourceResponseV1{}
	for range resources {
		metrics.IncReqCounter("add")
	}
	return &rsp, nil
}

func (a *api) UpdateResourceV1(ctx context.Context, req *desc.UpdateResourceRequestV1) (*desc.ResourceV1, error) {
	a.notifyProducer(ctx, "update")
	resource, err := a.repo.UpdateEntity(ctx, req.GetId(), req.GetFields().GetUserId(), req.GetFields().GetType(), req.GetFields().GetStatus())
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("request err: %v", err))
	}
	return mapResourceToResourceV1(resource), nil
}

func (a *api) notifyProducer(_ context.Context, operation string) {
	err := a.msgProducer.Send(producer.EventMessage{
		Timestamp: time.Now().Unix(),
		Operation: operation,
	})
	if err != nil {
		log.Err(err).Msgf("Issue during producer notification")
	}
}

func mapResourceToResourceV1(resource *models.Resource) *desc.ResourceV1 {
	res := desc.ResourceV1{
		Id:     resource.Id,
		UserId: resource.UserId,
		Type:   resource.Type,
		Status: resource.Status,
	}
	return &res
}

func NewOcpResourceApi(repo repo.Repo, producer producer.Producer) (desc.OcpResourceApiServer, error) {
	newFlusher := flusher.NewFlusher(chunkSize, repo)
	return &api{repo: repo, flusher: newFlusher, msgProducer: producer}, nil
}
