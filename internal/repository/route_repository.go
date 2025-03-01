package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/flohansen/walker/generated/database"
	"github.com/flohansen/walker/internal/controller"
)

//go:generate sqlc generate -f ../../sqlc.yaml

type RouteSQLRepository struct {
	Queries *database.Queries
}

func (r *RouteSQLRepository) CreateRoute(ctx context.Context, route controller.Route) (int32, error) {
	return r.Queries.CreateRoute(ctx, route.Name)
}

func (r *RouteSQLRepository) DeleteRoute(ctx context.Context, route controller.Route) error {
	_, err := r.Queries.DeleteRoute(ctx, route.ID)
	return err
}

func (r *RouteSQLRepository) GetRoute(ctx context.Context, id int32) (*controller.Route, error) {
	row, err := r.Queries.GetRoute(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &controller.Route{
		ID:   row.ID,
		Name: row.Name,
	}, nil
}

func (r *RouteSQLRepository) GetRoutes(ctx context.Context) ([]controller.Route, error) {
	rows, err := r.Queries.GetRoutes(ctx)
	if err != nil {
		return nil, err
	}

	var routes []controller.Route
	for _, row := range rows {
		routes = append(routes, controller.Route{
			ID:   row.ID,
			Name: row.Name,
		})
	}

	return routes, nil
}

func (r *RouteSQLRepository) UpdateRoute(ctx context.Context, route controller.Route) error {
	return r.Queries.UpdateRoute(ctx, database.UpdateRouteParams{
		ID:   route.ID,
		Name: route.Name,
	})
}
