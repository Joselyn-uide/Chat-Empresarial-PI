import json
import os
import secrets
import bcrypt
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel

app = FastAPI(title="Microservicio de Autenticación - Chat Empresarial")

DB_FILE = "usuarios.json"

# Modelos de entrada
class UsuarioRegistro(BaseModel):
    username: str
    password: str
    rol: str = "Empleado"

class UsuarioLogin(BaseModel):
    username: str
    password: str

# Funciones auxiliares para leer/guardar usuarios
def cargar_usuarios() -> dict:
    if not os.path.exists(DB_FILE):
        return {}
    with open(DB_FILE, "r", encoding="utf-8") as f:
        return json.load(f)

def guardar_usuarios(usuarios: dict):
    with open(DB_FILE, "w", encoding="utf-8") as f:
        json.dump(usuarios, f, indent=4)

# 1. Endpoint de Registro de Usuario
@app.post("/register", status_code=201)
def registrar_usuario(datos: UsuarioRegistro):
    username = datos.username.strip()
    if len(username) < 3 or len(datos.password) < 6:
        raise HTTPException(status_code=400, detail="Usuario (min 3) o contraseña (min 6) muy cortos")

    usuarios = cargar_usuarios()
    if username in usuarios:
        raise HTTPException(status_code=409, detail="El usuario ya existe")

    # Hashing + Salt con Bcrypt (nunca texto plano)
    salt = bcrypt.gensalt()
    password_hash = bcrypt.hashpw(datos.password.encode("utf-8"), salt).decode("utf-8")

    usuarios[username] = {
        "username": username,
        "password_hash": password_hash,
        "rol": datos.rol
    }
    guardar_usuarios(usuarios)

    return {"mensaje": f"Usuario '{username}' registrado exitosamente", "rol": datos.rol}

# 2. Endpoint de Inicio de Sesión (Login)
@app.post("/login")
def iniciar_sesion(datos: UsuarioLogin):
    usuarios = cargar_usuarios()
    usuario = usuarios.get(datos.username.strip())

    # Verificar existencia y comparar hash
    if not usuario or not bcrypt.checkpw(
        datos.password.encode("utf-8"),
        usuario["password_hash"].encode("utf-8")
    ):
        raise HTTPException(status_code=401, detail="Credenciales inválidas")

    # Generar token de sesión temporal
    token_sesion = secrets.token_hex(16)
    return {
        "mensaje": "Autenticación exitosa",
        "username": usuario["username"],
        "rol": usuario["rol"],
        "token": token_sesion
    }