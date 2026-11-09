// En este archivo escribimos las pruebas unitarias del historial
// para comprobar que guarda, exporta y valida correctamente.
package storage

import (
	"errors"
	"path/filepath"
	"testing"

	"chat-empresarial/internal/models"
)

// Probamos que no se pueda guardar un historial vacío
// y que devuelva nuestro error ErrHistorialVacio.
func TestHistorialVacioNoSeGuarda(t *testing.T) {
	h := NewHistorial()
	err := h.GuardarJSON(filepath.Join(t.TempDir(), "h.json"))
	if !errors.Is(err, ErrHistorialVacio) {
		t.Fatalf("esperábamos ErrHistorialVacio, pero obtuvimos: %v", err)
	}
}

// Probamos que un mensaje se guarde en JSON y que al cargarlo
// de nuevo recuperemos exactamente el mismo contenido.
func TestGuardarYCargarJSON(t *testing.T) {
	h := NewHistorial()

	// Creamos un mensaje usando el constructor del modelo del equipo
	m, err := models.NewMensaje("Darwin", "Hola equipo")
	if err != nil {
		t.Fatal(err)
	}
	if err := h.Agregar(m); err != nil {
		t.Fatal(err)
	}

	// Usamos una carpeta temporal para no dejar archivos basura
	ruta := filepath.Join(t.TempDir(), "historial.json")
	if err := h.GuardarJSON(ruta); err != nil {
		t.Fatal(err)
	}

	h2 := NewHistorial()
	if err := h2.CargarJSON(ruta); err != nil {
		t.Fatal(err)
	}
	if h2.Cantidad() != 1 || h2.ObtenerTodos()[0].GetContenido() != "Hola equipo" {
		t.Fatalf("el historial cargado no coincide: %+v", h2.ObtenerTodos())
	}
}

// Probamos que el historial rechace un mensaje nulo devolviendo un error.
func TestAgregarNulo(t *testing.T) {
	if err := NewHistorial().Agregar(nil); err == nil {
		t.Fatal("esperábamos un error al agregar un mensaje nulo")
	}
}