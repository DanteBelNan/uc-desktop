<script lang="ts">
    import logo from './assets/images/logo-universal.png'
    
    // Variables de estado para la conexion
    let roomID: string = "";
    let clientID: string = "";
    let connectionStatus: string = "Desconectado";
    let isConnected: boolean = false;

    // Funcion para intentar unirse a una sala
    // Por ahora solo cambia el estado visual, la logica de Go se integrara en el siguiente paso.
    function handleJoinRoom(): void {
        if (!roomID || !clientID) {
            alert("Por favor completa ambos campos para continuar.");
            return;
        }
        
        // Simulacion de inicio de conexion
        connectionStatus = "Conectando...";
        
        // TODO: Invocar funcion de Go para establecer el WebSocket
        console.log(`Intentando unirse a sala: ${roomID} con ID: ${clientID}`);
        
        // Simulacion de exito (temporal para v0.1.x)
        setTimeout(() => {
            isConnected = true;
            connectionStatus = `Conectado a la sala: ${roomID}`;
        }, 500);
    }

    // Funcion para desconectarse
    function handleDisconnect(): void {
        isConnected = false;
        connectionStatus = "Desconectado";
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
            <p>Sincronizando portapapeles en tiempo real...</p>
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
    }

    #logo {
        width: 120px;
        margin-bottom: 1.5rem;
    }

    .status-badge {
        padding: 0.5rem 1rem;
        border-radius: 20px;
        font-size: 0.9rem;
        margin-bottom: 2rem;
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
        gap: 0.5rem;
        text-align: left;
    }

    .input-group label {
        font-size: 0.85rem;
        color: #666;
    }

    input {
        padding: 0.8rem;
        border: 1px solid #ddd;
        border-radius: 6px;
        font-size: 1rem;
    }

    .btn-primary {
        background-color: #4a90e2;
        color: white;
        border: none;
        padding: 0.8rem;
        border-radius: 6px;
        cursor: pointer;
        font-weight: bold;
        transition: background 0.2s;
    }

    .btn-primary:hover {
        background-color: #357abd;
    }

    .btn-danger {
        background-color: #e74c3c;
        color: white;
        border: none;
        padding: 0.8rem;
        border-radius: 6px;
        cursor: pointer;
        font-weight: bold;
    }

    .active-session {
        text-align: center;
        margin-top: 2rem;
    }
</style>
