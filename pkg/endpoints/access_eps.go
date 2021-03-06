package endpoints

import (
	"context"
	"errors"

	"github.com/StevenRojas/goaccess-api/pkg/codec"
	e "github.com/StevenRojas/goaccess-api/pkg/errors"
	"github.com/StevenRojas/goaccess/pkg/entities"
	"github.com/StevenRojas/goaccess/pkg/service"
	"github.com/go-kit/kit/endpoint"
)

type AccessEndpoints struct {
	ListRoles             endpoint.Endpoint
	ListRolesByUser       endpoint.Endpoint
	AddRole               endpoint.Endpoint
	EditRole              endpoint.Endpoint
	DeleteRole            endpoint.Endpoint
	CloneRole             endpoint.Endpoint
	GetAllModules         endpoint.Endpoint
	GetAccessStructure    endpoint.Endpoint
	GetAccessByRole       endpoint.Endpoint
	GetAssignedModules    endpoint.Endpoint
	AssignModules         endpoint.Endpoint
	UnassignModules       endpoint.Endpoint
	GetAssignedSubModules endpoint.Endpoint
	AssignSubModules      endpoint.Endpoint
	UnassignSubModules    endpoint.Endpoint
	GetAssignedSections   endpoint.Endpoint
	AssignSections        endpoint.Endpoint
	UnassignSections      endpoint.Endpoint
}

func MakeAccessEndpoints(
	s service.AccessService,
	middlewares []endpoint.Middleware) AccessEndpoints {
	return AccessEndpoints{
		ListRoles:             wrapMiddlewares(makeListRoles(s), middlewares),
		ListRolesByUser:       wrapMiddlewares(makeListRolesByUser(s), middlewares),
		AddRole:               wrapMiddlewares(makeAddRole(s), middlewares),
		EditRole:              wrapMiddlewares(makeEditRole(s), middlewares),
		DeleteRole:            wrapMiddlewares(makeDeleteRole(s), middlewares),
		CloneRole:             wrapMiddlewares(makeCloneRole(s), middlewares),
		GetAllModules:         wrapMiddlewares(makeGetAllModules(s), middlewares),
		GetAssignedModules:    wrapMiddlewares(makeGetAssignedModules(s), middlewares),
		GetAccessStructure:    wrapMiddlewares(makeGetAccessStructure(s), middlewares),
		GetAccessByRole:       wrapMiddlewares(makeGetAccessByRole(s), middlewares),
		AssignModules:         wrapMiddlewares(makeAssignModules(s), middlewares),
		UnassignModules:       wrapMiddlewares(makeUnassignModules(s), middlewares),
		GetAssignedSubModules: wrapMiddlewares(makeGetAssignedSubModules(s), middlewares),
		AssignSubModules:      wrapMiddlewares(makeAssignSubModules(s), middlewares),
		UnassignSubModules:    wrapMiddlewares(makeUnassignSubModules(s), middlewares),
		GetAssignedSections:   wrapMiddlewares(makeGetAssignedSections(s), middlewares),
		AssignSections:        wrapMiddlewares(makeAssignSections(s), middlewares),
		UnassignSections:      wrapMiddlewares(makeUnassignSections(s), middlewares),
	}
}

func wrapMiddlewares(ep endpoint.Endpoint, middlewares []endpoint.Middleware) endpoint.Endpoint {
	for i := range middlewares {
		ep = middlewares[i](ep)
	}
	return ep
}

func makeListRoles(s service.AccessService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		roles, err := s.ListRoles(ctx)
		if err != nil {
			return nil, e.HTTPConflict("Unable to get a list of roles", err)
		}
		return &codec.MapString{List: roles}, nil
	}
}

func makeListRolesByUser(s service.AccessService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		userID, ok := request.(string)
		if !ok {
			return nil, e.HTTPBadRequest(errors.New("unable to cast the request to string"))
		}
		roles, err := s.ListRolesByUser(ctx, userID)
		if err != nil {
			return nil, e.HTTPConflict("Unable to get a list of roles by user", err)
		}
		return &codec.MapString{List: roles}, nil
	}
}

func makeAddRole(s service.AccessService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		role, ok := request.(entities.Role)
		if !ok {
			return nil, e.HTTPBadRequest(errors.New("unable to cast the request to RoleRequest"))
		}
		ID, err := s.AddRole(ctx, role.Name)
		if err != nil {
			return nil, e.HTTPConflict("Unable to add role", err)
		}
		return &codec.IDResponse{ID: ID}, nil
	}
}

func makeEditRole(s service.AccessService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		role, ok := request.(entities.Role)
		if !ok {
			return nil, e.HTTPBadRequest(errors.New("unable to cast the request to RoleRequest"))
		}
		if ok, _ := s.IsRoleExist(ctx, role.ID); !ok {
			return nil, e.HTTPNotFound("Role not found")
		}
		err := s.EditRole(ctx, role.ID, role.Name)
		if err != nil {
			return nil, e.HTTPConflict("Unable to edit role", err)
		}
		return &codec.EmptyResponse{}, nil
	}
}

func makeDeleteRole(s service.AccessService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		roleID, ok := request.(string)
		if !ok {
			return nil, e.HTTPBadRequest(errors.New("unable to cast the request to string"))
		}
		if ok, _ := s.IsRoleExist(ctx, roleID); !ok {
			return nil, e.HTTPNotFound("Role not found")
		}
		err := s.DeleteRole(ctx, roleID)
		if err != nil {
			return nil, e.HTTPConflict("Unable to delete the role", err)
		}
		return &codec.NoContentResponse{}, nil
	}
}

