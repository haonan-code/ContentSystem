package api

import (
	"contentsystem/internal/services"
	"github.com/gin-gonic/gin"
)

const (
	rootPath   = "/api/"
	noAuthPath = "/out/api/"
)

func CmsRouters(r *gin.Engine) {
	cmsApp := services.NewCmsApp()
	session := NewSessionAuth()

	root := r.Group(rootPath).Use(session.Auth)
	{
		root.GET("/cms/hello", cmsApp.Hello)
		root.POST("/cms/content/create", cmsApp.ContentCreate)
		root.POST("/cms/content/update", cmsApp.ContentUpdate)
		root.POST("/cms/content/delete", cmsApp.ContentDelete)
		root.POST("/cms/content/find", cmsApp.ContentFind)
	}

	noAuth := r.Group(noAuthPath)
	{
		// /out/api/cms/register
		noAuth.POST("/cms/register", cmsApp.Register)
		// /out/api/cms/login
		noAuth.POST("/cms/login", cmsApp.Login)
	}
}
