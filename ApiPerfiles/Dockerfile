# Usa la imagen de Python
FROM python:3.9

# Establece el directorio de trabajo
WORKDIR /app

# Copia los archivos necesarios
COPY ./app /app/app
COPY ./requirements.txt /app

# Instala las dependencias
RUN pip install --no-cache-dir -r /app/requirements.txt

# Variable de entorno para FastAPI
ENV PYTHONUNBUFFERED=1

# Exponer el puerto de la API
EXPOSE 8001

# Ejecuta el servidor de FastAPI
CMD ["uvicorn", "app.main:app", "--host", "0.0.0.0", "--port", "8001"]
