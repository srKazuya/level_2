// Package namegen provides utilities to generate filenames from URLs.
package namegen

import "strings"


func FileName(f string) string {
	name := strings.TrimPrefix(f, "https://")

	name = strings.ReplaceAll(name, "/", "_")

	if idx := strings.Index(name, "?"); idx != -1 {
		name = name[:idx]
	}
	if strings.HasSuffix(name, "_") || name == "" {
		name += "index.html"
	}
	
	return name
}
