package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pardnchiu/go-podrun/internal/database"
)

var (
	DB *database.SQLite
)

func GetAPIPodList(ctx *gin.Context) {
	containers, err := DB.ListPods(ctx.Request.Context())
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": containers})
}

func GetAPIUserList(ctx *gin.Context) {
	results, err := DB.ListPods(ctx.Request.Context())
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": results})
}
