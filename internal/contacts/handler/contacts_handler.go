package http

import (
	"github.com/hardikm9850/GoChat/internal/contacts/service"
	"net/http"
	dto "github.com/hardikm9850/GoChat/internal/contacts/handler/dto"
	"github.com/gin-gonic/gin"
)

type ContactsHandler struct {
	contactService service.ContactService
}

func NewContactsHandler(contactService service.ContactService) *ContactsHandler {
	return &ContactsHandler{
		contactService: contactService,
	}
}

// @Summary Sync contacts
// @Description Syncs a list of phone numbers with the server
// @Tags Sync contacts
// @Accept json
// @Produce json
// @Param request body SyncContactsRequest true "Sync Contacts Request"
// @Success 200 {array} string "List of synced contacts"
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 500 {object} map[string]string "Failed to sync contacts"
// @Router /contacts/sync [post]
func (h *ContactsHandler) SyncContacts(c *gin.Context) {
	var req dto.SyncContactsRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request payload",
		})
		return
	}

	userID := c.GetString("user_id")

	resp, err := h.contactService.SyncContacts(userID, req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to sync contacts",
		})
		return
	}

	c.JSON(http.StatusOK, resp)
}
