package server

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"user/adapter/apihook"
	"user/adapter/emailer"
	"user/adapter/storage/database"
	"user/config"
	"user/consts"
	"user/internal/controller"
	"user/internal/repo"
	"user/internal/usecase"
	"user/middlewares"
	"user/migrations"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type HTTPServer interface {
	Run() error
	ShutDown() error
}

type GinServer struct {
	Cfg    *config.EnvConfig
	Router *gin.Engine
}

func New() *GinServer {
	cfg, err := config.Load(consts.AppName)
	if err != nil {
		panic(err)
	}
	db, err := database.Connect(cfg.DB)
	if err != nil {
		log.Printf("unable to connect the database : %v", err)
		return nil
	}

	if err := migrations.Up(db); err != nil {
		log.Printf("migrations up failed : %v", err)
		return nil
	}
	// here initializing the router
	router := initRouter()
	if !cfg.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	api := router.Group("/api/:version/users")
	mailSrv := emailer.NewGmail(
		cfg.Gmail.Name,
		cfg.Gmail.Address,
		cfg.Gmail.Password,
	)

	apiSrv := apihook.HTTPAPI{}

	// Initialize controllers and use cases.
	userRepo := repo.NewUser(db, cfg)
	userUsecase := usecase.NewUser(userRepo, cfg, mailSrv, apiSrv)
	// Initialize middleware
	mw := middlewares.NewMiddlewares(cfg, userUsecase)
	api.Use(mw.APIVersionGuard(middlewares.VersionOptions{
		AcceptedVersions: cfg.AcceptedVersions,
	}))

	api.Use(mw.Authorize())
	{
		userCtrl := controller.NewUser(api, userUsecase, cfg)
		userCtrl.InitRoutes()
	}

	return &GinServer{
		Cfg:    cfg,
		Router: router,
	}
}

func (gs GinServer) Run() error {
	log.Print("Starting gin server in...", gs.Cfg.Port)
	return gs.launch()
}

func (gs GinServer) ShutDown() error {
	return nil
}

func initRouter() *gin.Engine {
	router := gin.Default()
	gin.SetMode(gin.DebugMode)

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "POST", "DELETE", "GET", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	return router
}

func (gs GinServer) launch() error {
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%v", gs.Cfg.Port),
		Handler:           gs.Router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	return srv.ListenAndServe()
}
