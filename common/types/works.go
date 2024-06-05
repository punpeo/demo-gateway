package types

type WorksListWhere struct {
	IsUpWall      int64   `json:"is_up_wall"`
	Page          int64   `json:"page"`
	PageSize      int64   `json:"page_size"`
	CampIds       []int64 `json:"camp_ids"`
	UserId        int64   `json:"user_id"`
	WorksIds      []int64 `json:"works_ids"`
	CampCourseId  int64   `json:"camp_course_id"`
	TermId        int64   `json:"term_id"`
	ChapterIds    []int64 `json:"chapter_ids"`
	CommentStatus int64   `json:"comment_status"`
}

type WorksList struct {
	Id        int64 `json:"id"`
	ChapterId int64 `json:"chapter_id"`
}
