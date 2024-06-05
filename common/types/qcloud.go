package types

type QcloudConfig struct {
	Cos Cos
}

type Cos struct {
	SecretId      string // 用户的 SecretId，
	SecretKey     string // 用户的 SecretKey，
	Region        string // 地区
	TempKeyTime   int64  // 临时密钥有效期（单位：秒）
	PublicPicture PublicPicture
}

type PublicPicture struct {
	BucketName string
	Domain     string
	AppId      string
}

type CosObjectInfo struct {
	Location string `json:"location"`
	Bucket   string `json:"bucket"`
	Key      string `json:"key"`
	ETag     string `json:"e_tag"`
}
