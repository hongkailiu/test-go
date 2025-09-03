package main

import (
	"context"
	"encoding/json"
	"net/http"
	
	"github.com/sirupsen/logrus"
)

func getRouter(ctx context.Context, graphService GraphService) *http.ServeMux {
	handler := http.NewServeMux()

	handler.HandleFunc("/graph", func(w http.ResponseWriter, r *http.Request) {
		registry, err := graphService.Get(ctx)
		if err != nil {
			logrus.WithError(err).Error("Error loading registry data")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(registry); err != nil {
			logrus.WithError(err).WithField("registry", registry).Error("failed to encode page")
		}
	})

	return handler
}
