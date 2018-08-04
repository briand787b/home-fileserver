package server

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sync"
	"testing"
)

var (
	whiteSpaceRegx *regexp.Regexp
	wsrOnce        sync.Once
)

// **************************************************************
// Why does the tesdata directory exist?  testdata directories
// are not analyzed by the go compiler unless tests are being
// run.  This prevents test fixtures from being included in the
// binary that `go build` outputs.
// **************************************************************

func cutWhiteSpace(b []byte, t *testing.T) []byte {
	wsrOnce.Do(func() {
		var err error
		whiteSpaceRegx, err = regexp.Compile(`\s`)
		if err != nil {
			t.Fatal("could not compile regexp: ", err)
		}
	})

	return whiteSpaceRegx.ReplaceAll(b, []byte(""))
}

func parseGoldenFile(t *testing.T) *bytes.Buffer {
	filepath := filepath.Join("testdata", t.Name()+".golden")
	bs, err := ioutil.ReadFile(filepath)
	if err != nil {
		t.Fatalf("could not read golden file %s: %s", t.Name(), err)
	}

	pwd, err := os.Getwd()
	if err != nil {
		t.Fatal("could not get working dir: ", err)
	}

	b := &bytes.Buffer{}
	if err := template.Must(template.New(
		filepath,
	).Parse(string(cutWhiteSpace(bs, t)))).Execute(b, pwd); err != nil {
		t.Fatal("could not execute parsed template: ", err)
	}

	t.Logf("%v\n", b.String())

	return b
}
