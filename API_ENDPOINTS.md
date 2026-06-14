# API Endpoints Documentation

## Base URL
- **Development**: `http://localhost:4100`
- **Production**: `https://api.cattfilemanager.com`

---

## Backend Endpoints (Go - Echo Framework)

### Auth Routes
**Group**: No authentication required (except logout)  
**Handler**: `internal/http/handlers/auth/authHandler.go`

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/login` | User login | No |
| POST | `/register` | User registration | No |
| POST | `/logout` | User logout | Yes |
| POST | `/refresh` | Refresh JWT token | No |

---

### Client Routes
**Group Prefix**: `/client`  
**Middleware**: `middlewares.IsAuthenticated` (all routes)  
**Handler**: `internal/http/handlers/clients/client.go`

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/client` | Get all clients | Yes |
| GET | `/client/search` | Search clients (query: `q`) | Yes |
| POST | `/client/new` | Create new client | Yes |
| PUT | `/client/edit/:id` | Update client by ID | Yes |
| DELETE | `/client/:id` | Delete client by ID | Yes |

---

### Project Routes
**Group Prefix**: `/project`  
**Middleware**: `middlewares.IsAuthenticated` (all routes)  
**Handler**: `internal/http/handlers/projects/project.go`

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/project` | Get all projects | Yes |
| GET | `/project/search` | Search projects (query: `q`) | Yes |
| POST | `/project/new` | Create new project | Yes |
| PUT | `/project/edit/:id` | Update project by ID | Yes |

---

### File Routes
**Group Prefix**: `/files`  
**Middleware**: None  
**Handler**: `internal/http/handlers/files/file.go`

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/files` | Get all files | No |
| GET | `/files/project/:id` | Get files by project ID | No |
| POST | `/files/upload` | Upload a file | No |

---

### TypeFile Routes
**Group Prefix**: `/type-files`  
**Middleware**: `middlewares.IsAuthenticated` (all routes)  
**Handler**: `internal/http/handlers/typeFile.go`

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/type-files` | Get all file types | Yes |
| POST | `/type-files/new` | Bulk create file types | Yes |

---

## Frontend Endpoints (Angular)

### API Route Definitions
**File**: `src/app/core/routes/apiUrls.ts`

### Auth Endpoints
| Method | Endpoint | Description | Service |
|--------|----------|-------------|---------|
| POST | `{BASE_URL}/login` | User login | `auth-state.service.ts` |
| POST | `{BASE_URL}/register` | User registration | `auth-state.service.ts` |
| POST | `{BASE_URL}/logout` | User logout | `auth-service.ts` |
| POST | `{BASE_URL}/refresh` | Refresh token | - |

---

### Users Endpoints
**Note**: These endpoints are defined in frontend but may not be implemented in backend yet.

| Method | Endpoint | Description | Service |
|--------|----------|-------------|---------|
| GET | `{BASE_URL}/users` | Get all users | - |
| GET | `{BASE_URL}/users/{id}` | Get user by ID | - |
| PUT | `{BASE_URL}/users/{id}` | Update user | - |
| DELETE | `{BASE_URL}/users/{id}` | Delete user | - |

---

### Files Endpoints
| Method | Endpoint | Description | Service |
|--------|----------|-------------|---------|
| GET | `{BASE_URL}/files` | Get all files | `files.ts` |
| GET | `{BASE_URL}/files/{id}` | Get file by ID | - |
| POST | `{BASE_URL}/files/upload` | Upload file | `file-service.ts` |
| GET | `{BASE_URL}/files/project/{id}` | Get files by project | `file-service.ts` |
| GET | `{BASE_URL}/files/by-ids` | Get files by multiple IDs | - |
| DELETE | `{BASE_URL}/files/{id}` | Delete file | `file-service.ts` |
| GET | `{BASE_URL}/files/{id}/download` | Download file | `file-service.ts` (empty) |
| GET | `{BASE_URL}/type-files` | Get file types | `file-service.ts` |
| POST | `{BASE_URL}/type-files` | Create file types | `file-service.ts` |

---

### Clients Endpoints
**Note**: Frontend defines `/client/delete/{id}` but backend uses `/client/:id` for DELETE.

| Method | Endpoint | Description | Service |
|--------|----------|-------------|---------|
| GET | `{BASE_URL}/client` | Get all clients | `client-service.ts` |
| GET | `{BASE_URL}/client/search` | Search clients (query: `q`) | `client-service.ts` |
| POST | `{BASE_URL}/client/new` | Create client | `client-service.ts` |
| GET | `{BASE_URL}/client/{id}` | Get client by ID | - |
| PUT | `{BASE_URL}/client/edit/{id}` | Update client | `client-service.ts` |
| DELETE | `{BASE_URL}/client/delete/{id}` | Delete client | - |

---

### Projects Endpoints
**Note**: Frontend defines `/project/delete/{id}` but backend uses `/project/:id` for DELETE (not implemented in backend yet).

| Method | Endpoint | Description | Service |
|--------|----------|-------------|---------|
| GET | `{BASE_URL}/project` | Get all projects | `project-service.ts` |
| GET | `{BASE_URL}/project/search` | Search projects (query: `q`) | `project-service.ts` |
| POST | `{BASE_URL}/project/new` | Create project | `project-service.ts` |
| GET | `{BASE_URL}/project/{id}` | Get project by ID | - |
| PUT | `{BASE_URL}/project/edit/{id}` | Update project | `project-service.ts` |
| DELETE | `{BASE_URL}/project/delete/{id}` | Delete project | - |

---

### Additional Frontend-Only Endpoints
These are defined in the frontend but may not have backend implementation:

| Method | Endpoint | Description | Service |
|--------|----------|-------------|---------|
| GET | `{BASE_URL}/files/{id}/download` | Download file | `file-service.ts` (empty) |
| GET | `{BASE_URL}/files/by-ids` | Get files by multiple IDs | - |

---

## Endpoint Mismatch Notes

1. **Client DELETE**: Frontend uses `/client/delete/{id}`, Backend uses `/client/:id`
2. **Project DELETE**: Frontend defines `/project/delete/{id}`, Backend doesn't have DELETE route implemented
3. **Users endpoints**: Defined in frontend but not implemented in backend
4. **Download endpoint**: Defined in frontend service but method is empty
5. **File by ID**: Frontend defines `/files/{id}` but backend doesn't have this endpoint

---

## Authentication Middleware

Backend routes protected with `middlewares.IsAuthenticated`:
- All `/client/*` routes
- All `/project/*` routes
- All `/type-files/*` routes
- `/logout` endpoint

Routes without authentication:
- `/login`, `/register`, `/refresh`
- All `/files/*` routes
