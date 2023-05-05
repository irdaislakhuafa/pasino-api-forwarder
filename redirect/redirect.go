package redirect

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Redirect interface {
	Redirect(ctx context.Context, w http.ResponseWriter, r *http.Request, path string, errHandlingResponse func(w http.ResponseWriter, err error))
	RedirectWebSocket(ctx context.Context, w http.ResponseWriter, r *http.Request)
}

type redirect struct {
	baseUrl      string
	webSocketUrl string
	client       *http.Client
}

func Init(ctx context.Context, baseUrl string, client *http.Client, webSocketUrl string) Redirect {
	return &redirect{
		baseUrl:      baseUrl,
		client:       client,
		webSocketUrl: webSocketUrl,
	}
}

func (redirect *redirect) Redirect(ctx context.Context, w http.ResponseWriter, r *http.Request, path string, errHandlingResponse func(w http.ResponseWriter, err error)) {
	log.Printf("INFO: %+v will redirected into %v", r.Host, (redirect.baseUrl + path))

	if errHandlingResponse == nil {
		errHandlingResponse = func(w http.ResponseWriter, err error) {
			log.Println(err)
		}
	}

	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		errHandlingResponse(w, err)
		return
	}

	request, err := http.NewRequest(r.Method, redirect.baseUrl+path, bytes.NewBuffer(reqBody))
	if err != nil {
		errHandlingResponse(w, err)
		return
	}

	response, err := redirect.client.Do(request)
	if err != nil {
		errHandlingResponse(w, err)
		return
	}

	resBody, err := io.ReadAll(response.Body)
	if err != nil {
		errHandlingResponse(w, err)
		return
	}

	w.Write(resBody)
	for k, v := range response.Header {
		head := ""
		for _, s := range v {
			head += s + ";"
		}
		w.Header().Set(k, head)
	}
	return
}

func (redirect *redirect) RedirectWebSocket(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	webSocketDialler := websocket.DefaultDialer
	headers := http.Header{}

	webSocketConnection, _, err := webSocketDialler.DialContext(ctx, redirect.webSocketUrl, headers)
	if err != nil {
		log.Fatalf("Failed to connect to WebSocket: %v", err)
	}
	defer webSocketConnection.Close()

	webSockerUpgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			// logic to allow origin here
			return true
		},
	}
	// Upgrade HTTP connection to WebSocket
	clientConnection, err := webSockerUpgrader.Upgrade(w, r, nil)
	if err != nil {
		message := fmt.Sprintf("Failed to upgrade HTTP connection to WebSockets: %v", err)
		log.Println(message)
		ReturnErrorResponse(w, message)
		return
	}
	defer clientConnection.Close()

	// Forwarding data from client to web socket server
	for {
		messageType, data, err := webSocketConnection.ReadMessage()
		if err != nil {
			message := fmt.Sprintf("Error reading from client: %v", err)
			log.Println(message)
			ReturnErrorResponse(w, message)
			return
		}

		// is message is close message
		if messageType == websocket.CloseMessage {
			message := "Receive close message"
			log.Println(message)

			err := webSocketConnection.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				message := fmt.Sprintf("Error writing close message: %v", err)
				log.Println(message)
				ReturnErrorResponse(w, message)
			}
			return
		} else {
			err = webSocketConnection.WriteMessage(messageType, data)
			if err != nil {
				message := fmt.Sprintf("Error writing to websocket server: %v", err)
				log.Println(message)
			}
			return
		}
	}

}

func ToStringJSON[T any](value T) string {
	jsonBytes, err := json.MarshalIndent(value, "", "\t")
	if err != nil {
		message := fmt.Sprintf("Error marshalling JSON: %+v", err)
		log.Println(message)
		return fmt.Sprintf(`{"error": "%+v"}`, message)
	}
	return string(jsonBytes)
}

func ReturnErrorResponse[T any](w http.ResponseWriter, message T) {
	w.Write([]byte(
		ToStringJSON(
			map[string]interface{}{
				"error": fmt.Sprintf("%+v", message),
			},
		),
	))
}
