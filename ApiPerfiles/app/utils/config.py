from dotenv import load_dotenv
import os

# Cargar las variables del archivo .env
load_dotenv()

# Acceder a las variables de entorno
MONGO_DETAILS = os.getenv("MONGO_DETAILS")
SECRET_KEY = os.getenv("SECRET_KEY")
ACCESS_TOKEN_EXPIRE_MINUTES = int(os.getenv("ACCESS_TOKEN_EXPIRE_MINUTES", 30))
LOGS_API_URL = os.getenv("LOGS_API_URL")
NOTIFICATIONS_API_URL = os.getenv("NOTIFICATIONS_API_URL")

