# Proyecto Integrador: Sistema de chat empresarial Cliente-Servidor con red segmentada

## Integrantes
* Joselyn Paulet Laguna Puetate 
* Jose Emiliano Velasco Sosa 
* Darwin Oswaldo Román Aguirre 

---

## Descripción del problema
Actualmente la empresa presenta dificultades en la comunicación interna entre los colaboradores que trabajan en sus dos sedes. Gran parte de la información operativa se transmite mediante llamadas, mensajes personales o correos electrónicos no centralizados, lo que provoca retrasos, pérdida de información y dificultad para mantener un registro ordenado de las conversaciones.

## Justificación del sistema seleccionado
Se desarrolló una aplicación de chat corporativo en arquitectura cliente-servidor para permitir a los empleados de ambas sedes comunicarse de forma directa, centralizada y sin depender de aplicaciones personales. El sistema utiliza el lenguaje Go por su alta eficiencia en el manejo de concurrencia mediante sockets y goroutines, y se integra con un microservicio de autenticación en FastAPI para proteger el acceso mediante hashes seguros.

## Objetivos del Proyecto
* **Objetivo General:** Desarrollar una aplicación de chat corporativo en arquitectura cliente/servidor, en el lenguaje Go, que permita la interacción concurrente y segura entre usuarios de dos unidades organizativas distintas, con verificación de identidad y sobre una infraestructura de red segmentada.
* **Objetivos Específicos:** 
  1.- Implementar funciones de mensajería bidireccional y gestión del historial básico.
  2.- Integrar el servicio de autenticación con FastAPI.
  3.- Controlar el acceso mediante hashing seguro.
  4.- Organizar el código bajo principios de Programación Orientada a Objetos (POO).
  5.-  Validar la conectividad mediante una red segmentada por VLANs y enrutamiento OSPF.

## Alcance del Proyecto
El proyecto abarca el diseño e implementación de un sistema que combina software y redes:
* **Software:** Lógica de chat cliente-servidor en Go y servicio de autenticación en FastAPI.
* **Infraestructura de Red:** Simulación en Cisco Packet Tracer de dos sedes interconectadas mediante VLANs, enlaces troncales (Trunk), Router-on-a-Stick, enrutamiento dinámico OSPF y redundancia RSTP.

## Exclusiones del Sistema
* No incluye despliegue en nubes comerciales de pago (AWS, Azure, Google Cloud); validación estrictamente local.
* No incluye llamadas de voz, videollamadas ni transmisión de archivos pesados (limitado a texto).
* No incluye aplicaciones móviles nativas ni pasarelas de mensajería pública externa.

---

## Video
Video grupal: https://drive.google.com/file/d/1BgLwMJoLMNJwATurPaoKMaBrOg6PJpEg/view?usp=drivesdk

---

## Requisitos de software
| Herramienta | Versión | Uso |
|---|---|---|
| Go | 1.21 o superior (probado con VERSION_AQUI) | Compilar y ejecutar el servidor y el cliente |
| Git | Cualquier versión reciente | Clonar el repositorio |
| Terminal | Bash, PowerShell o la terminal de GitHub Codespaces | Ejecutar los comandos |

No se necesitan librerías externas: el proyecto usa solo la librería estándar de Go.

## Estructura del código (`app-go/`)
app-go/
├── go.mod → Módulo chat-empresarial
├── cmd/
│ ├── server/main.go → Servidor TCP (una goroutine por cliente y reenvío a todos)
│ └── client/main.go → Cliente de consola con menú interactivo
├── internal/
│ ├── models/ → Estructuras Usuario y Mensaje
│ ├── handlers/cliente.go → Lógica de conexión del cliente (ClienteChat)
│ ├── services/ → Validaciones de entrada (mensaje, usuario, IP y puerto)
│ └── storage/ → Historial de la sesión en memoria y exportación a JSON
└── data/ → Historiales exportados en JSON (no se suben al repositorio)


## Cómo ejecutar el proyecto paso a paso
1. Clonar el repositorio y entrar a la carpeta del código:
   ```bash
   git clone https://github.com/Joselyn-uide/Chat-Empresarial-PI.git
   cd Chat-Empresarial-PI/app-go
   ```
2. Comprobar que Go está instalado:
   ```bash
   go version
   ```
3. **Terminal 1, iniciar el servidor:**
   ```bash
   go run ./cmd/server
   ```
   Debe mostrar `Esperando conexiones en el puerto :8080`.
4. **Terminal 2 (y 3 para un segundo usuario), iniciar el cliente:**
   ```bash
   go run ./cmd/client
   ```
   Presionar Enter para usar la IP `127.0.0.1` y el puerto `8080` por defecto, y escribir un nombre de usuario.
5. Usar el menú del cliente:
   - `1` Enviar mensaje
   - `2` Ver el historial de la sesión
   - `3` Guardar el historial en JSON (carpeta `data/`)
   - `4` Salir (envía el comando `SALIR`)

Para generar ejecutables en lugar de usar `go run`:
```bash
go build -o servidor ./cmd/server
go build -o cliente ./cmd/client
```

## Pruebas
```bash
go vet ./...      # análisis estático del código
go test ./...     # pruebas unitarias de validaciones e historial
```

## Declaración de uso de Inteligencia Artificial
En el desarrollo de este proyecto utilizamos las siguientes herramientas de Inteligencia Artificial como apoyo:
**Claude (Anthropic)** 
**ChatGPT (OpenAI)**
**ChatGPT (OpenAI)** 
||como herramienta de apoyo para:

- Orientación paso a paso en la estructura modular del código en Go, el manejo de errores y la concurrencia.
- Sugerencias de código para el cliente, el historial, las validaciones y el reenvío de mensajes del servidor, que revisamos, probamos y adaptamos al proyecto.

Todo el código fue revisado, ejecutado y validado por los integrantes del equipo, quienes entienden y pueden explicar su funcionamiento.
