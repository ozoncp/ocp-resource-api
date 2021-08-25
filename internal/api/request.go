package api

import (
	"context"
	desc "github.com/ozoncp/ocp-resource-api/pkg/ocp-resource-api"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	errEndpointNotFound = "endpoint not found"
)

type api struct {
	desc.UnimplementedOcpResourceApiServer
}

func (api) CreateResourceV1(_ context.Context, _ *desc.CreateResourceRequestV1) (*desc.ResourceV1, error) {
	log.Warn().Msgf("Requested not implemented method: CreateResourceV1")
	err := status.Error(codes.NotFound, errEndpointNotFound)
	return nil, err
}

func (_ *api) DescribeResourceV1(_ context.Context, _ *desc.DescribeResourceRequestV1) (*desc.ResourceV1, error) {
	log.Warn().Msgf("Requested not implemented method: DescribeResourceV1")
	err := status.Error(codes.NotFound, errEndpointNotFound)
	return nil, err
}

func (_ *api) ListResourcesV1(_ context.Context, _ *desc.ListResourcesRequestV1) (*desc.ListResourcesResponseV1, error) {
	log.Warn().Msgf("Requested not implemented method: ListResourcesV1")
	err := status.Error(codes.NotFound, errEndpointNotFound)
	return nil, err
}

func (_ *api) RemoveResourceV1(_ context.Context, _ *desc.RemoveResourceRequestV1) (*desc.RemoveResourceResponseV1, error) {
	log.Warn().Msgf("Requested not implemented method: RemoveResourceV1")
	err := status.Error(codes.NotFound, errEndpointNotFound)
	return nil, err
}

func NewOcpResourceApi() desc.OcpResourceApiServer {
	return &api{}
}
