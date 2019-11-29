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
	PhoneNumber  string             `json:"phone_number" bson:"phone_number"`
	Email        string             `json:"email" bson:"email"`
	InstitutePID primitive.ObjectID `json:"institute_pid" bson:"institute_pid"`
}

// CourseTarget struct
type CourseTarget struct {
	Tag  string `json:"tag" bson:"tag"`
	Desc string `json:"desc" bson:"desc"`
}

// Course struct
type Course struct {
	PID           primitive.ObjectID `json:"pid" bson:"_id,omitempty"`
	CourseUID     string             `json:"course_uid" bson:"course_uid"`
	CourseName    string             `json:"course_name" bson:"course_name"`
	CourseIntro   string             `json:"course_intro" bson:"course_intro"`
	CourseTargets []CourseTarget     `json:"course_targets" bson:"course_targets"`
	TeacherPID    primitive.ObjectID `json:"teacher_pid" bson:"teacher_pid"`
	AssistantPID  primitive.ObjectID `json:"assistant_pid" bson:"assistant_pid"`
	InstitutePID  primitive.ObjectID `json:"institute_pid" bson:"institute_pid"`
}

// Student struct
type Student struct {
	PID              primitive.ObjectID `json:"pid" bson:"_id,omitempty"`
	StudentName      string             `json:"student_name" bson:"student_name"`
	StudentImageName string             `json:"student_image_name" bson:"student_image_name"`
	StudentImageURL  string             `json:"student_image_url" bson:"student_image_url"`
	ParentWXID       string             `json:"parent_wxid" bson:"parent_wxid"`
	ParentName       string             `json:"parent_name" bson:"parent_name"`
	PhoneNumber      string             `json:"phone_number" bson:"phone_number"`
	Email            string             `json:"email" bson:"email"`
	BindingCode      string             `json:"binding_code" bson:"binding_code"`
	BindingExpire    int64              `json:"binding_expire" bson:"binding_expire"`
	TeacherPID       primitive.ObjectID `json:"teacher_pid" bson:"teacher_pid"`
}

// StudentMediaQueryReq struct
type StudentMediaQueryReq struct {
	StudentPID primitive.ObjectID `json:"student_pid" bson:"student_pid"`
	StartTS    int64              `json:"start_ts" bson:"start_ts"`
	EndTS      int64              `json:"end_ts" bson:"end_ts"`
}

// ParentWeChatLoginInfo struct
type ParentWeChatLoginInfo struct {
	AppID  string `json:"appid"`
	Secret string `json:"secret"`
	JSCode string `json:"js_code"`
}

const (
	CloudMediaTypeVideo  = "video"
	CloudMediaTypeImage  = "image"
	CloudMediaTypeOthers = "others"
)

var cloudMediaTypeMap = map[string]bool{
	CloudMediaTypeVideo:  true,
	CloudMediaTypeImage:  true,
	CloudMediaTypeOthers: true,
}

// CloudMedia struct
type CloudMedia struct {
	PID           primitive.ObjectID `json:"pid" bson:"_id,omitempty"`
	MediaType     string             `json:"media_type" bson:"media_type"`
	MediaName     string             `json:"media_name" bson:"media_name"`
	MediaURL      string             `json:"media_url" bson:"media_url"`
	RankScore     float64            `json:"rank_score" bson:"rank_score"`
	CreateTS      int64              `json:"create_ts" bson:"create_ts"`
	ContentLength int64              `json:"content_length" bson:"content_length"`
	StudentPID    primitive.ObjectID `json:"student_pid" bson:"student_pid"`
}

// AzureBlobProp struct
type AzureBlobProp struct {
	BlobName      string
	CreateTS      int64
	ContentLength int64
}
