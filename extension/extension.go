package extension

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
)

const FolderPattern string = "dynatrace_extension"
const InnerFileName string = "plugin.json"

type PackagedExtension struct {
	FilePath string
}

type ExtensionPayload struct {
	Payload string
}

type DynatraceExtensionMetadata struct {
	Name string `json:"name"`
}

func GetExtensionPayloadFromPackage(zipPackage []byte) (*ExtensionPayload, error) {
	reader, err := zip.NewReader(bytes.NewReader(zipPackage), int64(len(zipPackage)))
	if err != nil {
		return nil, err
	}

	pluginFile, err := reader.Open("plugin.json")
	if err != nil {
		return nil, err
	}

	pluginContentBytes, err := io.ReadAll(pluginFile)
	if err != nil {
		return nil, err
	}

	return &ExtensionPayload{
		Payload: string(pluginContentBytes),
	}, nil
}

func CreatePackagedExtension(payload string) (*PackagedExtension, error) {
	var metadata DynatraceExtensionMetadata
	err := json.Unmarshal([]byte(payload), &metadata)
	if err != nil {
		return nil, err
	}

	// create a zip file containing a plugin.json file with payload content
	zipDir, err := os.MkdirTemp(os.TempDir(), FolderPattern)
	if err != nil {
		return nil, err
	}

	zipFilePath := fmt.Sprintf("%s/%s.zip", zipDir, metadata.Name)

	zipFile, err := os.Create(zipFilePath)
	if err != nil {
		return nil, err
	}
	defer zipFile.Close()
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	entry, err := zipWriter.Create(InnerFileName)
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(entry, strings.NewReader(payload))
	if err != nil {
		return nil, err
	}
	return &PackagedExtension{FilePath: zipFilePath}, nil
}

func (p *PackagedExtension) Dispose() {
	parentDirectory := path.Dir(p.FilePath)
	os.RemoveAll(parentDirectory)
}
