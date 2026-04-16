# Gestrym Storage Backend - AI Memory

## 🏗️ Project Overview
This is the `Gestrym.Storage.Back` microservice, responsible for handling all file uploads and storage persistence logic across the Gestrym platform. It uses **MinIO** as the main storage infrastructure layer (S3-compatible object storage) and **PostgreSQL** (via GORM) to save file metadata.

## 📐 Architecture
The project strictly follows **Hexagonal Architecture** (Ports and Adapters):
- **Domain**: Interfaces (Ports) for the Repository and Storage Adapter, and the Domain `File` model (`src/storage/domain`).
- **Application**: The core logic (Use Cases) orchestrating inputs and interacting via Domain Ports (`src/storage/application/usecases`).
  - Added validations (e.g. file *size limit of 10MB* and *content-type checks* for standard image/pdf files).
  - Added support for concurrent multiple file uploads.
- **Infrastructure**: Adapters connecting to external libraries or I/O.
  - `MinioStorageAdapter` for physical file management.
  - `PostgresFileRepository` for metadata management.
  - HTTP definitions in `handlers/StorageHandler.go`.
- **Common Module**: Code shared across the backend service templates such as Config, Middleware, Models, Routes. 

> ⚠️ Existing models unrelated to User/Auth/Access (like Exercise, GymProfile) were stripped out to keep this context pure and lightweight as per the requirements.

## 📂 Folder Structure
```
.
├── deployment/       # Envoy / configuration files (e.g., env_local.yaml)
├── docs/             # Swagger auto-generated documentation
├── src/
│   ├── app.go        # Startup definitions, config and Gin initializer
│   ├── common/
│   │   ├── config/   # Environment setups, Database connections, and GORM Migrations 
│   │   ├── models/   # GORM Database Entities (e.g. File)
│   │   ├── routes/   # HTTP route bindings and DI container wiring 
│   ├── storage/
│   │   ├── domain/   # Interface ports
│   │   ├── application/
│   │   │   ├── usecases/   # UploadFile, GetFiles, DeleteFile
│   │   ├── infrastructure/
│   │   │   ├── adapters/   # Postgre & Minio Adapters
│   │   │   ├── http/
│   │   │   │   ├── handlers/ # Gin Handlers
├── Dockerfile        # Production Docker configuration
├── main.go           # Entrypoint
```

## 🧩 Modulo: Storage-Service
### 1. Model: `File` (`src/common/models/File.go`)
Central database representation containing:
`FileName`, `ContentType`, `Size`, `URL`, `Collection`, `EntityID`, and `EntityType`.

### 2. Microservice Integration Strategy
Other microservices (like Training or Auth) will communicate with this service either:
1. By HTTP REST interface, pushing the file as `multipart/form-data`, passing their contextual IDs (`EntityID` and `EntityType`).
2. Getting signed URLs to stream directly, or receiving the generated standard URL. The Storage service manages MinIO, issues Presigned URLs for external access securely.

### 3. API Endpoints
- **POST `/gestrym-storage/public/files/upload`**: Uploads files concurrently (expects multipart `files`, `collection`, `entityId`, `entityType`). Validate size (<= 10MB) and type.
- **GET `/gestrym-storage/public/files`**: Retrieve multiple file models (with ephemeral Pre-Signed URLs embedded) based on query params (`entityId`, `entityType`).
- **DELETE `/gestrym-storage/public/files/:id`**: Removes file from MinIO and deletes record from PostgreSQL.

## 📝 Rules for Future AI Interactions
- **Do Not Break existing Hexagonal pattern**. Dependency Injection happens at the `routes` level (`ServerRoutesDefinition.go`). Handlers, Adapters, and Repositories must always have clear interfaces in `domain/`.
- Models used for persistence must stay inside `src/common/models` and registered in `config/Migrations.go`. Keep models unique to their bounded context.
- Adhere to the pre-existing error handling shapes and HTTP status returns (`gin.H{...}`).
- Whenever you add an endpoint, document it with Swaggo formatted comments.

## 💻 Tech Stack
- Go 1.25+
- Gin-Gonic
- GORM (PostgreSQL)
- MinIO (minio-go/v7)
- Swagger (Swaggo)
