package cincinnati

import (
	"archive/tar"
	"compress/gzip"
	"io/fs"
	"net/http"
	"path/filepath"

	"github.com/blang/semver/v4"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/hongkailiu/test-go/pkg/util"
	"github.com/hongkailiu/test-go/pkg/version"
)

func GetHandler(r *gin.Engine, gb *GraphBuilder) http.Handler {
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

	prefix := "api/upgrades_info"
	r.GET(prefix+"/graph", getGraphHandler(gb))

	// ref. https://github.com/openshift/cincinnati/blob/49c914f804a1910ad3170b7724188643d404e49b/graph-builder/src/main.rs#L125
	r.GET(prefix+"/v1/graph", getGraphHandler(gb))

	r.GET(prefix+"/graph-data", func(c *gin.Context) {
		c.Header("Content-Type", "application/gzip")
		c.Header("Content-Disposition", `attachment; filename="graph-data.tar.gz"`)

		gz := gzip.NewWriter(c.Writer)
		defer func() {
			if err := gz.Close(); err != nil {
				logrus.WithError(err).Error("Error closing gzip writer")
			}
		}()

		tw := tar.NewWriter(gz)
		defer func() {
			if err := tw.Close(); err != nil {
				logrus.WithError(err).Error("Error closing tar writer")
			}
		}()

		for f := range graphDataFiles {
			file := filepath.Join(gb.graphDataDir, f)
			logrus.WithField("file", file).Warn("Add file to graph data")
			if err := util.AddFileToTar(tw, gb.graphDataDir, file); err != nil {
				logrus.WithError(err).WithField("file", file).Error("Error adding file to tar")
				c.Status(http.StatusInternalServerError)
				return
			}
		}

		for d := range graphDataDirs {
			dir := filepath.Join(gb.graphDataDir, d)
			if err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}

				if d.IsDir() {
					return nil
				}

				if d.Name() == "OWNERS" {
					return nil
				}

				if err := util.AddFileToTar(tw, gb.graphDataDir, path); err != nil {
					logrus.WithError(err).Error("Error adding file to tar")
					return err
				}

				return nil
			}); err != nil {
				c.String(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
				return
			}
		}

		c.Status(http.StatusOK)
	})

	return r.Handler()
}

var graphDataDirs = sets.New[string]("blocked-edges", "channels", "raw")
var graphDataFiles = sets.New[string]("LICENSE", "version")

func getGraphHandler(gb *GraphBuilder) gin.HandlerFunc {
	return func(c *gin.Context) {

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
	}
}
