package s3

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sidecar/internal/models"
	"sidecar/internal/rabbit"
	"sidecar/internal/state"
	"sidecar/internal/util"

	"github.com/minio/minio-go/v7"
)

func logDumpDirs() {
	for _, dir := range []string{"/tmp/dumps", "/tmp/dumps01"} {
		files, err := os.ReadDir(dir)
		if err != nil {
			log.Printf("[dumps] %s: %v", dir, err)
			continue
		}
		for _, f := range files {
			log.Printf("[dumps] %s/%s", dir, f.Name())
		}
	}
}

func UploadArtifacts(matchId int64) {
	logDumpDirs()
	log.Printf("Uploading artifacts for matchId %d", matchId)
	uploadFolder(util.LOG_FOLDER, models.ARTIFACT_TYPE_LOG, matchId)
	uploadFolder(util.REPLAY_FOLDER, models.ARTIFACT_TYPE_REPLAY, matchId)
	uploadFolder(util.DUMP_FOLDER, models.ARTIFACT_TYPE_DUMP, matchId)
	uploadFolder(util.DUMP_FOLDER_01, models.ARTIFACT_TYPE_DUMP, matchId)
	log.Printf("Artifacts for matchId %d successfully uploaded", matchId)
}

func uploadFolder(dir string, artifactType models.ArtifactType, matchId int64) {
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

func uploadFile(filePath string, artifactType models.ArtifactType, matchId int64) {
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

	if artifactType == models.ARTIFACT_TYPE_LOG {
		filename = fmt.Sprintf("%d.log", matchId)
		contentType = "text/plain"
		bucket = "logs"
	} else if artifactType == models.ARTIFACT_TYPE_REPLAY {
		filename = fmt.Sprintf("%d.dem", matchId)
		contentType = "application/octet-stream"
		bucket = "replays"
	} else if artifactType == models.ARTIFACT_TYPE_DUMP {
		filename = fmt.Sprintf("%d_%s", matchId, filepath.Base(filePath))
		contentType = "application/octet-stream"
		bucket = "dumps"
	}

	if bucket == "" || filename == "" {
		log.Printf("Unknown artifact type %s, skipping %s", artifactType, filePath)
		return
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
	evt := models.MatchArtifactUploadedEvent{
		MatchID:      matchId,
		ArtifactType: artifactType,
		LobbyType:    state.GlobalMatchInfo.LobbyType,
		Bucket:       bucket,
		Key:          filename,
	}
	rabbit.PublishArtifactUploadedEvent(&evt)

	err = os.Remove(filePath)
	if err != nil {
		log.Printf("Failed to remove file: %v", err)
	}
}
