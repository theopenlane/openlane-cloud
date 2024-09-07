package route

import (
	"net/http"

	echo "github.com/theopenlane/echox"
)

// registerOrganizationHandler registers the organization handler and route
func registerOrganizationHandler(router *Router) (err error) {
	path := "/organization"
	method := http.MethodPost
	name := "Organization"

	route := echo.Route{
		Name:   name,
		Method: method,
		Path:   path,
		Handler: func(c echo.Context) error {
			return router.Handler.OrganizationHandler(c)
		},
	}

	registerOperation := router.Handler.BindOrganizationHandler()

	if err := router.Addv1Route(path, method, registerOperation, route); err != nil {
		return err
	}

	return nil
}
