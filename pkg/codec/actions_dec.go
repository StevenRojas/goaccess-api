package codec

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	e "github.com/StevenRojas/goaccess-api/pkg/errors"
	"github.com/StevenRojas/goaccess/pkg/entities"
	"github.com/gorilla/mux"
)

// DecodeAssignActionsRequest decode request
func DecodeAssignActionsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	roleID, ok := vars["role_id"]
	if !ok {
		return "", e.HTTPBadRequestFromString("Role ID is missing")
	}
	moduleName, ok := vars["module"]
	if !ok {
		return "", e.HTTPBadRequestFromString("Module name is missing")
	}
	submoduleName, ok := vars["sub_module"]
	if !ok {
		return "", e.HTTPBadRequestFromString("Subodule name is missing")
	}

	defer r.Body.Close()
	actionsList := entities.ActionList{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&actionsList)
	if err != nil {
		return nil, e.HTTPBadRequestFromString("Validation error while assigning actions")
	}
	actionsList.RoleID = roleID
	actionsList.Module = moduleName
	actionsList.SubModule = submoduleName
	return actionsList, nil
}

// DecodeUnassignActionsRequest decode request
func DecodeUnassignActionsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	roleID, ok := vars["role_id"]
	if !ok {
		return "", e.HTTPBadRequestFromString("Role ID is missing")
	}
	moduleName, ok := vars["module"]
	if !ok {
		return "", e.HTTPBadRequestFromString("Module name is missing")
	}
	submoduleName, ok := vars["sub_module"]
	if !ok {
		return "", e.HTTPBadRequestFromString("Subodule name is missing")
	}
	actions, ok := vars["actions"]
	if !ok {
		return "", e.HTTPBadRequestFromString("Actions are missing")
	}
	actionsList := entities.ActionList{}
	actionsList.RoleID = roleID
	actionsList.Module = moduleName
	actionsList.SubModule = submoduleName
	actionsList.Actions = strings.Split(actions, ",")
	for i := range actionsList.Actions {
		actionsList.Actions[i] = strings.TrimSpace(actionsList.Actions[i])
	}
	return actionsList, nil
}

// DecodeAssignUnassignRoleRequest decode request
func DecodeAssignUnassignRoleRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	roleID, ok := vars["role_id"]
	if !ok {
		return "", e.HTTPBadRequestFromString("Role ID is missing")
	}
	userID, ok := vars["user_id"]
	if !ok {
		return "", e.HTTPBadRequestFromString("User ID is missing")
	}
	return map[string]string{
		"roleID": roleID,
		"userID": userID,
	}, nil
}

// DecodeGetAccessRequest decode request
func DecodeGetAccessRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	userID, ok := vars["user_id"]
	if !ok {
		return "", e.HTTPBadRequestFromString("User ID is missing")
	}
	return userID, nil
}

// DecodeGetActionsRequest decode request
func DecodeGetActionsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	userID, ok := vars["user_id"]
	if !ok {
		return "", e.HTTPBadRequestFromString("User ID is missing")
	}
	module, ok := vars["module"]
	if !ok {
		return "", e.HTTPBadRequestFromString("Role ID is missing")
	}
	return map[string]string{
		"userID": userID,
		"module": module,
	}, nil
}
