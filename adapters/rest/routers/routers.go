package routers

import (
	"github.com/jannotti-glaucio/planets/adapters/rest/routers/planets"
	"github.com/jannotti-glaucio/planets/adapters/rest/routers/users"
	"github.com/labstack/echo/v4"
)

//StartRouters ...
func StartRouters(Echo *echo.Echo) {
	users.SetupRoutes(Echo)
	planets.SetupRoutes(Echo)
}
