"""
Pruebas de seguridad PS-05 y PS-06 - Proyecto Integrador Fase 3.

"""
import argparse
import json
import time
import urllib.error
import urllib.request
from collections import Counter
from concurrent.futures import ThreadPoolExecutor
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


def ps05(base, usuario, total, hilos):
    print(f"\n=== PS-05: rafaga de {total} solicitudes simultaneas de login para '{usuario}' ===")
    def intento(i):
        codigo, _ = peticion(base, "POST", "/login",
                             {"username": usuario, "password": f"clave_incorrecta_{i}"})
        return codigo

    inicio = time.time()
    with ThreadPoolExecutor(max_workers=hilos) as ex:
        codigos = list(ex.map(intento, range(1, total + 1)))
    duracion = time.time() - inicio

    conteo = Counter(codigos)
    print(f"Solicitudes enviadas: {total} en {duracion:.2f} s ({hilos} hilos en paralelo)")
    for codigo, n in sorted(conteo.items()):
        nombre = {401: "Unauthorized (procesadas)", 429: "Too Many Requests (bloqueadas)"}.get(codigo, "")
        print(f"  HTTP {codigo} {nombre}: {n}")
    bloqueadas = conteo.get(429, 0)
    print(f"Esperado: como maximo 5 procesadas (limite 5/minuto) y el resto bloqueadas con 429.")
    return bloqueadas > 0 and conteo.get(401, 0) <= 5


def ps05_disponibilidad(base):
    print("\n=== PS-05: disponibilidad del servidor durante el bloqueo ===")
    codigo, _ = peticion(base, "GET", "/docs")
    print(f"GET /docs mientras la IP esta bloqueada en /login -> HTTP {codigo}")
    print("Esperado: HTTP 200 (el servidor sigue en linea; solo se limita /login).")
    return codigo == 200


def ps05_recuperacion(base, usuario):
    print("\n=== PS-05: recuperacion tras la ventana del limite ===")
    codigo, cuerpo = peticion(base, "POST", "/login",
                              {"username": usuario, "password": "clave_incorrecta_final"})
    print(f"Nuevo intento -> HTTP {codigo} | {cuerpo}")
    print("Esperado: HTTP 401 (la solicitud vuelve a ser procesada; ya no hay bloqueo 429).")
    return codigo == 401


def main():
    parser = argparse.ArgumentParser(description="Pruebas PS-05 y PS-06")
    parser.add_argument("--url", default="http://127.0.0.1:8000", help="direccion del microservicio FastAPI")
    parser.add_argument("--usuario", default="j.cardenas", help="usuario registrado para las pruebas")
    parser.add_argument("--total", type=int, default=200, help="solicitudes simultaneas para PS-05")
    parser.add_argument("--hilos", type=int, default=50, help="hilos en paralelo para la rafaga")
    parser.add_argument("--espera", type=int, default=65, help="segundos de espera antes de la recuperacion")
    args = parser.parse_args()

    print("Pruebas de seguridad - Proyecto Integrador Fase 3")
    print("Microservicio:", args.url, "| Fecha:", datetime.now().strftime("%Y-%m-%d %H:%M:%S"))

    try:
        r6 = ps06(args.url)
        bloqueo = ps05(args.url, args.usuario, args.total, args.hilos)
        disponible = ps05_disponibilidad(args.url)
        print(f"\nEsperando {args.espera} segundos para que se libere el limite...")
        time.sleep(args.espera)
        recupera = ps05_recuperacion(args.url, args.usuario)
    except urllib.error.URLError as e:
        print("\nNo se pudo conectar con el microservicio:", e.reason)
        print("Verifica que FastAPI este encendido en", args.url)
        return

    r5 = bloqueo and disponible and recupera
    print("\n=== Resumen ===")
    print("PS-05 (rafaga de 200 solicitudes: bloqueo, disponibilidad y recuperacion):", veredicto(r5))
    print("PS-06 (verificacion obligatoria del token): ", veredicto(r6))


if __name__ == "__main__":
    main()
