package main

import (
	"github.com/Azure/azure-storage-blob-go/azblob"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func cloudMediaRecycle() error {
	var err error
	defer func() {
		if err != nil {
			logging.Errormf(logModCloudMediaMgmt, err.Error())
		}
	}()

	// DB student images
	var students []*Student
	students, err = findStudent(primitive.NilObjectID)
	if err != nil {
		return err
	}
	var studentImageMap = map[string]bool{}
	for i := range students {
		studentImageMap[students[i].StudentImageName] = true
	}

	// DB cloud media
	var cloudMediaSlice []*CloudMedia
	cloudMediaSlice, err = findCloudMedia(primitive.NilObjectID)
	if err != nil {
		return err
	}
	var cloudMediaMap = map[string]bool{}
	for i := range cloudMediaSlice {
		cloudMediaMap[cloudMediaSlice[i].MediaName] = true
	}

	// azure blobs
	var azMediaBlobs []*azblob.BlobItem
	azMediaBlobs, err = azureStorageListBlobs(azMediaContainerURL, "")
	if err != nil {
		return err
	}

	// delete blobs which are not linked to DB image/cloudmedia entries
	for i := range azMediaBlobs {
		_, imageOK := studentImageMap[azMediaBlobs[i].Name]
		_, mediaOK := cloudMediaMap[azMediaBlobs[i].Name]
		if !imageOK && !mediaOK {
			err = azureStorageDeleteBlob(azMediaContainerURL, azMediaBlobs[i].Name)
			if err != nil {
				return err
			}
			logging.Infomf(logModCloudMediaMgmt, "[%s] is neither student image or cloud media ==> delete it", azMediaBlobs[i].Name)
		}
	}

	return nil
}
