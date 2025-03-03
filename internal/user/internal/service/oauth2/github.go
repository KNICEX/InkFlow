package oauth2

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/KNICEX/InkFlow/internal/user/internal/domain"
	"github.com/KNICEX/InkFlow/pkg/logx"
	"net/http"
	"time"
)

type GithubService struct {
	clientId       string
	clientSecret   string
	redirectDomain string
	redirectURI    string
	client         *http.Client
	logger         logx.Logger
}

func NewGithubService(clientId, clientSecret, redirectDomain string, logger logx.Logger) Service[domain.GithubInfo] {
	return &GithubService{
		clientId:       clientId,
		clientSecret:   clientSecret,
		redirectDomain: redirectDomain,
		logger:         logger,
	}
}
func (s *GithubService) AuthURL(ctx context.Context, state string) (string, error) {
	const urlPattern = "https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&state=%s"
	return fmt.Sprintf(urlPattern, s.clientId, s.redirectURI, state), nil
}

func (s *GithubService) VerifyCode(ctx context.Context, code string) (domain.GithubInfo, error) {
	const targetPattern = "https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&code=%s&redirect_uri=%s"
	targetURL := fmt.Sprintf(targetPattern, s.clientId, s.clientSecret, code, s.redirectURI)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, targetURL, nil)
	if err != nil {
		return domain.GithubInfo{}, err
	}
	req.Header.Set("Accept", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return domain.GithubInfo{}, err
	}
	defer resp.Body.Close()
	var res Result
	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return domain.GithubInfo{}, err
	}
	if res.AccessToken == "" {
		// code过期或者无效
		s.logger.WithCtx(ctx).Warn("github oauth2 service 获取accessToken失败", logx.Error(ErrCodeInvalid))
		return domain.GithubInfo{}, ErrCodeInvalid
	}
	user, err := s.getGithubUserInfo(ctx, res.AccessToken)
	if err != nil {
		return domain.GithubInfo{}, err
	}
	return domain.GithubInfo{
		AccessToken: res.AccessToken,
		Username:    user.Username,
		Id:          user.Id,
		AvatarUrl:   user.AvatarUrl,
	}, nil
}

// getGithubUserInfo
// 通过accessToken 获取github用户信息
func (s *GithubService) getGithubUserInfo(ctx context.Context, accessToken string) (User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.github.com/user", nil)
	if err != nil {
		return User{}, err
	}
	req.Header.Set("Authorization", "bearer "+accessToken)
	resp, err := s.client.Do(req)
	if err != nil {
		return User{}, err
	}
	defer resp.Body.Close()
	var user User
	if err = json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return User{}, err
	}
	return user, nil
}

type User struct {
	// github用户名
	Username  string `json:"login"`
	Id        int64  `json:"id"`
	AvatarUrl string `json:"avatar_url"`
}

type Result struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}
