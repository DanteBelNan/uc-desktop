package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "net/url"
    "os"
    "sync"

    "golang.design/x/clipboard"
    "github.com/gorilla/websocket"
    "github.com/wailsapp/wails/v2/pkg/runtime"
)

// Message representa la estructura estandar de mensajes que compartimos con uc-server.
type Message struct {
    Type    string `json:"type"`
    Payload string `json:"payload"`
    Sender  string `json:"sender_id"`
}

// App representa la estructura principal de la aplicacion de escritorio.
type App struct {
    ctx        context.Context
    conn       *websocket.Conn
    mu         sync.Mutex
    roomID     string
    clientID   string
    isJoined   bool
    autoSync   bool
}

// NewApp crea una nueva instancia de la aplicacion de escritorio.
func NewApp() *App {
    return &App{
        autoSync: true, // Sincronizacion automatica por defecto
    }
}

// startup se ejecuta cuando Wails inicia la aplicacion e inicializa el clipboard nativo.
func (a *App) startup(ctx context.Context) {
    a.ctx = ctx
    err := clipboard.Init()
    if err != nil {
        log.Fatalf("Fallo al inicializar el portapapeles: %v", err)
    }
}

// JoinRoom establece una conexion WebSocket con el servidor y se une a una sala.
func (a *App) JoinRoom(roomID string, clientID string) string {
    a.mu.Lock()
    defer a.mu.Unlock()

    if a.conn != nil {
        return "Ya estas conectado a una sala."
    }

    serverHost := os.Getenv("UC_SERVER_URL")
    if serverHost == "" {
        serverHost = "localhost:8080"
    }

    u := url.URL{Scheme: "ws", Host: serverHost, Path: "/ws"}
    conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
    if err != nil {
        return fmt.Sprintf("Error al conectar: %v", err)
    }

    a.conn = conn
    a.roomID = roomID
    a.clientID = clientID
    a.isJoined = true

    joinMsg := Message{Type: "JOIN", Payload: roomID, Sender: clientID}
    if err := a.conn.WriteJSON(joinMsg); err != nil {
        a.conn.Close()
        a.conn = nil
        return "Error al unirse"
    }

    go a.listenForMessages()
    go a.watchLocalClipboard()

    return "OK"
}

// watchLocalClipboard monitorea el portapapeles local.
func (a *App) watchLocalClipboard() {
    ch := clipboard.Watch(context.Background(), clipboard.FmtText)
    for data := range ch {
        a.mu.Lock()
        if !a.isJoined || a.conn == nil || !a.autoSync {
            a.mu.Unlock()
            if !a.isJoined { break }
            continue
        }
        
        content := string(data)
        if content != "" {
            a.sendMessage("UPDATE", content)
        }
        a.mu.Unlock()
    }
}

// listenForMessages escucha continuamente actualizaciones del servidor.
func (a *App) listenForMessages() {
    for {
        var msg Message
        a.mu.Lock()
        conn := a.conn
        if conn == nil {
            a.mu.Unlock()
            break
        }
        a.mu.Unlock()

        err := conn.ReadJSON(&msg)
        if err != nil {
            a.handleDisconnectEvent("Conexion perdida")
            break
        }

        switch msg.Type {
        case "UPDATE":
            if a.autoSync {
                clipboard.Write(clipboard.FmtText, []byte(msg.Payload))
            }
            runtime.EventsEmit(a.ctx, "clipboard_update", msg.Payload)
        case "USER_LIST":
            var users []string
            json.Unmarshal([]byte(msg.Payload), &users)
            runtime.EventsEmit(a.ctx, "user_list_updated", users)
        }
    }
}

// SetAutoSync activa o desactiva la sincronizacion automatica.
func (a *App) SetAutoSync(enabled bool) {
    a.mu.Lock()
    defer a.mu.Unlock()
    a.autoSync = enabled
}

// ManualPaste envia el contenido actual del portapapeles al servidor manualmente.
func (a *App) ManualPaste() {
    content := string(clipboard.Read(clipboard.FmtText))
    if content != "" {
        a.sendMessage("UPDATE", content)
    }
}

// ManualCopy copia un contenido especifico al portapapeles del sistema.
func (a *App) ManualCopy(content string) {
    clipboard.Write(clipboard.FmtText, []byte(content))
}

// sendMessage es una utilidad interna para enviar mensajes JSON.
func (a *App) sendMessage(msgType string, payload string) {
    // Nota: El llamador debe gestionar el lock si es necesario
    if a.conn != nil {
        msg := Message{Type: msgType, Payload: payload, Sender: a.clientID}
        a.conn.WriteJSON(msg)
    }
}

func (a *App) handleDisconnectEvent(reason string) {
    a.mu.Lock()
    a.conn = nil
    a.isJoined = false
    a.mu.Unlock()
    runtime.EventsEmit(a.ctx, "room_disconnected", reason)
}

// Disconnect limpia la sesion activa.
func (a *App) Disconnect() {
    a.mu.Lock()
    defer a.mu.Unlock()
    if a.conn != nil {
        a.conn.Close()
        a.conn = nil
    }
    a.isJoined = false
}
