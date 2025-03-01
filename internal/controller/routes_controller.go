package controller

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type Route struct {
	ID   int32
	Name string
}

type RoutesRepository interface {
	CreateRoute(ctx context.Context, route Route) (int32, error)
	GetRoute(ctx context.Context, id int32) (*Route, error)
	GetRoutes(ctx context.Context) ([]Route, error)
	UpdateRoute(ctx context.Context, route Route) error
	DeleteRoute(ctx context.Context, route Route) error
}

type RoutesController struct {
	mux  *http.ServeMux
	repo RoutesRepository
}

func NewRoutes(
	repo RoutesRepository,
) *RoutesController {
	h := &RoutesController{
		mux:  http.NewServeMux(),
		repo: repo,
	}

	h.mux.HandleFunc("GET /healthz", healthz)
	h.mux.HandleFunc("POST /routes", h.postRoutes)
	h.mux.HandleFunc("PUT /routes/{id}", h.putRoute)
	h.mux.HandleFunc("GET /routes", h.getRoutes)
	h.mux.HandleFunc("GET /routes/{id}", h.getRoute)
	h.mux.HandleFunc("DELETE /routes/{id}", h.deleteRoute)

	return h
}

func (h *RoutesController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

func healthz(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

type PostRoutesRequest struct {
	Name string `json:"name"`
}

type PostRoutesResponse struct {
	ID int32 `json:"id"`
}

func (h *RoutesController) postRoutes(w http.ResponseWriter, r *http.Request) {
	var req PostRoutesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := h.repo.CreateRoute(r.Context(), Route{
		Name: req.Name,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(PostRoutesResponse{
		ID: id,
	})
}

type PutRouteRequest struct {
	Name string `json:"name"`
}

type PutRouteResponse struct {
}

func (h *RoutesController) putRoute(w http.ResponseWriter, r *http.Request) {
	idVal := r.PathValue("id")
	if idVal == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idVal)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var req PutRouteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.repo.UpdateRoute(r.Context(), Route{
		ID:   int32(id),
		Name: req.Name,
	}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(PutRouteResponse{})
}

type GetRoutesResponse struct {
	Routes []Route `json:"routes,omitempty"`
}

func (h *RoutesController) getRoutes(w http.ResponseWriter, r *http.Request) {
	routes, err := h.repo.GetRoutes(r.Context())
	if err != nil {
		log.Printf("could not get routes: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(GetRoutesResponse{
		Routes: routes,
	})
}

type GetRouteResponse struct {
	Route Route `json:"route,omitempty"`
}

func (h *RoutesController) getRoute(w http.ResponseWriter, r *http.Request) {
	idVal := r.PathValue("id")
	if idVal == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idVal)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	route, err := h.repo.GetRoute(r.Context(), int32(id))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if route == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(GetRouteResponse{
		Route: *route,
	})
}

type DeleteRouteResponse struct {
}

func (h *RoutesController) deleteRoute(w http.ResponseWriter, r *http.Request) {
	idVal := r.PathValue("id")
	if idVal == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idVal)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.repo.DeleteRoute(r.Context(), Route{ID: int32(id)}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(DeleteRouteResponse{})
}
