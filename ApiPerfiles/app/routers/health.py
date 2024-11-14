from fastapi import APIRouter
from datetime import datetime, timedelta, timezone

router = APIRouter()

# Guardamos el tiempo de inicio del servicio
start_time = datetime.now(timezone.utc)  # Usamos timezone-aware datetime
service_version = "1.0.0"  # Definimos la versión del servicio

# Liveness check - /health/live
@router.get("/health/live", response_model=dict)
async def liveness_check():
    uptime = datetime.now(timezone.utc) - start_time
    response = {
        "status": "UP",
        "checks": [
            {
                "data": {
                    "from": datetime.now(timezone.utc).isoformat(),
                    "status": "ALIVE"
                },
                "name": "Liveness check",
                "status": "UP"
            }
        ],
        "details": {
            "version": service_version,
            "uptime": str(uptime)  # Convertimos timedelta a string legible
        }
    }
    return response

# Readiness check - /health/ready
@router.get("/health/ready", response_model=dict)
async def readiness_check():
    is_db_connected = True  # Simulamos conexión a base de datos
    status = "UP" if is_db_connected else "DOWN"
    
    uptime = datetime.now(timezone.utc) - start_time
    response = {
        "status": status,
        "checks": [
            {
                "data": {
                    "from": datetime.now(timezone.utc).isoformat(),
                    "status": "READY"
                },
                "name": "Readiness check",
                "status": status
            }
        ],
        "details": {
            "version": service_version,
            "uptime": str(uptime)
        }
    }
    return response

# General health status - /health
@router.get("/health", response_model=dict)
async def general_health():
    uptime = datetime.now(timezone.utc) - start_time
    response = {
        "status": "UP",
        "checks": [
            {
                "data": {
                    "from": datetime.now(timezone.utc).isoformat(),
                    "status": "READY"
                },
                "name": "Readiness check",
                "status": "UP"
            },
            {
                "data": {
                    "from": datetime.now(timezone.utc).isoformat(),
                    "status": "ALIVE"
                },
                "name": "Liveness check",
                "status": "UP"
            }
        ],
        "details": {
            "version": service_version,
            "uptime": str(uptime)
        }
    }
    return response
