package main

import (
	"errors"
	"fmt"
	"os"
	"testing"
)

const (
	MockDir  = "mockdir"
	MockFile = "mockdir/test"
)

func CreateMockArchive() error {
	err := os.MkdirAll(MockDir, 0777)
	if err != nil {
		return fmt.Errorf("error creating mock folder")
	}
	err = os.MkdirAll(MockDir+"/node_modules", 0777)
	if err != nil {
		return fmt.Errorf("error creating mock folder")
	}
	err = os.WriteFile(MockFile, []byte("test"), 0777)
	if err != nil {
		return fmt.Errorf("error creating mock file")
	}
	err = os.WriteFile(MockDir+"/node_modules/testfile", []byte("test"), 0777)
	if err != nil {
		return fmt.Errorf("error creating mock file")
	}
	return nil
}

func DeleteMockArchive() error {
	err := os.RemoveAll(MockDir)
	if err != nil {
		return err
	}
	return nil
}

func TestRemover_FormatSize(t *testing.T) {
	expected := "3.0B"
	got := formatSize(3)
	if got != expected {
		t.Errorf("expected %v, got %v", expected, got)
	}
}

func TestRemover_GetTotalSize(t *testing.T) {
	CreateMockArchive()
	defer DeleteMockArchive()
	expected := 4
	got := GetTotalSize(MockFile)
	if got != expected {
		t.Errorf("expected %v, got %v", expected, got)
	}
}

func TestRemover_MapAllPaths_NoMatches(t *testing.T) {
	expected := &PruneData{
		Size: "0.0B",
	}
	got := MapAllPaths(".")
	if !expected.Equals(expected) {
		t.Errorf("expected %v, got %v", expected, got)
	}
}

func TestRemover_MapAllPaths(t *testing.T) {
	CreateMockArchive()
	defer DeleteMockArchive()
	got := MapAllPaths(".")
	expected := PruneData{
		Size:  "4.0B",
		Paths: []string{fmt.Sprintf("mockdir%snode_modules", string(os.PathSeparator))},
	}
	if !got.Equals(&expected) {
		t.Errorf("expected %v, got %v", expected, got)
	}
}

func TestRemover_DeleteAll(t *testing.T) {
	CreateMockArchive()
	defer DeleteMockArchive()
	paths := MapAllPaths(".")
	DeleteAll(paths.Paths)
	_, err := os.Stat(MockDir + "/node_modules/testfile")
	if !errors.Is(err, os.ErrNotExist) {
		t.Errorf("expected folder to be clean but file exists")
	}
}
