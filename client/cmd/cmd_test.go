package cmd_test

import (
	"errors"
	"filestore/client/apicall"
	"filestore/client/cmd"
	"os"
	"strings"
	"testing"
)

func Add_Test(t *testing.T) {
	apicall.ServerUrl = "http://localhost:5000/v1/files"
	os.Create("test.txt")
	cmd.Addfile([]string{"test.txt"})
	if _, err := os.Stat("filestore/test.txt"); errors.Is(err, os.ErrNotExist) {
		t.Fatalf("Failed to create file")
	}
	op := apicall.List()
	if !strings.Contains(op, "test.txt") {
		t.Errorf("List gave wrong output")
	}
	op = apicall.WC()
	if !strings.Contains(op, "0") {
		t.Errorf("WC gave wrong output")
	}
	op = cmd.Remove([]string{"test.txt"})
	if !strings.Contains(op, "Deleted") {
		t.Errorf("Failed to delete file output")
	}

}
