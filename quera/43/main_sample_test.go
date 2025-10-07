package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// WsHandler -
type WsHandler struct {
	handler echo.HandlerFunc
	path string
	paramNames []string
	paramValues []string
  }
  
func (h *WsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	e := echo.New()
	c := e.NewContext(r, w)
	c.SetPath(h.path)
	c.SetParamNames(h.paramNames...)
	c.SetParamValues(h.paramValues...)

	forever := make(chan struct{})
	h.handler(c)
	<-forever
  }

func TestOneRoomTwoUser(t *testing.T) {
    // Create test server with the echo handler.
	h := WsHandler{
		handler: wsChatRoom,
		path: "/ws/chat/:roomId/user/:username",
		paramNames: []string{"roomId", "username"},
	}
	server := httptest.NewServer(http.HandlerFunc(h.ServeHTTP))
	defer server.Close()

	url := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"

	h.paramValues = []string{"room1", "user1"}
    ws1, _, err := websocket.DefaultDialer.Dial(url, nil)
    assert.Nil(t, err, err)
    defer ws1.Close()

	h.paramValues = []string{"room1", "user2"}
	ws2, _, err := websocket.DefaultDialer.Dial(url, nil)
    assert.Nil(t, err, err)
    defer ws2.Close()

    // Send message to server, read response and check to see if it's what we expect.
	err = ws1.WriteMessage(websocket.TextMessage, []byte("hello"))
	assert.Nil(t, err, err)
	_, msg, err := ws2.ReadMessage()
	assert.Nil(t, err, err)
	assert.Equal(t, string(msg), "user1: hello")

	err = ws2.WriteMessage(websocket.TextMessage, []byte("hellooo"))
	assert.Nil(t, err, err)
	err = ws2.WriteMessage(websocket.TextMessage, []byte("how are you?"))
	assert.Nil(t, err, err)
	_, msg, err = ws1.ReadMessage()
	assert.Nil(t, err, err)
	assert.Equal(t, string(msg), "user2: hellooo")
    _, msg, err = ws1.ReadMessage()
	assert.Nil(t, err, err)
	assert.Equal(t, string(msg), "user2: how are you?")
}

