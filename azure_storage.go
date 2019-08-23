package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"

	"github.com/Azure/azure-storage-blob-go/azblob"
)

var azureContainerURL *azblob.ContainerURL

func azureErrorInfo(err error) string {
	if err != nil {
		if serr, ok := err.(azblob.StorageError); ok { // This error is an Azure Service-specific
			return "[Azure Service Code]: " + string(serr.ServiceCode())
		}
		return err.Error()
	}
	return ""
}

func azureStorageInit(sc *ServerConfig) (*azblob.ContainerURL, error) {
	// config check
	if len(sc.AzureStorageAccount) == 0 || len(sc.AzureStorageAccessKey) == 0 {
		return nil, fmt.Errorf("Invalid azure storage account name or access key")
	}

	// create a default request pipeline using storage account name and account key.
	credential, credentailErr := azblob.NewSharedKeyCredential(sc.AzureStorageAccount, sc.AzureStorageAccessKey)
	if credentailErr != nil {
		return nil, fmt.Errorf("Invalid azure credentials with error: %s", credentailErr.Error())
	}
	pipeline := azblob.NewPipeline(credential, azblob.PipelineOptions{})

	// from the Azure portal, get storage account blob service URL endpoint.
	URL, _ := url.Parse(
		fmt.Sprintf("https://%s.blob.core.windows.net/%s", sc.AzureStorageAccount, sc.AzureStorageContainer))

	// create a ContainerURL object that wraps the container URL and a request pipeline to make requests.
	containerURL := azblob.NewContainerURL(*URL, pipeline)

	// check if container exists. If not, try to create
	if _, containerPropErr := containerURL.GetProperties(context.TODO(), azblob.LeaseAccessConditions{}); containerPropErr != nil {
		if serr, ok := containerPropErr.(azblob.StorageError); ok && serr.ServiceCode() == azblob.ServiceCodeContainerNotFound {
			if containerCreateErr := azureStorageCreateContainer(&containerURL); containerCreateErr != nil {
				return nil, containerCreateErr
			}
		} else {
			return nil, fmt.Errorf("Failed to connect to azure container: %s", azureErrorInfo(containerPropErr))
		}
	}

	return &containerURL, nil
}

func azureStorageCreateContainer(azureContainerURL *azblob.ContainerURL) error {
	if azureContainerURL == nil {
		return fmt.Errorf("Empty azure container URL object")
	}

	// create azure storage container
	if _, containerCreateErr := azureContainerURL.Create(context.TODO(), azblob.Metadata{}, azblob.PublicAccessBlob); containerCreateErr != nil {
		return fmt.Errorf("Failed to create azure storage container: %s", azureErrorInfo(containerCreateErr))
	}
	return nil
}

func azureStorageDeleteContainer(azureContainerURL *azblob.ContainerURL) error {
	if azureContainerURL == nil {
		return fmt.Errorf("Empty azure container URL object")
	}

	// delete azure storage container
	if _, containerDeleteErr := azureContainerURL.Delete(context.TODO(), azblob.ContainerAccessConditions{}); containerDeleteErr != nil {
		return fmt.Errorf("Failed to delete azure storage container: %s", azureErrorInfo(containerDeleteErr))
	}
	return nil
}

func azureStorageListBlobs(azureContainer *azblob.ContainerURL, prefix string) ([]*azblob.BlobItem, error) {
	if azureContainerURL == nil {
		return []*azblob.BlobItem{}, fmt.Errorf("Empty azure container URL object")
	}

	var blobItems = []*azblob.BlobItem{}
	for marker := (azblob.Marker{}); marker.NotDone(); {
		// get a result segment starting with the blob indicated by the current Marker.
		listBlob, listBlobErr := azureContainer.ListBlobsFlatSegment(context.TODO(), marker, azblob.ListBlobsSegmentOptions{
			Prefix: prefix,
		})
		if listBlobErr != nil {
			return []*azblob.BlobItem{}, fmt.Errorf("Failed to list blobs in azure container: %s", azureErrorInfo(listBlobErr))
		}

		// append blobs to results
		for bi := range listBlob.Segment.BlobItems {
			blobItems = append(blobItems, &listBlob.Segment.BlobItems[bi])
		}

		// update marker
		marker = listBlob.NextMarker
	}
	return blobItems, nil
}

func azureStorageIsBlobExist(azureContainerURL *azblob.ContainerURL, blobname string) bool {
	blobURL := azureContainerURL.NewBlobURL(blobname)
	blobResp, err := blobURL.GetProperties(context.TODO(), azblob.BlobAccessConditions{})
	if err != nil || blobResp == nil {
		return false
	}
	return true
}

func azureStorageUploadBlob(azureContainerURL *azblob.ContainerURL, blobname string) error {
	if azureContainerURL == nil {
		return fmt.Errorf("Empty azure container URL object")
	}

	// read blob file
	blobFile, fileErr := os.Open(blobname)
	if fileErr != nil {
		return fmt.Errorf("Failed to open file %s for blob upload", blobname)
	}
	defer blobFile.Close()

	// upload blob file
	blobURL := azureContainerURL.NewBlockBlobURL(blobname)
	_, blobUploadErr := azblob.UploadFileToBlockBlob(context.TODO(), blobFile, blobURL, azblob.UploadToBlockBlobOptions{
		BlockSize:   4 * 1024 * 1024,
		Parallelism: 16})
	if blobUploadErr != nil {
		return fmt.Errorf("Failed to upload blob %s: %s", blobname, azureErrorInfo(blobUploadErr))
	}

	return nil
}

func azureStorageDownloadBlob(azureContainerURL *azblob.ContainerURL, blobname string) error {
	if azureContainerURL == nil {
		return fmt.Errorf("Empty azure container URL object")
	}

	// download blob file
	blobURL := azureContainerURL.NewBlockBlobURL(blobname)
	downloadResp, downloadErr := blobURL.Download(context.TODO(), 0, azblob.CountToEnd, azblob.BlobAccessConditions{}, false)
	if downloadErr != nil {
		return fmt.Errorf("Failed to download blob %s: %s", blobname, azureErrorInfo(downloadErr))
	}
	bodyStream := downloadResp.Body(azblob.RetryReaderOptions{MaxRetryRequests: 5})
	downloadedData := bytes.Buffer{}
	_, dataErr := downloadedData.ReadFrom(bodyStream)
	if dataErr != nil {
		return fmt.Errorf("Failed to read downloaded blob data: %s", dataErr.Error())
	}

	// save blob file
	fileErr := ioutil.WriteFile(blobname, downloadedData.Bytes(), 0755)
	if fileErr != nil {
		return fmt.Errorf("Failed to save blob file %s: %s", blobname, fileErr.Error())
	}

	return nil
}

func azureStorageDeleteBlob(azureContainerURL *azblob.ContainerURL, blobname string) error {
	if azureContainerURL == nil {
		return fmt.Errorf("Empty azure container URL object")
	}

	// delete blob
	blobURL := azureContainerURL.NewBlockBlobURL(blobname)
	_, deleteErr := blobURL.Delete(context.TODO(), azblob.DeleteSnapshotsOptionInclude, azblob.BlobAccessConditions{})
	if deleteErr != nil {
		return fmt.Errorf("Failed to delete blob %s: %s", blobname, azureErrorInfo(deleteErr))
	}

	return nil
}
