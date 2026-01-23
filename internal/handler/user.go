package handler

// # NOT THIS PROJECT POINT, REMOVE IT FOR NOW
// func getAPIUserList(ctx *gin.Context) {
// 	results, err := DB.ListPods(ctx.Request.Context())
// 	if err != nil {
// 		ctx.String(http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, gin.H{"data": results})
// }

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
