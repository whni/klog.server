package main

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"null.v3"
	"time"
)

// AddressInfo struct
type AddressInfo struct {
	Street  string `json:"street" bson:"street"`
	Code    string `json:"code" bson:"code"`
	City    string `json:"city" bson:"city"`
	State   string `json:"state" bson:"state"`
	Country string `json:"country" bson:"country"`
}

// Institute struct
type Institute struct {
	PID           primitive.ObjectID `json:"pid" bson:"_id,omitempty"`
	InstituteUID  string             `json:"institute_uid" bson:"institute_uid"`
	InstituteName string             `json:"institute_name" bson:"institute_name"`
	Address       AddressInfo        `json:"address" bson:"address"`
}

// Teacher struct
type Teacher struct {
	PID          primitive.ObjectID `json:"pid" bson:"_id,omitempty"`
	TeacherUID   string             `json:"teacher_uid" bson:"teacher_uid"`
	TeacherName  string             `json:"teacher_name" bson:"teacher_name"`
	ClassName    string             `json:"class_name" bson:"class_name"`
	PhoneNumber  string             `json:"phone_number" bson:"phone_number"`
	Email        string             `json:"email" bson:"email"`
	InstitutePID primitive.ObjectID `json:"institute_pid" bson:"institute_pid"`
}

// Student struct
type Student struct {
	PID           int       `json:"PID"`
	StudentUID    string    `json:"studentUID"`
	FirstName     string    `json:"firstName"`
	LastName      string    `json:"lastName"`
	DateOfBirth   time.Time `json:"dateOfBirth"`
	MediaLocation string    `json:"mediaLocation"`
	ClassPID      null.Int  `json:"classPID"`
	CreateTS      time.Time `json:"createTS"`
	ModifyTS      time.Time `json:"modifyTS"`
}

// Parent struct
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

// VideoClip struct
type VideoClip struct {
	PID           int       `json:"PID"`
	VideoClipUID  string    `json:"videoUID"`
	VideoClipURL  string    `json:"videoClipURL"`
	VideoMetaURL  string    `json:"videoMetaURL"`
	VideoLifeTime int       `json:"videoLifeTime"`
	CreateTS      time.Time `json:"createTS"`
	ModifyTS      time.Time `json:"modifyTS"`
}

// VideoMeta struct
type VideoMeta struct {
	PID          int    `json:"PID"`
	VideoMetaUID string `json:"videoUID"`
	VideoMetaURL string `json:"videoMetaURL"`
	VideoClipPID string `json:"videoClipURL"`
}
