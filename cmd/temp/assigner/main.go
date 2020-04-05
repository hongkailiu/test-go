package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/hongkailiu/test-go/pkg/assigner"
)

func main() {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	s := assigner.NewService()

	r.GET("/", func(c *gin.Context) {
		s, err := s.GetStatus()
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, s)
	})

	ticker := time.NewTicker(5 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				handle(s)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	r.Run(":8080")
}

func handle(s assigner.Service) {
	c, err := s.GetConfig()
	if err != nil {
		logrus.WithError(err).Error("Failed to get config")
		return
	}
	g, err := s.GetGroup(c.GroupName)
	if err != nil {
		logrus.WithError(err).Error("Failed to get group")
		return
	}
	newGroup, err := determine(time.Now(), c)
	if err != nil {
		logrus.WithError(err).Error("Failed to determine group")
		return
	}
	if !isSame(g, newGroup) {
		if err := s.SetGroup(c.GroupName, newGroup); err != nil {
			logrus.WithError(err).Error("Failed to set group")
		}
	}
}

func determine(now time.Time, c assigner.Config) ([]string, error) {
	found := -1
	for i, a := range c.ScheduledActions {
		if a.At.Before(now) {
			if found == -1 {
				found = i
				continue
			}
			if c.ScheduledActions[found].At.Before(a.At) {
				found = i
			}
		}
	}
	if found == -1 {
		return nil, fmt.Errorf("not found")
	}
	return c.ScheduledActions[found].Members, nil
}

func isSame(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
