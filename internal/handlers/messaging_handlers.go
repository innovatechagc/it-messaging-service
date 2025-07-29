package handlers

import (
	"net/http"
	"strconv"

	"github.com/company/microservice-template/internal/auth"
	"github.com/company/microservice-template/internal/domain"
	"github.com/company/microservice-template/internal/services"
	"github.com/company/microservice-template/pkg/logger"
	"github.com/gin-gonic/gin"
)

type MessagingHandler struct {
	messagingService services.MessagingService
	fileService      services.FileService
	jwtManager       *auth.JWTManager
	logger           logger.Logger
}

func NewMessagingHandler(
	messagingService services.MessagingService,
	fileService services.FileService,
	jwtManager *auth.JWTManager,
	logger logger.Logger,
) *MessagingHandler {
	return &MessagingHandler{
		messagingService: messagingService,
		fileService:      fileService,
		jwtManager:       jwtManager,
		logger:           logger,
	}
}

// GetConversations godoc
// @Summary Lista conversaciones activas
// @Description Obtiene las conversaciones del usuario con filtros opcionales
// @Tags conversations
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param channel query string false "Canal de comunicación" Enums(whatsapp, web, messenger, instagram)
// @Param status query string false "Estado de la conversación" Enums(active, closed, archived)
// @Param limit query int false "Límite de resultados" default(20)
// @Param offset query int false "Offset para paginación" default(0)
// @Success 200 {object} domain.APIResponse{data=[]domain.Conversation}
// @Failure 401 {object} domain.APIResponse
// @Failure 500 {object} domain.APIResponse
// @Router /conversations [get]
func (h *MessagingHandler) GetConversations(c *gin.Context) {
	userID := h.getUserIDFromContext(c)
	if userID == "" {
		h.respondWithError(c, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid token")
		return
	}

	// Parse filters
	filters := domain.ConversationFilters{
		Channel: domain.Channel(c.Query("channel")),
		Status:  domain.ConversationStatus(c.Query("status")),
		Limit:   h.parseIntQuery(c, "limit", 20),
		Offset:  h.parseIntQuery(c, "offset", 0),
	}

	conversations, err := h.messagingService.GetConversations(c.Request.Context(), userID, filters)
	if err != nil {
		h.logger.Error("Failed to get conversations", err)
		h.respondWithError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get conversations")
		return
	}

	h.respondWithSuccess(c, http.StatusOK, "Conversations retrieved successfully", conversations)
}

// GetConversation godoc
// @Summary Obtiene detalles de una conversación
// @Description Trae detalles de una conversación incluyendo mensajes
// @Tags conversations
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path string true "ID de la conversación"
// @Success 200 {object} domain.APIResponse{data=domain.Conversation}
// @Failure 401 {object} domain.APIResponse
// @Failure 404 {object} domain.APIResponse
// @Failure 500 {object} domain.APIResponse
// @Router /conversations/{id} [get]
func (h *MessagingHandler) GetConversation(c *gin.Context) {
	userID := h.getUserIDFromContext(c)
	if userID == "" {
		h.respondWithError(c, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid token")
		return
	}

	conversationID := c.Param("id")
	if conversationID == "" {
		h.respondWithError(c, http.StatusBadRequest, "INVALID_REQUEST", "Conversation ID is required")
		return
	}

	conversation, err := h.messagingService.GetConversation(c.Request.Context(), conversationID, userID)
	if err != nil {
		h.logger.Error("Failed to get conversation", err)
		h.respondWithError(c, http.StatusNotFound, "NOT_FOUND", "Conversation not found")
		return
	}

	h.respondWithSuccess(c, http.StatusOK, "Conversation retrieved successfully", conversation)
}

// CreateConversation godoc
// @Summary Crea una nueva conversación
// @Description Crea una nueva conversación si no existe
// @Tags conversations
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param request body CreateConversationRequest true "Datos de la conversación"
// @Success 201 {object} domain.APIResponse{data=domain.Conversation}
// @Failure 400 {object} domain.APIResponse
// @Failure 401 {object} domain.APIResponse
// @Failure 500 {object} domain.APIResponse
// @Router /conversations [post]
func (h *MessagingHandler) CreateConversation(c *gin.Context) {
	userID := h.getUserIDFromContext(c)
	if userID == "" {
		h.respondWithError(c, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid token")
		return
	}

	var req CreateConversationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.respondWithError(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	conversation, err := h.messagingService.CreateConversation(c.Request.Context(), userID, req.Channel)
	if err != nil {
		h.logger.Error("Failed to create conversation", err)
		h.respondWithError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to create conversation")
		return
	}

	h.respondWithSuccess(c, http.StatusCreated, "Conversation created successfully", conversation)
}

// UpdateConversation godoc
// @Summary Actualiza estado de conversación
// @Description Actualiza el estado de una conversación (ej: cerrar conversación)
// @Tags conversations
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path string true "ID de la conversación"
// @Param request body UpdateConversationRequest true "Nuevo estado"
// @Success 200 {object} domain.APIResponse
// @Failure 400 {object} domain.APIResponse
// @Failure 401 {object} domain.APIResponse
// @Failure 404 {object} domain.APIResponse
// @Failure 500 {object} domain.APIResponse
// @Router /conversations/{id} [patch]
func (h *MessagingHandler) UpdateConversation(c *gin.Context) {
	userID := h.getUserIDFromContext(c)
	if userID == "" {
		h.respondWithError(c, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid token")
		return
	}

	conversationID := c.Param("id")
	if conversationID == "" {
		h.respondWithError(c, http.StatusBadRequest, "INVALID_REQUEST", "Conversation ID is required")
		return
	}

	var req UpdateConversationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.respondWithError(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	err := h.messagingService.UpdateConversationStatus(c.Request.Context(), conversationID, req.Status, userID)
	if err != nil {
		h.logger.Error("Failed to update conversation", err)
		h.respondWithError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to update conversation")
		return
	}

	h.respondWithSuccess(c, http.StatusOK, "Conversation updated successfully", nil)
}

// GetMessages godoc
// @Summary Lista mensajes de una conversación
// @Description Lista los mensajes de una conversación con paginación
// @Tags messages
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path string true "ID de la conversación"
// @Param limit query int false "Límite de resultados" default(50)
// @Param offset query int false "Offset para paginación" default(0)
// @Success 200 {object} domain.APIResponse{data=[]domain.Message}
// @Failure 401 {object} domain.APIResponse
// @Failure 404 {object} domain.APIResponse
// @Failure 500 {object} domain.APIResponse
// @Router /conversations/{id}/messages [get]
func (h *MessagingHandler) GetMessages(c *gin.Context) {
	userID := h.getUserIDFromContext(c)
	if userID == "" {
		h.respondWithError(c, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid token")
		return
	}

	conversationID := c.Param("id")
	if conversationID == "" {
		h.respondWithError(c, http.StatusBadRequest, "INVALID_REQUEST", "Conversation ID is required")
		return
	}

	pagination := domain.PaginationParams{
		Limit:  h.parseIntQuery(c, "limit", 50),
		Offset: h.parseIntQuery(c, "offset", 0),
		SortBy: "timestamp",
		Order:  "DESC",
	}

	messages, err := h.messagingService.GetMessages(c.Request.Context(), conversationID, userID, pagination)
	if err != nil {
		h.logger.Error("Failed to get messages", err)
		h.respondWithError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get messages")
		return
	}

	h.respondWithSuccess(c, http.StatusOK, "Messages retrieved successfully", messages)
}

// SendMessage godoc
// @Summary Envía un nuevo mensaje
// @Description Envía un nuevo mensaje (texto, archivo, IA, etc.)
// @Tags messages
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path string true "ID de la conversación"
// @Param request body services.SendMessageRequest true "Datos del mensaje"
// @Success 201 {object} domain.APIResponse{data=domain.Message}
// @Failure 400 {object} domain.APIResponse
// @Failure 401 {object} domain.APIResponse
// @Failure 500 {object} domain.APIResponse
// @Router /conversations/{id}/messages [post]
func (h *MessagingHandler) SendMessage(c *gin.Context) {
	userID := h.getUserIDFromContext(c)
	if userID == "" {
		h.respondWithError(c, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid token")
		return
	}

	conversationID := c.Param("id")
	if conversationID == "" {
		h.respondWithError(c, http.StatusBadRequest, "INVALID_REQUEST", "Conversation ID is required")
		return
	}

	var req services.SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.respondWithError(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	// Set conversation ID and sender ID from context
	req.ConversationID = conversationID
	req.SenderID = userID

	message, err := h.messagingService.SendMessage(c.Request.Context(), req)
	if err != nil {
		h.logger.Error("Failed to send message", err)
		h.respondWithError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to send message")
		return
	}

	h.respondWithSuccess(c, http.StatusCreated, "Message sent successfully", message)
}

// GetMessage godoc
// @Summary Consulta un mensaje individual
// @Description Obtiene los detalles de un mensaje específico
// @Tags messages
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path string true "ID del mensaje"
// @Success 200 {object} domain.APIResponse{data=domain.Message}
// @Failure 401 {object} domain.APIResponse
// @Failure 404 {object} domain.APIResponse
// @Failure 500 {object} domain.APIResponse
// @Router /messages/{id} [get]
func (h *MessagingHandler) GetMessage(c *gin.Context) {
	userID := h.getUserIDFromContext(c)
	if userID == "" {
		h.respondWithError(c, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid token")
		return
	}

	messageID := c.Param("id")
	if messageID == "" {
		h.respondWithError(c, http.StatusBadRequest, "INVALID_REQUEST", "Message ID is required")
		return
	}

	message, err := h.messagingService.GetMessage(c.Request.Context(), messageID, userID)
	if err != nil {
		h.logger.Error("Failed to get message", err)
		h.respondWithError(c, http.StatusNotFound, "NOT_FOUND", "Message not found")
		return
	}

	h.respondWithSuccess(c, http.StatusOK, "Message retrieved successfully", message)
}

// UploadAttachment godoc
// @Summary Sube un archivo adjunto
// @Description Sube un archivo y devuelve URL segura
// @Tags attachments
// @Accept multipart/form-data
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param file formData file true "Archivo a subir"
// @Success 200 {object} domain.APIResponse{data=UploadResponse}
// @Failure 400 {object} domain.APIResponse
// @Failure 401 {object} domain.APIResponse
// @Failure 500 {object} domain.APIResponse
// @Router /attachments/upload [post]
func (h *MessagingHandler) UploadAttachment(c *gin.Context) {
	userID := h.getUserIDFromContext(c)
	if userID == "" {
		h.respondWithError(c, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid token")
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		h.respondWithError(c, http.StatusBadRequest, "INVALID_REQUEST", "File is required")
		return
	}
	defer file.Close()

	uploadReq := services.UploadFileRequest{
		File:     file,
		Filename: header.Filename,
		Size:     header.Size,
		UserID:   userID,
	}

	result, err := h.fileService.UploadFile(c.Request.Context(), uploadReq)
	if err != nil {
		h.logger.Error("Failed to upload file", err)
		h.respondWithError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to upload file")
		return
	}

	response := UploadResponse{
		URL:      result.URL,
		Filename: result.Filename,
		Size:     result.Size,
		Type:     result.Type,
	}

	h.respondWithSuccess(c, http.StatusOK, "File uploaded successfully", response)
}

// GetAttachment godoc
// @Summary Obtiene detalles de un archivo adjunto
// @Description Devuelve los detalles de un archivo adjunto
// @Tags attachments
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path string true "ID del archivo adjunto"
// @Success 200 {object} domain.APIResponse{data=domain.Attachment}
// @Failure 401 {object} domain.APIResponse
// @Failure 404 {object} domain.APIResponse
// @Failure 500 {object} domain.APIResponse
// @Router /attachments/{id} [get]
func (h *MessagingHandler) GetAttachment(c *gin.Context) {
	userID := h.getUserIDFromContext(c)
	if userID == "" {
		h.respondWithError(c, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid token")
		return
	}

	attachmentID := c.Param("id")
	if attachmentID == "" {
		h.respondWithError(c, http.StatusBadRequest, "INVALID_REQUEST", "Attachment ID is required")
		return
	}

	attachment, err := h.messagingService.GetAttachment(c.Request.Context(), attachmentID, userID)
	if err != nil {
		h.logger.Error("Failed to get attachment", err)
		h.respondWithError(c, http.StatusNotFound, "NOT_FOUND", "Attachment not found")
		return
	}

	h.respondWithSuccess(c, http.StatusOK, "Attachment retrieved successfully", attachment)
}

// Helper methods

func (h *MessagingHandler) getUserIDFromContext(c *gin.Context) string {
	token, err := h.jwtManager.ExtractTokenFromHeader(c)
	if err != nil {
		return ""
	}

	claims, err := h.jwtManager.ValidateToken(token)
	if err != nil {
		return ""
	}

	return claims.UserID
}

func (h *MessagingHandler) parseIntQuery(c *gin.Context, key string, defaultValue int) int {
	if value := c.Query(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func (h *MessagingHandler) respondWithError(c *gin.Context, statusCode int, code, message string) {
	response := domain.APIResponse{
		Code:    code,
		Message: message,
		Data:    nil,
	}
	c.JSON(statusCode, response)
}

func (h *MessagingHandler) respondWithSuccess(c *gin.Context, statusCode int, message string, data interface{}) {
	response := domain.APIResponse{
		Code:    "SUCCESS",
		Message: message,
		Data:    data,
	}
	c.JSON(statusCode, response)
}

// Request/Response types

type CreateConversationRequest struct {
	Channel domain.Channel `json:"channel" binding:"required"`
}

type UpdateConversationRequest struct {
	Status domain.ConversationStatus `json:"status" binding:"required"`
}

type UploadResponse struct {
	URL      string                `json:"url"`
	Filename string                `json:"filename"`
	Size     int64                 `json:"size"`
	Type     domain.AttachmentType `json:"type"`
}