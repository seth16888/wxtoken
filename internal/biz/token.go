package biz

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/seth16888/wxcommon/domain"
	wxErr "github.com/seth16888/wxcommon/error"
	"github.com/seth16888/wxcommon/hc"
	"github.com/seth16888/wxcommon/helpers"
	"github.com/seth16888/wxcommon/paths"
	"github.com/seth16888/wxtoken/internal/cache"
	"github.com/seth16888/wxtoken/internal/consts"
	"github.com/seth16888/wxtoken/internal/entities"
	"go.uber.org/zap"
)

type WXAccessTokenRes struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type AccessTokenRes struct {
	AppId       string
	MpId        string
	AccessToken string
	ExpiresIn   uint64
	Deadline    int64
}

type GetWXStableAccessTokenReq struct {
	GrantType    string `json:"grant_type"`
	Appid        string `json:"appid"`
	Secret       string `json:"secret"`
	ForceRefresh bool   `json:"force_refresh"`
}

type AccessTokenRepo interface {
	Get(ctx context.Context, mpId string) (*entities.AccessToken, error)
	Save(ctx context.Context, entity *entities.AccessToken) (uint, error)
	GetNeedRefresh(ctx context.Context, deadline int64, limit int) ([]*entities.AccessToken, error)
}

type TokenUsecase struct {
	repo    AccessTokenRepo
	log     *zap.Logger
	cache   cache.Cache
	hc      *hc.Client
	appRepo AppRepo
}

func NewTokenUsecase(repo AccessTokenRepo, logger *zap.Logger,
	cache cache.Cache, appRepo AppRepo, hc *hc.Client) *TokenUsecase {
	return &TokenUsecase{
		repo:    repo,
		log:     logger,
		cache:   cache,
		appRepo: appRepo,
		hc:      hc,
	}
}

func (p *TokenUsecase) Get(ctx context.Context, appId string, mpId string) (*AccessTokenRes, error) {
	accessTokenRes := p.loadAccessToken(ctx, mpId)

	if accessTokenRes == nil || p.isNeedRefresh(accessTokenRes) { // 缓存和数据库中都不存在
    p.log.Debug("access_token not found in cache and db, or need refresh",
    zap.String("mpId", mpId), zap.Any(consts.RequestIdKey, ctx.Value(consts.RequestIdKey)))
		rt, wxErr := p.doFetch(ctx, appId, mpId, false) // 从微信获取
		if wxErr == nil {
			return rt, nil
		}
		// 获取失败
		if wxErr.ErrCode == 40001 { // app secret error
      p.log.Debug("app secret error", zap.String("mpId", mpId))
			cache.CacheRepo.Delete("secret_" + mpId)
			rt2, wxErr2 := p.doFetch(ctx, appId, mpId, false) // 从微信获取
			if wxErr2 != nil {                                // 获取失败
				return nil, fmt.Errorf("get error: %d %s", wxErr2.ErrCode, wxErr2.ErrMsg)
			}
			return rt2, nil
		} else {
			return nil, fmt.Errorf("get error: %d %s", wxErr.ErrCode, wxErr.ErrMsg)
		}
	}

	return accessTokenRes, nil
}

func (p *TokenUsecase) Refresh(ctx context.Context, appId string, mpId string) (*AccessTokenRes, error) {
	token, wxErr := p.doFetch(ctx, appId, mpId, false)
	if wxErr == nil {
		return token, nil
	}

	if wxErr.ErrCode == 40001 { // app secret error
    p.log.Debug("app secret error", zap.String("mpId", mpId))
		cache.CacheRepo.Delete("secret_" + mpId)
		rt2, wxErr2 := p.doFetch(ctx, appId, mpId, false) // 从微信获取
		if wxErr2 != nil {                                // 获取失败
			return nil, fmt.Errorf("get error: %d %s", wxErr2.ErrCode, wxErr2.ErrMsg)
		}
		return rt2, nil
	} else {
		return nil, fmt.Errorf("refresh error: %d %s", wxErr.ErrCode, wxErr.ErrMsg)
	}
}

// ForceRefresh 强制刷新
func (p *TokenUsecase) ForceRefresh(ctx context.Context, appId string, mpId string, force bool) (*AccessTokenRes, error) {
	token, wxErr := p.doFetch(ctx, appId, mpId, force)
	if wxErr == nil {
		return token, nil
	}

	if wxErr.ErrCode == 40001 { // app secret error
    p.log.Debug("app secret error", zap.String("mpId", mpId))
		cache.CacheRepo.Delete("secret_" + mpId)
		rt2, wxErr2 := p.doFetch(ctx, appId, mpId, force) // 从微信获取
		if wxErr2 != nil {                                // 获取失败
			return nil, fmt.Errorf("get error: %d %s", wxErr2.ErrCode, wxErr2.ErrMsg)
		}
		return rt2, nil
	} else {
		return nil, fmt.Errorf("force refresh error: %d %s", wxErr.ErrCode, wxErr.ErrMsg)
	}
}

