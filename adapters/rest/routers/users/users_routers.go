package users

import (
	"github.com/jannotti-glaucio/planets/adapters/rest/controllers/users"
	"github.com/jannotti-glaucio/planets/adapters/rest/middlewares/auth"
	"github.com/labstack/echo/v4"
)

//SetupRoutes ...
func SetupRoutes(Echo *echo.Echo) {
	//Public routes no authentication required
	routes := Echo.Group("/users")
	routes.POST("/auth", users.Auth())

	routes.POST("", users.Create(), auth.Authorize)
	routes.GET("", users.Index(), auth.Authorize)
	routes.GET("/:UUID", users.Show(), auth.Authorize)
	routes.PATCH("/:UUID", users.Update(), auth.Authorize)
	routes.DELETE("/:UUID", users.Destroy(), auth.Authorize)
}
