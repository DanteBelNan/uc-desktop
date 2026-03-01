<script lang="ts">
    import logo from './assets/images/logo-universal.png'
    import { JoinRoom, Disconnect, SetAutoSync, ManualCopy, ManualPaste } from '../wailsjs/go/main/App.js'
    import { EventsOn } from '../wailsjs/runtime/runtime.js'
    import { onMount } from 'svelte';
    
    // Estado
    let roomID: string = "";
    let clientID: string = "";
    let isConnected: boolean = false;
    let autoSync: boolean = true;
    let connectedUsers: string[] = [];
    let lastClipboardContent: string = "";

    onMount(() => {
        EventsOn("clipboard_update", (payload: string) => {
            lastClipboardContent = payload;
        });

        EventsOn("user_list_updated", (users: string[]) => {
            connectedUsers = users;
        });

        EventsOn("room_disconnected", (msg: string) => {
            isConnected = false;
        });
    });

    function handleJoinRoom(): void {
        if (!roomID || !clientID) return;
        JoinRoom(roomID, clientID).then(result => {
            if (result === "OK") isConnected = true;
        });
    }

    function toggleAutoSync(): void {
        autoSync = !autoSync;
        SetAutoSync(autoSync);
    }
</script>

<main>
    {#if isConnected}
        <!-- Top Bar solo con el panel de usuarios a la derecha -->
        <div class="top-bar">
            <div class="spacer"></div>
            <div class="users-panel">
                <h4>Conectados ({connectedUsers.length})</h4>
                <ul>
                    {#each connectedUsers as user}
                        <li class={user === clientID ? 'is-me' : ''}>{user} {user === clientID ? '(Tu)' : ''}</li>
                    {/each}
                </ul>
            </div>
        </div>

        <div class="session-content">
            <img alt="Logo" class="mini-logo" src="{logo}">
            <h2>Sala: {roomID}</h2>

            <div class="sync-controls-container">
                <div class="sync-controls">
                    <label class="switch">
                        <input type="checkbox" checked={autoSync} on:change={toggleAutoSync}>
                        <span class="slider round"></span>
                    </label>
                    <span>Sincronizacion Automatica: <strong>{autoSync ? 'ACTIVADA' : 'DESACTIVADA'}</strong></span>
                </div>
            </div>

            <div class="preview-card">
                <p>{lastClipboardContent || "Esperando contenido de la sala..."}</p>
            </div>

            <div class="action-buttons">
                <button class="btn-action" on:click={() => ManualCopy(lastClipboardContent)} disabled={!lastClipboardContent}>
                    Copiar de la Sala
                </button>
                <button class="btn-action" on:click={ManualPaste}>
                    Pegar a la Sala
                </button>
            </div>

            <button class="btn-disconnect" on:click={Disconnect}>Abandonar Sala</button>
        </div>
    {:else}
        <div class="login-view">
            <img alt="Logo" id="logo" src="{logo}">
            <h1>Universal Clipboard</h1>
            <div class="form">
                <input bind:value={roomID} placeholder="Nombre de la Sala" />
                <input bind:value={clientID} placeholder="Tu Nombre de Usuario" />
                <button class="btn-primary" on:click={handleJoinRoom}>Entrar</button>
            </div>
        </div>
    {/if}
</main>

<style>
    :global(body) { margin: 0; background: #e0e0e0; font-family: 'Nunito', sans-serif; }
    
    main { height: 100vh; display: flex; flex-direction: column; }

    .top-bar { display: flex; justify-content: space-between; padding: 1rem; }
    .spacer { flex: 1; }

    .users-panel {
        background: rgba(255, 255, 255, 0.8); backdrop-filter: blur(10px);
        padding: 0.8rem; border-radius: 12px; width: 160px;
        box-shadow: 0 4px 15px rgba(0,0,0,0.1); border: 1px solid rgba(255,255,255,0.4);
    }
    .users-panel h4 { margin: 0 0 0.5rem; font-size: 0.8rem; color: #555; text-transform: uppercase; text-align: left; }
    .users-panel ul { list-style: none; padding: 0; margin: 0; font-size: 0.85rem; text-align: left; }
    .users-panel li { padding: 4px 0; color: #333; border-bottom: 1px solid rgba(0,0,0,0.05); }
    .is-me { color: #2980b9; font-weight: bold; }

    .session-content { flex: 1; display: flex; flex-direction: column; align-items: center; justify-content: center; padding: 2rem; }
    .mini-logo { width: 60px; margin-bottom: 1rem; }

    .sync-controls-container {
        background: #f0f0f0; padding: 0.6rem 1.2rem; border-radius: 30px;
        box-shadow: inset 0 2px 4px rgba(0,0,0,0.05); margin-bottom: 1rem;
    }
    .sync-controls { display: flex; align-items: center; gap: 12px; font-size: 0.9rem; color: #444; }

    .preview-card {
        background: #ffffff; width: 100%; max-width: 450px; min-height: 120px;
        padding: 1.5rem; border-radius: 12px; margin: 1.5rem 0;
        box-shadow: 0 8px 20px rgba(0,0,0,0.08); display: flex; align-items: center; justify-content: center;
        word-break: break-all; color: #2c3e50; border: 1px solid #dcdde1;
    }

    .action-buttons { display: flex; gap: 1rem; margin-bottom: 2rem; }
    .btn-action { 
        background: #ffffff; border: 2px solid #2980b9; color: #2980b9; 
        padding: 0.8rem 1.4rem; border-radius: 8px; font-weight: bold; cursor: pointer;
        transition: all 0.2s;
    }
    .btn-action:hover:not(:disabled) { background: #2980b9; color: white; }
    .btn-action:disabled { opacity: 0.5; cursor: not-allowed; border-color: #bdc3c7; color: #7f8c8d; }

    /* Switch Style */
    .switch { position: relative; display: inline-block; width: 46px; height: 24px; }
    .switch input { opacity: 0; width: 0; height: 0; }
    .slider { position: absolute; cursor: pointer; top: 0; left: 0; right: 0; bottom: 0; background-color: #bdc3c7; transition: .4s; border-radius: 24px; }
    .slider:before { position: absolute; content: ""; height: 18px; width: 18px; left: 3px; bottom: 3px; background-color: white; transition: .4s; border-radius: 50%; }
    input:checked + .slider { background-color: #27ae60; }
    input:checked + .slider:before { transform: translateX(22px); }

    .btn-disconnect { background: none; border: none; color: #c0392b; cursor: pointer; text-decoration: underline; font-size: 0.9rem; margin-top: 1rem; }

    .login-view { height: 100vh; display: flex; flex-direction: column; align-items: center; justify-content: center; background: #e0e0e0; }
    #logo { width: 120px; margin-bottom: 2rem; }
    .form { display: flex; flex-direction: column; gap: 1rem; width: 300px; }
    input { padding: 0.8rem; border: 1px solid #bdc3c7; border-radius: 8px; font-size: 1rem; background: white; }
    .btn-primary { background: #2980b9; color: white; border: none; padding: 1rem; border-radius: 8px; font-weight: bold; cursor: pointer; }
</style>
