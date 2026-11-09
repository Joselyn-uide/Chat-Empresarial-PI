// Cliente de consola de nuestro Chat Empresarial.
// Para ejecutarlo, desde la carpeta app-go: go run ./cmd/client
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"chat-empresarial/internal/handlers"
	"chat-empresarial/internal/models"
	"chat-empresarial/internal/storage"
)

// En main solo decidimos cómo termina el programa.
// Toda la lógica la pusimos en ejecutar(), que devuelve un error en vez de usar panic.
func main() {
	if err := ejecutar(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func ejecutar() error {
	// Lector para leer lo que el usuario escribe en la consola
	teclado := bufio.NewReader(os.Stdin)

	fmt.Println("========================================")
	fmt.Println("     Cliente de Chat Empresarial")
	fmt.Println("========================================")

	// Pedimos la IP y el puerto; si el usuario solo presiona Enter usamos 127.0.0.1:8080
	ip := leer(teclado, "IP del servidor [127.0.0.1]: ", "127.0.0.1")
	puerto := leer(teclado, "Puerto [8080]: ", "8080")

	// Creamos el usuario con el constructor del modelo, que valida que el nombre no esté vacío
	usuario, err := models.NewUsuario(leer(teclado, "Nombre de usuario: ", ""), "")
	if err != nil {
		return err
	}

	// Creamos el historial de la sesión y el cliente que maneja la conexión
	historial := storage.NewHistorial()
	cliente, err := handlers.NewClienteChat(usuario, historial)
	if err != nil {
		return err
	}
	if err := cliente.Conectar(ip, puerto); err != nil {
		return err
	}
	fmt.Printf("Conectado a %s:%s como %s (%s)\n", ip, puerto, usuario.GetNombreUsuario(), usuario.GetRol())

	// Lanzamos una goroutine que escucha al servidor mientras el usuario usa el menú.
	// Si la conexión termina, nos avisa por el channel fin.
	fin := make(chan error, 1)
	go cliente.Escuchar(func(m *models.Mensaje) {
		fmt.Printf("\n<< %s\n> ", m.GetContenido())
	}, fin)

	for {
		// Revisamos si el servidor cerró la conexión, sin quedarnos bloqueados esperando
		select {
		case err := <-fin:
			if err != nil {
				return err
			}
			fmt.Println("\nEl servidor cerró la conexión.")
			return nil
		default:
		}

		mostrarMenu()
		switch leer(teclado, "> ", "") {
		case "1":
			// Enviar mensaje: si está vacío mostramos el error y el programa sigue
			if err := cliente.EnviarMensaje(leer(teclado, "Mensaje: ", "")); err != nil {
				fmt.Println("No se envió:", err)
			}
		case "2":
			// Consultar el historial de la sesión activa
			mostrarHistorial(historial)
		case "3":
			// Guardar el historial en un JSON con el nombre del usuario y la fecha
			ruta := fmt.Sprintf("data/historial_%s_%s.json",
				usuario.GetNombreUsuario(), time.Now().Format("20060102_150405"))
			if err := historial.GuardarJSON(ruta); err != nil {
				fmt.Println("No se guardó:", err)
			} else {
				fmt.Println("Historial guardado en", ruta)
			}
		case "4", handlers.ComandoSalir:
			// Salir: enviamos SALIR al servidor y cerramos el socket de forma ordenada
			if err := cliente.Desconectar(); err != nil {
				return err
			}
			fmt.Println("Sesión cerrada. ¡Hasta luego!")
			return nil
		case "":
			// Si el usuario solo presiona Enter, volvemos a mostrar el menú
		default:
			fmt.Println("Opción no válida. Elige un número del 1 al 4.")
		}
	}
}

// mostrarMenu imprime las opciones disponibles para el usuario.
func mostrarMenu() {
	fmt.Println("\n------------- MENÚ -------------")
	fmt.Println("1. Enviar mensaje")
	fmt.Println("2. Ver historial de la sesión")
	fmt.Println("3. Guardar historial en JSON")
	fmt.Println("4. Salir (SALIR)")
}

// mostrarHistorial recorre los mensajes guardados y los muestra con su hora y remitente.
func mostrarHistorial(h *storage.Historial) {
	mensajes := h.ObtenerTodos()
	if len(mensajes) == 0 {
		fmt.Println("(Todavía no hay mensajes en esta sesión)")
		return
	}
	fmt.Printf("------ Historial (%d mensajes) ------\n", len(mensajes))
	for _, m := range mensajes {
		fmt.Printf("[%s] %s: %s\n", m.GetFechaHora().Format("15:04:05"), m.GetRemitente(), m.GetContenido())
	}
}

// leer muestra un texto, lee una línea del teclado y, si viene vacía,
// devuelve el valor por defecto que le pasamos.
func leer(teclado *bufio.Reader, etiqueta, porDefecto string) string {
	fmt.Print(etiqueta)
	linea, err := teclado.ReadString('\n')
	if err != nil {
		return porDefecto
	}
	linea = strings.TrimSpace(linea)
	if linea == "" {
		return porDefecto
	}
	return linea
}