#!/bin/bash

# Script para configurar el entorno de pruebas locales
# IT Messaging Service - Local Testing Setup

set -e

echo "🚀 Configurando entorno de pruebas locales para IT Messaging Service"
echo "=================================================================="

# Colores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Función para imprimir mensajes con colores
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Verificar que Docker esté corriendo
print_status "Verificando Docker..."
if ! docker info > /dev/null 2>&1; then
    print_error "Docker no está corriendo. Por favor inicia Docker Desktop."
    exit 1
fi
print_success "Docker está corriendo"

# Verificar que Go esté instalado
print_status "Verificando Go..."
if ! command -v go &> /dev/null; then
    print_error "Go no está instalado. Por favor instala Go desde https://golang.org/"
    exit 1
fi
print_success "Go está instalado: $(go version)"

# Limpiar contenedores anteriores si existen
print_status "Limpiando contenedores anteriores..."
docker-compose down -v 2>/dev/null || true

# Crear directorio de uploads si no existe
print_status "Creando directorio de uploads..."
mkdir -p uploads
print_success "Directorio uploads creado"

# Instalar dependencias de Go
print_status "Instalando dependencias de Go..."
go mod download
go mod tidy
print_success "Dependencias instaladas"

# Levantar servicios de infraestructura (PostgreSQL, Redis, Vault)
print_status "Levantando servicios de infraestructura..."
docker-compose up -d postgres redis vault prometheus

# Esperar a que PostgreSQL esté listo
print_status "Esperando a que PostgreSQL esté listo..."
sleep 10

# Verificar conexión a PostgreSQL
print_status "Verificando conexión a PostgreSQL..."
max_attempts=30
attempt=1
while [ $attempt -le $max_attempts ]; do
    if docker-compose exec -T postgres pg_isready -U postgres > /dev/null 2>&1; then
        print_success "PostgreSQL está listo"
        break
    fi
    
    if [ $attempt -eq $max_attempts ]; then
        print_error "PostgreSQL no está respondiendo después de $max_attempts intentos"
        exit 1
    fi
    
    print_status "Intento $attempt/$max_attempts - Esperando PostgreSQL..."
    sleep 2
    ((attempt++))
done

# Generar tokens JWT para testing
print_status "Generando tokens JWT para testing..."
go run scripts/generate-jwt.go > jwt-tokens.txt
print_success "Tokens JWT generados en jwt-tokens.txt"

# Mostrar información de servicios
echo ""
echo "🎉 Entorno de pruebas configurado exitosamente!"
echo "=============================================="
echo ""
echo "📋 Servicios disponibles:"
echo "  • PostgreSQL: localhost:5433 (usuario: postgres, password: postgres)"
echo "  • Redis: localhost:6379"
echo "  • Vault: http://localhost:8200 (token: dev-token)"
echo "  • Prometheus: http://localhost:9090"
echo ""
echo "🔧 Para iniciar la aplicación:"
echo "  make run"
echo "  # o directamente:"
echo "  go run ."
echo ""
echo "📡 Endpoints de la API:"
echo "  • Health: http://localhost:8080/api/v1/health"
echo "  • Swagger: http://localhost:8080/swagger/index.html"
echo "  • API Base: http://localhost:8080/api/v1"
echo ""
echo "🔑 Tokens JWT para testing:"
echo "  Ver archivo: jwt-tokens.txt"
echo ""
echo "📮 Postman Collection:"
echo "  Importa el archivo: postman_collection.json"
echo ""
echo "🧪 Comandos útiles:"
echo "  make test          # Ejecutar tests"
echo "  make docker-logs   # Ver logs de contenedores"
echo "  make docker-down   # Detener todos los servicios"
echo ""

# Mostrar los primeros tokens generados
echo "🔑 Tokens JWT generados:"
echo "========================"
head -10 jwt-tokens.txt
echo ""
echo "Ver jwt-tokens.txt para todos los tokens disponibles"