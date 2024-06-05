package types

type DshRpcUserListRespItem struct {
	UserId    int64  `json:"user_id"`
	IsVip     int64  `json:"is_vip"`
	Nickname  string `json:"nickname"`
	Avatarurl string `json:"avatarurl"`
}
