package srcds

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// ArtifactType type of uploaded artifact
type ArtifactType string

// Define allowed values
const (
	ArtifactReplay ArtifactType = "replay"
	ArtifactLog    ArtifactType = "log"
)

func UploadArtifacts(matchId int64) {
	uploadFolder("./dota/logs", ArtifactLog, matchId)
	uploadFolder("./dota/replays", ArtifactReplay, matchId)
}

func uploadFolder(dir string, artifactType ArtifactType, matchId int64) {
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Printf("Failed to list files in logs: %v", err)
		return
	}
	for _, f := range files {
		if f.IsDir() {
			continue // skip subdirectories
		}
		filePath := filepath.Join(dir, f.Name())
		uploadFile(filePath, artifactType, matchId)
	}
}

func compressFile(filePath string) (string, error) {
	zipPath := filePath + ".zip"

	outFile, err := os.Create(zipPath)
	if err != nil {
		return "", fmt.Errorf("failed to create zip file: %w", err)
	}
	defer outFile.Close()

	zipWriter := zip.NewWriter(outFile)
	defer zipWriter.Close()

	fileToZip, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open source file: %w", err)
	}
	defer fileToZip.Close()

	w, err := zipWriter.Create(filepath.Base(filePath))
	if err != nil {
		return "", fmt.Errorf("failed to create zip entry: %w", err)
	}

	if _, err := io.Copy(w, fileToZip); err != nil {
		return "", fmt.Errorf("failed to write zip content: %w", err)
	}

	return zipPath, nil
}

func uploadFile(filePath string, artifactType ArtifactType, matchId int64) {
	log.Printf("Uploading %s %s", artifactType, filePath)

	ctx := context.Background()
	endpoint := os.Getenv("S3_ENDPOINT")
	accessKeyID := os.Getenv("S3_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("S3_ACCESS_KEY_SECRET")

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: false,
	})

	if err != nil {
		log.Fatalln(err)
	}

	// If it's a
	if artifactType == ArtifactReplay {
		filePath, err = compressFile(filePath)
		if err != nil {
			log.Fatalln(err)
		}
	}

	var bucket string
	var filename string
	var contentType string

	if artifactType == ArtifactLog {
		filename = fmt.Sprintf("%d.log", matchId)
		contentType = "text/plain"
		bucket = "logs"
	} else {
		filename = fmt.Sprintf("%d.dem.zip", matchId)
		contentType = "application/zip"
		bucket = "replays"
	}

	info, err := minioClient.FPutObject(ctx, bucket, filename, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Successfully uploaded %s of size %d", filePath, info.Size)

}
