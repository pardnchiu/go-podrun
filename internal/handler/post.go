package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pardnchiu/go-podrun/internal/model"
)

// # NOT THIS PROJECT POINT, REMOVE IT FOR NOW
// func PostAPIUserUpsert(ctx *gin.Context) {
// 	var user model.User
// 	if err := ctx.ShouldBindJSON(&user); err != nil {
// 		ctx.String(http.StatusBadRequest, err.Error())
// 		return
// 	}

// 	if err := DB.UpsertUser(ctx.Request.Context(), &user); err != nil {
// 		ctx.String(http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	ctx.String(http.StatusOK, "ok")
// }

func PostAPIPodUpsert(ctx *gin.Context) {
	var pod model.Pod
	if err := ctx.ShouldBindJSON(&pod); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	if err := DB.UpsertPod(ctx.Request.Context(), &pod); err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.String(http.StatusOK, "ok")
}

func PostAPIPodRecordUpdate(ctx *gin.Context) {
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

	if err := DB.UpdatePod(ctx.Request.Context(), &pod); err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.String(http.StatusOK, "ok")
}

func GetAPIPodRecordInsert(ctx *gin.Context) {
	var record model.Record
	if err := ctx.ShouldBindJSON(&record); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	if err := DB.InsertRecord(ctx.Request.Context(), &record); err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.String(http.StatusOK, "ok")
}
