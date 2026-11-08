package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
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

// manejarCliente gestiona la comunicación individual con cada cliente conectado
func manejarCliente(conexion net.Conn) {
	defer conexion.Close()

	direccionCliente := conexion.RemoteAddr().String()
	fmt.Println("¡Nuevo cliente conectado desde:", direccionCliente)

	// Enviamos un mensaje de bienvenida al cliente
	conexion.Write([]byte("¡Conectado exitosamente al Servidor de Chat Empresarial!\n"))

	// Lector para recibir los mensajes que envía el cliente por el socket TCP
	lector := bufio.NewReader(conn)

	for {
    //Lee lo que escribe el cliente linea por linea
		mensaje, err := lector.ReadString('\n')
		if err != nil {
			fmt.Println("Cliente desconectado:", direccionCliente)
			break
		}

		mensajeLimpio := strings.TrimSpace(mensaje)
		if mensajeLimpio == "SALIR" {
			fmt.Println("El cliente", direccionCliente, "solicitó salir del chat.")
			conn.Write([]byte("Conexión cerrada. ¡Hasta luego!\n"))
			break
		}

		fmt.Printf("Mensaje recibido de [%s]: %s\n", direccionCliente, mensajeLimpio)

		// Respuesta de confirmación al cliente
		conexion.Write([]byte("Servidor recibido -> " + mensajeLimpio + "\n"))
	}
}
