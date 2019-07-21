package main

import (
	"time"
)

// Student struct for student
type Student struct {
	PID           int64     `json:"DBID,string"`
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
