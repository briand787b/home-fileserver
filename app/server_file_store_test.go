package app

import (
	"fmt"
	"os"
	"testing"
)

const (
	fakeServerPath = "this_file_does_not_exist.json"
)

func TestMain(m *testing.M) {
	setupFS()
	exitStatus := m.Run()
	tearDownFS()
	os.Exit(exitStatus)
}

// prepares filesystem for test running
func setupFS() {
	removeTestFiles()
}

// cleans up filesystem after test running
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
func TestServerStoreGetState(t *testing.T) {
	if _, err := newFileServerStore(fakeServerPath); err != nil {
		t.Fatal("unable to instantiate new file server store: ", err)
	}
}

func TestFileServerStoreSaveStateFromBlank(t *testing.T) {
	fss, err := newFileServerStore(fakeServerPath)
	if err != nil {
		t.Fatal("unable to instantiate new file server store: ", err)
	}

	fmt.Println(fss)

	if err := fss.Save(&Server{}); err != nil {
		t.Fatal("could not save server to file: ", err)
	}
}
