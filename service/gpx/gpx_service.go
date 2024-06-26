package gpx

import (
	"context"
	"fmt"

	"api.backend.xjco2913/dao"
	"api.backend.xjco2913/service/sdto"
	"api.backend.xjco2913/service/sdto/errorx"
	"api.backend.xjco2913/util"
	"api.backend.xjco2913/util/zlog"
	"go.uber.org/zap"
)

type GPXService struct{}

var (
	gpxService GPXService
)

func Service() *GPXService {
	return &gpxService
}

// Store the gpx data as GEO type in mysql, and return route id
func (g *GPXService) ParseGPXData(ctx context.Context, in *sdto.ParseGPXDataInput) (*sdto.ParseGPXDataOutput, *errorx.ServiceErr) {
	gpxLonLatData, err := util.GPXToLonLat(in.GPXData)
	if err != nil {
		return nil, errorx.NewServicerErr(errorx.ErrExternal, "Invalid gpx format", nil)
	} else if len(gpxLonLatData) <= 1 {
		return nil, errorx.NewServicerErr(
			400,
			"Route data must have at least two points",
			nil,
		)
	}

	linestring := gpxLonLatData[0]
	for i := 1; i < len(gpxLonLatData); i++ {
		linestring += ", "
		linestring += gpxLonLatData[i]
	}
	// ST_GeomFromText('LINESTRING(?)')
	err = dao.DB.WithContext(ctx).Exec(
		fmt.Sprintf(
			"INSERT INTO GPSRoutes (path) VALUES (ST_GeomFromText('LINESTRING(%s)'));",
			linestring,
		),
	).Error
	if err != nil {
		zlog.Error("Error while store gpx route into mysql", zap.Error(err))
		return nil, errorx.NewInternalErr()
	}

	// Get last inserted route
	lastGPXRoute, err := dao.GetLastGPSRoute(ctx)
	if err != nil {
		zlog.Error("Error while get last inserted gps route", zap.Error(err))
		return nil, errorx.NewInternalErr()
	}

	return &sdto.ParseGPXDataOutput{
		RouteID: lastGPXRoute.ID,
	}, nil
}

// Store the [[lon, lat], [lon, lat]...] data into mysql and return route id
func (g *GPXService) ParseLonLatData(ctx context.Context, in *sdto.ParseLonLatDataInput) (*sdto.ParseLonLatDataOutput, *errorx.ServiceErr) {
	lonLatStrData := util.StrStrToGPXStr(in.LonLatData)
	if len(lonLatStrData) <= 1 {
		return nil, errorx.NewServicerErr(
			400,
			"Route data must have at least two points",
			nil,
		)
	}

	linestring := lonLatStrData[0]
	for i := 1; i < len(lonLatStrData); i++ {
		linestring += ", "
		linestring += lonLatStrData[i]
	}

	// ST_GeomFromText('LINESTRING(?)')
	err := dao.DB.WithContext(ctx).Exec(
		fmt.Sprintf(
			"INSERT INTO GPSRoutes (path) VALUES (ST_GeomFromText('LINESTRING(%s)'));",
			linestring,
		),
	).Error
	if err != nil {
		zlog.Error("Error while store gpx route into mysql", zap.Error(err))
		return nil, errorx.NewInternalErr()
	}

	// Get last inserted route
	lastGPXRoute, err := dao.GetLastGPSRoute(ctx)
	if err != nil {
		zlog.Error("Error while get last inserted gps route", zap.Error(err))
		return nil, errorx.NewInternalErr()
	}

	return &sdto.ParseLonLatDataOutput{
		RouteID: lastGPXRoute.ID,
	}, nil
}
