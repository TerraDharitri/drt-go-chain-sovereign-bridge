package server

import (
	"bytes"
	"context"
	"encoding/hex"
	"net"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/TerraDharitri/drt-go-chain-core/data/sovereign"
	"github.com/TerraDharitri/drt-go-chain-sovereign-bridge/testscommon"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

func TestNewServerHandler(t *testing.T) {
	t.Parallel()

	t.Run("nil gin handler", func(t *testing.T) {
		handler, err := NewServerHandler(nil, grpc.NewServer())
		require.Equal(t, errNilGinHandler, err)
		require.Nil(t, handler)
	})
	t.Run("nil grpc handler", func(t *testing.T) {
		handler, err := NewServerHandler(gin.New(), nil)
		require.Equal(t, errNilGRPCHandler, err)
		require.Nil(t, handler)
	})
	t.Run("should work", func(t *testing.T) {
		handler, err := NewServerHandler(gin.New(), grpc.NewServer())
		require.Nil(t, err)
		require.NotNil(t, handler)
	})
}

func TestServerRequestsHandler_ServeHTTP(t *testing.T) {
	t.Parallel()

	// Mock Gin handler
	ginHandler := gin.New()
	ginHandler.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello from HTTP")
	})

	// Mock grpc handler
	wasSendCalled := false
	expectedResponse := &sovereign.BridgeOperationsResponse{
		TxHashes: []string{"hash1", "hash2", "hash3"}, // Sample data for TxHashes field
	}
	grpcServer := grpc.NewServer()
	mockGRPCServer := &testscommon.MockBridgeTxSenderServer{
		SendCalled: func(ctx context.Context, req *sovereign.BridgeOperations) (*sovereign.BridgeOperationsResponse, error) {
			wasSendCalled = true
			return expectedResponse, nil
		},
	}
	sovereign.RegisterBridgeTxSenderServer(grpcServer, mockGRPCServer)

	buffer := bufconn.Listen(1024 * 1024)
	go func() {
		if err := grpcServer.Serve(buffer); err != nil {
			require.Fail(t, "server exited with error", "error", err)
		}
	}()

	handler, err := NewServerHandler(ginHandler, grpcServer)
	require.Nil(t, err)

	// Test HTTP request
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Content-Type", "text/plain")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)
	require.Contains(t, w.Body.String(), "Hello from HTTP")

	// Test GRPC request
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		ctx := context.Background()
		dialOptWithCtx := grpc.WithContextDialer(
			func(context.Context, string) (net.Conn, error) {
				bytesReq, err := hex.DecodeString("00000000500a4e0a06686173685f3412140a076f70486173683112096272696467654f703112140a076f70486173683212096272696467654f70321a0d6167677265676174656453696722096c6561646572536967")
				require.Nil(t, err)

				grpcReq, _ := http.NewRequest("POST", "/sovereign.BridgeTxSender/Send", bytes.NewReader(bytesReq))
				grpcReq.Header.Set("Content-Type", "application/grpc")
				grpcReq.ProtoMajor = 2
				grpcReq.Proto = "HTTP/2.0"
				grpcW := httptest.NewRecorder()
				handler.ServeHTTP(grpcW, grpcReq)
				require.Equal(t, http.StatusOK, grpcW.Code)

				return buffer.Dial()
			},
		)
		dialOptCredentials := grpc.WithTransportCredentials(insecure.NewCredentials())
		clientConn, err := grpc.DialContext(ctx, "bufnet", dialOptWithCtx, dialOptCredentials)
		if err != nil {
			require.Fail(t, "client failed to dial", "error", err)
		}

		defer func() {
			err = clientConn.Close()
			wg.Done()
			require.Nil(t, err)
		}()

		client := sovereign.NewBridgeTxSenderClient(clientConn)
		resp, err := client.Send(ctx, &sovereign.BridgeOperations{})
		require.Nil(t, err)
		require.Equal(t, expectedResponse.TxHashes, resp.TxHashes)
	}()

	wg.Wait()
	grpcServer.Stop()
	require.True(t, wasSendCalled)
}
