package main

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"

	"github.com/jannotti-glaucio/planets/adapters/persistence/mongodb"

	"github.com/jannotti-glaucio/planets/adapters/rest/routers"
	"github.com/jannotti-glaucio/planets/core/tools/communication"
	"github.com/jannotti-glaucio/planets/core/tools/providers/cache"
	"github.com/jannotti-glaucio/planets/core/tools/providers/logger"
	"github.com/jannotti-glaucio/planets/core/tools/providers/tracer"
)

func init() {
	viper.SetConfigFile(`./config.yml`)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Duration(time.Second*10))
	defer cancel()

	log := logger.New()

	connection := mongodb.New(ctx)
	if connection.Error != nil {
		log.Error(ctx, "Error connecting to database", connection.Error)
		return
	}

	memcacheHost := viper.GetString("memcache.host") + ":" + viper.GetString("memcache.port")
	_, err := cache.New(memcacheHost)
	if err != nil {
		log.Error(ctx, "Error connecting to memcache", err)
		return
	}

	jaegerTracer := tracer.New("planets")
	if jaegerTracer.Error != nil {
		log.Error(ctx, "Error connecting to jaeger", jaegerTracer)
	}
	defer jaegerTracer.Closer.Close()

	Echo := echo.New()
	routers.StartRouters(Echo)
	Echo.HTTPErrorHandler = customHTTPErrorHandler
	Echo.Logger.Fatal(Echo.Start(":" + viper.GetString("web.port")))
}

func customHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}

	comm := communication.New()
	if code == 404 {
		resp := comm.Response(code, "end_point_not_found")
		c.JSON(resp.Status, resp)
		return
	}

	resp := comm.Response(code, "Unexpected")
	c.JSON(resp.Status, resp)
}
