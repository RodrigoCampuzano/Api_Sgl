package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type FileHandler struct {
	uploadPath string
	maxSize    int64
}

func NewFileHandler(uploadPath string, maxSizeMB int) *FileHandler {
	return &FileHandler{
		uploadPath: uploadPath,
		maxSize:    int64(maxSizeMB) * 1024 * 1024,
	}
}

type UploadResponse struct {
	URL      string `json:"url"`
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
	Type     string `json:"type"`
}

// UploadFile godoc
// @Summary      Upload de archivo
// @Description  Sube archivos (imágenes, PDFs) al servidor
// @Tags         files
// @Accept       multipart/form-data
// @Produce      json
// @Param        file  formData  file  true  "Archivo a subir"
// @Param        type  formData  string  false  "Tipo: invoice, damage_photo, vehicle_photo"
// @Success      200   {object}  UploadResponse
// @Security     Bearer
// @Router       /api/v1/files/upload [post]
func (h *FileHandler) UploadFile(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Archivo no proporcionado"})
		return
	}
	defer file.Close()

	// Validar tamaño
	if header.Size > h.maxSize {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Archivo demasiado grande. Máximo %dMB", h.maxSize/(1024*1024)),
		})
		return
	}

	// Validar extensión
	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowed := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".pdf":  true,
		".xml":  true,
	}

	if !allowed[ext] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Tipo de archivo no permitido. Solo: JPG, PNG, PDF, XML",
		})
		return
	}

	// Crear directorio
	uploadType := c.PostForm("type")
	if uploadType == "" {
		uploadType = "general"
	}

	now := time.Now()
	yearMonth := now.Format("2006/01")
	targetDir := filepath.Join(h.uploadPath, yearMonth, uploadType)

	if err := os.MkdirAll(targetDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating directory"})
		return
	}

	// Generar nombre único
	uniqueFilename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	targetPath := filepath.Join(targetDir, uniqueFilename)

	// Guardar archivo
	out, err := os.Create(targetPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving file"})
		return
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error writing file"})
		return
	}

	// URL pública
	publicURL := fmt.Sprintf("/uploads/%s/%s/%s", yearMonth, uploadType, uniqueFilename)

	c.JSON(http.StatusOK, UploadResponse{
		URL:      publicURL,
		Filename: header.Filename,
		Size:     header.Size,
		Type:     header.Header.Get("Content-Type"),
	})
}

// ServeFile sirve archivos subidos
func (h *FileHandler) ServeFile(c *gin.Context) {
	filepath := c.Param("filepath")
	fullPath := filepath + filepath

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	c.File(fullPath)
}
