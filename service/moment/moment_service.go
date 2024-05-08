package moment

import (
	"context"
	"errors"
	"time"

	"api.backend.xjco2913/dao"
	"api.backend.xjco2913/dao/minio"
	"api.backend.xjco2913/dao/model"
	"api.backend.xjco2913/service/gpx"
	"api.backend.xjco2913/service/sdto"
	"api.backend.xjco2913/service/sdto/errorx"
	"api.backend.xjco2913/util"
	"api.backend.xjco2913/util/zlog"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
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
		zlog.Error("Error while create new moment", zap.Error(err))
		return errorx.NewInternalErr()
	}

	return nil
}

func (m *MomentService) CreateWithImage(ctx context.Context, in *sdto.CreateMomentImageInput) *errorx.ServiceErr {
	imageName := uuid.New()
	err := minio.UploadMomentImage(ctx, imageName.String(), in.ImageData)
	if err != nil {
		zlog.Error("Error while store moment image into minio", zap.Error(err))
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
		zlog.Error("Error while create new moment", zap.Error(err))

		// Async remove image in minio
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
		zlog.Error("Error while create new moment", zap.Error(err))
		return errorx.NewInternalErr()
	}

	// Async upload video
	go func() {
		err := minio.UploadMomentVideo(ctx, videoNameStr, in.VideoData)
		if err != nil {
			// Error, remove moment record in mysql
			zlog.Error("Error while async upload moment video", zap.Error(err))
			dao.DeleteMomentByID(ctx, newMomentID)
		}
	}()

	return nil
}

func (m *MomentService) CreateWithGPX(ctx context.Context, in *sdto.CreateMomentGPXInput) *errorx.ServiceErr {
	gpxResp, sErr := gpx.Service().ParseGPXData(ctx, &sdto.ParseGPXDataInput{
		GPXData: in.GPXData,
	})
	if sErr != nil {
		return sErr
	}

	momentId := uuid.New()
	momentIdStr := momentId.String()
	_, err := dao.CreateNewMoment(ctx, &model.Moment{
		AuthorID: in.UserID,
		Content:  &in.Content,
		RouteID:  &gpxResp.RouteID,
		MomentID: momentIdStr,
	})
	if err != nil {
		zlog.Error("Error while create new moment", zap.Error(err))
		return errorx.NewInternalErr()
	}

	return nil
}

func (m *MomentService) Feed(ctx context.Context, in *sdto.FeedMomentInput) (*sdto.FeedMomentOutput, *errorx.ServiceErr) {
	moments, err := dao.GetMomentsByTime(ctx, MOMENT_FEED_LIMIT, time.UnixMilli(in.LatestTime))
	if err != nil {
		zlog.Error("Error while get moment feed by latest time", zap.Error(err))
		return nil, errorx.NewInternalErr()
	}

	var nextTime int64
	if len(moments) == 0 {
		nextTime = time.Now().UnixMilli()
	} else {
		nextTime = moments[len(moments)-1].CreatedAt.UnixMilli()
	}

	res := &sdto.FeedMomentOutput{
		GPXRouteText:  make(map[int][][]string),
		AuthorInfoMap: make(map[string]*model.User),
	}
	for i, moment := range moments {
		// Get author info
		author, err := dao.GetUserByID(ctx, moment.AuthorID)
		if err != nil {
			zlog.Error("Error while get moment author info", zap.Error(err), zap.String("momentID", moment.MomentID))
			return nil, errorx.NewInternalErr()
		}

		// Get author avatar url
		if author.AvatarURL != nil {
			url, err := minio.GetUserAvatarUrl(ctx, *author.AvatarURL)
			if err != nil {
				zlog.Error("Error while get author avatar url", zap.Error(err))
				return nil, errorx.NewInternalErr()
			}

			author.AvatarURL = &url
		}

		res.AuthorInfoMap[moment.MomentID] = author

		if moment.ImageURL != nil {
			url, err := minio.GetMomentImageUrl(ctx, *moment.ImageURL)
			if err != nil {
				zlog.Error("Error while get moment image url from minio", zap.Error(err))
				return nil, errorx.NewInternalErr()
			}

			moment.ImageURL = &url
		}
		if moment.VideoURL != nil {
			url, err := minio.GetMomentImageUrl(ctx, *moment.VideoURL)
			if err != nil {
				zlog.Error("Error while get moment video url from minio", zap.Error(err))
				return nil, errorx.NewInternalErr()
			}

			moment.VideoURL = &url
		}
		if moment.RouteID != nil {
			path, err := dao.GetPathAsText(ctx, *moment.RouteID)
			if err != nil {
				zlog.Error("Error while get GPX route from mysql", zap.Error(err))
				return nil, errorx.NewInternalErr()
			}

			pathText, err := util.GPXRoute(path)
			if err != nil {
				zlog.Error("Error while parse gpx route to text", zap.String("path", path))
				return nil, errorx.NewInternalErr()
			}

			res.GPXRouteText[i] = util.GPXStrTo2DString(pathText)
		}
	}

	res.Moments = moments
	res.NextTime = nextTime

	return res, nil
}

