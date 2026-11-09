// Package handlers
// En este paquete organizamos toda la lógica de comunicación del cliente
// por sockets TCP: conectar, enviar, escuchar y desconectar.
// Decidimos devolver siempre los fallos con el tipo error de Go
// y no usar panic, para que el programa nunca se caiga de golpe.
package handlers

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	"chat-empresarial/internal/models"
	"chat-empresarial/internal/storage"
)

// ComandoSalir es la instrucción que nuestro servidor reconoce para cerrar la sesión.
const ComandoSalir = "SALIR"

// RemitenteServidor lo usamos para identificar en el historial
// los mensajes que nos llegan desde el servidor.
const RemitenteServidor = "Servidor"

// Definimos nuestros propios errores para poder compararlos con errors.Is.
var (
	ErrDireccionInvalida = errors.New("la IP y el puerto no pueden estar vacíos")
	ErrNoConectado       = errors.New("el cliente no está conectado al servidor")
)

// ClienteChat representa la conexión de un usuario con el servidor.
// Dejamos sus atributos en minúscula (privados) para aplicar encapsulamiento:
// solo se pueden usar a través de sus métodos.
type ClienteChat struct {
	usuario   *models.Usuario
	conexion  net.Conn
	historial *storage.Historial
}

// NewClienteChat es nuestro constructor. Validamos que reciba un usuario
// y un historial antes de crear el cliente.
func NewClienteChat(usuario *models.Usuario, historial *storage.Historial) (*ClienteChat, error) {
	if usuario == nil {
		return nil, errors.New("el cliente necesita un usuario válido")
	}
	if historial == nil {
		return nil, errors.New("el cliente necesita un historial")
	}
	return &ClienteChat{usuario: usuario, historial: historial}, nil
}

// Conectar abre el socket TCP hacia la IP y el puerto que escribe el usuario.
// Pusimos un límite de 5 segundos para que no se quede esperando para siempre.
func (c *ClienteChat) Conectar(ip, puerto string) error {
	ip, puerto = strings.TrimSpace(ip), strings.TrimSpace(puerto)
	if ip == "" || puerto == "" {
		return ErrDireccionInvalida
	}

	direccion := net.JoinHostPort(ip, puerto)
	conexion, err := net.DialTimeout("tcp", direccion, 5*time.Second)
	if err != nil {
		return fmt.Errorf("no se pudo conectar a %s: %w", direccion, err)
	}
	c.conexion = conexion
	return nil
}

// EnviarMensaje valida el texto usando el constructor NewMensaje del modelo,
// lo envía por el socket y lo guarda en el historial de la sesión.
func (c *ClienteChat) EnviarMensaje(texto string) error {
	if c.conexion == nil {
		return ErrNoConectado
	}

	// Si el mensaje está vacío, NewMensaje nos devuelve el error y no se envía nada
	mensaje, err := models.NewMensaje(c.usuario.GetNombreUsuario(), texto)
	if err != nil {
		return err
	}

	// Enviamos el mensaje como una línea de texto, que es lo que entiende el servidor
	linea := fmt.Sprintf("[%s] %s\n", mensaje.GetRemitente(), mensaje.GetContenido())
	if _, err := c.conexion.Write([]byte(linea)); err != nil {
		return fmt.Errorf("error al enviar el mensaje: %w", err)
	}
	return c.historial.Agregar(mensaje)
}

// Escuchar lee las respuestas del servidor línea por línea.
// Este método lo ejecutamos en una goroutine para poder recibir mensajes
// mientras el usuario usa el menú. Cada respuesta la guardamos en el
// historial y la mostramos en pantalla con la función mostrar.
// Cuando termina, avisamos por el channel fin (nil si fue un cierre normal).
func (c *ClienteChat) Escuchar(mostrar func(*models.Mensaje), fin chan<- error) {
	if c.conexion == nil {
		fin <- ErrNoConectado
		return
	}

	lector := bufio.NewReader(c.conexion)
	for {
		linea, err := lector.ReadString('\n')
		if err != nil {
			// Si el servidor cerró la conexión de forma normal, no es un error grave
			if errors.Is(err, io.EOF) || errors.Is(err, net.ErrClosed) {
				fin <- nil
				return
			}
			fin <- fmt.Errorf("se perdió la conexión con el servidor: %w", err)
			return
		}

		// Si llega una línea vacía la ignoramos sin detener el cliente
		respuesta, err := models.NewMensaje(RemitenteServidor, linea)
		if err != nil {
			continue
		}
		if err := c.historial.Agregar(respuesta); err != nil {
			fin <- err
			return
		}
		mostrar(respuesta)
	}
}

// Desconectar envía el comando SALIR al servidor y cierra el socket de forma ordenada.
func (c *ClienteChat) Desconectar() error {
	if c.conexion == nil {
		return ErrNoConectado
	}

	// Si el aviso falla igual cerramos el socket, por eso solo lo informamos
	if _, err := c.conexion.Write([]byte(ComandoSalir + "\n")); err != nil {
		fmt.Println("Aviso: no se pudo notificar la salida al servidor:", err)
	}

	// Esperamos un momento para que alcance a llegar la despedida del servidor
	time.Sleep(300 * time.Millisecond)

	err := c.conexion.Close()
	c.conexion = nil
	if err != nil && !errors.Is(err, net.ErrClosed) {
		return fmt.Errorf("error al cerrar la conexión: %w", err)
	}
	return nil
}