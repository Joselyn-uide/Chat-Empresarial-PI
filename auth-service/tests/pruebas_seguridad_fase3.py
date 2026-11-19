"""
Pruebas de seguridad PS-06 y PS-07 - Proyecto Integrador Fase 3.

"""
import argparse
import json
import time
import urllib.error
import urllib.request
from datetime import datetime


def peticion(base, metodo, ruta, cuerpo=None, token=None):
    """Envia una solicitud HTTP y devuelve (codigo, cuerpo de la respuesta)."""
    datos = json.dumps(cuerpo).encode() if cuerpo is not None else None
    req = urllib.request.Request(base + ruta, data=datos, method=metodo)
    req.add_header("Content-Type", "application/json")
    if token is not None:
        req.add_header("Authorization", f"Bearer {token}")
    try:
        with urllib.request.urlopen(req, timeout=10) as r:
            return r.status, r.read().decode().strip()
    except urllib.error.HTTPError as e:
        return e.code, e.read().decode().strip()


def veredicto(ok):
    return "CUMPLE" if ok else "NO CUMPLE"


def ps06(base):
    print("\n=== PS-06: acceso a /chat/mensajes sin token y con token alterado ===")
    cod_a, cuerpo_a = peticion(base, "GET", "/chat/mensajes")
    print(f"[a] Sin cabecera Authorization -> HTTP {cod_a} | {cuerpo_a}")
    cod_b, cuerpo_b = peticion(base, "GET", "/chat/mensajes", token="token_falso_manipulado_123")
    print(f"[b] Con token falso/alterado   -> HTTP {cod_b} | {cuerpo_b}")
    ok = cod_a in (401, 403) and cod_b in (401, 403)
    print("Esperado: ambos rechazados (401 o 403) y sin devolver el historial.")
    print("Resultado PS-06:", veredicto(ok))
    return ok


def ps07(base, usuario, intentos):
    print(f"\n=== PS-07: {intentos} intentos de login fallidos seguidos para '{usuario}' ===")
    primer_429 = None
    for i in range(1, intentos + 1):
        codigo, _ = peticion(base, "POST", "/login",
                             {"username": usuario, "password": f"clave_incorrecta_{i}"})
        print(f"Intento {i:02d} -> HTTP {codigo}")
        if codigo == 429 and primer_429 is None:
            primer_429 = i
    if primer_429 is None:
        print("Resultado PS-07: NO CUMPLE (no aparecio HTTP 429; revisar el limite en FastAPI).")
        return False
    print(f"El limite se activo en el intento {primer_429} (esperado: intento 6, limite de 5/minuto).")
    return True


def recuperacion(base, usuario):
    print("\n=== PS-07 (continuacion): recuperacion del servicio tras la ventana del limite ===")
    codigo, cuerpo = peticion(base, "POST", "/login",
                              {"username": usuario, "password": "clave_incorrecta_final"})
    print(f"Nuevo intento -> HTTP {codigo} | {cuerpo}")
    ok = codigo == 401
    print("Esperado: HTTP 401 (la solicitud vuelve a ser procesada; ya no hay bloqueo 429).")
    print("Resultado PS-07:", veredicto(ok))
    return ok


def main():
    parser = argparse.ArgumentParser(description="Pruebas PS-06 y PS-07")
    parser.add_argument("--url", default="http://127.0.0.1:8000", help="direccion del microservicio FastAPI")
    parser.add_argument("--usuario", default="j.cardenas", help="usuario registrado para las pruebas")
    parser.add_argument("--intentos", type=int, default=8, help="intentos de login para PS-07")
    parser.add_argument("--espera", type=int, default=65, help="segundos de espera antes de la recuperacion")
    args = parser.parse_args()

    print("Pruebas de seguridad - Proyecto Integrador Fase 3")
    print("Microservicio:", args.url, "| Fecha:", datetime.now().strftime("%Y-%m-%d %H:%M:%S"))

    try:
        r6 = ps06(args.url)
        limite = ps07(args.url, args.usuario, args.intentos)
        print(f"\nEsperando {args.espera} segundos para que se libere el limite...")
        time.sleep(args.espera)
        r7 = limite and recuperacion(args.url, args.usuario)
    except urllib.error.URLError as e:
        print("\nNo se pudo conectar con el microservicio:", e.reason)
        print("Verifica que FastAPI este encendido en", args.url)
        return

    print("\n=== Resumen ===")
    print("PS-06 (verificacion obligatoria del token):", veredicto(r6))
    print("PS-07 (limitacion de tasa y recuperacion): ", veredicto(r7))


if __name__ == "__main__":
    main()