// loadAccessToken 从缓存和数据库中加载
func (p *TokenUsecase) loadAccessToken(ctx context.Context, mpId string) *AccessTokenRes {
	if token, err := p.cache.Get(mpId); err == nil {
		dataStr := fmt.Sprintf("%v", token)
		tk := &AccessTokenRes{}
		err = json.Unmarshal([]byte(dataStr), tk)
		if err != nil {
      p.log.Debug("AccessTokenRes Unmarshal error", zap.Error(err))
			return nil
		}

		if tk.Deadline <= time.Now().Unix()+10 { // 过期
			cache.CacheRepo.Delete(mpId)
			return nil
		}
		p.log.Debug("get access_token from cache", zap.String("mpId", mpId))
		return tk
	}

	row, err := p.repo.Get(ctx, mpId)
	if err == nil {
		if row.Deadline <= time.Now().Unix()+10 { // 过期
			return nil
		}
		p.log.Debug("get access_token from db", zap.String("mpId", mpId))
		// 保存到 缓存
		token := AccessTokenRes{
			AppId:       row.AppId,
			MpId:        row.MpId,
			AccessToken: row.AccessToken,
			Deadline:    row.Deadline,
			ExpiresIn:   row.ExpiresIn}

		data, err := json.Marshal(token)
		if err != nil {
			p.log.Error("token marshal error", zap.Error(err))
			return nil
		}

		p.cache.Set(mpId, string(data))

		return &token
	}
	p.log.Debug("access_token not found in cache and db", zap.Error(err))
	return nil
}

// isNeedRefresh 是否需要重新生成
func (p *TokenUsecase) isNeedRefresh(accessToken *AccessTokenRes) bool {
	startTime := time.Now().Unix() + 290
	return accessToken.Deadline <= startTime
}

// doFetch 从微信获取
func (p *TokenUsecase) doFetch(ctx context.Context, appId string,
  mpId string, force bool) (*AccessTokenRes, *wxErr.WXError) {
	secret := ""
	secretCached, err := cache.CacheRepo.Get("secret_" + mpId)
	if err != nil {
		app, err := p.appRepo.GetSecret(ctx, appId)
		if err != nil {
			p.log.Debug("get secret error",zap.Error(err), zap.String("mpId", mpId))
			return nil, &wxErr.WXError{ErrCode: 500, ErrMsg: err.Error()}
		}
		secret = app.AppSecret
		// 存入缓存
		cache.CacheRepo.Set("secret_"+mpId, secret)
	} else {
		secret = secretCached.(string)
	}

	accessToken, wxErr := p.fetchFromWX(ctx, appId, mpId, secret, force) // 从微信获取
	return accessToken, wxErr
}

// fetchFromWX 从微信获取
func (p *TokenUsecase) fetchFromWX(ctx context.Context, appId string,
  mpId string, secret string, force bool) (*AccessTokenRes, *wxErr.WXError) {
	bodyJson := GetWXStableAccessTokenReq{
		GrantType:    "client_credential",
		Appid:        mpId,
		Secret:       secret,
		ForceRefresh: force,
	}

	reader, err := helpers.BuildRequestBody[GetWXStableAccessTokenReq](bodyJson)
	if err != nil {
		p.log.Error("build request body error", zap.Error(err))
		return nil, &wxErr.WXError{ErrCode: 500, ErrMsg: err.Error()}
	}

	wxDomain := domain.GetWXAPIDomain()
	url := fmt.Sprintf("https://%s%s",
		wxDomain,
		paths.Path_Get_Stable_Access_Token,
	)
	p.log.Debug("url", zap.String("url", url))

	resp, err := p.hc.Post(url, "application/json", reader)
	if err != nil {
		p.log.Error("access WX server error", zap.Error(err))
		return nil, &wxErr.WXError{ErrCode: 500, ErrMsg: err.Error()}
	}
	result, wxErr2 := helpers.BuildHttpResponse[WXAccessTokenRes](resp, err)
	if wxErr2 != nil {
		p.log.Error("build http response error", zap.Error(wxErr2))
		return nil, &wxErr.WXError{ErrCode: wxErr2.ErrCode, ErrMsg: wxErr2.ErrMsg}
	}

	// 计算deadline
	token := &AccessTokenRes{
		AppId:       appId,
		MpId:        mpId,
		AccessToken: result.AccessToken,
		ExpiresIn:   uint64(result.ExpiresIn),
		Deadline:    time.Now().Unix() + int64(result.ExpiresIn),
	}

	// 保存到缓存
	data, err := json.Marshal(token)
	if err != nil {
		p.log.Error("token marshal error", zap.Error(err))
		return nil, &wxErr.WXError{ErrCode: 500, ErrMsg: err.Error()}
	}

	cache.CacheRepo.Set(mpId, string(data))

	// 保存到数据库
	row := &entities.AccessToken{
		AppId:       appId,
		MpId:        mpId,
		AccessToken: token.AccessToken,
		Deadline:    token.Deadline,
		ExpiresIn:   token.ExpiresIn,
	}
	if _, err := p.repo.Save(ctx, row); err != nil {
		p.log.Error("update access_token error", zap.Error(err))
		return nil, &wxErr.WXError{ErrCode: 500, ErrMsg: err.Error()}
	}

	return token, nil
}
