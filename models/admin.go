package models

type TimeAdminDash struct {
	Holidays  int `json:"holidays"`
	LieuHours int `json:"lieu_hours"`
	FreeTime  int `json:"free_time"`
	SickDays  int `json:"sick_days"`
}

type HolidayAdminDash struct {
	Pending int `json:"pending"`
}

type LieuAdminDash struct {
	Pending int `json:"pending"`
}

type FreeTimeAdminDash struct {
	Pending int `json:"pending"`
}

type SickAdminDash struct {
	Pending int `json:"pending"`
}
