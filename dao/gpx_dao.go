package dao

import (
	"context"

	"api.backend.xjco2913/dao/model"
	"api.backend.xjco2913/dao/query"
)

func GetLastGPSRoute(ctx context.Context) (*model.GPSRoute, error) {
	g := query.Use(DB).GPSRoute

	last, err := g.WithContext(ctx).Last()
	if err != nil {
		return nil, err
	}

	return last, nil
}