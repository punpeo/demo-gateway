package types

type DshApiGetChapterListByIdRespItem struct {
	Id           int64  `json:"id"`
	CampCourseId int64  `json:"camp_course_id"`
	CampId       int64  `json:"camp_id"`
	StageId      int64  `json:"stage_id"`
	TaskId       int64  `json:"task_id"`
	TaskName     string `json:"task_name"`
	StageName    string `json:"stage_name"`
}

type DshBaseApiPushMsgExtra struct {
	WorksId int64  `json:"works_id"`
	Title   string `json:"title"`
	Message string `json:"message"`
	UserId  int64  `json:"user_id"`
}

type DshApiGetAdminListRespItem struct {
	Id       int64  `json:"id"`
	RealName string `json:"real_name"`
	NickName string `json:"nick_name"`
	Avatar   string `json:"avatar"`
	Status   int64  `json:"status"`
}

type DshApiCheckWorksSubmitResp struct {
	Code int64  `json:"code"`
	Msg  string `json:"msg"`
}

type DshApiGetLantingMissionListResp struct {
	MissionList []DshApiGetLantingMissionListRespMissionListItem `json:"mission_list"`
}

type DshApiGetLantingMissionListRespMissionListItem struct {
	Word  string `json:"word"`
	Desc  string `json:"desc"`
	Image string `json:"image"`
	Id    int64  `json:"id"`
}

type DshApiGetCampUserInfoResp struct {
	ConsultantNow       int64 `json:"consultant_now"`
	CommentConsultantId int64 `json:"comment_consultant_id"`
}

type DshApiGetCampInfoResp struct {
	Id             int64  `json:"id"`
	Name           string `json:"name"`
	ShowPage       int64  `json:"show_page"`
	ShowPageBranch int64  `json:"show_page_branch"`
}

type DshApiGetCampListByIdsRespItem struct {
	Id         int64  `json:"id"`
	Name       string `json:"name"`
	ShowPage   int64  `json:"show_page"`
	ModuleName string `json:"module_name"`
}

type DshApiIsBuyReqProductItem struct {
	ProductType int64 `json:"product_type"`
	ProductId   int64 `json:"product_id"`
}

type DshApiIsBuyRespItem struct {
	ProductType int64 `json:"product_type"`
	ProductId   int64 `json:"product_id"`
	Status      int64 `json:"status"`
	PayTime     int64 `json:"pay_time"`
}

type DshApiGetCampTermUserListResp struct {
	List []struct {
		Id     int64 `json:"id"`
		UserId int64 `json:"user_id"`
	} `json:"list"`
	PageSize int64 `json:"pageSize"`
}

type DshCampUserApplyUnlockLog struct {
	ID             int64 `json:"id"`               // id自增
	UserId         int64 `json:"user_id"`          // 用户id
	CampId         int64 `json:"camp_id"`          // 课程id
	TermId         int64 `json:"term_id"`          // 期数id
	UnlockDayStart int64 `json:"unlock_day_start"` // 解锁开始章节-序号
	UnlockDayEnd   int64 `json:"unlock_day_end"`   // 解锁结束章节-序号
	UnitId         int64 `json:"unit_id"`          // 单元id
	CreateTime     int64 `json:"create_time"`      // 创建时间
}

type DshCampClassUserListReqUserListItem struct {
	CampId uint32 `json:"camp_id"`
	UserId uint32 `json:"user_id"`
}
type DshCampClassUserListRespItem struct {
	CampId        uint32 `json:"camp_id"`
	UserId        uint32 `json:"user_id"`
	CampTerm      uint32 `json:"camp_term"`
	CampTime      uint32 `json:"camp_time"`
	ServiceExpire uint32 `json:"service_expire"`
	UnlockType    uint32 `json:"unlock_type"`
	ShowPage      uint32 `json:"show_page"`
}

type DshUserCamplearnProgressListRespItem struct {
	ProductType           int    `json:"product_type"`
	ProductId             int    `json:"product_id"`
	Name                  string `json:"name"`
	Cover                 string `json:"cover"`
	ThumbCover            string `json:"thumb_cover"`
	ShowPage              int    `json:"show_page"`
	ShowPageBranch        int    `json:"show_page_branch"`
	ChapterCompletedCount int    `json:"chapter_completed_count"`
	ChapterCount          int    `json:"chapter_count"`
}
