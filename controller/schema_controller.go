package controller

import "github.com/gin-gonic/gin"

type SchemaController interface {
	GetSchema(c *gin.Context)
}
