package cincinnati

import (
	"net/http"

	"github.com/blang/semver/v4"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/hongkailiu/test-go/pkg/version"
)



func releaseMode() bool {
	return gin.Mode() == gin.ReleaseMode
}

func GetHandler(r *gin.Engine,gb *GraphBuilder) http.Handler {
	r.GET("/version", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"name":    version.Name,
			"version": version.Version,
		})
	})

	r.GET("/healthz", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	r.GET("/readyz", func(c *gin.Context) {
		if ready := gb.Ready(); ready {
			c.Status(http.StatusOK)
		} else {
			c.Status(http.StatusServiceUnavailable)
		}
	})

	r.GET("/upgrades_info/v1/graph", func(c *gin.Context) {
		channel := c.Query("channel")
		if channel == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"kind":  "missing_params",
				"value": "mandatory client parameters missing: channel",
			})
		}

		arch := c.DefaultQuery("arch", "amd64")
		params := GraphParams{
			Channel: channel,
			Arch:    arch,
		}

		if versionStr := c.Query("version"); versionStr != "" {
			v, err := semver.Parse(versionStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"kind":  "invalid_params",
					"value": "invalid parameter version: " + versionStr,
				})

				return
			}

			params.Version = v
		}

		if id := c.Query("id"); id != "" {
			params.Id = id
		}

		graph, err := gb.Build(params)
		if err != nil {
			logrus.WithError(err).Error("Error generating graph")
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": http.StatusText(http.StatusInternalServerError),
			})

			return
		}

		c.JSON(http.StatusOK, graph)
	})

	return r.Handler()
}
