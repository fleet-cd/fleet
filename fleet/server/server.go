package server

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/tgs266/fleet/config"
	"github.com/tgs266/fleet/fleet/services"
	"github.com/tgs266/fleet/rest-gen/generated/com/fleet/ships"
)

type Server struct {
	Router *gin.Engine
	Logger zerolog.Logger
	Config config.Config
}

func New(logger zerolog.Logger, config config.Config) *Server {
	logger.Info().Msg("creating fleet server")
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = ioutil.Discard
	router := gin.New()
	server := &Server{
		Router: router,
		Logger: logger,
		Config: config,
	}
	router.Use(Middleware(logger))
	router.Use(ErrorMiddleware())

	shipSvc := services.ShipService{}
	ships.RegisterShipServiceRoutes(router, ships.ShipServiceHandler{
		Handler: &shipSvc,
	})

	return server
}

func (server *Server) createPath(path string) string {
	return fmt.Sprintf("%s/%s", strings.TrimRight(server.Config.Server.BasePath, "/"), strings.TrimLeft(path, "/"))
}

func (server *Server) Start() {
	host := server.Config.Server.Host + ":" + server.Config.Server.Port
	server.Logger.Info().Msg("starting fleet server at http://" + host)
	server.Router.Run(host)
}
