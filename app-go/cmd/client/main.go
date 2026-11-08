package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	// Nos conectamos al servidor en localhost y el puerto 8080
	direccion := "127.0.0.1:8080"
	
	fmt.Println("Conectándose al servidor de chat en", direccion, "...")
	conexion, err := net.Dial("tcp", direccion)
	if err != nil {
		fmt.Println("No se pudo conectar al servidor:", err)
		os.Exit(1)
	}
	defer conexion.Close()

	fmt.Println("¡Conectado exitosamente al chat!")
	fmt.Println("Escribe tu mensaje y presiona Enter (Escribe 'SALIR' para desconectarte):")

	// Canal para leer lo que escribe el usuario desde su teclado en la terminal
	lectorTerminal := bufio.NewReader(os.Stdin)
	
	// Canal para leer las respuestas que envía el servidor
	lectorServidor := bufio.NewReader(conexion)

	// Goroutine para escuchar de forma concurrente lo que el servidor responde
	go func() {
		for {
			respuesta, err := lectorServidor.ReadString('\n')
			if err != nil {
				fmt.Println("\n[Conexión cerrada por el servidor]")
				os.Exit(0)
			}
			fmt.Print(strings.TrimSpace(respuesta), "\n> ")
		}
	}()

	// Bucle principal para enviar mensajes al servidor
	fmt.Print("> ")
	for {
		texto, err := lectorTerminal.ReadString('\n')
		if err != nil {
			fmt.Println("Error al leer la entrada:", err)
			break
		}

		mensaje := strings.TrimSpace(texto)
		if mensaje == "" {
			continue
		}

		// Enviamos el mensaje al servidor por la conexión TCP
		_, err = conexion.Write([]byte(mensaje + "\n"))
		if err != nil {
			fmt.Println("Error al enviar el mensaje al servidor.")
			break
		}

		if mensaje == "SALIR" {
			fmt.Println("Cerrando sesión de chat...")
			break
		}
	}
}