func makeCloneRole(s service.AccessService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		role, ok := request.(entities.Role)
		if !ok {
			return nil, e.HTTPBadRequest(errors.New("unable to cast the request to RoleRequest"))
		}
		if ok, _ := s.IsRoleExist(ctx, role.ID); !ok {
			return nil, e.HTTPNotFound("Role not found")
		}
		ID, err := s.CloneRole(ctx, role.ID, role.Name)
		if err != nil {
			return nil, e.HTTPConflict("Unable to clone role", err)
		}
		return &codec.IDResponse{ID: ID}, nil
	}
}

func makeGetAllModules(s service.AccessService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		modules, err := s.ModulesList(ctx)
		if err != nil {
			return nil, e.HTTPConflict("Unable to get list of modules", err)
		}
		return &codec.StringList{List: modules}, nil
	}
}

func makeGetAccessStructure(s service.AccessService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		module, ok := request.(string)
		if !ok {
			return nil, e.HTTPBadRequest(errors.New("unable to cast the request to string"))
		}
		structure, err := s.ModuleStructure(ctx, module)
		if err != nil {
			return nil, e.HTTPConflict("Unable to get module structure", err)
		}
		if structure == nil {
			return nil, e.HTTPNotFound("Module not found")
		}
		return structure, nil
	}
}

func makeGetAccessByRole(s service.AccessService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		roleID, ok := request.(string)
		if !ok {
			return nil, e.HTTPBadRequest(errors.New("unable to cast the request to string"))
		}
		modules, err := s.GetRoleAccessList(ctx, roleID)
		if err != nil {
			return nil, e.HTTPConflict("Unable to get modules by role", err)
		}
		return modules, nil
	}
}

func makeGetAssignedModules(s service.AccessService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		roleID, ok := request.(string)
		if !ok {
			return nil, e.HTTPBadRequest(errors.New("unable to cast the request to string"))
		}
		modules, err := s.ModulesListByRole(ctx, roleID)
		if err != nil {
			return nil, e.HTTPConflict("Unable to get list of modules", err)
		}
		return &codec.StringList{List: modules}, nil
	}
}

func makeAssignModules(s service.AccessService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		moduleList, ok := request.(entities.ModuleList)
		if !ok {
			return nil, e.HTTPBadRequest(errors.New("unable to cast the request to module list"))
		}
		err := s.AssignModules(ctx, moduleList.RoleID, moduleList.Modules)
		if err != nil {
			return nil, e.HTTPConflict("Unable to assign modules", err)
		}
		return &codec.EmptyResponse{}, nil
	}
}

func makeUnassignModules(s service.AccessService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		moduleList, ok := request.(entities.ModuleList)
		if !ok {
			return nil, e.HTTPBadRequest(errors.New("unable to cast the request to module list"))
		}
		err := s.UnassignModules(ctx, moduleList.RoleID, moduleList.Modules)
		if err != nil {
			return nil, e.HTTPConflict("Unable to unassign modules", err)
		}
		return &codec.NoContentResponse{}, nil
	}
}

func makeGetAssignedSubModules(s service.AccessService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		roleID, ok := request.(string)
		if !ok {
			return nil, e.HTTPBadRequest(errors.New("unable to cast the request to string"))
		}
		modules, err := s.SubModulesListByRole(ctx, roleID)
		if err != nil {
			return nil, e.HTTPConflict("Unable to get list of modules", err)
		}
		return &codec.MapStringList{List: modules}, nil
	}
}

func makeAssignSubModules(s service.AccessService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		submoduleList, ok := request.(entities.SubModuleList)
		if !ok {
			return nil, e.HTTPBadRequest(errors.New("unable to cast the request to submodule list"))
		}
		err := s.AssignSubModules(ctx, submoduleList.RoleID, submoduleList.Module, submoduleList.SubModules)
		if err != nil {
			return nil, e.HTTPConflict("Unable to assign submodules", err)
		}
		return &codec.EmptyResponse{}, nil
	}
}

func makeUnassignSubModules(s service.AccessService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		submoduleList, ok := request.(entities.SubModuleList)
		if !ok {
			return nil, e.HTTPBadRequest(errors.New("unable to cast the request to submodule list"))
		}
		err := s.UnassignSubModules(ctx, submoduleList.RoleID, submoduleList.Module, submoduleList.SubModules)
		if err != nil {
			return nil, e.HTTPConflict("Unable to unassign submodules", err)
		}
		return &codec.NoContentResponse{}, nil
	}
}

func makeGetAssignedSections(s service.AccessService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		roleID, ok := request.(string)
		if !ok {
			return nil, e.HTTPBadRequest(errors.New("unable to cast the request to string"))
		}
		modules, err := s.SectionsListByRole(ctx, roleID)
		if err != nil {
			return nil, e.HTTPConflict("Unable to get list of modules", err)
		}
		return &codec.MapOfMapStringList{List: modules}, nil
	}
}

func makeAssignSections(s service.AccessService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		sectionList, ok := request.(entities.SectionList)
		if !ok {
			return nil, e.HTTPBadRequest(errors.New("unable to cast the request to section list"))
		}
		err := s.AssignSections(ctx, sectionList.RoleID, sectionList.Module, sectionList.SubModule, sectionList.Sections)
		if err != nil {
			return nil, e.HTTPConflict("Unable to assign sections", err)
		}
		return &codec.EmptyResponse{}, nil
	}
}

func makeUnassignSections(s service.AccessService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		sectionList, ok := request.(entities.SectionList)
		if !ok {
			return nil, e.HTTPBadRequest(errors.New("unable to cast the request to section list"))
		}
		err := s.UnassignSections(ctx, sectionList.RoleID, sectionList.Module, sectionList.SubModule, sectionList.Sections)
		if err != nil {
			return nil, e.HTTPConflict("Unable to unassign sections", err)
		}
		return &codec.NoContentResponse{}, nil
	}
}
