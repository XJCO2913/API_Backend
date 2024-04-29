package comment

import (
	"context"

	"api.backend.xjco2913/dao"
	"api.backend.xjco2913/dao/model"
	"api.backend.xjco2913/service/sdto"
	"api.backend.xjco2913/service/sdto/errorx"
	"api.backend.xjco2913/util/zlog"
	"go.uber.org/zap"
)

type CommentService struct{}

var (
	commentService CommentService
)

func Service() *CommentService {
	return &commentService
}

func (s *CommentService) Create(ctx context.Context, input *sdto.CreateCommentInput) *errorx.ServiceErr {
	err := dao.CreateNewComment(ctx, &model.Comment{
		AuthorID: input.AuthorID,
		MomentID: input.MomentID,
		Content:  input.Content,
	})
	if err != nil {
		zlog.Error("Error while create new comment", zap.String("authorID", input.AuthorID), zap.String("momentID", input.MomentID), zap.Error(err))
		return errorx.NewInternalErr()
	}

	return nil
}
