package file

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLoggerWithLFSHook(t *testing.T) {
	file, err := ioutil.TempFile("", "test-log-file")
	if err != nil {
		log.Fatal(err)
	}
	filename := file.Name()
	defer os.Remove(filename)
	l := NewLoggerWithLFSHook(filename)
	content := "filename: " + filename
	l.Println(content)

	dat, err := ioutil.ReadFile(filename)
	assert.Nil(t, err)
	assert.True(t, strings.Contains(string(dat), content), "log file should contain the expected line")
}
