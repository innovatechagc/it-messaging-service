# Messaging Service

Servicio de mensajería en tiempo real para chatbots con soporte para múltiples canales de comunicación, archivos adjuntos y eventos pub/sub.

## 🚀 Características

- ✅ **Conversaciones multi-canal** (WhatsApp, Web, Messenger, Instagram)
- ✅ **Mensajes en tiempo real** con tipos de contenido variados
- ✅ **Archivos adjuntos** con clasificación automática
- ✅ **Caché con Redis** para alto rendimiento
- ✅ **Eventos pub/sub** para integraciones
- ✅ **Autenticación JWT** con middleware de seguridad
- ✅ **API REST** completamente documentada con Swagger
- ✅ **Base de datos PostgreSQL** con índices optimizados
- ✅ **Logging estructurado** con Zap
- ✅ **Health checks** y métricas con Prometheus
- ✅ **Docker** y despliegue en Google Cloud Run

## 📋 Entidades Principales

### Conversation
- `id`: UUID único
- `user_id`: ID del usuario
- `channel`: Canal de comunicación (whatsapp, web, messenger, instagram)
- `status`: Estado (active, closed, archived)
- `created_at`, `updated_at`: Timestamps

### Message
- `id`: UUID único
- `conversation_id`: Referencia a conversación
- `sender_type`: Tipo de remitente (user, bot, system)
- `sender_id`: ID del remitente
- `content`: Contenido del mensaje
- `content_type`: Tipo de contenido (text, image, video, audio, file)
- `metadata`: Datos adicionales en JSONB
- `timestamp`: Fecha y hora del mensaje

### Attachment
- `id`: UUID único
- `message_id`: Referencia al mensaje
- `url`: URL del archivo
- `type`: Tipo de archivo (image, video, file, audio)
- `size`: Tamaño en bytes
- `filename`: Nombre original del archivo

## 🛠 API Endpoints

### Base path: `/api/v1/messaging`

#### 🔁 Conversaciones
| Método | Ruta | Descripción |
|--------|------|-------------|
| `GET` | `/conversations` | Lista conversaciones activas |
| `GET` | `/conversations/:id` | Detalles de una conversación |
| `POST` | `/conversations` | Crea nueva conversación |
| `PATCH` | `/conversations/:id` | Actualiza estado de conversación |

#### ✉️ Mensajes
| Método | Ruta | Descripción |
|--------|------|-------------|
| `GET` | `/conversations/:id/messages` | Lista mensajes con paginación |
| `POST` | `/conversations/:id/messages` | Envía nuevo mensaje |
| `GET` | `/messages/:id` | Consulta mensaje individual |

#### 📎 Archivos Adjuntos
| Método | Ruta | Descripción |
|--------|------|-------------|
| `POST` | `/attachments/upload` | Sube archivo y devuelve URL |
| `GET` | `/attachments/:id` | Detalles de archivo adjunto |

## 🚀 Inicio Rápido

### Prerrequisitos

- Go 1.21+
- PostgreSQL 13+
- Redis 6+ (opcional, para caché y eventos)
- Docker y Docker Compose

### Instalación

1. **Clona el repositorio:**
```bash
git clone <repository-url>
cd messaging-service
```

2. **Configura el entorno:**
```bash
cp .env.example .env
# Edita .env con tus configuraciones
```

3. **Inicia la base de datos:**
```bash
docker-compose up -d postgres redis
```

4. **Ejecuta las migraciones:**
```bash
psql -h localhost -U postgres -d messaging_service -f scripts/init-messaging.sql
```

5. **Inicia el servicio:**
```bash
go mod tidy
go run main.go
```

### Con Docker Compose

```bash
docker-compose up -d
```

## 📖 Uso de la API

### Autenticación

Todas las rutas requieren un token JWT en el header:
```bash
Authorization: Bearer <your-jwt-token>
```

### Ejemplos de uso

#### Crear conversación
```bash
curl -X POST http://localhost:8080/api/v1/messaging/conversations \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"channel": "web"}'
```

#### Enviar mensaje
```bash
curl -X POST http://localhost:8080/api/v1/messaging/conversations/{id}/messages \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "sender_type": "user",
    "content": "Hola, necesito ayuda",
    "content_type": "text",
    "metadata": {"priority": "normal"}
  }'
```

#### Subir archivo
```bash
curl -X POST http://localhost:8080/api/v1/messaging/attachments/upload \
  -H "Authorization: Bearer <token>" \
  -F "file=@imagen.jpg"
```

## ⚙️ Configuración

### Variables de Entorno Principales

```bash
# Base de datos
DB_HOST=localhost
DB_NAME=messaging_service
DB_USER=postgres
DB_PASSWORD=your_password

# Redis (opcional)
REDIS_ENABLED=true
REDIS_HOST=localhost

# JWT
JWT_SECRET=your-secret-key
JWT_ISSUER=messaging-service

# Archivos
FILE_STORAGE_LOCAL_PATH=./uploads
FILE_STORAGE_MAX_SIZE=10485760

# Eventos
EVENTS_PROVIDER=redis
EVENTS_TOPIC=message.events
```

## 🔧 Funcionalidades Técnicas

### Caché con Redis
- Conversaciones recientes cacheadas por 30 minutos
- Mensajes cacheados por 10 minutos
- Invalidación automática en actualizaciones

### Eventos Pub/Sub
Cuando se recibe un mensaje nuevo, se publica un evento:
```json
{
  "type": "message.received",
  "conversation_id": "uuid",
  "message": { /* objeto mensaje completo */ },
  "timestamp": "2025-01-22T10:30:00Z"
}
```

### Seguridad
- Middleware JWT en todas las rutas protegidas
- Validación de propiedad de recursos por usuario
- Sanitización de archivos subidos
- Límites de tamaño de archivo configurables

## 📊 Monitoreo

### Health Checks
- `GET /api/v1/health` - Estado general
- `GET /api/v1/ready` - Readiness para tráfico

### Métricas Prometheus
- Requests HTTP por endpoint
- Duración de requests
- Errores por tipo
- Métricas de base de datos y Redis

### Logs Estructurados
```json
{
  "level": "info",
  "timestamp": "2025-01-22T10:30:00Z",
  "message": "Message sent",
  "message_id": "uuid",
  "conversation_id": "uuid",
  "user_id": "user123"
}
```

## 🧪 Testing

```bash
# Tests unitarios
go test ./...

# Tests con coverage
go test -cover ./...

# Tests de integración
go test -tags=integration ./tests/integration/...
```

## 🚢 Deployment

### Docker
```bash
docker build -t messaging-service .
docker run -p 8080:8080 messaging-service
```

### Google Cloud Run
```bash
gcloud run deploy messaging-service \
  --source . \
  --platform managed \
  --region us-central1 \
  --set-env-vars="DB_HOST=your-db-host,REDIS_HOST=your-redis-host"
```

## 📚 Documentación API

Una vez iniciado el servicio, la documentación Swagger estará disponible en:
`http://localhost:8080/swagger/index.html`

## 🤝 Contribución

1. Fork el proyecto
2. Crea una rama para tu feature (`git checkout -b feature/nueva-funcionalidad`)
3. Commit tus cambios (`git commit -m 'Agrega nueva funcionalidad'`)
4. Push a la rama (`git push origin feature/nueva-funcionalidad`)
5. Abre un Pull Request

## 📄 Licencia

Este proyecto está bajo la Licencia MIT - ver el archivo [LICENSE](LICENSE) para detalles.