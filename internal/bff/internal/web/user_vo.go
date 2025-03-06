package web

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

type LoginEmailPwdReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=30"`
}

type EditAccountNameReq struct {
	AccountName string `json:"accountName" binding:"required,min=1,max=30"`
}

type EditProfileReq struct {
	Nickname string `json:"nickname" binding:"required,min=1,max=30"`
	Birthday string `json:"birthday"`
	AboutMe  string `json:"aboutMe" binding:"max=1024"`
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
