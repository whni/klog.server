package main

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	TeacherKey   string             `json:"teacher_key" bson:"teacher_key"`
	TeacherName  string             `json:"teacher_name" bson:"teacher_name"`
	ClassName    string             `json:"class_name" bson:"class_name"`
	PhoneNumber  string             `json:"phone_number" bson:"phone_number"`
	Email        string             `json:"email" bson:"email"`
	InstitutePID primitive.ObjectID `json:"institute_pid" bson:"institute_pid"`
}

// Student struct
type Student struct {
	PID             primitive.ObjectID `json:"pid" bson:"_id,omitempty"`
	StudentUID      string             `json:"student_uid" bson:"student_uid"`
	StudentName     string             `json:"student_name" bson:"student_name"`
	StudentImageURL string             `json:"student_image_url" bson:"student_image_url"`
	ParentWXID      string             `json:"parent_wxid" bson:"parent_wxid"`
	ParentName      string             `json:"parent_name" bson:"parent_name"`
	PhoneNumber     string             `json:"phone_number" bson:"phone_number"`
	Email           string             `json:"email" bson:"email"`
	TeacherPID      primitive.ObjectID `json:"teacher_pid" bson:"teacher_pid"`
}

// VideoMeta struct
type VideoMeta struct {
	PID          primitive.ObjectID `json:"pid" bson:"_id,omitempty"`
	VideoMetaUID string             `json:"videometa_uid" bson:"videometa_uid"`
	StorageURL   string             `json:"storage_url" bson:"storage_url"`
	Timestamp    int64              `json:"timestamp" bson:"timestamp"`
	Lifetime     int64              `json:"lifetime" bson:"lifetime"`
	StudentPID   primitive.ObjectID `json:"student_pid" bson:"student_pid"`
}
