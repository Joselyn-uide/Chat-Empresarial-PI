// En este archivo probamos las validaciones de entrada con varios casos
// correctos e incorrectos, usando pruebas en forma de tabla.
package services

import (
	"errors"
	"strings"
	"testing"
)

// Probamos ValidarMensaje con mensajes correctos, vacíos, largos y con caracteres raros.
func TestValidarMensaje(t *testing.T) {
	casos := []struct {
		nombre   string
		entrada  string
		esperado string
		err      error
	}{
		{"mensaje normal", "Hola equipo", "Hola equipo", nil},
		{"quita espacios de los bordes", "   Hola   ", "Hola", nil},
		{"acepta tildes y ñ", "Reunión mañana", "Reunión mañana", nil},
		{"vacío", "", "", ErrMensajeVacio},
		{"solo espacios", "     ", "", ErrMensajeVacio},
		{"solo tabulaciones", "\t\t", "", ErrMensajeVacio},
		{"demasiado largo", strings.Repeat("a", MaxLongitudMensaje+1), "", ErrMensajeMuyLargo},
		{"justo en el límite", strings.Repeat("a", MaxLongitudMensaje), strings.Repeat("a", MaxLongitudMensaje), nil},
		{"carácter de control", "hola\x00mundo", "", ErrCaracteresInvalidos},
		{"salto de línea en medio", "hola\nmundo", "", ErrCaracteresInvalidos},
	}
	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			obtenido, err := ValidarMensaje(c.entrada)
			if !errors.Is(err, c.err) {
				t.Fatalf("esperábamos el error %v, pero obtuvimos %v", c.err, err)
			}
			if obtenido != c.esperado {
				t.Fatalf("esperábamos %q, pero obtuvimos %q", c.esperado, obtenido)
			}
		})
	}
}

// Probamos nombres de usuario válidos e inválidos.
func TestValidarUsuario(t *testing.T) {
	validos := []string{"Darwin", "jose_v", "ana.m", "Joselyn-L", "  Ana  "}
	for _, u := range validos {
		if _, err := ValidarUsuario(u); err != nil {
			t.Errorf("el usuario %q debería ser válido, pero obtuvimos: %v", u, err)
		}
	}
	invalidos := []string{"", "ab", "nombre con espacios", "usuario@correo", strings.Repeat("x", MaxLongitudUsuario+1)}
	for _, u := range invalidos {
		if _, err := ValidarUsuario(u); !errors.Is(err, ErrUsuarioInvalido) {
			t.Errorf("el usuario %q debería ser inválido, pero obtuvimos: %v", u, err)
		}
	}
}

// Probamos direcciones IP con formato correcto e incorrecto.
func TestValidarIP(t *testing.T) {
	for _, ip := range []string{"127.0.0.1", "192.168.10.5", "localhost", "::1"} {
		if _, err := ValidarIP(ip); err != nil {
			t.Errorf("la IP %q debería ser válida, pero obtuvimos: %v", ip, err)
		}
	}
	for _, ip := range []string{"", "999.1.1.1", "hola", "127.0.0"} {
		if _, err := ValidarIP(ip); !errors.Is(err, ErrIPInvalida) {
			t.Errorf("la IP %q debería ser inválida, pero obtuvimos: %v", ip, err)
		}
	}
}

// Probamos puertos dentro y fuera del rango permitido.
func TestValidarPuerto(t *testing.T) {
	for _, p := range []string{"8080", "1", "65535", " 9000 "} {
		if _, err := ValidarPuerto(p); err != nil {
			t.Errorf("el puerto %q debería ser válido, pero obtuvimos: %v", p, err)
		}
	}
	for _, p := range []string{"", "0", "65536", "abc", "-5", "80.5"} {
		if _, err := ValidarPuerto(p); !errors.Is(err, ErrPuertoInvalido) {
			t.Errorf("el puerto %q debería ser inválido, pero obtuvimos: %v", p, err)
		}
	}
}