package moment

import (
	"context"
	"fmt"
	"time"

	"api.backend.xjco2913/dao"
	"api.backend.xjco2913/dao/minio"
	"api.backend.xjco2913/dao/model"
	"api.backend.xjco2913/service/sdto"
	"api.backend.xjco2913/service/sdto/errorx"
	"api.backend.xjco2913/util"
	"api.backend.xjco2913/util/zlog"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	MOMENT_FEED_LIMIT = 10
)

type MomentService struct{}

var (
	momentService MomentService
)

func Service() *MomentService {
	return &momentService
}

func (m *MomentService) Create(ctx context.Context, in *sdto.CreateMomentInput) *errorx.ServiceErr {
	momentId := uuid.New()
	momentIdStr := momentId.String()
	_, err := dao.CreateNewMoment(ctx, &model.Moment{
		AuthorID: in.UserID,
		Content:  &in.Content,
		MomentID: momentIdStr,
	})
	if err != nil {
		zlog.Error("error while create new moment", zap.Error(err))
		return errorx.NewInternalErr()
	}

	return nil
}

func (m *MomentService) CreateWithImage(ctx context.Context, in *sdto.CreateMomentImageInput) *errorx.ServiceErr {
	imageName := uuid.New()
	err := minio.UploadMomentImage(ctx, imageName.String(), in.ImageData)
	if err != nil {
		zlog.Error("error while store moment image into minio", zap.Error(err))
		return errorx.NewInternalErr()
	}

	imageNameStr := imageName.String()
	momentId := uuid.New()
	momentIdStr := momentId.String()
	_, err = dao.CreateNewMoment(ctx, &model.Moment{
		AuthorID: in.UserID,
		Content:  &in.Content,
		ImageURL: &imageNameStr,
		MomentID: momentIdStr,
	})
	if err != nil {
		zlog.Error("error while create new moment", zap.Error(err))

		// async remove image in minio
		go func() {
			minio.RemoveObjectFromMoment(ctx, imageNameStr)
		}()

		return errorx.NewInternalErr()
	}

	return nil
}

func (m *MomentService) CreateWithVideo(ctx context.Context, in *sdto.CreateMomentVideoInput) *errorx.ServiceErr {
	videoName := uuid.New()
	videoNameStr := videoName.String()

	momentId := uuid.New()
	momentIdStr := momentId.String()
	newMomentID, err := dao.CreateNewMoment(ctx, &model.Moment{
		AuthorID: in.UserID,
		Content:  &in.Content,
		VideoURL: &videoNameStr,
		MomentID: momentIdStr,
	})
	if err != nil {
		zlog.Error("error while create new moment", zap.Error(err))
		return errorx.NewInternalErr()
	}

	// async upload video
	go func() {
		err := minio.UploadMomentVideo(ctx, videoNameStr, in.VideoData)
		if err != nil {
			// error, remove moment record in mysql
			zlog.Error("error while async upload moment video", zap.Error(err))
			dao.DeleteMomentByID(ctx, newMomentID)
		}
	}()

	return nil
}

func (m *MomentService) CreateWithGPX(ctx context.Context, in *sdto.CreateMomentGPXInput) *errorx.ServiceErr {
	gpxLonLatData, err := util.GPXToLonLat(in.GPXData)
	if err != nil {
		return errorx.NewServicerErr(errorx.ErrExternal, "invalid gpx format", nil)
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
		zlog.Error("error while store gpx route into mysql", zap.Error(err))
		return errorx.NewInternalErr()
	}

	// get last inserted route
	lastGPXRoute, err := dao.GetLastGPSRoute(ctx)
	if err != nil {
		zlog.Error("error while get last inserted gps route", zap.Error(err))
		return errorx.NewInternalErr()
	}
	momentId := uuid.New()
	momentIdStr := momentId.String()
	_, err = dao.CreateNewMoment(ctx, &model.Moment{
		AuthorID: in.UserID,
		Content:  &in.Content,
		RouteID:  &lastGPXRoute.ID,
		MomentID: momentIdStr,
	})
	if err != nil {
		zlog.Error("error while create new moment", zap.Error(err))
		return errorx.NewInternalErr()
	}

	return nil
}

func (m *MomentService) Feed(ctx context.Context, in *sdto.FeedMomentInput) (*sdto.FeedMomentOutput, *errorx.ServiceErr) {
	moments, err := dao.GetMomentsByTime(ctx, MOMENT_FEED_LIMIT, time.UnixMilli(in.LatestTime))
	if err != nil {
		zlog.Error("error while get moment feed by latest time", zap.Error(err))
		return nil, errorx.NewInternalErr()
	}

	var nextTime int64
	if len(moments) == 0 {
		nextTime = time.Now().UnixMilli()
	} else {
		nextTime = moments[len(moments)-1].CreatedAt.UnixMilli()
	}

	res := &sdto.FeedMomentOutput{
		GPXRouteText: make(map[int]string),
	}
	for i, moment := range moments {
		if moment.ImageURL != nil {
			url, err := minio.GetMomentImageUrl(ctx, *moment.ImageURL)
			if err != nil {
				zlog.Error("error while get moment image url from minio", zap.Error(err))
				return nil, errorx.NewInternalErr()
			}

			moment.ImageURL = &url
		}
		if moment.VideoURL != nil {
			url, err := minio.GetMomentImageUrl(ctx, *moment.VideoURL)
			if err != nil {
				zlog.Error("error while get moment video url from minio", zap.Error(err))
				return nil, errorx.NewInternalErr()
			}

			moment.VideoURL = &url
		}
		if moment.RouteID != nil {
			path, err := dao.GetPathAsText(ctx, *moment.RouteID)
			if err != nil {
				zlog.Error("error while get GPX route from mysql", zap.Error(err))
				return nil, errorx.NewInternalErr()
			}

			pathText, err := util.GPXRoute(path)
			if err != nil {
				zlog.Error("error while parse gpx route to text", zap.String("path", path))
				return nil, errorx.NewInternalErr()
			}

			res.GPXRouteText[i] = pathText
		}
	}

	res.Moments = moments
	res.NextTime = nextTime

	return res, nil
}
