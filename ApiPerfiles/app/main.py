from fastapi import FastAPI
from app.routers import profile, health
import yaml
from fastapi.middleware.cors import CORSMiddleware

app = FastAPI()

# Registrar el router para perfiles de usuario
app.include_router(profile.router, prefix="/profile", tags=["User Profiles"])
app.include_router(health.router)

# Cargar el archivo OpenAPI YAML
def load_openapi_yaml():
    with open("openapi.yml", "r") as file:
        openapi_schema = yaml.safe_load(file)
    return openapi_schema

# Sobreescribir el esquema OpenAPI
def custom_openapi():
    if app.openapi_schema:
        return app.openapi_schema
    app.openapi_schema = load_openapi_yaml()
    return app.openapi_schema

app.openapi = custom_openapi

origins = [
    "http://localhost",
    "http://127.0.0.1",
]

# Añadir el middleware de CORS
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],  # Permitir estos orígenes
    allow_credentials=True,
    allow_methods=["*"],  # Permitir todos los métodos (GET, POST, etc.)
    allow_headers=["*"],  # Permitir todos los headers
)

# Ruta de ejemplo
@app.get("/")
async def root():
    return {"message": "API funcionando correctamente"}
