package app

import (
	"config-service/controller"
	"config-service/exception"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(configController controller.ConfigController, schemaController controller.SchemaController) *gin.Engine {
	// Create Gin engine
	router := gin.Default()

	// Global error handling (replace PanicHandler)
	router.Use(func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				exception.ErrorHandler(c.Writer, c.Request, r)
				c.Abort()
			}
		}()
		c.Next()
	})

	// Config routes
	configs := router.Group("/configs")
	{
		configs.POST("/:schema/:name", configController.CreateConfig)
		configs.PUT("/:schema/:name", configController.UpdateConfig)
		configs.POST("/:schema/:name/rollback", configController.RollbackConfig)
		configs.GET("/:schema/:name", configController.FetchConfig)
		configs.GET("/:schema/:name/versions", configController.ListVersions)
	}

	// Swagger UI
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Display individual schema
	router.GET("/schemas/:name", schemaController.GetSchema)

	return router
}
