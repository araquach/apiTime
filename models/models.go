package models

import (
	"gorm.io/datatypes"
	"time"
)

type Time struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	StaffId     int            `json:"staff_id"`
	HolidayEnt  float32        `json:"holiday_ent"`
	SaturdayEnt float32        `json:"saturday_ent"`
	FreeTimeEnt float32        `json:"free_time_ent"`
	Schedule    datatypes.JSON `json:"schedule"`
}

type TimeDash struct {
	Holidays        float32 `json:"holidays"`
	HolidaysPending float32 `json:"holidays_pending"`
	Lieu            float32 `json:"lieu"`
	LieuPending     float32 `json:"lieu_pending"`
	FreeTime        float32 `json:"free_time"`
	FreeTimePending float32 `json:"free_time_pending"`
	Sick            float32 `json:"sick"`
	SickPending     float32 `json:"sick_pending"`
}

type SHCats struct {
	ID      uint   `json:"id" gorm:"primaryKey"`
	CatName string `json:"cat_name"`
	Tag     string `json:"tag"`
	Info    string `json:"info"`
	Filter  string `json:"filter"`
	Auth    uint   `json:"auth"`
}

type Holiday struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at" sql:"index"`
	StaffId     int        `json:"staff_id"`
	Requested   float32    `json:"requested"`
	Description string     `json:"description"`
	DateFrom    time.Time  `json:"date_from"`
	DateTo      time.Time  `json:"date_to"`
	Saturday    float32    `json:"saturday"`
	Approved    uint       `json:"approved"`
}

// HolidayDashboard Not in DB
type HolidayDash struct {
	Entitlement      float32 `json:"entitlement"`
	TotalBooked      float32 `json:"total_booked"`
	TotalPending     float32 `json:"total_pending"`
	Remaining        float32 `json:"remaining"`
	RemainingPending float32 `json:"remaining_pending"`
	SatRemaining     float32 `json:"sat_remaining"`
	SatPending       float32 `json:"sat_pending"`
}

type Sick struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at" sql:"index"`
	StaffId     int        `json:"staff_id"`
	Hours       float32    `json:"hours"`
	Description string     `json:"description"`
	DateFrom    time.Time  `json:"date_from"`
	DateTo      time.Time  `json:"date_to"`
	Deducted    uint       `json:"deducted"`
}

type SickDash struct {
	SickDays float32 `json:"sick_days"`
	Pending  float32 `json:"pending"`
}

type Lieu struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at" sql:"index"`
	StaffId     int        `json:"staff_id"`
	Hours       float32    `json:"hours"`
	Description string     `json:"description"`
	RequestDate time.Time  `json:"request_date"`
	Approved    uint       `json:"approved"`
}

type LieuDash struct {
	Used    float32 `json:"used"`
	Pending float32 `json:"pending"`
}

type FreeTime struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at" sql:"index"`
	StaffId     int        `json:"staff_id"`
	Hours       float32    `json:"hours"`
	Description string     `json:"description"`
	RequestDate time.Time  `json:"request_date"`
	Approved    uint       `json:"approved"`
}

type FreeTimeDash struct {
	Entitlement      float32 `json:"entitlement"`
	Used             float32 `json:"used"`
	Pending          float32 `json:"pending"`
	Remaining        float32 `json:"remaining"`
	RemainingPending float32 `json:"remaining_pending"`
}
