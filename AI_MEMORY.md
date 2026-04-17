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
`FileName`, `ContentType`, `Size`, `URL`, `CollectionID` (grouping UUID), `Service` (origin service), and `IsActive`.

### 2. Microservice Integration Strategy
Other microservices (like Training or Auth) will communicate with this service either:
1. By HTTP REST interface, pushing the file as `multipart/form-data`. They will receive a `collection_id` which they should store in their own databases.
2. They can later retrieve files by sending the `collection_id`.
3. The Storage service manages MinIO, issues Presigned URLs for external access securely.

### 3. API Endpoints
- **POST `/gestrym-storage/public/files/upload`**: Uploads files concurrently. Returns a single `collection_id`. Expects multipart `files`, and optional `collectionId`, `service`.
- **GET `/gestrym-storage/public/files/collection`**: Retrieve files by `collectionId`. If only one file exists, returns the object; otherwise, an array.
- **DELETE `/gestrym-storage/public/files/:id`**: Logical deletion. Sets `is_active = false` in PostgreSQL. File is NOT removed from MinIO.

## 📝 Rules for Future AI Interactions
- **Logical Deletion**: Never delete records or physical files. Always use the `is_active` flag.
- **Collections**: Use `collection_id` for grouping files related to a single entity field (e.g., gallery of images for one exercise).
- **Return Shapes**: Upload endpoints should prioritize returning the `collection_id` for easy integration with other microservices.
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
