package freqwords_test

import (
	"filestore/server/freqwords"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommonWords(t *testing.T) {
	os.Mkdir("testfreqword", 775)
	_, err := os.Create("testfreqword/newfile.txt")
	if err != nil {
		t.Error(err)
	}
	b := []byte("Hello World Helo Hello Helo Helo")
	err = ioutil.WriteFile("testfreqword/newfile.txt", b, 0775)
	if err != nil {
		t.Error(err)
	}
	freqwords.Order = "asc"
	obj := freqwords.CommonWords("testfreqword/", 2)
	assert.Equal(t, 2, len(obj))
	for _, v := range obj {
		if v.Word == "Hello" {
			assert.Equal(t, 2, v.Count)
		}
		if v.Word == "World" {
			assert.Equal(t, 1, v.Count)
		}
	}

}
