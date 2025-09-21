package filename

import "path/filepath"

func MyName(f string) string {
	ext := filepath.Ext(f)
	name := f[:len(f)-len(ext)]
	return name
}
