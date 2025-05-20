package handler

import (
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/elginbrian/FILKOMPLAIN-BE/pkg/response"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type SystemInfo struct {
	DB        *gorm.DB
	StartTime time.Time
}

var systemInfo = SystemInfo{
	StartTime: time.Now(),
}

func SetupSystemInfo(db *gorm.DB) {
	systemInfo.DB = db
}

func WelcomeHandler(c *fiber.Ctx) error {
	dbStatus := "offline"
	dbName := "unknown"
	migrationStatus := "unknown"
	
	if systemInfo.DB != nil {
		sqlDB, err := systemInfo.DB.DB()
		if err == nil {
			if err = sqlDB.Ping(); err == nil {
				dbStatus = "connected"
			}
		}
		
		var result string
		if err := systemInfo.DB.Raw("SELECT current_database()").Scan(&result).Error; err == nil {
			dbName = result
		}
		
		var count int64
		if err := systemInfo.DB.Raw("SELECT count(*) FROM information_schema.tables WHERE table_name = 'migrations'").Count(&count).Error; err == nil {
			if count > 0 {
				var migrationsCount int64
				systemInfo.DB.Table("migrations").Count(&migrationsCount)
				migrationStatus = "completed - " + strconv.FormatInt(migrationsCount, 10) + " migrations applied"
			} else {
				var tableCount int64
				if err := systemInfo.DB.Raw("SELECT count(*) FROM information_schema.tables WHERE table_name IN ('users', 'reports')").Count(&tableCount).Error; err == nil {
					if tableCount > 0 {
						migrationStatus = "schema exists without migration tracking"
					} else {
						migrationStatus = "no schema detected"
					}
				}
			}
		}
	}
	
	env := os.Getenv("GO_ENV")
	if env == "" {
		env = os.Getenv("APP_ENV")
	}
	if env == "" {
		if os.Getenv("PRODUCTION") == "true" {
			env = "production"
		} else {
			env = "development"
		}
	}
	
	uptime := time.Since(systemInfo.StartTime).String()
	
	return c.JSON(response.Response{
		Success: true,
		Message: "Welcome to FILKOMPLAIN API",
		Data: fiber.Map{
			"status":  "online",
			"version": "1.0.0",
			"system": fiber.Map{
				"environment": env,
				"go_version":  runtime.Version(),
				"uptime":      uptime,
			},
			"database": fiber.Map{
				"status":    dbStatus,
				"name":      dbName,
				"migration": migrationStatus,
			},
		},
	})
}
