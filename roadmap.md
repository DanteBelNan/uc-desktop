# Roadmap de Desarrollo: uc-desktop (v0.1.x)

Este roadmap detalla los pasos especificos para la implementacion del cliente de escritorio del Universal Clipboard utilizando Wails y Svelte.

## 🚀 Fase 1: MVP Cliente de Escritorio (Sincronizacion de Texto Plano)

1. **Setup Inicial (v0.1.0)**
    - 1.1 Inicializacion del proyecto Wails con template Svelte-TS. [COMPLETADO]
    - 1.2 Configuracion de Git y repositorio remoto en GitHub. [PENDIENTE]
    - 1.3 Definicion de la interfaz basica de usuario (UI) para gestion de salas. [PENDIENTE]

2. **Conectividad (v0.1.1)**
    - 2.1 Implementacion de cliente WebSocket en Go. [PENDIENTE]
    - 2.2 Integracion de eventos de WebSocket con el frontend de Svelte. [PENDIENTE]
    - 2.3 Validacion de conexion con `uc-server`. [PENDIENTE]

3. **Sincronizacion (v0.1.2)**
    - 3.1 Integracion con la API de portapapeles del sistema operativo. [PENDIENTE]
    - 3.2 Implementacion de escucha activa de cambios en el portapapeles local. [PENDIENTE]
    - 3.3 Sincronizacion automatica: Local -> Servidor -> Otros clientes. [PENDIENTE]

4. **Testing y Calidad (v0.1.3)**
    - 4.1 Tests unitarios para el cliente WebSocket y la logica de sincronizacion. [PENDIENTE]
    - 4.2 Documentacion de uso del cliente en `uc-desktop/README.md`. [PENDIENTE]

---

## 🔒 Futuras Versiones
- **v0.2.x:** Implementacion de cifrado E2EE (AES-256).
- **v0.3.x:** Soporte para imagenes y archivos pequeños.
- **v0.4.x:** Minimizar a bandeja de sistema (System Tray) y notificaciones nativas.
