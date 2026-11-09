// Package storage
// En este paquete implementamos el almacenamiento temporal del historial
// de la sesión activa del chat. Decidimos guardar los mensajes en memoria
// y además dar la opción de exportarlos a un archivo JSON.
package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"chat-empresarial/internal/models"
)

// ErrHistorialVacio lo definimos como un error propio para avisar
// que no hay mensajes que guardar, en lugar de usar panic.
var ErrHistorialVacio = errors.New("el historial está vacío")

// Historial es la estructura donde guardamos los mensajes de la sesión.
// Dejamos el atributo mensajes en minúscula (privado) para aplicar
// encapsulamiento: desde afuera solo se puede usar con sus métodos.
// Agregamos un mutex porque la goroutine que escucha al servidor y el
// menú del usuario usan el historial al mismo tiempo.
type Historial struct {
	mu       sync.RWMutex
	mensajes []*models.Mensaje
}

// NewHistorial es nuestro constructor: crea un historial vacío y listo para usar.
func NewHistorial() *Historial {
	return &Historial{mensajes: make([]*models.Mensaje, 0)}
}

// Agregar guarda un mensaje nuevo en el historial.
// Validamos que el mensaje no sea nulo antes de guardarlo.
func (h *Historial) Agregar(m *models.Mensaje) error {
	if m == nil {
		return errors.New("no se puede guardar un mensaje nulo")
	}
	// Bloqueamos el mutex mientras escribimos para evitar conflictos entre goroutines
	h.mu.Lock()
	defer h.mu.Unlock()
	h.mensajes = append(h.mensajes, m)
	return nil
}

// ObtenerTodos devuelve los mensajes de la sesión.
// Entregamos una copia para que nadie pueda modificar la lista original.
func (h *Historial) ObtenerTodos() []*models.Mensaje {
	h.mu.RLock()
	defer h.mu.RUnlock()
	copia := make([]*models.Mensaje, len(h.mensajes))
	copy(copia, h.mensajes)
	return copia
}

// Cantidad nos permite consultar cuántos mensajes hay en la sesión.
func (h *Historial) Cantidad() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.mensajes)
}

// GuardarJSON exporta el historial a un archivo JSON en la ruta indicada.
// Cada posible fallo lo devolvemos como error con un mensaje claro.
func (h *Historial) GuardarJSON(ruta string) error {
	mensajes := h.ObtenerTodos()
	if len(mensajes) == 0 {
		return ErrHistorialVacio
	}

	// Creamos la carpeta (por ejemplo data/) si todavía no existe
	if err := os.MkdirAll(filepath.Dir(ruta), 0o755); err != nil {
		return fmt.Errorf("no se pudo crear la carpeta: %w", err)
	}

	archivo, err := os.Create(ruta)
	if err != nil {
		return fmt.Errorf("no se pudo crear el archivo %s: %w", ruta, err)
	}
	defer archivo.Close()

	// Configuramos el JSON con sangría para que sea fácil de leer
	codificador := json.NewEncoder(archivo)
	codificador.SetIndent("", "  ")
	codificador.SetEscapeHTML(false) // así "->" no se guarda como "\u003e"
	if err := codificador.Encode(mensajes); err != nil {
		return fmt.Errorf("no se pudo escribir el historial en JSON: %w", err)
	}
	return nil
}

// CargarJSON nos permite recuperar un historial que guardamos antes.
func (h *Historial) CargarJSON(ruta string) error {
	datos, err := os.ReadFile(ruta)
	if err != nil {
		return fmt.Errorf("no se pudo leer el archivo %s: %w", ruta, err)
	}

	var mensajes []*models.Mensaje
	if err := json.Unmarshal(datos, &mensajes); err != nil {
		return fmt.Errorf("el archivo %s no contiene un JSON válido: %w", ruta, err)
	}

	h.mu.Lock()
	defer h.mu.Unlock()
	h.mensajes = mensajes
	return nil
}