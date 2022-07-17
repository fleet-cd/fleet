package server

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/tgs266/fleet/config"
	authService "github.com/tgs266/fleet/fleet/auth"
	cargoService "github.com/tgs266/fleet/fleet/cargo"
	healthService "github.com/tgs266/fleet/fleet/health"
	productService "github.com/tgs266/fleet/fleet/products"
	securedAuthService "github.com/tgs266/fleet/fleet/secured_auth"
	shipService "github.com/tgs266/fleet/fleet/ships"
	userService "github.com/tgs266/fleet/fleet/users"
	"github.com/tgs266/fleet/rest-gen/generated/com/fleet/auth"
	"github.com/tgs266/fleet/rest-gen/generated/com/fleet/cargo"
	"github.com/tgs266/fleet/rest-gen/generated/com/fleet/health"
	"github.com/tgs266/fleet/rest-gen/generated/com/fleet/products"
	"github.com/tgs266/fleet/rest-gen/generated/com/fleet/ships"
	"github.com/tgs266/fleet/rest-gen/generated/com/fleet/users"
	"github.com/tgs266/rest-gen/runtime/errors"
)

type Server struct {
	Router *gin.Engine
	Logger zerolog.Logger
	Config *config.Config
}

func New(logger zerolog.Logger, config *config.Config) *Server {
	logger.Info().Msg("creating fleet server")
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = ioutil.Discard
	router := gin.New()
	server := &Server{
		Router: router,
		Logger: logger,
		Config: config,
	}
	corsCfg := cors.DefaultConfig()
	corsCfg.AllowOrigins = []string{"http://localhost:3000"}
	corsCfg.AllowCredentials = true
	router.Use(cors.New(corsCfg))
	router.Use(RecoveryMiddleware(func(ctx *gin.Context, err error) {
		e := errors.GetError(err)
		ctx.JSON(e.Code(), e)
	}))
	router.Use(Middleware(logger))
	router.Use(ErrorMiddleware())

	shipSvc := shipService.ShipService{}
	productSvc := productService.ProductService{}
	cargoSvc := cargoService.CargoService{}
	healthSvc := healthService.HealthService{}
	authSvc := authService.AuthService{}
	securedAuthSvc := securedAuthService.SecuredAuthService{}
	userSvc := userService.UserService{}
	ships.RegisterShipServiceRoutes(router, ships.ShipServiceHandler{
		Handler: &shipSvc,
	})
	products.RegisterProductServiceRoutes(router, products.ProductServiceHandler{
		Handler: &productSvc,
	})
	cargo.RegisterCargoServiceRoutes(router, cargo.CargoServiceHandler{
		Handler: &cargoSvc,
	})
	health.RegisterHealthServiceRoutes(router, health.HealthServiceHandler{
		Handler: &healthSvc,
	})
	auth.RegisterAuthServiceRoutes(router, auth.AuthServiceHandler{
		Handler: &authSvc,
	})
	auth.RegisterSecuredAuthServiceRoutes(router, auth.SecuredAuthServiceHandler{
		Handler: &securedAuthSvc,
	})
	users.RegisterUserServiceRoutes(router, users.UserServiceHandler{
		Handler: &userSvc,
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
