package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"testing"
)

const (
	MOCKDIR  = "MOCKDIR"
	MOCKFILE = "MOCKDIR/test"
	PERM_ALL = 0777
	NO_PERM  = 0000
)

func CreateMockArchive(permission int) error {
	var perm = (fs.FileMode(permission))
	err := os.MkdirAll(MOCKDIR, perm)
	if err != nil {
		return fmt.Errorf("error creating mock folder")
	}
	err = os.MkdirAll(MOCKDIR+"/node_modules", perm)
	if err != nil {
		return fmt.Errorf("error creating mock folder")
	}
	err = os.WriteFile(MOCKFILE, []byte("test"), perm)
	if err != nil {
		return fmt.Errorf("error creating mock file")
	}
	err = os.WriteFile(MOCKDIR+"/node_modules/testfile", []byte("test"), perm)
	if err != nil {
		return fmt.Errorf("error creating mock file")
	}
	return nil
}

func DeleteMockArchive() error {
	err := os.RemoveAll(MOCKDIR)
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
	CreateMockArchive(0777)
	defer DeleteMockArchive()
	expected := 4
	got := GetTotalSize(MOCKFILE)
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
	CreateMockArchive(PERM_ALL)
	defer DeleteMockArchive()
	got := MapAllPaths(".")
	expected := PruneData{
		Size:  "4.0B",
		Paths: []string{fmt.Sprintf("MOCKDIR%snode_modules", string(os.PathSeparator))},
	}
	if !got.Equals(&expected) {
		t.Errorf("expected %v, got %v", expected, got)
	}
}

func TestRemover_DeleteAll(t *testing.T) {
	CreateMockArchive(PERM_ALL)
	defer DeleteMockArchive()
	paths := MapAllPaths(".")
	DeleteAll(paths.Paths)
	_, err := os.Stat(MOCKDIR + "/node_modules/testfile")
	if !errors.Is(err, os.ErrNotExist) {
		t.Errorf("expected folder to be clean but file exists")
	}
}

func TestNoPermission(t *testing.T) {
	CreateMockArchive(NO_PERM)
	defer DeleteMockArchive()
	paths := MapAllPaths(".")
	DeleteAll(paths.Paths)
}
