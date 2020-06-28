package utils

import (
	"path/filepath"
	"testing"

	"github.com/jfrog/jfrog-client-go/utils/io/fileutils"
)

func TestVcsDetails(t *testing.T) {
	path := initVcsTestDir(t)
	vcsDetals := NewVcsDetals()
	revision, url, err := vcsDetals.GetVcsDetails(filepath.Join(path))
	if err != nil {
		t.Error(err)
	}
	if url != "https://github.com/jfrog/jfrog-cli.git" {
		t.Errorf("TestGitManager() error, want %s, got %s", url, "https://github.com/jfrog/jfrog-cli.git")
	}
	if revision != "d63c5957ad6819f4c02a817abe757f210d35ff92" {
		t.Errorf("TestGitManager() error, want %s, got %s", url, "d63c5957ad6819f4c02a817abe757f210d35ff92")
	}
}

func initVcsTestDir(t *testing.T) string {
	testdataSrc := filepath.Join("testdata", "vcs")
	testdataTarget := filepath.Join("testdata", "tmp")
	err := fileutils.CopyDir(testdataSrc, testdataTarget, true, nil)
	if err != nil {
		t.Error(err)
	}
	if found, err := fileutils.IsDirExists(filepath.Join(testdataTarget, "gitdata"), false); found {
		if err != nil {
			t.Error(err)
		}
		err := fileutils.RenamePath(filepath.Join(testdataTarget, "gitdata"), filepath.Join(testdataTarget, ".git"))
		if err != nil {
			t.Error(err)
		}
	}
	if found, err := fileutils.IsDirExists(filepath.Join(testdataTarget, "OtherGit", "gitdata"), false); found {
		if err != nil {
			t.Error(err)
		}
		err := fileutils.RenamePath(filepath.Join(testdataTarget, "OtherGit", "gitdata"), filepath.Join(testdataTarget, "OtherGit", ".git"))
		if err != nil {
			t.Error(err)
		}
	}
	path, err := filepath.Abs(testdataTarget)
	if err != nil {
		t.Error(err)
	}
	return path
}
