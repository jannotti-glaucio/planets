package planets

import (
	"github.com/jannotti-glaucio/planets/adapters/rest/controllers/planets"
	"github.com/jannotti-glaucio/planets/adapters/rest/middlewares/auth"
	"github.com/labstack/echo/v4"
)

//SetupRoutes ...
func SetupRoutes(Echo *echo.Echo) {
	//Public routes no authentication required
	routes := Echo.Group("/planets")

	routes.GET("/:UUID", planets.Show())
	routes.GET("", planets.Index())
	routes.POST("", planets.Create(), auth.Authorize)
	routes.PATCH("/:UUID", planets.Update(), auth.Authorize)
	routes.DELETE("/:UUID", planets.Destroy(), auth.Authorize)
}
