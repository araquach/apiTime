package models

import (
	"time"
)

type SHCats struct {
	ID      uint   `json:"id" gorm:"primaryKey"`
	CatName string `json:"cat_name"`
	Tag     string `json:"tag"`
	Info    string `json:"info"`
	Filter  string `json:"filter"`
	Auth    uint   `json:"auth"`
}

type Holiday struct {
	ID              uint       `json:"id" gorm:"primaryKey"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at" sql:"index"`
	StaffId         int        `json:"staff_id"`
	HoursRequested  float32    `json:"hours_requested"`
	Description     string     `json:"description"`
	RequestDateFrom time.Time  `json:"request_date_from"`
	RequestDateTo   time.Time  `json:"request_date_to"`
	Saturday        float32    `json:"saturday"`
	Approved        uint       `json:"approved"`
}

type Sick struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at" sql:"index"`
	StaffId     int        `json:"staff_id"`
	SickHours   float32    `json:"sick_hours"`
	Description string     `json:"description"`
	SickFrom    time.Time  `json:"sick_from"`
	SickTo      time.Time  `json:"sick_to"`
	Deducted    bool       `json:"deducted"`
}

type Lieu struct {
	ID            uint       `json:"id" gorm:"primaryKey"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at" sql:"index"`
	StaffId       int        `json:"staff_id"`
	LieuHours     float32    `json:"lieu_hours"`
	Description   string     `json:"description"`
	DateRegarding time.Time  `json:"date_regarding"`
	Approved      uint       `json:"approved"`
}

type FreeTime struct {
	ID            uint       `json:"id" gorm:"primaryKey"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at" sql:"index"`
	StaffId       int        `json:"staff_id"`
	FreeTimeHours float32    `json:"free_time_hours"`
	Description   string     `json:"description"`
	DateRegarding time.Time  `json:"date_regarding"`
	Approved      uint       `json:"approved"`
}
