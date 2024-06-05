package types

type CampCourseChaptersInfo struct {
	TermInfo struct {
		CampId          int64 `json:"campId"`
		CourseId        int64 `json:"courseId"`
		TermId          int64 `json:"termId"`
		CampTime        int64 `json:"campTime"`
		IsAuto          int64 `json:"isAuto"`
		LearnInBuy      int64 `json:"learnInBuy"`
		IsFreeTestStudy int64 `json:"isFreeTestStudy"`
		ShowPage        int64 `json:"showPage"`
		UnlockConfig    struct {
			ConnectUnlockCampId int64 `json:"connectUnlockCampId"`
			SkipWeek            int64 `json:"skipWeek"`
			UnlockType          int64 `json:"unlockType"`
			UnlockList          []struct {
				ChapterId  int64 `json:"chapterId"`
				StageId    int64 `json:"stageId"`
				UnlockTime int64 `json:"unlockTime"`
			} `json:"unlockList"`
			RegularDayUnlock     int64 `json:"regularDayUnlock"`
			RegularDayUnlockDays []int `json:"regularDayUnlockDays"`
			DayLimit             int64 `json:"dayLimit"`
			PassMode             int64 `json:"passMode"`
			UnitPassLimit        int64 `json:"unitPassLimit"`
			PassType             int64 `json:"passType"`
			HomeworkLimit        int64 `json:"homeworkLimit"`
			ServiceCycle         int64 `json:"serviceCycle"`
			WeekUnlockDay        int64 `json:"weekUnlockDay"`
			WeekUnlockHour       int64 `json:"weekUnlockHour"`
			IsAllPreStudy        int64 `json:"isAllPreStudy"`
			UnlockUnitNum        int64 `json:"unlockUnitNum"`
			AutoWeekUnlockNum    int64 `json:"autoWeekUnlockNum"`
			AdvanceStudyTime     struct {
				Day  int64 `json:"day"`
				Time int64 `json:"time"`
			} `json:"advanceStudyTime"`
		} `json:"unlockConfig"`
	} `json:"termInfo"`
	Chapters []struct {
		ChapterId int64 `json:"chapterId"`
		StageId   int64 `json:"stageId"`
		TaskId    int64 `json:"taskId"`
		UnitId    int64 `json:"unitId"`
		UnitSort  int64 `json:"unitSort"`
		Day       int64 `json:"day"`
		IsGuide   int64 `json:"isGuide"`
		IsFree    int64 `json:"isFree"`
	} `json:"chapters"`
}

type UserCampCompleteTime struct {
	CompleteTime int64 `json:"completeTime"`
}

type UserStudyRecord struct {
	Id         int64 `json:"id"`
	UserId     int64 `json:"user_id"`
	CampId     int64 `json:"camp_id"`
	TaskId     int64 `json:"task_id"`
	Day        int64 `json:"day"`
	Status     int64 `json:"status"`
	StudyTime  int64 `json:"study_time"`
	UpdateTime int64 `json:"update_time"`
	ChapterId  int64 `json:"chapter_id"`
}

type UnlockChapterItem struct {
	ChapterId  int64 `json:"chapter_id"`
	TaskId     int64 `json:"task_id"`
	UnlockTime int64 `json:"unlock_time"`
}
