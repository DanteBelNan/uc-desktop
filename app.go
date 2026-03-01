package main

import (
    "context"
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
    ctx      context.Context
    conn     *websocket.Conn
    mu       sync.Mutex
    roomID   string
    clientID string
    isJoined bool
}

// NewApp crea una nueva instancia de la aplicacion de escritorio.
func NewApp() *App {
    return &App{}
}

// startup se ejecuta cuando Wails inicia la aplicacion e inicializa el clipboard nativo.
func (a *App) startup(ctx context.Context) {
    a.ctx = ctx
    // Inicializar el acceso nativo al portapapeles (necesario en Linux/Cgo)
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

    // Mensaje de union
    joinMsg := Message{Type: "JOIN", Payload: roomID, Sender: clientID}
    if err := a.conn.WriteJSON(joinMsg); err != nil {
        a.conn.Close()
        a.conn = nil
        return "Error al unirse"
    }

    // Iniciar servicios asincronos (basados en eventos, no polling)
    go a.listenForMessages()
    go a.watchLocalClipboard()

    return "OK"
}

// watchLocalClipboard utiliza un canal nativo para detectar cambios en el portapapeles.
// Esto elimina el "polling" y el parpadeo de foco en sistemas Linux.
func (a *App) watchLocalClipboard() {
    // Escuchar solo cambios de texto plano (UTF-8)
    ch := clipboard.Watch(context.Background(), clipboard.FmtText)
    
    for data := range ch {
        a.mu.Lock()
        if !a.isJoined || a.conn == nil {
            a.mu.Unlock()
            break
        }
        
        content := string(data)
        if content != "" {
            log.Printf("Cambio detectado nativamente. Enviando...")
            updateMsg := Message{
                Type:    "UPDATE",
                Payload: content,
                Sender:  a.clientID,
            }
            err := a.conn.WriteJSON(updateMsg)
            if err != nil {
                log.Printf("Error al enviar actualizacion: %v", err)
            }
        }
        a.mu.Unlock()
    }
}

// listenForMessages escucha continuamente actualizaciones del servidor central.
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
            a.mu.Lock()
            a.conn = nil
            a.isJoined = false
            a.mu.Unlock()
            runtime.EventsEmit(a.ctx, "room_disconnected", "Conexion perdida")
            break
        }

        if msg.Type == "UPDATE" {
            // Actualizar localmente SIN disparar el Watcher (evitar bucles)
            clipboard.Write(clipboard.FmtText, []byte(msg.Payload))
            runtime.EventsEmit(a.ctx, "clipboard_update", msg.Payload)
        }
    }
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
