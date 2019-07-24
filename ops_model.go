package main

import (
	"time"
)

// Institute struct
type Institute struct {
	PID           int       `json:"PID"`
	InstituteUID  string    `json:"instituteUID"`
	InstituteName string    `json:"instituteName"`
	Address       string    `json:"address"`
	CountryCode   string    `json:"countryCode"`
	CreateTS      time.Time `json:"createTS"`
	ModifyTS      time.Time `json:"modifyTS"`
}

// Class struct
type Class struct {
	PID          int       `json:"PID"`
	ClassUID     string    `json:"classUID"`
	ClassName    string    `json:"className"`
	Location     string    `json:"location"`
	InstitutePID int       `json:"institutePID"`
	CreateTS     time.Time `json:"createTS"`
	ModifyTS     time.Time `json:"modifyTS"`
}

// Teacher struct
type Teacher struct {
	PID          int       `json:"PID"`
	TeacherUID   string    `json:"teacherUID"`
	FirstName    string    `json:"firstName"`
	LastName     string    `json:"lastName"`
	DateOfBirth  time.Time `json:"dateOfBirth"`
	Address      string    `json:"address"`
	PhoneNumber  string    `json:"phoneNumber"`
	Email        string    `json:"email"`
	InstitutePID int       `json:"institutePID"`
	CreateTS     time.Time `json:"createTS"`
	ModifyTS     time.Time `json:"modifyTS"`
}

// Student struct
type Student struct {
	PID           int       `json:"PID"`
	StudentUID    string    `json:"studentUID"`
	FirstName     string    `json:"firstName"`
	LastName      string    `json:"lastName"`
	DateOfBirth   time.Time `json:"dateOfBirth"`
	MediaLocation string    `json:"mediaLocation"`
	ClassPID      int       `json:"classPID"`
	CreateTS      time.Time `json:"createTS"`
	ModifyTS      time.Time `json:"modifyTS"`
}

// Parent struct for parent
type Parent struct {
	PID         int       `json:"PID"`
	ParentUID   string    `json:"parentUID"`
	FirstName   string    `json:"firstName"`
	LastName    string    `json:"lastName"`
	DateOfBirth time.Time `json:"dateOfBirth"`
	Address     string    `json:"address"`
	PhoneNumber string    `json:"phoneNumber"`
	Email       string    `json:"email"`
	Occupation  string    `json:"occupation"`
	CreateTS    time.Time `json:"createTS"`
	ModifyTS    time.Time `json:"modifyTS"`
}
