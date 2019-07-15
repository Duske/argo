package ipfs

import (
	"fmt"
	wfv1 "github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	"github.com/argoproj/argo/workflow/common"
	"github.com/argoproj/pkg/json"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"os"
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
	log.Infof("IPFS Load by hash: %s", inputArtifact.IPFS.Hash)
	nodeIP := os.Getenv(common.EnvVarDownwardAPINodeIP)
	if nodeIP == "" {
		return fmt.Errorf("empty envvar %s", common.EnvVarDownwardAPINodeIP)
	}
	endpoint := "/ip4/" + nodeIP + "/tcp/5001"
	return common.RunCommand("ipfs", "--api", endpoint, "get", "-o", path, inputArtifact.IPFS.Hash)
}

func (h *IPFSDriver) Save(path string, outputArtifact *wfv1.Artifact) error {
	fileArg := "file=@" + path
	storageEndpoint := outputArtifact.IPFS.StorageEndpoint
	if storageEndpoint == "" {
		return fmt.Errorf("Please define a storageEndpoint ")
	}
	outputJSON, err := common.RunCommandWithOutput("curl", "-F", fileArg, "http://ipfs-cluster:9094/add")
	if err != nil {
		return err
	}
	var result map[string]interface{}
	err = json.Unmarshal([]byte(outputJSON), &result)
	if err != nil {
		log.Errorf("Could not decode curl output")
		return err
	}
	cidData := result["cid"].(map[string]interface{})
	cid := cidData["/"].(string)
	outputArtifact.IPFS.Hash = cid
	log.Infof("IPFS save by cid: %s on %s", cid, storageEndpoint)
	log.Infof("Storage: %s", )
	return nil
}
