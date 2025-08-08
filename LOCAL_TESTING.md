# Guía de Pruebas Locales - IT Messaging Service

Esta guía te ayudará a configurar y ejecutar pruebas locales del servicio de mensajería.

## 🚀 Inicio Rápido

### 1. Configuración Automática

Ejecuta el script de configuración automática:

```bash
./scripts/local-test-setup.sh
```

Este script:
- ✅ Verifica que Docker y Go estén instalados
- ✅ Levanta PostgreSQL, Redis, Vault y Prometheus
- ✅ Instala dependencias de Go
- ✅ Genera tokens JWT para testing
- ✅ Crea directorios necesarios

### 2. Iniciar la Aplicación

```bash
# Opción 1: Usando Makefile
make run

# Opción 2: Directamente con Go
go run .

# Opción 3: Con Docker Compose (aplicación completa)
make docker-dev
```

### 3. Verificar que Todo Funciona

```bash
# Health check
curl http://localhost:8080/api/v1/health

# Readiness check
curl http://localhost:8080/api/v1/ready
```

## 🔧 Configuración Manual

Si prefieres configurar manualmente:

### Prerrequisitos

- Docker y Docker Compose
- Go 1.21+
- Make (opcional)

### Variables de Entorno

El archivo `.env` ya está configurado para desarrollo local. Las variables principales son:

```env
ENVIRONMENT=development
PORT=8080
DB_HOST=localhost
DB_PORT=5433
JWT_SECRET=dev-jwt-secret-key-change-in-production
```

### Servicios de Infraestructura

```bash
# Levantar solo la infraestructura
docker-compose up -d postgres redis vault prometheus

# Ver logs
docker-compose logs -f

# Detener servicios
docker-compose down -v
```

## 🔑 Autenticación JWT

### Generar Tokens de Prueba

```bash
go run scripts/generate-jwt.go
```

Esto genera tokens para usuarios de prueba:
- `user1@example.com` (rol: user)
- `user2@example.com` (rol: user)  
- `admin@example.com` (rol: admin)

### Usar Tokens en Requests

```bash
# Ejemplo con curl
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
     http://localhost:8080/api/v1/messaging/conversations
```

## 📮 Pruebas con Postman

### Importar Colección

1. Abre Postman
2. Importa el archivo `postman_collection.json`
3. La colección incluye:
   - Variables de entorno preconfiguradas
   - Todos los endpoints de la API
   - Ejemplos de requests
   - Scripts para capturar IDs automáticamente

### Configurar Variables

En Postman, configura estas variables de colección:
- `base_url`: `http://localhost:8080`
- `jwt_token`: Copia un token del archivo `jwt-tokens.txt`

### Flujo de Pruebas Recomendado

1. **Health Checks**: Verifica que la API esté funcionando
2. **Create Conversation**: Crea una conversación (guarda el ID automáticamente)
3. **Send Messages**: Envía diferentes tipos de mensajes
4. **Upload Files**: Prueba la subida de archivos
5. **Get Data**: Consulta conversaciones y mensajes

## 🧪 Tipos de Mensajes Soportados

### Mensaje de Texto
```json
{
    "type": "text",
    "content": "Hola, este es un mensaje de prueba",
    "metadata": {
        "source": "postman_test"
    }
}
```

### Mensaje con IA
```json
{
    "type": "ai_response",
    "content": "¿En qué puedo ayudarte hoy?",
    "ai_context": {
        "model": "gpt-3.5-turbo",
        "prompt": "Responde de manera amigable",
        "confidence": 0.95
    }
}
```

### Mensaje con Archivo
```json
{
    "type": "file",
    "content": "Archivo adjunto",
    "attachments": [
        {
            "url": "http://localhost:8080/uploads/example.pdf",
            "filename": "documento.pdf",
            "type": "document",
            "size": 1024
        }
    ]
}
```

## 📊 Endpoints Disponibles

### Health & Status
- `GET /api/v1/health` - Health check
- `GET /api/v1/ready` - Readiness check

### Conversaciones
- `GET /api/v1/messaging/conversations` - Listar conversaciones
- `POST /api/v1/messaging/conversations` - Crear conversación
- `GET /api/v1/messaging/conversations/{id}` - Obtener conversación
- `PATCH /api/v1/messaging/conversations/{id}` - Actualizar estado

### Mensajes
- `GET /api/v1/messaging/conversations/{id}/messages` - Listar mensajes
- `POST /api/v1/messaging/conversations/{id}/messages` - Enviar mensaje
- `GET /api/v1/messaging/messages/{id}` - Obtener mensaje

### Archivos
- `POST /api/v1/messaging/attachments/upload` - Subir archivo
- `GET /api/v1/messaging/attachments/{id}` - Obtener detalles de archivo

### Documentación
- `GET /swagger/index.html` - Documentación Swagger

## 🔍 Debugging y Logs

### Ver Logs de la Aplicación
```bash
# Si ejecutas con go run
# Los logs aparecen en la consola

# Si usas Docker
docker-compose logs -f app
```

### Ver Logs de Base de Datos
```bash
docker-compose logs -f postgres
```

### Conectar a PostgreSQL
```bash
# Usando psql
psql -h localhost -p 5433 -U postgres -d messaging_service

# Usando Docker
docker-compose exec postgres psql -U postgres -d messaging_service
```

### Conectar a Redis
```bash
# Usando redis-cli
redis-cli -h localhost -p 6379

# Usando Docker
docker-compose exec redis redis-cli
```

## 🛠️ Comandos Útiles

```bash
# Ejecutar tests
make test

# Ejecutar tests con cobertura
make test-coverage

# Limpiar y reiniciar todo
make dev-reset

# Ver métricas
curl http://localhost:8080/metrics

# Formatear código
make format

# Ejecutar linter
make lint
```

## 🐛 Solución de Problemas

### La aplicación no inicia
1. Verifica que Docker esté corriendo
2. Verifica que los puertos no estén ocupados:
   ```bash
   lsof -i :8080  # Puerto de la app
   lsof -i :5433  # Puerto de PostgreSQL
   lsof -i :6379  # Puerto de Redis
   ```

### Error de conexión a base de datos
1. Verifica que PostgreSQL esté corriendo:
   ```bash
   docker-compose ps postgres
   ```
2. Verifica la conexión:
   ```bash
   docker-compose exec postgres pg_isready -U postgres
   ```

### Tokens JWT inválidos
1. Regenera los tokens:
   ```bash
   go run scripts/generate-jwt.go
   ```
2. Verifica que uses el token completo (incluyendo el Bearer)

### Archivos no se suben
1. Verifica que el directorio `uploads` exista y tenga permisos
2. Verifica el tamaño del archivo (máximo 10MB por defecto)

## 📚 Recursos Adicionales

- [Documentación Swagger](http://localhost:8080/swagger/index.html) (cuando la app esté corriendo)
- [Prometheus Metrics](http://localhost:9090) (métricas de la aplicación)
- [Vault UI](http://localhost:8200) (manejo de secretos)

## 🤝 Contribuir

Para contribuir al proyecto:
1. Ejecuta las pruebas: `make test-all`
2. Verifica el formato: `make format`
3. Ejecuta el linter: `make lint`
4. Asegúrate de que todos los endpoints funcionen con Postman