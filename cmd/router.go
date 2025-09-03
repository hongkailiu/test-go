package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

func getRouter(ctx context.Context, graphService GraphService) *http.ServeMux {
	handler := http.NewServeMux()

	handler.HandleFunc("/graph", func(w http.ResponseWriter, r *http.Request) {
		graph, err := graphService.Get(ctx)
		if err != nil {
			logrus.WithError(err).Error("Error loading registry data")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(graph); err != nil {
			logrus.WithError(err).WithField("graph", graph).Error("failed to encode graph")
		}
	})

	handler.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		_, err := graphService.Get(ctx)
		if err != nil {
			logrus.WithError(err).Error("Error loading registry data")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = fmt.Fprintf(w, "OK")
		logrus.WithError(err).Error("failed to write response")
	})

	return handler
}
