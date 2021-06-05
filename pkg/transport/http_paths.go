package transport

const accessPrefix = "/api/access"
const authorizationPrefix = "/api/auth"
const initializationPrefix = "/api/init"
const appPrefix = "/api/app"

var appPaths = map[string]string{
	"login":  "/login",  // POST
	"logout": "/logout", // POST
}

var accessPaths = map[string]string{
	"listRoles":   "/roles",                //GET
	"rolesByUser": "/roles/user/{user_id}", //GET

	"addRole":    "/roles",                 //POST
	"editRole":   "/roles/{role_id}",       //PUT
	"deleteRole": "/roles/{role_id}",       //DELETE
	"cloneRole":  "/roles/{role_id}/clone", //POST

	"getAccessStructure": "/modules/{module}/structure", //GET
	"getAccessByRole":    "/modules/{role_id}/access",   //GET

	"getAllModules":      "/modules",                     //GET
	"getAssignedModules": "/modules/{role_id}",           //GET
	"assignModules":      "/modules/{role_id}",           //POST
	"unassignModules":    "/modules/{role_id}/{modules}", //DELETE

	"getAssignedSubModules": "/submodules/{role_id}",                                //GET
	"assignSubModules":      "/submodules/{role_id}/modules/{module}",               //POST
	"unassignSubModules":    "/submodules/{role_id}/modules/{module}/{sub_modules}", //DELETE

	"getAssignedSections": "/sections/{role_id}",                                                     //GET
	"assignSections":      "/sections/{role_id}/modules/{module}/submodules/{sub_module}",            //POST
	"unassignSections":    "/sections/{role_id}/modules/{module}/submodules/{sub_module}/{sections}", //DELETE
}

var authorizationPaths = map[string]string{
	"listUsers":   "/users",                //GET
	"usersByRole": "/users/role/{role_id}", //GET

	"assignRole":      "/users/{user_id}/role/{role_id}",                                       //POST
	"unassignRole":    "/users/{user_id}/role/{role_id}",                                       //DELETE
	"assignActions":   "/actions/{role_id}/modules/{module}/submodules/{sub_module}",           //POST
	"unassignActions": "/actions/{role_id}/modules/{module}/submodules/{sub_module}/{actions}", //DELETE

	"getAccessList": "/users/{user_id}/access",                   //GET
	"getActionList": "/users/{user_id}/modules/{module}/actions", //GET
}

var initializationPaths = map[string]string{
	"initDB": "/reset-db", //POST
}

func getAccessPath(path string) string {
	if p, ok := accessPaths[path]; ok {
		return accessPrefix + p
	}
	panic("Undefined path")
}

func getActionsPath(path string) string {
	if p, ok := authorizationPaths[path]; ok {
		return authorizationPrefix + p
	}
	panic("Undefined path")
}

func getInitPath(path string) string {
	if p, ok := initializationPaths[path]; ok {
		return initializationPrefix + p
	}
	panic("Undefined path")
}

func getAppPaths(path string) string {
	if p, ok := appPaths[path]; ok {
		return appPrefix + p
	}
	panic("Undefined path")
}
