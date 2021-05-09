package www

import (
	"io"
	"mime"
	"net/http"
	"path/filepath"
)

func serveStatic(fs http.FileSystem, path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		f, err := fs.Open(path)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		defer func() { _ = f.Close() }()

		ctype := mime.TypeByExtension(filepath.Ext(path))
		if ctype != "" {
			w.Header().Set("Content-Type", ctype)
		}

		_, _ = io.Copy(w, f)
	}
}
