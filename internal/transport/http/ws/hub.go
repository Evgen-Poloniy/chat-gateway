package ws

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"

	"github.com/Evgen-Poloniy/chat-gateway/internal/config"
	"github.com/Evgen-Poloniy/chat-gateway/internal/dto"
	"github.com/Evgen-Poloniy/chat-gateway/internal/entity"
	errs "github.com/Evgen-Poloniy/chat-gateway/pkg/errors"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

// ResolverService represents interface for work with user ids resolve
type ResolverService interface {
	// SubscribeOnEventChannel subscribes on broker channel once fro all time of work application.
	SubscribeOnEventChannel(ctx context.Context)

	// ResolveEvent returns event from channel.
	ResolveEvent(ctx context.Context) (*entity.EventMessage, error)
}

var _serverAddress string

// wsHub manages WebSocket connections
type Hub struct {
	resolver      ResolverService
	dispatchConf  *config.DispatchConfig
	upgrader      websocket.Upgrader
	serverAddress string
	mu            sync.RWMutex
	conns         map[uuid.UUID]map[*websocket.Conn]*Client
	logger        *logrus.Logger

	ctx    context.Context
	cancel context.CancelFunc
}

func NewHub(
	resolver ResolverService,
	dispatchConf *config.DispatchConfig,
	serverAddress string,
	logger *logrus.Logger,
) *Hub {
	ctx, cancel := context.WithCancel(context.Background())

	return &Hub{
		resolver:     resolver,
		dispatchConf: dispatchConf,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		serverAddress: serverAddress,
		conns:         make(map[uuid.UUID]map[*websocket.Conn]*Client),
		logger:        logger,
		ctx:           ctx,
		cancel:        cancel,
	}
}

// register is used only in WebSocketUpgrade for add new ws connection with user
func (h *Hub) register(client *Client) {
	h.mu.Lock()
	if _, exists := h.conns[client.userID]; !exists {
		h.conns[client.userID] = make(map[*websocket.Conn]*Client)
	}
	h.conns[client.userID][client.conn] = client
	h.mu.Unlock()
}

// unregister is used only in WebSocketUpgrade for close ws connection user with user
func (h *Hub) unregister(client *Client) error {
	h.mu.Lock()
	if userConns, exists := h.conns[client.userID]; exists {
		if _, connExists := userConns[client.conn]; connExists {
			delete(userConns, client.conn)

			if len(userConns) == 0 {
				delete(h.conns, client.userID)
			}

			close(client.send)
		}
	}
	h.mu.Unlock()

	if err := client.conn.Close(); err != nil {
		if errors.Is(err, net.ErrClosed) {
			return nil
		}
		return err
	}
	return nil
}

func (h *Hub) StartDispatchMessage() {
	h.resolver.SubscribeOnEventChannel(h.ctx)

	for {
		event, err := h.resolver.ResolveEvent(h.ctx)
		if err != nil {
			if h.ctx.Err() != nil {
				return
			}

			logError(h.logger, uuid.NewString(), h.serverAddress, err)

			if appErr, exists := errors.AsType[*errs.AppError](err); exists {
				if appErr.Code == errs.CodeFailedEventChannel {
					return
				}
			}

			continue
		}

		message := dto.DispatchMessage{
			MessageID: event.MessageID,
			ChatID:    event.ChatID,
			SenderID:  event.SenderID,
			Text:      event.Text,
			CreatedAt: event.CreatedAt,
		}

		h.broadcastMessage(event.UserIDs, message)
	}
}

// broadcastMessage sends message to personal client channels.
func (h *Hub) broadcastMessage(userIDs uuid.UUIDs, message dto.DispatchMessage) {
	h.mu.RLock()

	var targets []*Client
	for _, id := range userIDs {
		if clients, exists := h.conns[id]; exists {
			for _, client := range clients {
				targets = append(targets, client)
			}
		}
	}
	h.mu.RUnlock()

	for _, client := range targets {
		select {
		case client.send <- message:
		default:
			id := uuid.NewString()

			logWarn(h.logger, id,
				h.serverAddress,
				fmt.Sprintf("the client with user_id '%s' has been disconnected because of bad network", client.userID.String()),
			)

			if err := client.conn.Close(); err != nil {
				logError(h.logger, id, h.serverAddress, err)
			}
		}
	}
}

// writeMessage writes (sends) message to client.
func (h *Hub) writeMessage(client *Client) {
	for {
		select {
		case <-h.ctx.Done():
			return
		case message, ok := <-client.send:
			if !ok {
				return
			}

			if err := client.conn.WriteJSON(message); err != nil {
				logError(h.logger, uuid.NewString(), h.serverAddress, fmt.Errorf("websocket error: %v", err))
				return
			}
		}
	}
}

// readMessage reads message from client.
func (h *Hub) readMessage(client *Client) {
	// // Настраиваем ограничения для защиты от DOS-атак
	// c.conn.SetReadLimit(512 * 1024) // Ограничение на размер сообщения (например, 512 КБ)

	// // Настраиваем таймаут на чтение (Pong Timeout).
	// // Если клиент не ответит на наш Ping в течение этого времени, сокет закроется.
	// c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	// c.conn.SetPongHandler(func(string) error {
	// 	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	// 	return nil
	// })

	for {
		messageType, payload, err := client.conn.ReadMessage()
		if err != nil {
			var closeCode int
			var code string
			var message string
			var closeErr *websocket.CloseError
			if errors.As(err, &closeErr) {
				closeCode = closeErr.Code
			}

			switch closeCode {
			case websocket.CloseNormalClosure:
				closeCode = websocket.CloseNormalClosure
				message = "websocket connection closed normally"
			case websocket.CloseGoingAway:
				closeCode = websocket.CloseGoingAway
				message = "websocket connection going away"
			case websocket.CloseAbnormalClosure:
				closeCode = websocket.CloseAbnormalClosure
				code = "websocket_connection_unexpectedly_closed"
				message = fmt.Sprintf("websocket connection closed unexpectedly: %v", err)
			case websocket.CloseInternalServerErr:
				closeCode = websocket.CloseInternalServerErr
				code = "websocket_internal_server_error"
				message = fmt.Sprintf("websocket connection error: %v", err)
			default:
				closeCode = websocket.CloseInternalServerErr
				code = "websocket_unknown_error"
				message = fmt.Sprintf("websocket connection error: %v", err)
			}

			switch closeCode {
			case websocket.CloseNormalClosure, websocket.CloseGoingAway:
				h.logger.WithFields(map[string]interface{}{
					"id":         uuid.NewString(),
					"address":    h.serverAddress,
					"code":       "normal_closure",
					"close_code": closeCode,
				}).Info(message)

			default:
				h.logger.WithFields(map[string]interface{}{
					"id":         uuid.NewString(),
					"address":    h.serverAddress,
					"close_code": closeCode,
					"code":       code,
				}).Error(message)
			}

			return
		}

		if messageType != websocket.TextMessage {
			continue
		}

		logInfo(h.logger, uuid.NewString(), h.serverAddress, string(payload))
	}
}

func (h *Hub) ShutdownDispatchMessage(ctx context.Context) {
	h.cancel()
}
