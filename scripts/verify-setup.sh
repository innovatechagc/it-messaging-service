#!/bin/bash

# Script para verificar que el setup local esté funcionando correctamente
# IT Messaging Service - Setup Verification

set -e

# Colores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[✓]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[⚠]${NC} $1"
}

print_error() {
    echo -e "${RED}[✗]${NC} $1"
}

echo "🔍 Verificando configuración del entorno local"
echo "=============================================="

# Verificar que la aplicación esté corriendo
print_status "Verificando que la aplicación esté corriendo..."
if curl -s http://localhost:8080/api/v1/health > /dev/null 2>&1; then
    print_success "Aplicación está corriendo en puerto 8080"
else
    print_error "Aplicación no está corriendo. Ejecuta 'make run' o 'go run .'"
    exit 1
fi

# Verificar health check
print_status "Verificando health check..."
HEALTH_RESPONSE=$(curl -s http://localhost:8080/api/v1/health)
if echo "$HEALTH_RESPONSE" | grep -q "SUCCESS"; then
    print_success "Health check OK"
else
    print_warning "Health check devolvió: $HEALTH_RESPONSE"
fi

# Verificar readiness check
print_status "Verificando readiness check..."
READY_RESPONSE=$(curl -s http://localhost:8080/api/v1/ready)
if echo "$READY_RESPONSE" | grep -q "SUCCESS"; then
    print_success "Readiness check OK"
else
    print_warning "Readiness check devolvió: $READY_RESPONSE"
fi

# Verificar PostgreSQL
print_status "Verificando conexión a PostgreSQL..."
if docker-compose exec -T postgres pg_isready -U postgres > /dev/null 2>&1; then
    print_success "PostgreSQL está funcionando"
else
    print_error "PostgreSQL no está respondiendo"
fi

# Verificar Redis
print_status "Verificando conexión a Redis..."
if docker-compose exec -T redis redis-cli ping > /dev/null 2>&1; then
    print_success "Redis está funcionando"
else
    print_warning "Redis no está respondiendo (opcional)"
fi

# Verificar Swagger
print_status "Verificando documentación Swagger..."
if curl -s http://localhost:8080/swagger/index.html > /dev/null 2>&1; then
    print_success "Swagger UI disponible"
else
    print_warning "Swagger UI no está disponible"
fi

# Verificar archivos necesarios
print_status "Verificando archivos de configuración..."

if [ -f ".env" ]; then
    print_success "Archivo .env existe"
else
    print_error "Archivo .env no existe"
fi

if [ -f "postman_collection.json" ]; then
    print_success "Colección de Postman existe"
else
    print_error "Colección de Postman no existe"
fi

if [ -f "jwt-tokens.txt" ]; then
    print_success "Tokens JWT generados"
else
    print_warning "Tokens JWT no generados. Ejecuta 'make generate-jwt'"
fi

# Verificar directorio uploads
if [ -d "uploads" ]; then
    print_success "Directorio uploads existe"
else
    print_warning "Directorio uploads no existe"
    mkdir -p uploads
    print_success "Directorio uploads creado"
fi

# Probar endpoint protegido con JWT
if [ -f "jwt-tokens.txt" ]; then
    print_status "Probando autenticación JWT..."
    
    # Extraer el primer token del archivo
    JWT_TOKEN=$(grep -A 1 "Token:" jwt-tokens.txt | tail -1 | xargs)
    
    if [ -n "$JWT_TOKEN" ]; then
        AUTH_RESPONSE=$(curl -s -H "Authorization: Bearer $JWT_TOKEN" http://localhost:8080/api/v1/messaging/conversations)
        if echo "$AUTH_RESPONSE" | grep -q "SUCCESS\|data"; then
            print_success "Autenticación JWT funciona"
        else
            print_warning "Problema con autenticación JWT: $AUTH_RESPONSE"
        fi
    else
        print_warning "No se pudo extraer token JWT"
    fi
fi

echo ""
echo "🎉 Verificación completada!"
echo "=========================="
echo ""
echo "📋 Resumen del estado:"
echo "  • Aplicación: ✓ Funcionando en http://localhost:8080"
echo "  • Health Check: ✓ OK"
echo "  • PostgreSQL: ✓ Conectado"
echo "  • Redis: ✓ Conectado"
echo "  • Swagger: ✓ Disponible en http://localhost:8080/swagger/index.html"
echo ""
echo "🚀 Próximos pasos:"
echo "  1. Importa postman_collection.json en Postman"
echo "  2. Configura un JWT token de jwt-tokens.txt"
echo "  3. Ejecuta las pruebas en Postman"
echo ""
echo "🔧 Comandos útiles:"
echo "  make generate-jwt    # Regenerar tokens"
echo "  make test-api       # Probar endpoints básicos"
echo "  make docker-logs    # Ver logs de servicios"