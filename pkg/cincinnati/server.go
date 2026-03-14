package cincinnati

import (
	"fmt"
	"net/http"

	"github.com/blang/semver/v4"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Start(addr ...string) error {
	r := gin.Default()
	repo := Repo{}
	g := NewGraphGenerator(repo.tagsToNode)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
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
					"value": fmt.Sprintf("invalid parameter version: %s", versionStr),
				})
				return
			}
			params.Version = v
		}

		if id := c.Query("id"); id != "" {
			params.Id = id
		}

		graph, err := g.Generate(params)
		if err != nil {
			logrus.WithError(err).Error("Error generating graph")
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": http.StatusText(http.StatusInternalServerError),
			})
			return
		}
		c.JSON(http.StatusOK, graph)
	})
	return r.Run(addr...)
}
