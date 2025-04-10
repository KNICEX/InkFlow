package web

import (
	"github.com/KNICEX/InkFlow/internal/user"
	"time"
)

// SendCodeReq 发送验证码请求
// 手机号和邮箱二选一，都没有则返回参数错误
type SendCodeReq struct {
	Phone string `json:"phone"`
	Email string `json:"email"`
}

type LoginAccountPwdReq struct {
	Account  string `json:"account" binding:"required,min=1,max=30"`
	Password string `json:"password" binding:"required,min=6,max=30"`
}

type LoginEmailReq struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required"`
}

type RegisterByEmailReq struct {
	Account  string `json:"account" binding:"required,min=1,max=30"`
	Username string `json:"username" binding:"required,min=1,max=30"`
	Email    string `json:"email" binding:"required,email"`
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required,min=6,max=30"`
}

type LoginEmailPwdReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=30"`
}

type EditProfileReq struct {
	Avatar   string   `json:"avatar"`
	Banner   string   `json:"banner"`
	Username string   `json:"username" binding:"required,min=1,max=30"`
	Birthday string   `json:"birthday"`
	Links    []string `json:"links"`
	AboutMe  string   `json:"aboutMe" binding:"max=1024"`
}

type ProfileReq struct {
	Uid     int64  `json:"uid"`
	Account string `json:"account" binding:"max=30"`
}

type SmsResetPwdReq struct {
	Code        string `json:"code" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required,min=6,max=30"`
}

type EmailResetPwdReq struct {
	Code        string `json:"code" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required,min=6,max=30"`
}

type ChangePwdReq struct {
	OldPassword string `json:"oldPassword" binding:"required,min=6,max=30"`
	NewPassword string `json:"newPassword" binding:"required,min=6,max=30"`
}

type Oauth2Callback struct {
	Code  string `json:"code" binding:"required"`
	State string `json:"state" binding:"required"`
}

type UserVO struct {
	Id int64 `json:"id"`
	//Email     string    `json:"email"`
	//Phone     string    `json:"phone"`
	Account   string    `json:"account"`
	Username  string    `json:"username"`
	Avatar    string    `json:"avatar"`
	Banner    string    `json:"banner"`
	Birthday  string    `json:"birthday"`
	AboutMe   string    `json:"aboutMe"`
	Followers int64     `json:"followers"`
	Following int64     `json:"following"`
	Followed  bool      `json:"followed"`
	Links     []string  `json:"links"`
	CreatedAt time.Time `json:"createdAt"`
}

type FollowInfo struct {
	Followers int64 `json:"followers"`
	Following int64 `json:"following"`
	Followed  bool  `json:"followed"`
}

func userToVO(u user.User) UserVO {
	birthday := ""
	if !u.Birthday.IsZero() {
		birthday = u.Birthday.Format(time.DateOnly)
	}
	return UserVO{
		Id:       u.Id,
		Account:  u.Account,
		Username: u.Username,
		Birthday: birthday,
		AboutMe:  u.AboutMe,
		Avatar:   u.Avatar,
		Banner:   u.Banner,
		//Followers: u.Followers,
		//Following: u.Following,
		//Followed:  u.Followed,
		Links:     u.Links,
		CreatedAt: u.CreatedAt,
	}
}

type FollowListReq struct {
	MaxId int64 `json:"maxId" form:"maxId"`
	Limit int   `json:"limit" form:"limit" binding:"required"`
}

type DashboardInfo struct {
	InkCount       int64 `json:"inkCount"`
	CommentCount   int64 `json:"commentCount"`
	FollowerCount  int64 `json:"followerCount"`
	FollowingCount int64 `json:"followingCount"`
	FavoriteCount  int64 `json:"favoriteCount"`
	ViewCount      int64 `json:"viewCount"`
	LikeCount      int64 `json:"likeCount"`
}
