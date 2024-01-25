package roles

import (
	"errors"

	"github.com/codebarz/employee-service/rpc/proto/rolepb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func PBToRole(pb *rolepb.CreateRoleRequest) (*Role, error) {
	return &Role{
		Title: pb.GetTitle(),
		Level: int32(pb.GetLevel()),
	}, nil
}

func PBToNewRole(pb *rolepb.CreateRoleRequest) (*NewRole, error) {
	return &NewRole{
		Title: pb.GetTitle(),
		Level: int(pb.GetLevel()),
	}, nil
}

func RoleToPB(role *Role) (*rolepb.Role, error) {
	if role == nil {
		return nil, errors.New("role is nil")
	}
	var pbCreatedAt *timestamppb.Timestamp
	var pbUpdatedAt *timestamppb.Timestamp

	if !role.CreatedAt.IsZero() {
		pbCreatedAt = timestamppb.New(role.CreatedAt)
	}

	if !role.UpdatedAt.IsZero() {
		pbUpdatedAt = timestamppb.New(role.UpdatedAt)
	}

	return &rolepb.Role{
		Id:        role.ID,
		Title:     role.Title,
		Level:     int32(role.Level),
		CreatedAt: pbCreatedAt,
		UpdatedAt: pbUpdatedAt,
	}, nil
}

func RolesToPB(roles []Role) (*rolepb.Roles, error) {
	pbs := []*rolepb.Role{}
	for _, role := range roles {
		pb, err := RoleToPB(&role)
		if err != nil {
			return nil, err
		}
		pbs = append(pbs, pb)
	}
	return &rolepb.Roles{
		Roles: pbs,
	}, nil
}
