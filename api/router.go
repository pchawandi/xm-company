package api

import (
	"context"
	"time"

	"github.com/pchawandi/xm-company/database"
	"github.com/pchawandi/xm-company/middleware"
	"golang.org/x/time/rate"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func ContextMiddleware(companyRepository CompanyRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("appCtx", companyRepository)
		c.Next()
	}
}

func NewRouter(ctx context.Context, db database.Database, logger *zap.Logger) *gin.Engine {
	companyRepository := NewCompanyRepository(ctx, db, logger)
	userRepository := NewUserRepository(ctx, db)

	r := gin.New()
	r.Use(ContextMiddleware(companyRepository))
	r.Use(middleware.Logger(logger))

	r.Use(middleware.RateLimiter(rate.Every(1*time.Minute), 600)) // 600 requests per minute

	v1 := r.Group("/api/v1")
	{
		// companies management routes
		v1.POST("/companies", middleware.JWTAuth(), companyRepository.Create)
		v1.GET("/companies/:id", companyRepository.Get)
		v1.PATCH("/companies/:id", middleware.JWTAuth(), companyRepository.Patch)
		v1.DELETE("/companies/:id", middleware.JWTAuth(), companyRepository.Delete)

		// user management routes
		v1.POST("/users/login", userRepository.LoginHandler)
		v1.POST("/users/register", userRepository.RegisterHandler)
	}

	return r
}
