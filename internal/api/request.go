package api

import (
	"context"
	"fmt"
	"github.com/ozoncp/ocp-resource-api/internal/models"
	"github.com/ozoncp/ocp-resource-api/internal/repo"
	desc "github.com/ozoncp/ocp-resource-api/pkg/ocp-resource-api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//TODO add a tests as improvements
type api struct {
	desc.UnimplementedOcpResourceApiServer
	repo repo.Repo
}

func (a *api) CreateResourceV1(ctx context.Context, req *desc.CreateResourceRequestV1) (*desc.ResourceV1, error) {
	createdResource, err := a.repo.AddEntity(ctx, models.NewResourceNotPersisted(req.GetUserId(), req.GetType(), req.GetStatus()))
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("request err: %v", err))
	}
	rsp := mapResourceToResourceV1(createdResource)
	return rsp, nil
}

func (a *api) DescribeResourceV1(ctx context.Context, req *desc.DescribeResourceRequestV1) (*desc.ResourceV1, error) {
	resource, err := a.repo.DescribeEntity(ctx, req.GetResourceId())
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("request err: %v", err))
	}
	rsp := mapResourceToResourceV1(resource)
	return rsp, nil
}

func (a *api) ListResourcesV1(ctx context.Context, req *desc.ListResourcesRequestV1) (*desc.ListResourcesResponseV1, error) {
	resourcesPtr, err := a.repo.ListEntities(ctx, req.GetLimit(), req.GetOffset())
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("request err: %v", err))
	}
	resources := *resourcesPtr
	resourcesV1Slice := make([]*desc.ResourceV1, 0, len(resources))
	for i := range resources {
		resourcesV1Slice = append(resourcesV1Slice, mapResourceToResourceV1(&resources[i]))
	}
	rsp := desc.ListResourcesResponseV1{Resources: resourcesV1Slice}
	return &rsp, nil
}

func (a *api) RemoveResourceV1(ctx context.Context, req *desc.RemoveResourceRequestV1) (*desc.RemoveResourceResponseV1, error) {
	err := a.repo.RemoveEntity(ctx, req.GetResourceId())
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("request err: %v", err))
	}
	rsp := desc.RemoveResourceResponseV1{}
	return &rsp, nil
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

func NewOcpResourceApi(repo *repo.Repo) (desc.OcpResourceApiServer, error) {
	return &api{repo: *repo}, nil
}
