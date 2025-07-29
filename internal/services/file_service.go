package services

import (
	"context"
	"fmt"
	"io"
	"mime"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/company/microservice-template/internal/config"
	"github.com/company/microservice-template/internal/domain"
	"github.com/company/microservice-template/pkg/logger"
	"github.com/google/uuid"
)

type FileService interface {
	UploadFile(ctx context.Context, req UploadFileRequest) (*UploadFileResponse, error)
	DeleteFile(ctx context.Context, url string) error
	GetFileInfo(ctx context.Context, url string) (*FileInfo, error)
}

type UploadFileRequest struct {
	File     io.Reader
	Filename string
	Size     int64
	UserID   string
}

type UploadFileResponse struct {
	URL      string
	Filename string
	Size     int64
	Type     domain.AttachmentType
}

type FileInfo struct {
	URL      string
	Filename string
	Size     int64
	Type     domain.AttachmentType
	Exists   bool
}

type localFileService struct {
	config *config.FileStorageConfig
	logger logger.Logger
}

func NewLocalFileService(config *config.FileStorageConfig, logger logger.Logger) FileService {
	return &localFileService{
		config: config,
		logger: logger,
	}
}

func (s *localFileService) UploadFile(ctx context.Context, req UploadFileRequest) (*UploadFileResponse, error) {
	// Validate file size
	if req.Size > s.config.MaxFileSize {
		return nil, fmt.Errorf("file size exceeds maximum allowed size of %d bytes", s.config.MaxFileSize)
	}

	// Generate unique filename
	ext := filepath.Ext(req.Filename)
	uniqueFilename := fmt.Sprintf("%s_%s%s", uuid.New().String(), time.Now().Format("20060102_150405"), ext)
	
	// Create user directory
	userDir := filepath.Join(s.config.LocalPath, req.UserID)
	if err := os.MkdirAll(userDir, 0755); err != nil {
		s.logger.Error("Failed to create user directory", err)
		return nil, fmt.Errorf("failed to create upload directory: %w", err)
	}

	// Full file path
	filePath := filepath.Join(userDir, uniqueFilename)

	// Create and write file
	file, err := os.Create(filePath)
	if err != nil {
		s.logger.Error("Failed to create file", err)
		return nil, fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Copy file content
	written, err := io.Copy(file, req.File)
	if err != nil {
		s.logger.Error("Failed to write file content", err)
		// Clean up partial file
		os.Remove(filePath)
		return nil, fmt.Errorf("failed to write file: %w", err)
	}

	// Determine file type
	fileType := s.determineFileType(req.Filename)

	// Generate URL (relative path for local storage)
	url := fmt.Sprintf("/uploads/%s/%s", req.UserID, uniqueFilename)

	response := &UploadFileResponse{
		URL:      url,
		Filename: req.Filename,
		Size:     written,
		Type:     fileType,
	}

	s.logger.Info("File uploaded successfully", map[string]interface{}{
		"filename":     req.Filename,
		"size":         written,
		"type":         fileType,
		"user_id":      req.UserID,
		"unique_name":  uniqueFilename,
	})

	return response, nil
}

func (s *localFileService) DeleteFile(ctx context.Context, url string) error {
	// Convert URL to file path
	filePath := filepath.Join(s.config.LocalPath, strings.TrimPrefix(url, "/uploads/"))
	
	if err := os.Remove(filePath); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file not found")
		}
		s.logger.Error("Failed to delete file", err)
		return fmt.Errorf("failed to delete file: %w", err)
	}

	s.logger.Info("File deleted successfully", map[string]interface{}{
		"url": url,
	})

	return nil
}

func (s *localFileService) GetFileInfo(ctx context.Context, url string) (*FileInfo, error) {
	// Convert URL to file path
	filePath := filepath.Join(s.config.LocalPath, strings.TrimPrefix(url, "/uploads/"))
	
	stat, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return &FileInfo{
				URL:    url,
				Exists: false,
			}, nil
		}
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	// Extract filename from path
	filename := filepath.Base(filePath)
	
	// Determine file type
	fileType := s.determineFileType(filename)

	return &FileInfo{
		URL:      url,
		Filename: filename,
		Size:     stat.Size(),
		Type:     fileType,
		Exists:   true,
	}, nil
}

func (s *localFileService) determineFileType(filename string) domain.AttachmentType {
	ext := strings.ToLower(filepath.Ext(filename))
	mimeType := mime.TypeByExtension(ext)

	switch {
	case strings.HasPrefix(mimeType, "image/"):
		return domain.AttachmentTypeImage
	case strings.HasPrefix(mimeType, "video/"):
		return domain.AttachmentTypeVideo
	case strings.HasPrefix(mimeType, "audio/"):
		return domain.AttachmentTypeAudio
	default:
		return domain.AttachmentTypeFile
	}
}

// NoOpFileService for when file storage is disabled
type noOpFileService struct{}

func NewNoOpFileService() FileService {
	return &noOpFileService{}
}

func (s *noOpFileService) UploadFile(ctx context.Context, req UploadFileRequest) (*UploadFileResponse, error) {
	return nil, fmt.Errorf("file storage is disabled")
}

func (s *noOpFileService) DeleteFile(ctx context.Context, url string) error {
	return fmt.Errorf("file storage is disabled")
}

func (s *noOpFileService) GetFileInfo(ctx context.Context, url string) (*FileInfo, error) {
	return nil, fmt.Errorf("file storage is disabled")
}