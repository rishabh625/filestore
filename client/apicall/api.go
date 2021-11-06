package apicall

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"
)

var ServerUrl string

// CopyFileData ... Struct to hold File Name and SHA256 Hash of File Content
type CopyFileData struct {
	FileName string
	Hash     string
}

// Hexists ... Takes Filename and its content's Hash as input, Checks whether File Content is repeated or File Already Exists Querying Server. return Siggnature (?file content is repeated,?file is duplicated)
func Hexists(hash, file string) (bool, bool) {
	req, err := http.NewRequest(http.MethodGet, ServerUrl+"/hexist", nil)
	if err != nil {
		return false, true
	}
	req.Header.Set("file", file)
	req.Header.Set("hash", hash)
	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := netClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return false, true
	}
	if resp.StatusCode == http.StatusOK {
		return true, true
	}
	if resp.StatusCode == http.StatusConflict {
		return false, false
	}
	return false, true
}

// CopyCall ... Takes Filename and its content's Hash as input, Calls /copy URI to store file without Sending its content returns true if file is copied else returns false
func CopyCall(filename, hash string) bool {
	payloadstr := CopyFileData{
		FileName: filename,
		Hash:     hash,
	}
	marshaldata, err := json.Marshal(payloadstr)
	if err != nil {
		return false
	}
	payload := strings.NewReader(string(marshaldata))

	req, err := http.NewRequest(http.MethodPost, ServerUrl+"/copy", payload)
	if err != nil {
		req.Header.Set("hash", hash)
	}
	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}
	resp, _ := netClient.Do(req)
	if resp.StatusCode == http.StatusOK {
		return true
	}
	return false

}

// AddCall ... Asks Server to store file By calling /store input if 200 responsecode recieved from server returns true (file stored succesfully) else returns false File Data is sent as multipart form to avoid netwrk issues
func AddCall(filepath string, hash string) bool {
	r, w := io.Pipe()
	m := multipart.NewWriter(w)
	go func() {
		defer w.Close()
		defer m.Close()
		part, err := m.CreateFormFile("file", filepath)
		if err != nil {
			return
		}
		file, err := os.Open(filepath)
		if err != nil {
			return
		}
		defer file.Close()
		if _, err = io.Copy(part, file); err != nil {
			return
		}
	}()
	req, _ := http.NewRequest(http.MethodPost, ServerUrl+"/store", r)
	req.Header.Set("Content-Type", m.FormDataContentType())
	req.Header.Set("hash", hash)
	var netClient = &http.Client{
		Timeout: time.Minute * 1,
	}
	resp, _ := netClient.Do(req)
	if resp.StatusCode == http.StatusOK {
		return true
	}
	return false
}

// List ... Lists All files in filestore by Querying Server at /ls , returns error as string or list as string of json
func List() string {
	req, _ := http.NewRequest(http.MethodGet, ServerUrl+"/ls", nil)
	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}
	resp, _ := netClient.Do(req)
	if resp.StatusCode == http.StatusOK {
		body := &bytes.Buffer{}
		_, err := body.ReadFrom(resp.Body)
		resp.Body.Close()
		if err != nil {
			return fmt.Sprintf(err.Error())
		} else {
			return fmt.Sprintf(body.String())
		}
	}
	return fmt.Sprintf("Failed to list files")
}

// Remove ... Asks Server to Delete input file returns by sending HTTP Verb Delete
func Remove(file string) string {
	payloadstr := CopyFileData{
		FileName: file,
	}
	marshaldata, err := json.Marshal(payloadstr)
	if err != nil {
		return fmt.Sprintf("Failed to Delete file %s", file)
	}
	payload := strings.NewReader(string(marshaldata))
	req, _ := http.NewRequest(http.MethodDelete, ServerUrl+"/rm", payload)
	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}
	resp, _ := netClient.Do(req)
	body := &bytes.Buffer{}
	_, err = body.ReadFrom(resp.Body)
	resp.Body.Close()
	if err != nil {
		resp.Body.Close()
		return err.Error()
	}
	if resp.StatusCode == http.StatusOK {
		return fmt.Sprintf("File Deleted Successfully %s ", file)
	} else {
		return fmt.Sprintf("Failed to Delete file ", file, "Reason: ", body.String())
	}
}

//WC ... returns number of word present in filestore as string or returns error as string
func WC() string {
	req, _ := http.NewRequest(http.MethodGet, ServerUrl+"/wc", nil)
	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}
	resp, _ := netClient.Do(req)
	if resp.StatusCode == http.StatusAccepted {
		body := &bytes.Buffer{}
		_, err := body.ReadFrom(resp.Body)
		resp.Body.Close()
		if err != nil {
			return err.Error()
		} else {
			return body.String()
		}
	}
	return "Failed to give word count"
}
