package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/example/worklog-system/internal/config"
	"github.com/example/worklog-system/internal/db"
	"github.com/example/worklog-system/internal/handlers"
	"github.com/example/worklog-system/internal/middleware"
	"github.com/example/worklog-system/internal/seed"
)

func main() {
	cfg := config.Load()
	dbConn := db.Connect(cfg)
	seed.Run(dbConn)

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins: cfg.CORSOrigins,
		AllowMethods: []string{"GET","POST","PUT","DELETE","OPTIONS"},
		AllowHeaders: []string{"Authorization","Content-Type"},
		AllowCredentials: true,
	}))

	jwt := middleware.NewJWT(cfg.JWTSecret)
	authH := &handlers.AuthHandler{DB: dbConn, JWT: jwt}
	projH := &handlers.ProjectHandler{DB: dbConn}
	tsH := &handlers.TimesheetHandler{DB: dbConn}
	repH := &handlers.ReportHandler{DB: dbConn}

	api := r.Group("/api")
	{
		api.POST("/login", authH.Login)

		auth := api.Group("")
		auth.Use(jwt.AuthRequired())
		{
			admin := auth.Group("/admin")
			admin.Use(middleware.RequireRole("admin"))
			{
				admin.GET("/projects", projH.List)
				admin.POST("/projects", projH.Create)
				admin.PUT("/projects/:id", projH.Update)
				admin.DELETE("/projects/:id", projH.Delete)
				admin.POST("/reports/project-totals", repH.ProjectTotals)
				admin.GET("/reports/project-totals.csv", repH.ProjectTotalsCSV) // CSV 导出
			}

			auth.GET("/projects", projH.List)
			auth.GET("/timesheets/mine", tsH.ListMine)
			auth.POST("/timesheets", tsH.Create)
			auth.DELETE("/timesheets/:id", tsH.Delete)
		}
	}

	log.Println("server starting on :8080")
	if err := r.Run(":8080"); err != nil { panic(err) }
}
