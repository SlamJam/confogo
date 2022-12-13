package main

type Container interface {
	GetValueAtPath(path Path) (string, error)
	// GetUnnesessaryPaths(paths []Path) []Path
}
