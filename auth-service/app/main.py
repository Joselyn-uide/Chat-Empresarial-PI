import json
import os
import secrets
import bcrypt
from fastapi import FastAPI, HTTPException, Depends
from fastapi.security import HTTPBearer, HTTPAuthorizationCredentials
from pydantic import BaseModel

app = FastAPI(title="Microservicio de Autenticación y Control de Acceso (RBAC) - Chat Empresarial")

DB_FILE = "usuarios.json"
ROLES_VALIDOS ={"Empleado", "Administrador"}

# Almacén en memoria de sesiones activas: { token: {"username": ..., "rol": ...} }
sesiones_activas = {}

# Esquema de seguridad nativo para Swagger UI (Bearer Token)
security = HTTPBearer()

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

# Dependencia de Autorización RBAC
def obtener_sesion_actual(credentials: HTTPAuthorizationCredentials = Depends(security)):
    token = credentials.credentials.strip()
    sesion = sesiones_activas.get(token)
    if not sesion:
        raise HTTPException(status_code=401, detail="Sesión inválida o expirada")
    return sesion

def requerir_roles(roles_permitidos: set):
    def verificador(sesion: dict = Depends(obtener_sesion_actual)):
        if sesion["rol"] not in roles_permitidos:
            raise HTTPException(
                status_code=403,
                detail=f"Acceso denegado: el rol '{sesion['rol']}' no tiene privilegios suficientes para este recurso administrativo."
            )
        return sesion
    return verificador

# 1. Endpoint de Registro de Usuario
@app.post("/register", status_code=201)
def registrar_usuario(datos: UsuarioRegistro):
    username = datos.username.strip()
    rol = datos.rol.strip()
    
    if len(username) < 3 or len(datos.password) < 6:
        raise HTTPException(status_code=400, detail="Usuario (min 3) o contraseña (min 6) muy cortos")

    if rol not in ROLES_VALIDOS:
        raise HTTPException(
            status_code=400,
            detail=f"Rol inválido. Los roles permitidos según el diseño de seguridad son: {', '.join(sorted(ROLES_VALIDOS))}"
        )
    
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

    return {"mensaje": f"Usuario '{username}' registrado exitosamente", "rol": rol}

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
    sesiones_activas[token_sesion] = {
        "username": usuario["username"],
        "rol": usuario["rol"]
    }
    return {
        "mensaje": "Autenticación exitosa",
        "username": usuario["username"],
        "rol": usuario["rol"],
        "token": token_sesion
    }
    
# 3. Recurso Operativo (Permitido para Empleados y Administradores)
@app.get("/chat/mensajes")
def acceder_chat_operativo(sesion: dict = Depends(requerir_roles({"Empleado", "Administrador"}))):
    return {
        "estado": "Acceso permitido",
        "usuario": sesion["username"],
        "rol": sesion["rol"],
        "recurso": "Historial de chat activo y envío de mensajes corporativos"
    }

# 4. Módulo Administrativo Exclusivo (Restringido solo a Administrador)
@app.post("/admin/usuarios")
def gestionar_cuentas_administrativas(sesion: dict = Depends(requerir_roles({"Administrador"}))):
    return {
        "estado": "Acceso administrativo concedido",
        "usuario": sesion["username"],
        "rol": sesion["rol"],
        "mensaje": "Acción administrativa ejecutada con éxito (Gestión de cuentas/auditoría)"
    }