package types

type UserInfo struct {
	UserId    int64  `json:"user_id"`
	Nickname  string `json:"nickname"`
	Avatarurl string `json:"avatarurl"`
	Phone     string `json:"phone"`
}

//type GormConfig struct {
//	Dsn             string
//	Tracing         bool
//	MaxIdle         int
//	MaxOpen         int
//	ConnMaxIdleTime int
//}
