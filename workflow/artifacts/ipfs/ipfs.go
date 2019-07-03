package ipfs

import (
	wfv1 "github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"time"
)

// HTTPArtifactDriver is the artifact driver for a HTTP URL
type IPFSDriver struct{}

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func String(length int) string {
	return StringWithCharset(length, charset)
}

// Load download artifacts from an HTTP URL
func (h *IPFSDriver) Load(inputArtifact *wfv1.Artifact, path string) error {
	// Download the file to a local file path
	log.Infof("IPFS Load by hash: %s", inputArtifact.IPFS.Hash)
	return nil
}

func (h *IPFSDriver) Save(path string, outputArtifact *wfv1.Artifact) error {
	var hash = String(256)
	outputArtifact.IPFS.Hash = hash
	log.Infof("IPFS Load by hash: %s", hash)
	return nil
}
