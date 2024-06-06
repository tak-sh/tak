package fileutils

import (
	"os"
	"path/filepath"
)

// FindUpwardFrom recursively searches upward in the filesystem from start
// until it finds a file/directory with the target name or reaches end.
//
// If start is not specified, it's set to the cwd.
//
// If end is not specified, it's set to the root.
//
// Returns the absolute path to the found directory. If the directory could
// not be found, an empty string is returned.
func FindUpwardFrom(name, start, end string) string {
	if start == "" {
		start, _ = os.Getwd()
	}

	if end == "" {
		end = string(os.PathSeparator)
	}

	if end == start {
		return ""
	}

	if filepath.Ext(start) != "" {
		return FindUpwardFrom(name, filepath.Dir(start), end)
	}
	entries, _ := os.ReadDir(start)
	if len(entries) > 0 {
		for _, v := range entries {
			if v.Name() == name {
				return filepath.Join(start, v.Name())
			}
		}
	}

	return FindUpwardFrom(name, filepath.Dir(start), end)
}
