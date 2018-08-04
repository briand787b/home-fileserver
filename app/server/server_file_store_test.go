package server

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

const (
	fakeServerPath = "this_file_does_not_exist.json"
	testServerName = "test_server"
	testServerIP   = "127.0.0.1"
)

var (
	validTestServer *Server
)

func TestMain(m *testing.M) {
	setupFS()
	exitStatus := m.Run()
	tearDownFS()
	os.Exit(exitStatus)
}

// prepares for running fileStore tests
func setupFS() {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println("cannot get working dir: ", err)
		os.Exit(1)
	}

	validTestServer = &Server{
		Name:       testServerName,
		IP:         testServerIP,
		WorkingDir: pwd,
	}

	removeTestFiles()
}

// cleans up after running fileStore tests
func tearDownFS() {
	removeTestFiles()
}

// removes test-only files
func removeTestFiles() {
	if _, err := os.Stat(fakeServerPath); err == nil {
		if err := os.Remove(fakeServerPath); err != nil {
			fmt.Printf("cannot remove %s: %s", fakeServerPath, err)
			os.Exit(1)
		}
	}
}

// Verify that no error is generated when
// instantiating new file-based server store
func TestGetNewFileStore(t *testing.T) {
	if _, err := newFileStore(fakeServerPath); err != nil {
		t.Fatal("unable to instantiate new file server store: ", err)
	}
}

func TestFileStoreSaveFailsNoName(t *testing.T) {
	fss, err := newFileStore(fakeServerPath)
	if err != nil {
		t.Fatal("unable to instantiate new file server store: ", err)
	}

	if err := fss.add(&Server{
		IP:         "0.0.0.0",
		WorkingDir: "/valid/path",
	}, fakeServerPath); err == nil {
		t.Fatal("should not be able to save server with no name")
	}
}

func TestFileStoreSaveFailsBadIPv4(t *testing.T) {
	fss, err := newFileStore(fakeServerPath)
	if err != nil {
		t.Fatal("unable to instantiate new file server store: ", err)
	}

	badIPs := [][]string{
		[]string{"", "IP is empty string"},
		[]string{"0.0.0", "not enough dots"},
		[]string{"0.0.0.cat", "not a number"},
		[]string{"0.0.0.256", "number too large"},
		[]string{"0.0.0.-1", "number too small"},
	}

	for _, ip := range badIPs {
		if err := fss.add(&Server{
			Name:       testServerName,
			IP:         ip[0],
			WorkingDir: "/valid",
		}, fakeServerPath); err == nil {
			t.Fatal("should not be able to save because: ", ip[1])
		}
	}
}

func TestFileStoreSaveFailsNoWorkingDir(t *testing.T) {
	fss, err := newFileStore(fakeServerPath)
	if err != nil {
		t.Fatal("unable to instantiate new file server store: ", err)
	}

	if err := fss.add(&Server{
		IP:   "0.0.0.0",
		Name: "ValidName",
	}, fakeServerPath); err == nil {
		t.Fatal("should not be able to save server with no working dir")
	}
}

func TestFileStoreAddServerFromCleanState(t *testing.T) {
	fss, err := newFileStore(fakeServerPath)
	if err != nil {
		t.Fatal("unable to instantiate new file server store: ", err)
	}

	if err := fss.add(validTestServer, fakeServerPath); err != nil {
		t.Fatal("could not save single validTestServer")
	}

	actual, err := ioutil.ReadFile(fakeServerPath)
	if err != nil {
		t.Fatal("could not open written JSON file: ", err)
	}

	tActual := cutWhiteSpace(actual, t)
	t.Logf("tActual: %v\n", string(tActual))

	goldenFile := parseGoldenFile(t).Bytes()

	// validate saved JSON is correct
	if bytes.Compare(goldenFile, tActual) != 0 {
		t.Fatal("written JSON file does not match golden file")
	}
}

func TestFileStoreAddSecondServer(t *testing.T) {
	fss, err := newFileStore(fakeServerPath)
	if err != nil {
		t.Fatal("unable to instantiate new file server store: ", err)
	}

	// first add, starts at 0
	if err := fss.add(validTestServer, fakeServerPath); err != nil {
		t.Fatal("could not save single validTestServer")
	}

	// second add, starts at 1
	if err := fss.add(&Server{
		Name:       testServerName + "1",
		IP:         "127.0.0.2",
		WorkingDir: "working_dir",
	}, fakeServerPath); err != nil {
		t.Fatal("could not save single validTestServer")
	}

	actual, err := ioutil.ReadFile(fakeServerPath)
	if err != nil {
		t.Fatal("could not open written JSON file: ", err)
	}

	tActual := cutWhiteSpace(actual, t)
	t.Logf("tActual: %v\n", string(tActual))

	goldenFile := parseGoldenFile(t).Bytes()

	// validate saved JSON is correct
	if bytes.Compare(goldenFile, tActual) != 0 {
		t.Fatal("written JSON file does not match golden file")
	}
}

func TestFileStoreRemoveFirstServer(t *testing.T) {
	fss, err := newFileStore(fakeServerPath)
	if err != nil {
		t.Fatal("unable to instantiate new file server store: ", err)
	}

	// first add, starts at 0
	if err := fss.add(validTestServer, fakeServerPath); err != nil {
		t.Fatal("could not save single validTestServer")
	}

	// second add, starts at 1
	if err := fss.add(&Server{
		Name:       testServerName + "1",
		IP:         "127.0.0.2",
		WorkingDir: "working_dir",
	}, fakeServerPath); err != nil {
		t.Fatal("could not save single validTestServer")
	}

	if err := fss.removeServer(validTestServer.Name, fakeServerPath); err != nil {
		t.Fatal("could not remove first server: ", err)
	}

	actual, err := ioutil.ReadFile(fakeServerPath)
	if err != nil {
		t.Fatal("could not open written JSON file: ", err)
	}

	tActual := cutWhiteSpace(actual, t)
	t.Logf("tActual: %v\n", string(tActual))

	goldenFile := parseGoldenFile(t).Bytes()

	// validate saved JSON is correct
	if bytes.Compare(goldenFile, tActual) != 0 {
		t.Fatal("written JSON file does not match golden file")
	}
}
