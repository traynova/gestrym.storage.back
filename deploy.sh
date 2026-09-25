#!/usr/bin/env bash
set -eo pipefail

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$DIR"

BOLD='\033[1m'
GREEN='\033[0;32m'
CYAN='\033[0;36m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

SERVICE_NAME="gestrym-storage"
IMAGE_BASE="jhonnier/back-storage"

VERSION="${1:-}"

if [ -z "$VERSION" ]; then
    if [ -f stack.env ] && grep -q '^VERSION=' stack.env; then
        CURRENT_VER=$(grep '^VERSION=' stack.env | cut -d'=' -f2 || true)
    fi
    VERSION="${CURRENT_VER:-1.0.0}"
fi
VERSION="${VERSION#:}"

echo -e "${CYAN}${BOLD}🚀 Desplegando ${SERVICE_NAME} - Versión: ${VERSION}${NC}"

if [ ! -f stack.env ]; then
    if [ -f stack.env.example ]; then
        echo -e "${YELLOW}Creando stack.env a partir de stack.env.example...${NC}"
        cp stack.env.example stack.env
    else
        echo -e "${RED}Error: stack.env no encontrado.${NC}"
        exit 1
    fi
fi

set_env_var() {
    local key="$1"
    local value="$2"
    if grep -q "^${key}=" stack.env; then
        sed -i.bak "s|^${key}=.*|${key}=${value}|" stack.env && rm -f stack.env.bak
    else
        echo "${key}=${value}" >> stack.env
    fi
}

set_env_var "VERSION" "${VERSION}"
set_env_var "STORAGE_IMAGE" "${IMAGE_BASE}:${VERSION}"

if ! docker network inspect gestrym-network >/dev/null 2>&1; then
    echo -e "${YELLOW}Creando red externa gestrym-network...${NC}"
    docker network create gestrym-network
fi

echo -e "${BOLD}Descargando imagen ${IMAGE_BASE}:${VERSION}...${NC}"
docker compose --env-file stack.env pull || echo -e "${YELLOW}Aviso: No se pudo descargar de Docker Hub, usando build local...${NC}"

docker compose --env-file stack.env up -d --remove-orphans

echo -e "\n${BOLD}Estado del contenedor:${NC}"
docker compose --env-file stack.env ps

echo -e "\n${GREEN}${BOLD}✅ ${SERVICE_NAME} desplegado con éxito en versión ${VERSION}.${NC}"
