package cmd_test

import (
	"errors"
	"filestore/client/apicall"
	"filestore/client/cmd"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestAdd Tests Function Usecase plz start server at localhost:5000 before using this
func TestAdd(t *testing.T) {
	apicall.ServerUrl = "http://localhost:5000/files/v1"
	_, err := os.Create("../../test.txt")
	if err != nil {
		t.Fatalf("%s %+v", "Failed to initiate tests ", err)
	}
	cmd.Addfile([]string{"test.txt"})
	if _, err := os.Stat("../../filestore/test.txt"); errors.Is(err, os.ErrNotExist) {
		t.Fatalf("%s %s", "Failed to create file", "filestore/test.txt")
	}
	op := apicall.List()
	if !strings.Contains(op, "test.txt") {
		t.Errorf("List gave wrong output")
	}
	op = apicall.WC()
	if !strings.Contains(op, "count") {
		t.Errorf("WC gave wrong output %s", op)
	}
	op = cmd.Remove([]string{"test.txt"})
	if !strings.Contains(op, "Deleted") {
		t.Errorf("Failed to delete file output")
	}

	b := []byte("Hello World")
	// write the whole body at once
	err = ioutil.WriteFile("test.txt", b, 0644)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	cmd.Updatefile([]string{"test.txt"})
	readbytes, err := ioutil.ReadFile("../../filestore/test.txt")
	if err != nil {
		t.Fatalf("%+v", err)
	}
	assert.Equal(t, "Hello World", string(readbytes))
}

func TestValidateFlags(t *testing.T) {
	limit, order := cmd.ValidateFlags("non", "wrong")
	assert.Equal(t, "10", limit)
	assert.Equal(t, "desc", order)
}
