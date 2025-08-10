package main

import (
	"log"
	"strings"

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
		AllowOrigins:     cfg.CORSOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	jwt := middleware.NewJWT(cfg.JWTSecret)
	authH := &handlers.AuthHandler{DB: dbConn, JWT: jwt}
	projH := &handlers.ProjectHandler{DB: dbConn}
	tsH := &handlers.TimesheetHandler{DB: dbConn}
	repH := &handlers.ReportHandler{DB: dbConn}
	acctH := &handlers.AccountHandler{DB: dbConn}

	api := r.Group("/api")
	{
		api.POST("/login", authH.Login)

		auth := api.Group("")
		auth.Use(jwt.AuthRequired())
		{
			admin := auth.Group("/admin")
			admin.Use(middleware.RequireRole("admin"))
			{
				admin.GET("/allprojects", projH.AllList)
				admin.POST("/projects", projH.Create)
				admin.PUT("/projects/:id", projH.Update)
				admin.DELETE("/projects/:id", projH.Delete)
				admin.POST("/reports/project-totals", repH.ProjectTotals)
				// CSV 导出
				admin.GET("/reports/project-totals.csv", repH.ProjectTotalsCSV)
				usersH := &handlers.UsersHandler{DB: dbConn}
				admin.POST("/users", usersH.Create)
				admin.GET("/users", usersH.List)
				// admin.DELETE("/users/:id", usersH.Delete)
				// 管理员重置用户密码
				admin.POST("/users/reset-password", usersH.ResetPassword)
			}
			auth.GET("/projects", projH.List)
			auth.GET("/timesheets/mine", tsH.ListMine)
			auth.POST("/timesheets", tsH.Create)
			auth.DELETE("/timesheets/:id", tsH.Delete)
			// 新增：更新本人工时
			auth.PUT("/timesheets/:id", tsH.Update)

			// 用户自助改密
			auth.POST("/change-password", acctH.ChangePassword)
		}
	}

	// API 路由...
	r.Static("/assets", "./frontend/dist/assets")
	r.StaticFile("/favicon.ico", "./frontend/dist/favicon.ico")

	r.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api/") {
			c.JSON(404, gin.H{"error": "not found"})
			return
		}
		c.File("./frontend/dist/index.html")
	})

	log.Println("server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
