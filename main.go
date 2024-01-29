package main

import (
	"Glue-API/controller"
	"Glue-API/docs"
	"log"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//	@title			Glue-API
//	@version		v1.0
//	@description	This is a GlueAPI server.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	윤여천
//	@contact.url	http://www.ablecloud.io
//	@contact.email	support@ablecloud.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@BasePath	/api/v1

//	@securityDefinitions.basic	BasicAuth

//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization
//	@description				Description for what is this security definition being used

func main() {
	// programmatically set swagger info

	docs.SwaggerInfo.Title = "Glue API"
	docs.SwaggerInfo.Description = "This is a GlueAPI server."
	docs.SwaggerInfo.Version = "1.0"
	//docs.SwaggerInfo.Host = ".swagger.io"
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	r := gin.Default()
	r.ForwardedByClientIP = true
	r.SetTrustedProxies(nil)
	c := controller.NewController()

	v1 := r.Group("/api/v1")
	{
		glue := v1.Group("/glue")
		{
			glue.GET("", c.GlueStatus)
			glue.GET("/version", c.GlueVersion)
		}
		pool := v1.Group("/pool")
		{
			pool.GET("", c.ListPools)
			pool.GET("/:pool_name", c.ListImages)
			pool.DELETE("/:pool_name", c.PoolDelete)
		}
		service := v1.Group("/service")
		{
			service.GET("", c.ServiceLs)
			service.POST("/:service_name", c.ServiceControl)
		}
		fs := v1.Group("/gluefs")
		{
			fs.GET("", c.FsStatus)
			fs.POST("/:fs_name", c.FsCreate)
			fs.DELETE("/:fs_name", c.FsDelete)
			fs.GET("/info/:fs_name", c.FsGetInfo)
			fs.GET("/list", c.FsList)
		}
		nfs := v1.Group("/nfs")
		{
			nfs.GET("", c.NfsClusterLs)
			nfs.GET("/:cluster_id", c.NfsClusterInfo)
			nfs.POST("/:cluster_id/:port", c.NfsClusterCreate)
			nfs.DELETE("/:cluster_id", c.NfsClusterDelete)
			nfs.GET("/export/:cluster_id", c.NfsExportDetailed)
			nfs.POST("/export/:cluster_id", c.NfsExportCreate)
			nfs.PUT("/export/:cluster_id", c.NfsExportUpdate)
			nfs.DELETE("/export/:cluster_id/:export_id", c.NfsExportDelete)
		}
		iscsi := v1.Group("/iscsi")
		{
			iscsi.POST("", c.IscsiServiceCreate)
			iscsi.GET("/target", c.IscsiTargetList)
			iscsi.POST("/target/:iqn_id/:hostname", c.IscsiTargetCreate)
			iscsi.DELETE("/target/:iqn_id", c.IscsiTargetDelete)
			iscsi.POST("/disk/:image_name", c.IscsiDiskCreate)

		}
		mirror := v1.Group("/mirror")
		{
			mirror.GET("", c.MirrorStatus) //Get Mirroring Status
			//Todo
			mirror.POST("", c.MirrorSetup) //Setup Mirroring
			//mirror.PATCH("", c.MirrorUpdate)  //Configure Mirroring
			mirror.DELETE("", c.MirrorDelete) //Unconfigure Mirroring
			//
			mirrorimage := mirror.Group("/image")
			{
				mirrorimage.GET("", c.MirrorImageList)                             //List Mirroring Images
				mirrorimage.GET("/:mirrorPool/:imageName", c.MirrorImageInfo)      //Get Image Mirroring Status
				mirrorimage.POST("/:mirrorPool/:imageName", c.MirrorImageSetup)    //Setup Image Mirroring
				mirrorimage.PATCH("/:mirrorPool/:imageName", c.MirrorImageUpdate)  //Config Image Mirroring
				mirrorimage.DELETE("/:mirrorPool/:imageName", c.MirrorImageDelete) //Unconfigure Mirroring

				mirrorimage.GET("/promote/:mirrorPool/:imageName", c.MirrorImagestatus)   //Promote Image
				mirrorimage.POST("/promote/:mirrorPool/:imageName", c.MirrorImagePromote) //
				mirrorimage.DELETE("/promote/:mirrorPool/:imageName", c.MirrorImageDemote)
			}
			//
			//
		}
		/*
			admin := v1.Group("/admin")
			{
				admin.Use(auth())
				admin.POST("/auth", c.Auth)
			}
		*/
		r.Any("/version", c.Version)
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run(":8080")
}

/*
func auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if len(c.GetHeader("Authorization")) == 0 {
			httputil.NewError(c, http.StatusUnauthorized, errors.New("Authorization is required Header"))
			c.Abort()
		}
		c.Next()
	}
}
*/
