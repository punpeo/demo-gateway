package types

import "github.com/golang-jwt/jwt/v4"

type LogingAccount struct {
	AdminId    int64
	Username   string
	RealName   string
	ExpireTime int64
}

type AccountDataAuth struct {
	Type int // 数据权限类型：0-无权限，1-有全部权限，2-有部分用户数据权限
	List []AdminData
}

type AdminData struct {
	Id    int    // 账号ID
	Phone string // 手机号码
}

type JwtClaims struct {
	jwt.RegisteredClaims
	LoginAccount    LogingAccount   `json:"login_account"`
	AdminId         int             `json:"admin_id"`
	AccountDataAuth AccountDataAuth `json:"account_data_auth"`
}
