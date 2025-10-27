package s3

import (
	"log"
	"os"
	"sidecar/internal/util"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var MinioClient *minio.Client

func InitS3Client() {
	endpoint := util.StripProtocol(os.Getenv("S3_ENDPOINT"))
	accessKeyID := os.Getenv("S3_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("S3_ACCESS_KEY_SECRET")

	var err error
	// Initialize minio client object.
	MinioClient, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: false,
	})

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("S3 client initialized")
}
