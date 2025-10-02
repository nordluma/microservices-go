package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/nordluma/microservices-go/pkg/discovery"
)

type Handler struct {
	register discovery.Registry
}

func NewHandler(register discovery.Registry) *Handler {
	return &Handler{register: register}
}

func (h *Handler) Register(w http.ResponseWriter, req *http.Request) {
	serviceName := req.FormValue("serviceName")
	if serviceName == "" {
		log.Printf("missing serviceName")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	instanceID := req.FormValue("instanceId")
	if instanceID == "" {
		log.Printf("missing instanceID")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	hostPort := req.FormValue("hostPort")
	if hostPort == "" {
		log.Printf("missing hostPort")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := h.register.Register(req.Context(), instanceID, serviceName, hostPort)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Deregister(w http.ResponseWriter, req *http.Request) {
	instanceID := req.FormValue("instanceId")
	if instanceID == "" {
		log.Printf("missing instanceId")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	serviceName := req.FormValue("serviceName")
	if serviceName == "" {
		log.Printf("missing serviceName")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := h.register.Deregister(req.Context(), instanceID, serviceName)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Discover(w http.ResponseWriter, req *http.Request) {
	serviceName := req.FormValue("serviceName")
	if serviceName == "" {
		log.Printf("missing serviceName")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	instances, err := h.register.Discover(req.Context(), serviceName)
	if err != nil {
		log.Println(err)
		if errors.Is(err, discovery.ErrNotFound) {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = json.NewEncoder(w).Encode(instances); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *Handler) HealthCheck(w http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()
	instanceID := query.Get("instanceId")
	if instanceID == "" {
		log.Println("missing instanceId")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	serviceName := query.Get("serviceName")
	if serviceName == "" {
		log.Println("missing serviceName")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.register.HealthCheck(instanceID, serviceName); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
