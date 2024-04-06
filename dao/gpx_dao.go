package dao

import (
	"context"
	"fmt"

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

// SELECT ST_ASTEXT(path) from GPSRoutes;
func GetPathAsText(ctx context.Context, routeId int32) (string, error) {
	stat := fmt.Sprintf(
		"SELECT ST_ASTEXT(path) FROM GPSRoutes WHERE id = %v",
		routeId,
	)

	var path string
	err := DB.WithContext(ctx).Raw(stat).Scan(&path).Error
	if err != nil {
		return "", err
	}

	return path, nil
}