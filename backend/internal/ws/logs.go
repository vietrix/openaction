package ws

import (
	"bufio"
	"context"
	"io"
	"net/http"
	"time"

	"nhooyr.io/websocket"
)

type LogStreamer interface {
	OpenLog(ctx context.Context, logPath string) (io.ReadCloser, error)
}

type LogHandler struct {
	Open func(ctx context.Context, logPath string) (io.ReadCloser, error)
}

func (h *LogHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logPath := r.URL.Query().Get("path")
	if logPath == "" {
		http.Error(w, "missing log path", http.StatusBadRequest)
		return
	}

	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})
	if err != nil {
		return
	}
	defer conn.Close(websocket.StatusNormalClosure, "")

	reader, err := h.Open(r.Context(), logPath)
	if err != nil {
		_ = conn.Close(websocket.StatusInternalError, "failed to open log")
		return
	}
	defer reader.Close()

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		_ = conn.Write(ctx, websocket.MessageText, []byte(line))
		cancel()
		time.Sleep(30 * time.Millisecond)
	}
}
