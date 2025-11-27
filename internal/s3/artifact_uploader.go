package s3

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sidecar/internal/util"

	"github.com/minio/minio-go/v7"
)

// ArtifactType type of uploaded artifact
type ArtifactType string

// Define allowed values
const (
	ArtifactReplay ArtifactType = "replay"
	ArtifactLog    ArtifactType = "log"
)

func UploadArtifacts(matchId int64) {
	log.Printf("Uploading artifacts for matchId %d", matchId)
	uploadFolder(util.LOG_FOLDER, ArtifactLog, matchId)
	uploadFolder(util.REPLAY_FOLDER, ArtifactReplay, matchId)
	log.Printf("Artifacts for matchId %d successfully uploaded", matchId)
}

func uploadFolder(dir string, artifactType ArtifactType, matchId int64) {
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Printf("Failed to list files in folder: %v", err)
		return
	}
	for _, f := range files {
		if f.IsDir() {
			continue // skip subdirectories
		}
		if f.Name() == "discarded" {
			continue // skip discarded replays(failed match)
		}
		filePath := filepath.Join(dir, f.Name())
		uploadFile(filePath, artifactType, matchId)
	}
}

func uploadFile(filePath string, artifactType ArtifactType, matchId int64) {
	log.Printf("Uploading %s %s", artifactType, filePath)

	ctx := context.Background()

	// If it's a
	//if artifactType == ArtifactReplay {
	//	outPath := filePath + ".zip"
	//	err := util.CompressFile(filePath, outPath)
	//	log.Printf("Zipped replay file: %s", filePath)
	//	if err != nil {
	//		log.Fatalln(err)
	//	}
	//	filePath = outPath
	//}

	var bucket string
	var filename string
	var contentType string

	if artifactType == ArtifactLog {
		filename = fmt.Sprintf("%d.log", matchId)
		contentType = "text/plain"
		bucket = "logs"
	} else if artifactType == ArtifactReplay {
		filename = fmt.Sprintf("%d.dem", matchId)
		contentType = "application/octet-stream"
		bucket = "replays"
	}

	log.Printf("Uploading %s to bucket %s", filename, bucket)

	info, err := MinioClient.FPutObject(
		ctx,
		bucket,
		filename,
		filePath,
		minio.PutObjectOptions{ContentType: contentType},
	)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Successfully uploaded %s of size %d", filePath, info.Size)

	err = os.Remove(filePath)
	if err != nil {
		log.Printf("Failed to remove file: %v", err)
	}
}
