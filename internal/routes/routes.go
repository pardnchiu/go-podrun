package routes

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pardnchiu/go-podrun/internal/database"
	"github.com/pardnchiu/go-podrun/internal/model"
)

func New(db *database.SQLite) error {
	r := gin.Default()

	r.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "adsf")
	})
	r.GET("/pod/list", func(ctx *gin.Context) {
		containers, err := db.ListContainers(ctx.Request.Context())
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"data": containers})
	})
	r.POST("/pod/upsert", func(ctx *gin.Context) {
		var pod model.Pod
		if err := ctx.ShouldBindJSON(&pod); err != nil {
			ctx.String(http.StatusBadRequest, err.Error())
			return
		}

		if err := db.UpsertPod(ctx.Request.Context(), &pod); err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		ctx.String(http.StatusOK, "ok")
	})
	r.POST("/pod/update/:uid", func(ctx *gin.Context) {
		var pod model.Pod
		if err := ctx.ShouldBindJSON(&pod); err != nil {
			ctx.String(http.StatusBadRequest, err.Error())
			return
		}

		pod.UID = ctx.Param("uid")

		if pod.UID == "" {
			ctx.String(http.StatusBadRequest, "uid is required")
			return
		}

		if err := db.UpdatePod(ctx.Request.Context(), &pod); err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		ctx.String(http.StatusOK, "ok")
	})
	r.NoRoute(func(c *gin.Context) {
		select {}
	})

	log.Println("start on :8080")
	if err := r.Run(":8080"); err != nil {
		return err
	}

	return nil
}
