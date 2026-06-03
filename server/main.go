package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"carbon-tracker/config"
	"carbon-tracker/middleware"
	"carbon-tracker/models"
	"carbon-tracker/routes"
	"carbon-tracker/services"
)

func main() {
	cfg := config.LoadConfig()

	db, err := gorm.Open(mysql.Open(cfg.DSN()), &gorm.Config{})
	if err != nil {
		log.Fatal("数据库连接失败: ", err)
	}

	db.AutoMigrate(&models.User{}, &models.EmissionFactor{}, &models.CarbonRecord{}, &models.ReductionGoal{}, &models.RegionalAverage{})

	authService := services.NewAuthService(db, cfg)
	carbonService := services.NewCarbonService(db)
	goalService := services.NewGoalService(db)

	authHandler := routes.NewAuthHandler(authService)
	carbonHandler := routes.NewCarbonHandler(carbonService)
	goalHandler := routes.NewGoalHandler(goalService)

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	api := r.Group("/api")
	{
		api.POST("/register", authHandler.Register)
		api.POST("/login", authHandler.Login)

		api.GET("/factors", carbonHandler.GetFactors)

		auth := api.Group("/")
		auth.Use(middleware.AuthMiddleware(cfg))
		{
			auth.POST("/records", carbonHandler.CreateRecord)
			auth.GET("/records", carbonHandler.GetDailyRecords)
			auth.DELETE("/records/:id", carbonHandler.DeleteRecord)

			auth.GET("/trend/weekly", carbonHandler.GetWeeklyTrend)
			auth.GET("/trend/monthly", carbonHandler.GetMonthlyTrend)
			auth.GET("/summary/category", carbonHandler.GetCategorySummary)

			auth.POST("/goals", goalHandler.CreateGoal)
			auth.GET("/goals/active", goalHandler.GetActiveGoal)
			auth.GET("/goals/history", goalHandler.GetGoalHistory)
		}
	}

	log.Printf("服务器启动在端口 %s", cfg.ServerPort)
	r.Run(":" + cfg.ServerPort)
}
