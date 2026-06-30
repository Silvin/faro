// Package uploads recibe imágenes (multipart) y las guarda en disco local.
// Las imágenes se sirven estáticamente desde /files/ (ver server).
// Nota: en producción (filesystem efímero) migrar a object storage (S3/R2).
package uploads

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
)

const maxUploadBytes = 10 << 20 // 10 MB

var allowedExt = map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true}

// Handler gestiona la subida de archivos a un directorio local.
type Handler struct {
	dir            string
	requireSession func(http.Handler) http.Handler
}

func New(dir string, requireSession func(http.Handler) http.Handler) *Handler {
	return &Handler{dir: dir, requireSession: requireSession}
}

// Routes se monta en /uploads (POST). Requiere sesión.
func (h *Handler) Routes() http.Handler {
	r := chi.NewRouter()
	r.Use(h.requireSession)
	r.Post("/", h.handleUpload)
	return r
}

func (h *Handler) handleUpload(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadBytes)
	if err := r.ParseMultipartForm(maxUploadBytes); err != nil {
		writeError(w, http.StatusBadRequest, "too_large", "Archivo inválido o demasiado grande (máx 10MB)")
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		writeError(w, http.StatusBadRequest, "validation_error", "Falta el archivo (campo 'file')")
		return
	}
	defer file.Close()

	if ct := header.Header.Get("Content-Type"); !strings.HasPrefix(ct, "image/") {
		writeError(w, http.StatusBadRequest, "invalid_type", "Solo se permiten imágenes")
		return
	}
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !allowedExt[ext] {
		writeError(w, http.StatusBadRequest, "invalid_type", "Formato no permitido (jpg, png, webp)")
		return
	}

	name := randomHex(16) + ext
	dst, err := os.Create(filepath.Join(h.dir, name))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal", "No se pudo guardar la imagen")
		return
	}
	defer dst.Close()
	if _, err := io.Copy(dst, file); err != nil {
		writeError(w, http.StatusInternalServerError, "internal", "No se pudo guardar la imagen")
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{"url": "/files/" + name})
}

func randomHex(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, map[string]string{"code": code, "message": message})
}
