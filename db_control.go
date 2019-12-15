package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	DBCollectionInstitute          = "institute"
	DBCollectionTeacher            = "teacher"
	DBCollectionCourse             = "course"
	DBCollectionCourseRecord       = "course_record"
	DBCollectionCourseComment      = "course_comment"
	DBCollectionRelative           = "relative"
	DBCollectionStudent            = "student"
	DBCollectionCloudMedia         = "cloudmedia"
	DBCollectionTemplate           = "template"
	DBCollectionStudentRelativeRef = "student_relative_ref"
	DBCollectionStudentCourseRef   = "student_course_ref"
)

var dbPool *mongo.Database

func dbPoolInit(sc *ServerConfig) (*mongo.Database, error) {
	dbURL := fmt.Sprintf("mongodb://%s:%s@%s/%s", sc.DBUsername, sc.DBPassword, sc.DBHostAddress, sc.DBName)
	dbClientOptions := options.Client().ApplyURI(dbURL)
	dbClientOptions.SetConnectTimeout(5 * time.Second)

	// Connect to MongoDB
	dbClient, connErr := mongo.Connect(context.TODO(), dbClientOptions)
	if connErr != nil {
		return nil, connErr
	}

	// Check the connection
	pingErr := dbClient.Ping(context.TODO(), nil)
	if pingErr != nil {
		return nil, pingErr
	}

	database := dbClient.Database(sc.DBName)
	return database, nil
}