func (m *MomentService) GetLikesByMomentId(ctx context.Context, momentId string) (*sdto.GetLikesOutput, *errorx.ServiceErr) {
	likes, err := dao.GetLikeByMomentId(ctx, momentId)
	if err != nil {
		zlog.Error("Error while get moment likes", zap.Error(err))
		return nil, errorx.NewInternalErr()
	}

	personLikes := []sdto.MomentUser{}
	for _, like := range likes {
		likeId := like.UserID

		personLike, err := dao.GetUserByID(ctx, likeId)
		if err != nil {
			zlog.Error("Error while get user liked by id", zap.Error(err))
			return nil, errorx.NewInternalErr()
		}

		var url string
		if personLike.AvatarURL != nil && *personLike.AvatarURL != "" {
			url, err = minio.GetUserAvatarUrl(ctx, *personLike.AvatarURL)
			if err != nil {
				zlog.Error("Error while get user avatar url", zap.Error(err))
				return nil, errorx.NewInternalErr()
			}
		}

		personLikes = append(personLikes, sdto.MomentUser{
			Name:      personLike.Username,
			AvatarUrl: url,
		})
	}

	return &sdto.GetLikesOutput{
		PersonLikes: personLikes,
	}, nil
}

func (m *MomentService) GetCommentListByMomentId(ctx context.Context, momentId string) (*sdto.GetCommentListOutput, *errorx.ServiceErr) {
	commentlist := []sdto.MomentComment{}

	commentModels, err := dao.GetCommentsByMomentId(ctx, momentId)
	if err != nil {
		zlog.Error("Error while get moment comments", zap.Error(err))
		return nil, errorx.NewInternalErr()
	}

	for _, commentModel := range commentModels {
		var comment sdto.MomentComment

		author, err := dao.GetUserByID(ctx, commentModel.AuthorID)
		if err != nil {
			zlog.Error("Error while get comment user", zap.Error(err))
			return nil, errorx.NewInternalErr()
		}

		var url string
		if author.AvatarURL != nil && *author.AvatarURL != "" {
			url, err = minio.GetUserAvatarUrl(ctx, *author.AvatarURL)
			if err != nil {
				zlog.Error("Error while get user avatar url", zap.Error(err))
				return nil, errorx.NewInternalErr()
			}
		}

		comment.Id = commentModel.CommentID
		comment.CreatedAt = *commentModel.CreatedAt
		comment.Message = commentModel.Content
		comment.Author.Name = author.Username
		comment.Author.AvatarUrl = url

		commentlist = append(commentlist, comment)
	}

	return &sdto.GetCommentListOutput{
		CommentList: commentlist,
	}, nil
}

func (m *MomentService) IsLiked(ctx context.Context, momentId, userId string) (bool, *errorx.ServiceErr) {
	likeModel, err := dao.GetLikeByIDs(ctx, userId, momentId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}

		zlog.Error("Error while get like by two ids", zap.Error(err))
		return false, errorx.NewInternalErr()
	}

	return likeModel != nil, nil
}
