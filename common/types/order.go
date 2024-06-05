package types

type GetUserLantingOrder struct {
	ProductType int64 `json:"product_type"`
	ProductId   int64 `json:"product_id"`
	ShowPage    int64 `json:"show_page"`
	PayTime     int64 `json:"pay_time"`
	CampTime    int64 `json:"camp_time"`
}
