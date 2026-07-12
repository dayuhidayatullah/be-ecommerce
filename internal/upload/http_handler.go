package upload

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"be-ecommerce/internal/utils"
)

// UploadImage handles file uploads
func UploadImage(w http.ResponseWriter, r *http.Request) {
	// Parse our multipart form, 10 << 20 specifies a maximum upload of 10 MB files.
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "File size too large or invalid form", err.Error())
		return
	}

	// FormFile returns the first file for the given key `image`
	file, handler, err := r.FormFile("image")
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Error retrieving the file", err.Error())
		return
	}
	defer file.Close()

	// Only allow certain extensions
	ext := strings.ToLower(filepath.Ext(handler.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		utils.RespondError(w, http.StatusBadRequest, "Invalid file format. Only JPG, JPEG, and PNG are allowed.", "invalid_format")
		return
	}

	// Create a unique filename
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	
	// Ensure uploads directory exists
	uploadDir := "./uploads"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.Mkdir(uploadDir, 0755)
	}

	filePath := filepath.Join(uploadDir, filename)

	// Create a temporary file within our uploads directory
	dst, err := os.Create(filePath)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to save file", err.Error())
		return
	}
	defer dst.Close()

	// Copy the uploaded file to the created file on the filesystem
	if _, err := io.Copy(dst, file); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to write file", err.Error())
		return
	}

	// Construct URL (assuming server runs on localhost:8080 or based on request Host)
	// For production, you'd use a config variable for the domain.
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	fileURL := fmt.Sprintf("%s://%s/uploads/%s", scheme, r.Host, filename)

	resp := map[string]string{
		"url":       fileURL,
		"file_name": filename,
	}

	utils.RespondSuccess(w, http.StatusOK, "File uploaded successfully", resp)
}
