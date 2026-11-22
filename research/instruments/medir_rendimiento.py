import argparse
import csv
import json
import socket
import time
import urllib.error
import urllib.request
from datetime import datetime, timezone
from pathlib import Path


def medir_fastapi_http(url: str, metodo: str, payload: dict | None, timeout_s: float) -> tuple[float, int, str, int]:
    cuerpo_bytes = b""
    headers = {"Content-Type": "application/json"}
    if payload is not None:
        cuerpo_bytes = json.dumps(payload).encode("utf-8")

    req = urllib.request.Request(url=url, data=cuerpo_bytes if metodo.upper() != "GET" else None, headers=headers, method=metodo.upper())
    t0 = time.perf_counter()
    try:
        with urllib.request.urlopen(req, timeout=timeout_s) as resp:
            _ = resp.read()
            latencia_ms = (time.perf_counter() - t0) * 1000.0
            codigo = resp.getcode() or 200
            exito = 1 if 200 <= codigo < 300 else 0
            return round(latencia_ms, 3), len(cuerpo_bytes), str(codigo), exito
    except urllib.error.HTTPError as e:
        _ = e.read()
        latencia_ms = (time.perf_counter() - t0) * 1000.0
        codigo = e.code
        # En escenarios de ráfaga/límite (slowapi), 429 es bloqueo esperado de control
        exito = 1 if codigo in (200, 401, 429) else 0
        return round(latencia_ms, 3), len(cuerpo_bytes), str(codigo), exito
    except Exception:
        latencia_ms = (time.perf_counter() - t0) * 1000.0
        return round(latencia_ms, 3), len(cuerpo_bytes), "ERR_HTTP", 0


def medir_go_tcp(host: str, port: int, mensaje: str, timeout_s: float) -> tuple[float, int, str, int]:
    datos = (mensaje.rstrip("\n") + "\n").encode("utf-8")
    t0 = time.perf_counter()
    try:
        with socket.create_connection((host, port), timeout=timeout_s) as sock:
            sock.settimeout(timeout_s)
            sock.sendall(datos)
            respuesta = sock.recv(4096)
            latencia_ms = (time.perf_counter() - t0) * 1000.0
            if respuesta:
                return round(latencia_ms, 3), len(datos), "OK_TCP", 1
            return round(latencia_ms, 3), len(datos), "EMPTY_TCP", 0
    except Exception:
        latencia_ms = (time.perf_counter() - t0) * 1000.0
        return round(latencia_ms, 3), len(datos), "ERR_TCP", 0


def ejecutar_campana(
    componente: str,
    escenario: str,
    n_muestras: int,
    pausa_ms: float,
    salida_csv: Path,
    url: str,
    metodo: str,
    payload_str: str,
    host: str,
    port: int,
    mensaje_tcp: str,
    timeout_s: float,
) -> None:
    salida_csv.parent.mkdir(parents=True, exist_ok=True)
    payload = json.loads(payload_str) if payload_str else None
    existe_archivo = salida_csv.exists()

    campos = [
        "id_medicion",
        "timestamp_utc",
        "componente",
        "escenario_condicion",
        "latencia_ms",
        "tamano_carga_bytes",
        "codigo_estado",
        "operacion_exitosa",
    ]

    with salida_csv.open("a", newline="", encoding="utf-8") as f:
        writer = csv.DictWriter(f, fieldnames=campos)
        if not existe_archivo:
            writer.writeheader()

        for i in range(1, n_muestras + 1):
            if componente == "fastapi_auth":
                lat_ms, tam_b, cod_est, exito = medir_fastapi_http(url, metodo, payload, timeout_s)
            elif componente == "go_tcp_chat":
                lat_ms, tam_b, cod_est, exito = medir_go_tcp(host, port, mensaje_tcp, timeout_s)
            else:
                raise ValueError(f"Componente no soportado: {componente}")

            writer.writerow(
                {
                    "id_medicion": f"M-{componente[:4].upper()}-{i:04d}",
                    "timestamp_utc": datetime.now(timezone.utc).isoformat(),
                    "componente": componente,
                    "escenario_condicion": escenario,
                    "latencia_ms": lat_ms,
                    "tamano_carga_bytes": tam_b,
                    "codigo_estado": cod_est,
                    "operacion_exitosa": exito,
                }
            )

            if pausa_ms > 0 and i < n_muestras:
                time.sleep(pausa_ms / 1000.0)


def main() -> None:
    parser = argparse.ArgumentParser(description="Medidor repetible de latencia y estado para el PI.")
    parser.add_argument("--componente", choices=["fastapi_auth", "go_tcp_chat"], required=True)
    parser.add_argument("--escenario", default="normal_secuencial")
    parser.add_argument("--n", type=int, default=30, help="Número de mediciones repetidas (recomendado >= 30)")
    parser.add_argument("--pausa-ms", type=float, default=50.0, help="Pausa entre mediciones en ms")
    parser.add_argument("--salida", type=Path, default=Path("research/data/mediciones_pi.csv"))
    parser.add_argument("--url", default="http://127.0.0.1:8000/login")
    parser.add_argument("--metodo", default="POST")
    parser.add_argument("--payload", default='{"username":"admin","password":"test"}')
    parser.add_argument("--host", default="127.0.0.1")
    parser.add_argument("--port", type=int, default=8080)
    parser.add_argument("--mensaje-tcp", default="PING_METRICA_PI")
    parser.add_argument("--timeout", type=float, default=5.0)
    args = parser.parse_args()

    ejecutar_campana(
        componente=args.componente,
        escenario=args.escenario,
        n_muestras=args.n,
        pausa_ms=args.pausa_ms,
        salida_csv=args.salida,
        url=args.url,
        metodo=args.metodo,
        payload_str=args.payload,
        host=args.host,
        port=args.port,
        mensaje_tcp=args.mensaje_tcp,
        timeout_s=args.timeout,
    )
    print(f"Mediciones guardadas en: {args.salida}")


if __name__ == "__main__":
    main()