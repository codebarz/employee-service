package role

import (
	"context"

	"github.com/codebarz/employee-service/entities/roles"
	"github.com/codebarz/employee-service/rpc/proto/rolepb"
	"github.com/go-kit/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Handler struct {
	l       log.Logger
	service Service
	rolepb.UnimplementedRoleServiceServer
}

func NewGRPCHandler(l log.Logger, s Service) *Handler {
	return &Handler{l: l, service: s}
}

func (h *Handler) CreateRole(ctx context.Context, rb *rolepb.CreateRoleRequest) (*rolepb.Role, error) {
	newRoleReq, err := roles.PBToNewRole(rb)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	role, err := h.service.Create(ctx, newRoleReq)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return roles.RoleToPB(role)
}

func (h *Handler) QueryRoles(ctx context.Context, rb *rolepb.QueryRoleRequest) (*rolepb.Roles, error) {
	r, err := h.service.Query(ctx)

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return roles.RolesToPB(r)
}

func (h *Handler) QueryRoleByID(ctx context.Context, rb *rolepb.QueryRoleByIDRequest) (*rolepb.Role, error) {
	r, err := h.service.QueryByID(ctx, rb.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return roles.RoleToPB(r)
}

func (h *Handler) DeleteRole(ctx context.Context, rb *rolepb.DeleteRoleRequest) (*emptypb.Empty, error) {
	err := h.service.Delete(ctx, rb.GetId())

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (h *Handler) UpdateRole(ctx context.Context, rb *rolepb.UpdateRoleRequest) (*rolepb.Role, error) {
	updatedRole := roles.UpdateRole{
		Id:    rb.Id,
		Title: &rb.Title,
		Level: &rb.Level,
	}
	r, err := h.service.Update(ctx, updatedRole)

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return roles.RoleToPB(r)
}
