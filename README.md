## Guía de Configuración de Twilio y Nodemailer

En este proyecto usamos **Twilio** y **Nodemailer** para enviar notificaciones SMS y por correo electrónico. Aquí se explica cómo configurarlos para que el sistema funcione correctamente.

### Configuración de Twilio

Twilio es la herramienta que utilizamos para enviar SMS. Para configurarlo:

1. **Se crea una cuenta en [Twilio](https://www.twilio.com/)**.
2. **Se obtiene el SID de la cuenta y el Token de Autenticación** desde el panel de Twilio.
3. **Se configura el número de Twilio** para el envío de SMS.
4. Se deben agregar las variables de entorno en el archivo `.env`:
   ```plaintext
   TWILIO_ACCOUNT_SID=tu_account_sid
   TWILIO_AUTH_TOKEN=tu_auth_token
   TWILIO_PHONE_NUMBER=tu_numero_de_twilio
### Configuración de Nodemailer

Nodemailer nos permite enviar correos electrónicos. Para configurarlo:

1. **Configuramos las credenciales de correo electrónico** que usamos para enviar notificaciones (por ejemplo, una cuenta de Gmail).
2. Agregamos las variables de entorno en `.env`:
   ```plaintext
   EMAIL_HOST=smtp.ejemplo.com
   EMAIL_PORT=587
   EMAIL_USER=tu_usuario
   EMAIL_PASSWORD=tu_contraseña
3. Configuramos el transporter en services/notificationService.js para que use estas credenciales.

---

# Estructura General del Sistema de Notificaciones

Nuestro sistema de notificaciones tiene tres componentes principales:

1. **Modelo de Notificación** (`models/notification.js`): Define la estructura de las notificaciones en la base de datos.
2. **Rutas de Notificación** (`routes/notification.js`): Define las rutas API para enviar y consultar notificaciones.
3. **Servicio de Notificación** (`services/notificationService.js`): Contiene la lógica de envío a través de **Twilio** y **Nodemailer** para manejar notificaciones SMS y por correo electrónico.

---

## Funcionamiento de Cada Componente

### 1. Modelo de Notificación

En `models/notification.js` definimos un esquema de Mongoose para las notificaciones, el cual contiene:

- **`recipient`**: El destinatario (número de teléfono para SMS o correo electrónico).
- **`channels`**: Un array que especifica los tipos de canales (ej. `["sms"]` o `["email"]`).
- **`message`**: Contenido de la notificación.
- **`status`**: Estado actual de la notificación (puede ser "pendiente", "enviado", o "fallido").
- **`createdAt`**: Fecha de creación.

Cada vez que se envía una notificación, se crea una instancia de este modelo y se guarda en la base de datos.

### 2. Rutas de Notificación

En `routes/notification.js` definimos las rutas principales para manejar las notificaciones:

- **POST /notifications/send**: Permite enviar una nueva notificación.
  - Recibe un JSON con el tipo (`tipo`), mensaje (`mensaje`), y destinatario (`destinatario`).
  - Valida los datos y llama al servicio de notificación (`sendNotification`) para enviar la notificación.
  - Responde con un mensaje de éxito o un error, según el resultado.

- **GET /notifications**: Permite obtener una lista de notificaciones con filtros y paginación.
  - Los filtros incluyen `recipient`, `channels`, `status`, y rango de fechas (`startDate` y `endDate`).
  - Maneja paginación con los parámetros `page` y `limit`.

### 3. Servicio de Notificación

El servicio de notificación en `services/notificationService.js` es el núcleo del sistema, donde se realiza el envío mediante **Twilio** y **Nodemailer**.

#### Configuración de Twilio y Nodemailer:

- **Twilio**: Usamos las credenciales `accountSid`, `authToken`, y `number` para enviar SMS.
- **Nodemailer**: Configuramos un transporter con las credenciales de correo para enviar emails.

#### Función `sendNotification`:

1. **Creación y Guardado de Notificación**: 
   - Crea una entrada en la base de datos con estado "pendiente".

2. **Envío de Notificación**:
   - **SMS**: Si el tipo es `sms`, utiliza `client.messages.create` de Twilio. Si se envía correctamente, actualiza el estado a "enviado".
   - **Email**: Si el tipo es `email`, usa `transporter.sendMail` de Nodemailer. Al enviarse exitosamente, cambia el estado a "enviado".

3. **Manejo de Errores**: Si ocurre un error, actualiza el estado a "fallido" y registra el error.

---

## Ejemplo de Flujo Completo

1. **Solicitud de Envío**: Se envía una solicitud `POST` a `/notifications/send` con el tipo de notificación, mensaje y destinatario.
2. **Creación en Base de Datos**: La notificación se guarda con estado "pendiente".
3. **Envío de Notificación**:
   - **SMS**: Twilio envía el mensaje al número especificado.
   - **Email**: Nodemailer envía el correo al destinatario.
4. **Actualización de Estado**: El estado se actualiza a "enviado" o "fallido".
5. **Respuesta al Cliente**: La API responde con el resultado del envío.

---
