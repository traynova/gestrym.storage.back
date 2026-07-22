# 📦 Gestrym Storage Backend (`gestrym.storage.back`)

![Go Version](https://img.shields.io/badge/Go-1.25%2B-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![Gin Framework](https://img.shields.io/badge/Gin-Gonic-008080?style=for-the-badge&logo=gin&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-GORM-4169E1?style=for-the-badge&logo=postgresql&logoColor=white)
![MinIO](https://img.shields.io/badge/MinIO-S3--Compatible-C72C48?style=for-the-badge&logo=minio&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-Containers-2496ED?style=for-the-badge&logo=docker&logoColor=white)

Microservicio dedicado al almacenamiento de archivos y gestión de persisencia de metadatos dentro de la plataforma **Gestrym**. Utiliza **MinIO** como capa de almacenamiento físico de objetos (S3 compatible) y **PostgreSQL** para la persistencia de metadatos con **GORM**.

---

## 🏛️ Arquitectura

El proyecto está diseñado siguiendo **Arquitectura Hexagonal (Puertos y Adaptadores)**:

```
src/
├── app.go                       # Inicializador de servidor Gin y configuraciones
├── common/
│   ├── config/                  # DB connection, variables de entorno, AutoMigrations GORM
│   ├── models/                  # Entidades GORM (ej. File)
│   └── routes/                  # Definición de rutas HTTP y DI (Dependency Injection)
└── storage/
    ├── domain/                  # Puertos e interfaces del dominio (File, Repository, StorageAdapter)
    ├── application/
    │   └── usecases/            # Casos de uso (UploadFile, GetFiles, DeleteFile)
    └── infrastructure/
        ├── adapters/            # PostgresFileRepository & MinioStorageAdapter
        └── http/
            └── handlers/        # Controladores/Handlers Gin HTTP
```

---

## 🚀 Características Principales

- **Carga Concurrente de Archivos**: Procesamiento paralelo y almacenamiento eficiente de múltiples archivos.
- **Validaciones de Seguridad y Tamaño**: Límite de tamaño máximo por archivo (**10 MB**) y verificación de tipos de contenido permitidos (`image/*`, `application/pdf`).
- **Agrupación por Colecciones**: Generación y asociación de un `collection_id` único para agrupar archivos pertenecientes a una misma entidad (ej. galerías de imágenes).
- **Borrado Lógico (Soft Delete)**: Desactivación por bandera `is_active = false` en base de datos sin eliminar físicamente de MinIO.
- **Integración Microservicios**: Entrega de URLs Firmadas (Presigned URLs) para un acceso temporal seguro a los recursos almacenados.

---

## 🛠️ Tecnologías

- **Lenguaje**: Go (v1.25+)
- **Framework Web**: [Gin-Gonic](https://github.com/gin-gonic/gin)
- **ORM / Base de Datos**: [GORM](https://gorm.io/) con PostgreSQL
- **Almacenamiento de Objetos**: [MinIO Go SDK](https://github.com/minio/minio-go) (S3 Compatible)
- **Documentación API**: Swagger / Swaggo

---

## 📡 API Endpoints

### Archivos (`/gestrym-storage/public/files`)

| Método | Endpoint | Descripción |
| :--- | :--- | :--- |
| `POST` | `/gestrym-storage/public/files/upload` | Carga de archivos (`multipart/form-data`). Retorna un `collection_id`. |
| `GET` | `/gestrym-storage/public/files/collection` | Obtiene los archivos/presigned URLs mediante un `collectionId`. |
| `DELETE` | `/gestrym-storage/public/files/:id` | Eliminación lógica del archivo por su ID (`is_active = false`). |

---

## 📋 Requisitos Previos e Instalación

### Requisitos
- **Go** >= 1.25
- **Docker** y **Docker Compose**
- Instancia activa de **PostgreSQL** y **MinIO**

### Ejecución Local

1. **Clonar el repositorio:**
   ```bash
   git clone https://github.com/tu-usuario/gestrym.storage.back.git
   cd gestrym.storage.back
   ```

2. **Instalar dependencias:**
   ```bash
   go mod download
   ```

3. **Ejecutar con Docker Compose:**
   ```bash
   docker-compose up -d
   ```

4. **Ejecutar localmente con Go:**
   ```bash
   go run main.go
   ```

---

## 📑 Documentación Swagger

La documentación interactiva de la API está integrada con Swaggo:

- **Ruta Swagger UI**: `http://localhost:<PUERTO>/swagger/index.html`

---

## 📜 Licencia

Desarrollado como parte del ecosistema de microservicios para la plataforma **Gestrym**.
