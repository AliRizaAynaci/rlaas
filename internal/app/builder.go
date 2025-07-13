package app

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/AliRizaAynaci/rlaas/internal/app/health"
	"github.com/AliRizaAynaci/rlaas/internal/auth"
	"github.com/AliRizaAynaci/rlaas/internal/check"
	"github.com/AliRizaAynaci/rlaas/internal/config"
	"github.com/AliRizaAynaci/rlaas/internal/database"
	"github.com/AliRizaAynaci/rlaas/internal/middleware"
	"github.com/AliRizaAynaci/rlaas/internal/project"
	"github.com/AliRizaAynaci/rlaas/internal/rule"
	"github.com/AliRizaAynaci/rlaas/internal/service"
	"github.com/AliRizaAynaci/rlaas/internal/user"
)

func New() *fiber.App {
	cfg := config.Load()

	/* ------------ DB ------------ */
	db, err := database.Connect(cfg.DSN)
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}

	if err := database.Migrate(db,
		&user.User{},
		&project.Project{},
		&rule.Rule{},
	); err != nil {
		log.Fatalf("db migrate: %v", err)
	}

	/* ------------ Services ------------ */
	userSvc := user.NewService(user.NewGormRepo(db))
	projSvc := project.NewService(project.NewGormRepo(db))
	ruleSvc := rule.NewService(rule.NewGormRepo(db), db)
	rateCfgSvc := service.NewRateConfigService(db)

	/* ------------ Handlers ------------ */
	userHdl := user.NewHandler(userSvc)
	projHdl := project.NewHandler(projSvc)
	ruleHdl := rule.NewHandler(ruleSvc)
	checkH := check.NewHandler(rateCfgSvc)
	healthH := health.New(db)

	/* ------------ Fiber ------------ */
	app := fiber.New()
	srv := app.Server()
	srv.ReadTimeout = 10 * time.Second
	srv.WriteTimeout = 15 * time.Second
	app.Use(middleware.RequestLogger())

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "https://rlaas.tech",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS,PATCH",
		AllowHeaders:     "Accept, Authorization, Content-Type, X-CSRF-Token",
		AllowCredentials: true,
	}))

	/* ------------ Public routes ------------ */
	app.Get("/healthz", healthH.Liveness) // liveness
	app.Get("/readyz", healthH.Readiness) // readiness

	app.Get("/auth/google/login", auth.Login)
	app.Get("/auth/google/callback", auth.Callback(userSvc))
	app.Get("/logout", auth.Logout)
	app.Post("/check", checkH.Handle)

	/* ------------ Protected routes ------------ */
	api := app.Group("/", middleware.Auth())
	api.Get("/me", userHdl.Me)

	/* --- Projects --- */
	api.Post("/projects", projHdl.Create)
	api.Get("/projects", projHdl.List)
	api.Delete("/projects/:pid", projHdl.Delete)

	/* --- Nested Rules --- */
	rules := api.Group("/projects/:pid/rules")
	rules.Get("/", ruleHdl.List)
	rules.Post("/", ruleHdl.Create)
	rules.Put("/:rid", ruleHdl.Update)
	rules.Delete("/:rid", ruleHdl.Delete)

	return app
}
