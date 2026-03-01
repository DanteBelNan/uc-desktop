package main

import (
    "context"
    "fmt"
    "log"
    "net/url"
    "sync"

    "github.com/gorilla/websocket"
    "github.com/wailsapp/wails/v2/pkg/runtime"
    "os"
)

// Message representa la estructura estandar de mensajes que compartimos con uc-server.
type Message struct {
    Type    string `json:"type"`
    Payload string `json:"payload"`
    Sender  string `json:"sender_id"`
}

// App representa la estructura principal de la aplicacion de escritorio.
type App struct {
    ctx      context.Context
    conn     *websocket.Conn
    mu       sync.Mutex
    roomID   string
    clientID string
}

// NewApp crea una nueva instancia de la aplicacion de escritorio.
func NewApp() *App {
    return &App{}
}

// startup se ejecuta cuando Wails inicia la aplicacion.
func (a *App) startup(ctx context.Context) {
    a.ctx = ctx
}

// JoinRoom establece una conexion WebSocket con el servidor y se une a una sala.
// roomID: ID de la sala a la que unirse.
// clientID: ID unico del dispositivo.
func (a *App) JoinRoom(roomID string, clientID string) string {
    a.mu.Lock()
    defer a.mu.Unlock()

    if a.conn != nil {
        return "Ya estas conectado a una sala."
    }

    // Obtener la URL del servidor desde el entorno o usar localhost por defecto
    serverHost := os.Getenv("UC_SERVER_URL")
    if serverHost == "" {
        serverHost = "localhost:8080"
    }

    u := url.URL{Scheme: "ws", Host: serverHost, Path: "/ws"}
    log.Printf("Conectando a %s", u.String())

    conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
    if err != nil {
        return fmt.Sprintf("Error al conectar: %v", err)
    }

    a.conn = conn
    a.roomID = roomID
    a.clientID = clientID

    // Enviar mensaje JOIN inicial
    joinMsg := Message{
        Type:    "JOIN",
        Payload: roomID,
        Sender:  clientID,
    }
    if err := a.conn.WriteJSON(joinMsg); err != nil {
        a.conn.Close()
        a.conn = nil
        return fmt.Sprintf("Error al unirse a la sala: %v", err)
    }

    // Iniciar la escucha de mensajes en segundo plano
    go a.listenForMessages()

    return "OK"
}

// listenForMessages escucha continuamente actualizaciones del servidor.
func (a *App) listenForMessages() {
    for {
        var msg Message
        err := a.conn.ReadJSON(&msg)
        if err != nil {
            log.Printf("Conexion cerrada o error: %v", err)
            a.mu.Lock()
            a.conn = nil
            a.mu.Unlock()
            // Notificar al frontend la desconexion
            runtime.EventsEmit(a.ctx, "room_disconnected", "Conexion perdida con el servidor")
            break
        }

        // Si es una actualizacion del portapapeles, emitir evento al frontend
        if msg.Type == "UPDATE" {
            log.Printf("Nueva actualizacion recibida de %s", msg.Sender)
            runtime.EventsEmit(a.ctx, "clipboard_update", msg.Payload)
        }
    }
}

// Disconnect cierra la conexion activa con el servidor.
func (a *App) Disconnect() {
    a.mu.Lock()
    defer a.mu.Unlock()
    if a.conn != nil {
        a.conn.Close()
        a.conn = nil
    }
}
