package service

import (
	"context"

	v1 "github.com/seth16888/wxtoken/api/v1"
	"github.com/seth16888/wxtoken/internal/biz"
	"go.uber.org/zap"
	"google.golang.org/grpc/status"
)

type WXTokenService struct {
  v1.UnimplementedTokenServer

  log *zap.Logger
  uc *biz.TokenUsecase
}

// NewWXTokenService 创建一个WXTokenService实例
func NewWXTokenService(log *zap.Logger,uc *biz.TokenUsecase) *WXTokenService {
  return &WXTokenService{log: log,uc: uc}
}

func (s *WXTokenService) GetAccessToken(ctx context.Context, in *v1.GetTokenRequest) (*v1.GetTokenReply, error) {
	s.log.Debug("GetAccessToken", zap.String("mpId", in.GetMpId()), zap.String("appId", in.GetAppId()))
	token, err := s.uc.Get(ctx, in.GetAppId(), in.GetMpId())
	if err != nil {
		return nil, status.Error(10400, err.Error())
	}
	return &v1.GetTokenReply{
		AccessToken: token.AccessToken,
		ExpiresIn:   token.ExpiresIn,
	}, nil
}

func (s *WXTokenService) RefreshAccessToken(ctx context.Context, in *v1.RefreshTokenRequest) (*v1.RefreshTokenReply, error) {
	s.log.Debug("RefreshAccessToken", zap.String("mpId", in.GetMpId()), zap.String("appId", in.GetAppId()))
	token, err := s.uc.Refresh(ctx, in.GetAppId(), in.GetMpId())
	if err != nil {
		return nil, status.Error(10401, err.Error())
	}
	return &v1.RefreshTokenReply{
		AppId:       token.AppId,
		MpId:        token.MpId,
		Deadline:    token.Deadline,
		AccessToken: token.AccessToken,
		ExpiresIn:   token.ExpiresIn,
	}, nil
}

func (s *WXTokenService) ForceRefreshAccessToken(ctx context.Context, in *v1.ForceRefreshTokenRequest) (*v1.ForceRefreshTokenReply, error) {
	s.log.Debug("ForceRefreshAccessToken", zap.String("mpId", in.GetMpId()), zap.String("appId", in.GetAppId()))
	token, err := s.uc.ForceRefresh(ctx, in.GetAppId(), in.GetMpId(), in.GetForceRefresh())
	if err != nil {
		return nil, status.Error(10402, err.Error())
	}
	return &v1.ForceRefreshTokenReply{
		AppId:       token.AppId,
		MpId:        token.MpId,
		Deadline:    token.Deadline,
		AccessToken: token.AccessToken,
		ExpiresIn:   token.ExpiresIn,
	}, nil
}
