package account

import (
	"github.com/Stardome-Team/Service-Template/pkg/logset"
	"github.com/gin-gonic/gin"
)

const (
	// AccountEndpoint is the endpoint called to register a new user
	AccountEndpoint = "/account"
)

type handler struct {
	service Service
	logger  logset.Logger
}

// AuthenticationRequest request model sent to authenticate players
type AuthenticationRequest struct {
	UserName string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Controller contains interface for authentication handlers
type Controller interface {
	account(c *gin.Context)
}

// CreateHandlers sets up routing to the HTTP request
func CreateHandlers(r *gin.RouterGroup, s Service, l logset.Logger) {
	h := &handler{service: s, logger: l}

	registerHandlers(r, h)
}

func registerHandlers(r *gin.RouterGroup, ctr Controller) {
	r.POST(AccountEndpoint, ctr.account)
}

func (h *handler) account(c *gin.Context) {

	var request AuthenticationRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		h.logger.Errorf("invalid request: %v", err)
		return
	}
}
