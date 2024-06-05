package types

type StudyStatusUpdateEvent struct {
	CampCourseId int64 `json:"camp_course_id"`
	CampId       int64 `json:"camp_id"`
	CampTerm     int64 `json:"camp_term"`
	ChapterId    int64 `json:"chapter_id"`
	Day          int64 `json:"day"`
	Status       int64 `json:"status"`
	StudyTime    int64 `json:"study_time"`
	TaskId       int64 `json:"task_id"`
	UserId       int64 `json:"user_id"`
	Version      int64 `json:"version"`
	UnitId       int64 `json:"unit_id"`
	IsRebuild    int64 `json:"is_rebuild"`
	SubMessage   []struct {
		CampCourseId int64 `json:"camp_course_id"`
		CampId       int64 `json:"camp_id"`
		CampTerm     int64 `json:"camp_term"`
		ChapterId    int64 `json:"chapter_id"`
		Day          int64 `json:"day"`
		Status       int64 `json:"status"`
		StudyTime    int64 `json:"study_time"`
		TaskId       int64 `json:"task_id"`
		UserId       int64 `json:"user_id"`
		Version      int64 `json:"version"`
		UnitId       int64 `json:"unit_id"`
		IsRebuild    int64 `json:"is_rebuild"`
	} `json:"sub_message"`
}

type DshActionProductPayEvent struct {
	UserId      int64  `json:"user_id"`
	OrderNo     string `json:"order_no"`
	ProductType int64  `json:"product_type"`
	ProductId   int64  `json:"product_id"`
}

type DshCampTermChangeUserEvent struct {
	UserId   int64 `json:"user_id"`
	CampId   int64 `json:"camp_id"`
	CampTerm int64 `json:"term"`
}
