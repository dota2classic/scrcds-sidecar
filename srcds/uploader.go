package srcds

import (
	"io/ioutil"
	"log"
)

// ArtifactType type of uploaded artifact
type ArtifactType string

// Define allowed values
const (
	ArtifactReplay ArtifactType = "replay"
	ArtifactLog    ArtifactType = "log"
)

func UploadArtifacts() {
	uploadFolder("./dota/logs", ArtifactLog)
	uploadFolder("./dota/replays", ArtifactReplay)
}

func uploadFolder(dir string, artifactType ArtifactType) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Printf("Failed to list files in logs: %v", err)
		return
	}
	for _, f := range files {
		log.Printf("Uploading %s %s", artifactType, f.Name())
	}
}
