package routes

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pardnchiu/go-podrun/internal/database"
	"github.com/pardnchiu/go-podrun/internal/handler"
	"github.com/pardnchiu/go-podrun/internal/utils"
)

func New(db *database.SQLite) error {
	if handler.DB == nil {
		handler.DB = db
	}

	r := gin.Default()

	ip, err := utils.GetLocalIP()
	if err != nil {
		return err
	}

	r.SetTrustedProxies([]string{
		"127.0.0.1",
		ip,
	})

	// * Pod > GET
	r.GET("/api/pod/list", handler.GetAPIPodList)

	// * Pod > POST
	r.POST("/api/pod/upsert", handler.PostAPIPodUpsert)
	r.POST("/api/pod/update/:uid", handler.PostAPIPodRecordUpdate)
	r.POST("/api/pod/record/insert", handler.GetAPIPodRecordInsert)

	// # NOT THIS PROJECT POINT, REMOVE IT FOR NOW
	// // * User > POST
	// r.POST("/api/user/upsert", handler.PostAPIUserUpsert)

	// * Other
	r.GET("/api/health", func(ctx *gin.Context) {
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
