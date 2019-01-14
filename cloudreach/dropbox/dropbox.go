package dropbox

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
)

const dropboxUploadUrl = "https://content.dropboxapi.com/2/files/upload"

func AscendFile(filepath string, destination string, token string) error {
	params := map[string]interface{}{
		"path": destination,
		"mode": "add",
		"autorename": true,
	}
	jsonData, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("could not convert req params to JSON: %s", err)
	}

	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	defer bodyWriter.Close()

	// this step is very important
	fileWriter, err := bodyWriter.CreateFormFile("uploadfile", filepath)
	if err != nil {
		fmt.Println("error writing to buffer")
		return fmt.Errorf("error writing file %s to buffer: %s", filepath, err)
	}

	// open file handle
	fh, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("error opening file %s: %s", filepath, err)
	}
	defer fh.Close()

	//iocopy
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		return fmt.Errorf("could not copy file %s to writer: %s", filepath, err)
	}

	bodyWriter.Close()

	r, err := http.NewRequest("POST", dropboxUploadUrl, bodyBuf)

	r.Header.Add("Authorization", "Bearer " + token)
	r.Header.Add("Content-Type", "application/octet-stream" )
	r.Header.Add("Dropbox-API-Arg", string(jsonData))

	response, err := http.DefaultClient.Do(r)
	if err != nil {
		return fmt.Errorf("error posting to dropbox: %s", err)
	}
	if response.StatusCode > 399 {
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return fmt.Errorf("could not read dropbox error response: %s", err)
		}
		return fmt.Errorf("error posting to dropbox: %s", body)
	}
	return nil
}

