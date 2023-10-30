package controller

type AppController struct {
	User  interface{ UserEndpointsHandler }
	Admin interface{ AdminEndpointHandler }
}
