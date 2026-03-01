<script lang="ts">
    import logo from './assets/images/logo-universal.png'
    import { JoinRoom, Disconnect } from '../wailsjs/go/main/App.js'
    import { EventsOn } from '../wailsjs/runtime/runtime.js'
    import { onMount } from 'svelte';
    
    // Variables de estado para la conexion
    let roomID: string = "";
    let clientID: string = "";
    let connectionStatus: string = "Desconectado";
    let isConnected: boolean = false;
    let lastClipboardContent: string = "(Esperando sincronizacion...)";

    // Configurar listeners de eventos al montar el componente
    onMount(() => {
        EventsOn("clipboard_update", (payload: string) => {
            lastClipboardContent = payload;
            console.log("Nueva actualizacion de portapapeles:", payload);
        });

        EventsOn("room_disconnected", (msg: string) => {
            isConnected = false;
            connectionStatus = "Desconectado: " + msg;
        });
    });

    // Funcion para intentar unirse a una sala real
    function handleJoinRoom(): void {
        if (!roomID || !clientID) {
            alert("Por favor completa ambos campos para continuar.");
            return;
        }
        
        connectionStatus = "Conectando...";
        
        JoinRoom(roomID, clientID).then(result => {
            if (result === "OK") {
                isConnected = true;
                connectionStatus = `Conectado a la sala: ${roomID}`;
            } else {
                connectionStatus = "Error: " + result;
                isConnected = false;
            }
        });
    }

    // Funcion para desconectarse
    function handleDisconnect(): void {
        Disconnect().then(() => {
            isConnected = false;
            connectionStatus = "Desconectado";
        });
    }
</script>

<main>
    <img alt="Universal Clipboard Logo" id="logo" src="{logo}">
    
    <div class="status-badge {isConnected ? 'connected' : 'disconnected'}">
        {connectionStatus}
    </div>

    {#if !isConnected}
        <div class="form-container">
            <div class="input-group">
                <label for="roomID">ID de la Sala</label>
                <input 
                    id="roomID" 
                    type="text" 
                    bind:value={roomID} 
                    placeholder="Ej: mi-sala-secreta"
                    autocomplete="off"
                />
            </div>

            <div class="input-group">
                <label for="clientID">ID del Dispositivo</label>
                <input 
                    id="clientID" 
                    type="text" 
                    bind:value={clientID} 
                    placeholder="Ej: Laptop-Trabajo"
                    autocomplete="off"
                />
            </div>

            <button class="btn-primary" on:click={handleJoinRoom}>
                Unirse a la Sala
            </button>
        </div>
    {:else}
        <div class="active-session">
            <h3>Sincronizacion Activa</h3>
            <div class="clipboard-preview">
                <strong>Ultimo contenido recibido:</strong>
                <p>{lastClipboardContent}</p>
            </div>
            <button class="btn-danger" on:click={handleDisconnect}>
                Abandonar Sala
            </button>
        </div>
    {/if}
</main>

<style>
    main {
        display: flex;
        flex-direction: column;
        align-items: center;
        padding: 2rem;
        font-family: 'Nunito', sans-serif;
        text-align: center;
    }

    #logo {
        width: 100px;
        margin-bottom: 1rem;
    }

    .status-badge {
        padding: 0.5rem 1rem;
        border-radius: 20px;
        font-size: 0.8rem;
        margin-bottom: 1.5rem;
        font-weight: bold;
    }

    .connected {
        background-color: #d4edda;
        color: #155724;
    }

    .disconnected {
        background-color: #f8d7da;
        color: #721c24;
    }

    .form-container {
        width: 100%;
        max-width: 300px;
        display: flex;
        flex-direction: column;
        gap: 1.2rem;
    }

    .input-group {
        display: flex;
        flex-direction: column;
        gap: 0.4rem;
        text-align: left;
    }

    .input-group label {
        font-size: 0.8rem;
        color: #666;
    }

    input {
        padding: 0.6rem;
        border: 1px solid #ddd;
        border-radius: 6px;
        font-size: 0.9rem;
    }

    .btn-primary {
        background-color: #4a90e2;
        color: white;
        border: none;
        padding: 0.7rem;
        border-radius: 6px;
        cursor: pointer;
        font-weight: bold;
    }

    .btn-danger {
        background-color: #e74c3c;
        color: white;
        border: none;
        padding: 0.7rem;
        border-radius: 6px;
        cursor: pointer;
        font-weight: bold;
        margin-top: 1rem;
    }

    .clipboard-preview {
        background: #f9f9f9;
        padding: 1rem;
        border-radius: 8px;
        border: 1px dashed #ccc;
        margin: 1rem 0;
        max-width: 400px;
        word-break: break-all;
    }

    h3 { margin-bottom: 0.5rem; color: #333; }
</style>
