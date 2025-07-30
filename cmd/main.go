package main

import (
	"contentsystem/internal/api"
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	api.CmsRouters(r)

	err := r.Run(":8081")
	if err != nil {
		fmt.Printf("r run error = %v", err)
		return
	}

}
