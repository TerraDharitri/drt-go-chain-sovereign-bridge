package server

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

type serverRequestsHandler struct {
	ginHandler  *gin.Engine
	grpcHandler *grpc.Server
}

// NewServerHandler creates a server wrapper over a gin and grpc server. This wrapper can serve both http and grpc requests
// based on header content type
func NewServerHandler(ginHandler *gin.Engine, grpcHandler *grpc.Server) (*serverRequestsHandler, error) {
	if ginHandler == nil {
		return nil, errNilGinHandler
	}
	if grpcHandler == nil {
		return nil, errNilGRPCHandler
	}

	return &serverRequestsHandler{
		ginHandler:  ginHandler,
		grpcHandler: grpcHandler,
	}, nil
}

// ServeHTTP will server the http request(http or grpc)
func (h *serverRequestsHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	contentType := req.Header.Get("Content-Type")

	if strings.HasPrefix(contentType, "application/grpc") {
		log.Trace("server handling grpc request")
		h.grpcHandler.ServeHTTP(w, req)
		return
	}

	log.Trace("server handling http request")
	h.ginHandler.ServeHTTP(w, req)
}
