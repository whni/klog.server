package main

import (
	"time"
)

// Institute struct
type Institute struct {
	PID           int    `json:"PID"`
	InstituteUID  string `json:"instituteUID"`
	InstituteName string `json:"instituteName"`
	Address       string `json:"address"`
	CountryCode   string `json:"countryCode"`
	CreateTS      string `json:"createTS"`
	ModifyTS      string `json:"modifyTS"`
}

// Class struct
type Class struct {
	PID          int    `json:"PID"`
	ClassUID     string `json:"classUID"`
	ClassName    string `json:"className"`
	Location     string `json:"location"`
	InstitutePID int    `json:"institutePID"`
	CreateTS     string `json:"createTS"`
	ModifyTS     string `json:"modifyTS"`
}

struct Parent struct {
	
}

// Student struct
type Student struct {
	pid           int64     `json:"DBID,string"`
	StudentID     string    `json:"studentID"`
	FirstName     string    `json:"firstName"`
	LastName      string    `json:"lastName"`
	DateOfBirth   time.Time `json:"dateOfBirth"`
	ClassID       string    `json:"classID"`
	ParentIDs     []string  `json:"parentIDs"`
	MediaLocation string    `json:"mediaLocation"`
}

// Parent struct for parent
type Parent struct {
	DBID        int64     `json:"DBID,string"`
	ParentID    string    `json:"parentID"`
	FirstName   string    `json:"firstName"`
	LastName    string    `json:"lastName"`
	DateOfBirth time.Time `json:"dateOfBirth"`
	StudentIDs  []string  `json:"studentIDs"`
	Address     string    `json:"address"`
	PhoneNumber string    `json:"phoneNumber"`
	Email       string    `json:"email"`
}
