package employees

import (
	"errors"
	"fmt"

	"github.com/codebarz/employee-service/rpc/proto/employeepb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func PBToNewEmployee(pb *employeepb.CreateEmployeeRequest) (*NewEmployee, error) {
	return &NewEmployee{
		FirstName: pb.GetFirstName(),
		LastName:  pb.GetLastName(),
		Role:      pb.GetRole(),
		Email:     pb.GetEmail(),
	}, nil
}

func EmployeeToPB(employee *Employee) (*employeepb.Employee, error) {
	if employee == nil {
		return nil, errors.New("employee is nil")
	}

	var pbCreatedAt *timestamppb.Timestamp
	var pbUpdatedAt *timestamppb.Timestamp

	if !employee.CreatedAt.IsZero() {
		pbCreatedAt = timestamppb.New(employee.CreatedAt)
	}

	if !employee.UpdatedAt.IsZero() {
		pbUpdatedAt = timestamppb.New(employee.UpdatedAt)
	}

	return &employeepb.Employee{
		Id:        employee.ID,
		FirstName: employee.FirstName,
		LastName:  employee.LastName,
		Role:      employee.Role,
		Email:     employee.Email,
		CreatedAt: pbCreatedAt,
		UpdatedAt: pbUpdatedAt,
	}, nil
}

func EmployeesToPB(employees []Employee) (*employeepb.Employees, error) {
	fmt.Print(employees)
	pbs := []*employeepb.Employee{}
	for _, e := range employees {
		pb, err := EmployeeToPB(&e)
		if err != nil {
			return nil, err
		}
		pbs = append(pbs, pb)
	}
	return &employeepb.Employees{
		Employees: pbs,
	}, nil
}
