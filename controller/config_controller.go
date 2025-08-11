package controller

import (
	"github.com/gin-gonic/gin"
)

type ConfigController interface {
	CreateConfig(ctx *gin.Context)
	UpdateConfig(ctx *gin.Context)
	RollbackConfig(ctx *gin.Context)
	FetchConfig(ctx *gin.Context)
	ListVersions(ctx *gin.Context)
}
