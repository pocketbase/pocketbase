package osutils

import (
	"os"
	"strconv"
	"testing"
)

func TestIsProbablyGoRun(t *testing.T) {
	scenarios := []struct {
		arg0     string
		runDirs  []string
		expected bool
	}{
		{"", nil, false},
		{"/a/b", nil, false},
		{"/a/b", []string{""}, false},
		{"/a/b", []string{"/b/"}, false},
		{"/a/b", []string{"/a/"}, true},
		{"/a/b", []string{"", "/b/", "/a/"}, true},
	}

	originalArgs := os.Args
	defer func() {
		os.Args = originalArgs
	}()

	originalRunDirs := runDirs
	defer func() {
		runDirs = originalRunDirs
	}()

	for i, s := range scenarios {
		t.Run(strconv.Itoa(i)+"_"+s.arg0, func(t *testing.T) {
			os.Args = []string{s.arg0}
			runDirs = s.runDirs

			result := IsProbablyGoRun()

			if result != s.expected {
				t.Fatalf("Expected %v, got %v", s.expected, result)
			}
		})
	}
}
