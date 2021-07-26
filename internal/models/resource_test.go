package models

import (
	"reflect"
	"testing"
)

func TestNewResource(t *testing.T) {
	type args struct {
		id           uint64
		userId       uint64
		resourceType uint64
		status       uint64
	}
	tests := []struct {
		name string
		args args
		want Resource
	}{
		{
			name: "base",
			args: args{id: 10, resourceType: 20, status: 30, userId: 40},
			want: Resource{Id: 10, Type: 20, Status: 30, UserId: 40},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewResource(tt.args.id, tt.args.userId, tt.args.resourceType, tt.args.status); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewResource() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResource_String(t *testing.T) {
	type fields struct {
		Id     uint64
		UserId uint64
		Type   uint64
		Status uint64
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "base",
			fields: fields{Id: 10, UserId: 20, Type: 30, Status: 40},
			want:   "Resource {Id: 10, UserId: 20, Type: 30, Status: 40}"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Resource{
				Id:     tt.fields.Id,
				UserId: tt.fields.UserId,
				Type:   tt.fields.Type,
				Status: tt.fields.Status,
			}
			if got := r.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
