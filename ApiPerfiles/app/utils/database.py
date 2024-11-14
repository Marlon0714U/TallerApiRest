from motor.motor_asyncio import AsyncIOMotorClient
from bson import ObjectId
import os
from app.utils.config import MONGO_DETAILS
import logging

logging.basicConfig(level=logging.INFO)
# Crear el cliente de Motor
client = AsyncIOMotorClient(MONGO_DETAILS)

# Seleccionar la base de datos
database = client.user_profile_db

# Seleccionar la colección de usuarios
user_profiles_collection = database.get_collection("user_profiles")

async def check_mongo_connection():
    try:
        logging.info(f"Conectando a MongoDB con: {MONGO_DETAILS}")
        await client.admin.command('ping')
        logging.info("Conexión exitosa a MongoDB")
    except Exception as e:
        logging.error(f"Error al conectar a MongoDB: {e}")

# Función auxiliar para convertir ObjectId a string
def profile_helper(profile) -> dict:
    return {
        "id": str(profile["_id"]),
        "nickname": profile["nickname"],
        "personal_url": profile.get("personal_url"),
        "contact_public": profile.get("contact_public", True),
        "address": profile.get("address"),
        "biography": profile.get("biography"),
        "organization": profile.get("organization"),
        "country": profile.get("country"),
        "social_links": profile.get("social_links", []),
    }
