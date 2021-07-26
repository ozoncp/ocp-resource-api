package models

import "fmt"

type Resource struct {
	Id     uint64
	UserId uint64
	Type   uint64
	Status uint64
}

func NewResource(id uint64, userId uint64, resourceType uint64, status uint64) Resource {
	return Resource{
		Id:     id,
		UserId: userId,
		Type:   resourceType,
		Status: status,
	}
}

func (r *Resource) String() string {
	return fmt.Sprintf("Resource {Id: %d, UserId: %d, Type: %d, Status: %d}", r.Id, r.UserId, r.Type, r.Status)
}
