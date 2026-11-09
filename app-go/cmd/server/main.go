package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
)

// Guardamos las conexiones de todos los clientes conectados para poder
// reenviarles los mensajes (broadcast). Usamos un mutex porque cada cliente
// se atiende en su propia goroutine y todas usan este mapa al mismo tiempo.
var (
	clientes   = make(map[net.Conn]bool)
	clientesMu sync.Mutex
)

func main() {
	// Definimos el puerto en el que va a escuchar (funcionar) nuestro servidor de chat
	puerto := ":8080"

	servidor, err := net.Listen("tcp", puerto)
	if err != nil {
		fmt.Println("Error al iniciar el servidor:", err)
		os.Exit(1)
	}
	defer servidor.Close()

	fmt.Println("========================================")
	fmt.Println(" Servidor de Chat Empresarial Iniciado")
	fmt.Println(" Esperando conexiones en el puerto", puerto)
	fmt.Println("========================================")

	// Ciclo infinito para aceptar múltiples conexiones de clientes de forma concurrente
	for {
		conexion, err := servidor.Accept()
		if err != nil {
			fmt.Println("Error al aceptar la conexión del cliente:", err)
			continue
		}

		// Usamos una goroutine para atender a cada cliente de forma concurrente e independiente
		go manejarCliente(conexion)
	}
}

// agregarCliente registra una nueva conexión en la lista de clientes.
func agregarCliente(conexion net.Conn) {
	clientesMu.Lock()
	defer clientesMu.Unlock()
	clientes[conexion] = true
}

// quitarCliente elimina una conexión de la lista cuando el cliente se va.
func quitarCliente(conexion net.Conn) {
	clientesMu.Lock()
	defer clientesMu.Unlock()
	delete(clientes, conexion)
}

// difundir envía un mensaje a todos los clientes conectados, excepto al remitente.
// Si falla el envío a uno de ellos, lo informamos y seguimos con los demás.
func difundir(remitente net.Conn, mensaje string) {
	clientesMu.Lock()
	defer clientesMu.Unlock()
	for conexion := range clientes {
		if conexion == remitente {
			continue
		}
		if _, err := conexion.Write([]byte(mensaje + "\n")); err != nil {
			fmt.Println("No se pudo enviar el mensaje a", conexion.RemoteAddr(), ":", err)
		}
	}
}

// manejarCliente gestiona la comunicación individual con cada cliente conectado
func manejarCliente(conexion net.Conn) {
	direccionCliente := conexion.RemoteAddr().String()
	agregarCliente(conexion)
	defer func() {
		quitarCliente(conexion)
		conexion.Close()
		difundir(conexion, "*** Un usuario salió del chat ***")
	}()

	fmt.Println("¡Nuevo cliente conectado desde:", direccionCliente)

	// Enviamos un mensaje de bienvenida al cliente y avisamos a los demás
	conexion.Write([]byte("¡Conectado exitosamente al Servidor de Chat Empresarial!\n"))
	difundir(conexion, "*** Un nuevo usuario se unió al chat ***")

	// Lector para recibir los mensajes que envía el cliente por el socket TCP
	lector := bufio.NewReader(conexion)

	for {
		//Lee lo que escribe el cliente linea por linea
		mensaje, err := lector.ReadString('\n')
		if err != nil {
			fmt.Println("Cliente desconectado:", direccionCliente)
			break
		}

		mensajeLimpio := strings.TrimSpace(mensaje)
		if mensajeLimpio == "" {
			continue
		}
		if mensajeLimpio == "SALIR" {
			fmt.Println("El cliente", direccionCliente, "solicitó salir del chat.")
			conexion.Write([]byte("Conexión cerrada. ¡Hasta luego!\n"))
			break
		}

		fmt.Printf("Mensaje recibido de [%s]: %s\n", direccionCliente, mensajeLimpio)

		// Respuesta de confirmación al cliente que envió el mensaje
		conexion.Write([]byte("Servidor recibido -> " + mensajeLimpio + "\n"))

		// Reenviamos el mensaje a todos los demás clientes conectados
		difundir(conexion, mensajeLimpio)
	}
}