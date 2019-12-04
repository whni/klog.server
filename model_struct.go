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
	BindingCode      string             `json:"binding_code" bson:"binding_code"`
	BindingExpire    int64              `json:"binding_expire" bson:"binding_expire"`
}

// Relative struct
type Relative struct {
	PID          primitive.ObjectID `json:"pid" bson:"_id,omitempty"`
	RelativeName string             `json:"relative_name" bson:"relative_name"`
	RelativeWXID string             `json:"relative_wxid" bson:"relative_wxid"`
	PhoneNumber  string             `json:"phone_number" bson:"phone_number"`
	Email        string             `json:"email" bson:"email"`
}

// CourseRecord struct
type CourseRecord struct {
	PID        primitive.ObjectID `json:"pid" bson:"_id,omitempty"`
	StudentPID primitive.ObjectID `json:"student_pid" bson:"student_pid"`
	CoursePID  primitive.ObjectID `json:"course_pid" bson:"course_pid"`
	TargetTag  string             `json:"target_tag" bson:"target_tag"`
	RecordTS   int64              `json:"record_ts" bson:"record_ts"`
	IsMakeUp   bool               `json:"is_makeup" bson:"is_makeup"`
}

// CourseComment struct
type CourseComment struct {
	PID               primitive.ObjectID `json:"pid" bson:"_id,omitempty"`
	CourseRecordPID   primitive.ObjectID `json:"course_record_pid" bson:"course_record_pid"`
	CommentPersonPID  primitive.ObjectID `json:"comment_person_pid" bson:"comment_person_pid"`
	CommentPersonType string             `json:"comment_person_type" bson:"comment_person_type"`
	CommentTS         int64              `json:"comment_ts" bson:"comment_ts"`
	CommentBody       string             `json:"comment_body" bson:"comment_body"`
}

// CloudMedia struct
type CloudMedia struct {
	PID             primitive.ObjectID `json:"pid" bson:"_id,omitempty"`
	StudentPID      primitive.ObjectID `json:"student_pid" bson:"student_pid"`
	CourseRecordPID primitive.ObjectID `json:"course_record_pid" bson:"course_record_pid"`
	MediaType       string             `json:"media_type" bson:"media_type"`
	MediaName       string             `json:"media_name" bson:"media_name"`
	MediaURL        string             `json:"media_url" bson:"media_url"`
	RankScore       float64            `json:"rank_score" bson:"rank_score"`
	MediaTags       []string           `json:"media_tags" bson:"media_tags"`
	CreateTS        int64              `json:"create_ts" bson:"create_ts"`
	ContentLength   int64              `json:"content_length" bson:"content_length"`
}

// StudentRelativeRef struct
type StudentRelativeRef struct {
	PID          primitive.ObjectID `json:"pid" bson:"_id,omitempty"`
	StudentPID   primitive.ObjectID `json:"student_pid" bson:"student_pid"`
	RelativePID  primitive.ObjectID `json:"relative_pid" bson:"relative_pid"`
	Relationship string             `json:"relationship" bson:"relationship"`
	IsMain       bool               `json:"is_main" bson:"is_main"`
}

// StudentRelativeBindInfo struct
type StudentRelativeBindInfo struct {
	StudentPID   primitive.ObjectID `json:"student_pid" bson:"student_pid"`
	RelativeWXID string             `json:"relative_wxid" bson:"relative_wxid"`
	BindingCode  string             `json:"binding_code" bson:"binding_code"`
	Relationship string             `json:"relationship" bson:"relationship"`
}

// StudentCourseRef struct
type StudentCourseRef struct {
	PID        primitive.ObjectID `json:"pid" bson:"_id,omitempty"`
	StudentPID primitive.ObjectID `json:"student_pid" bson:"student_pid"`
	CoursePID  primitive.ObjectID `json:"course_pid" bson:"course_pid"`
}

// AzureBlobProp struct
type AzureBlobProp struct {
	BlobName      string
	CreateTS      int64
	ContentLength int64
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
	CommentPersonTypeTeacher  = "teacher"
	CommentPersonTypeRelative = "relative"
)

var CommentPersonTypeMap = map[string]bool{
	CommentPersonTypeTeacher:  true,
	CommentPersonTypeRelative: true,
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
