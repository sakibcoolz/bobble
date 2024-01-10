package main

import (
	"bobble/utils"
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	syMap = new(sync.Map)
)

func init() {
	func() {
		for _, v := range utils.Cache(nil, http.StatusOK, "sync/getall") {
			for key, val := range utils.ToMap(v) {
				syMap.Store(key, val)
			}
		}
	}()
}

func main() {
	r := gin.Default()
	hostport := os.Getenv("HOSTPORT")
	r.GET("/set", func(ctx *gin.Context) {
		var requestBodyMap map[string]string
		if err := ctx.ShouldBindJSON(&requestBodyMap); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		for key, val := range requestBodyMap {
			syMap.Store(key, val)
		}

		go utils.Cache(requestBodyMap, http.StatusCreated, "sync/set")

		ctx.String(http.StatusCreated, "")
	})

	r.GET("/get/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")

		val, ok := syMap.Load(id)
		if !ok {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "not found"})

			return
		}

		ctx.JSON(http.StatusOK, gin.H{id: val})
	})

	r.GET("/delete/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")

		syMap.Delete(id)

		go utils.Cache(nil, http.StatusOK, fmt.Sprintf("sync/delete/%s", id))

		ctx.JSON(http.StatusOK, gin.H{id: "deleted"})
	})

	back := r.Group("/sync")

	back.GET("/set", func(ctx *gin.Context) {
		var requestBodyMap map[string]string
		if err := ctx.ShouldBindJSON(&requestBodyMap); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		for key, val := range requestBodyMap {
			syMap.Store(key, val)
		}

		ctx.String(http.StatusCreated, "")
	})

	back.GET("/delete/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")

		syMap.Delete(id)

		ctx.JSON(http.StatusOK, gin.H{id: "deleted"})
	})

	back.GET("/getall", func(ctx *gin.Context) {

		m := map[string]string{}
		syMap.Range(func(key, value interface{}) bool {
			m[fmt.Sprint(key)] = value.(string)
			return true
		})

		ctx.JSON(http.StatusOK, m)
	})

	r.Run(hostport)
}
