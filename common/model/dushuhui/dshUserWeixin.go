package model

import (
	"context"
	"errors"
	"gorm.io/gorm"
)

// 微信用户表
type DshUserWeixin struct {
	Id           int64  `gorm:"column:id;type:int(10);primary_key;AUTO_INCREMENT;comment:自增id" json:"id"`
	UserId       int64  `gorm:"column:user_id;type:int(10);default:0;comment:用户id;NOT NULL" json:"user_id"`
	Unionid      string `gorm:"column:unionid;type:varchar(50);comment:微信unionid;NOT NULL" json:"unionid"`
	Openid       string `gorm:"column:openid;type:varchar(50);comment:小程序openid;NOT NULL" json:"openid"`
	OpenidPublic string `gorm:"column:openid_public;type:varchar(50);comment:公众号openid;NOT NULL" json:"openid_public"`
	Subscribe    int64  `gorm:"column:subscribe;type:tinyint(1);default:0;comment:是否关注读书会公众号：0-否，1-是;NOT NULL" json:"subscribe"`
	Nickname     string `gorm:"column:nickname;type:varchar(30);comment:昵称;NOT NULL" json:"nickname"`
	Avatarurl    string `gorm:"column:avatarurl;type:varchar(255);comment:头像;NOT NULL" json:"avatarurl"`
	Gender       int64  `gorm:"column:gender;type:tinyint(1);default:0;comment:性别：1-男，2-女，0-未知;NOT NULL" json:"gender"`
	City         string `gorm:"column:city;type:varchar(20);comment:城市;NOT NULL" json:"city"`
	Province     string `gorm:"column:province;type:varchar(20);comment:省份;NOT NULL" json:"province"`
	Country      string `gorm:"column:country;type:varchar(200);comment:国家;NOT NULL" json:"country"`
	CreateTime   int64  `gorm:"column:create_time;type:int(10);default:0;comment:第一次授权时间;NOT NULL" json:"create_time"`
	UpdateTime   int64  `gorm:"column:update_time;type:int(10);default:0;comment:更新时间;NOT NULL" json:"update_time"`
}

func (m *DshUserWeixin) TableName() string {
	return "dsh_user_weixin"
}

// 获取单条数据
func (m *DshUserWeixin) GetInfo(db *gorm.DB, ctx context.Context, id int64) (resp *DshUserWeixin, err error) {
	resp = &DshUserWeixin{}
	err = db.Table(m.TableName()).WithContext(ctx).Where("id =? ", id).First(&resp).Error
	//查空
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	//err = errors.New("111")
	return
}
