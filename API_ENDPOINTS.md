# API Endpoints

## Base URL
- **Development**: `http://localhost:4100`
- **Production**: `https://api.cattfilemanager.com`

---

## Backend Endpoints (Go — Echo Framework)

### Auth Routes
**Group**: No prefix  
**Middleware**: `IsAuthenticated` only on `/logout`  
**Handler**: `internal/http/handlers/auth/authHandler.go`

| Method | Endpoint | Auth |
|--------|----------|------|
| POST | `/login` | No |
| POST | `/register` | No |
| POST | `/logout` | Yes |
| POST | `/refresh` | No |
| POST | `/forgot-password` | No |
| POST | `/recovery-password` | No |

---

### Client Routes
**Group Prefix**: `/client`  
**Middleware**: `IsAuthenticated` (all)  
**Handler**: `internal/http/handlers/clients/client.go`

| Method | Endpoint | Auth |
|--------|----------|------|
| GET | `/client` | Yes |
| GET | `/client/search` | Yes |
| POST | `/client/new` | Yes |
| PUT | `/client/edit/:id` | Yes |
| DELETE | `/client/:id` | Yes |
| GET | `/client/with-projects` | Yes |

---

### Project Routes
**Group Prefix**: `/project`  
**Middleware**: `IsAuthenticated` (all)  
**Handler**: `internal/http/handlers/projects/project.go`

| Method | Endpoint | Auth |
|--------|----------|------|
| GET | `/project` | Yes |
| GET | `/project/search` | Yes |
| POST | `/project/new` | Yes |
| PUT | `/project/edit/:id` | Yes |

> No DELETE endpoint implemented yet.

---

### File Routes
**Group Prefix**: `/files`  
**Middleware**: `IsAuthenticated` (all)  
**Handler**: `internal/http/handlers/files/file.go`

| Method | Endpoint | Auth |
|--------|----------|------|
| GET | `/files` | Yes |
| GET | `/files/project/:id` | Yes |
| POST | `/files/upload` | Yes |
| DELETE | `/files/:id` | Yes |
| GET | `/files/search` | Yes |
| GET | `/files/:id/download` | Yes |
| GET | `/files/:id/preview` | Yes |
| PUT | `/files/edit/:id` | Yes |

---

### User Routes
**Group Prefix**: `/users`  
**Middleware**: `IsAuthenticated` (all)  
**Handler**: `internal/http/handlers/users/user.go`

| Method | Endpoint | Auth |
|--------|----------|------|
| GET | `/users` | Yes |
| GET | `/users/search` | Yes |
| GET | `/users/:id` | Yes |
| PUT | `/users/edit/:id` | Yes |
| DELETE | `/users/:id` | Yes |

---

### TypeFile Routes
**Group Prefix**: `/type-files`  
**Middleware**: `IsAuthenticated` (all)  
**Handler**: `internal/http/handlers/typeFile.go`

| Method | Endpoint | Auth |
|--------|----------|------|
| GET | `/type-files` | Yes |
| POST | `/type-files/new` | Yes |

---

## Total: 31 live endpoints

---

## Frontend API Routes (`apiUrls.ts`)

All endpoints use `withCredentials: true` and are defined in `src/app/core/routes/apiUrls.ts`:

### Auth
| Key | Method | Path |
|-----|--------|------|
| `auth.login` | POST | `{BASE_URL}/login` |
| `auth.register` | POST | `{BASE_URL}/register` |
| `auth.logout` | POST | `{BASE_URL}/logout` |
| `auth.refreshToken` | POST | `{BASE_URL}/refresh` |
| `auth.forgotPassword` | POST | `{BASE_URL}/forgot-password` |
| `auth.recoveryPassword` | POST | `{BASE_URL}/recovery-password` |

### Users
| Key | Method | Path |
|-----|--------|------|
| `users.getAll` | GET | `{BASE_URL}/users` |
| `users.getById` | GET | `{BASE_URL}/users/{id}` |
| `users.update` | PUT | `{BASE_URL}/users/edit/{id}` |
| `users.delete` | DELETE | `{BASE_URL}/users/{id}` |
| `users.search` | GET | `{BASE_URL}/users/search?query=` |

### Files
| Key | Method | Path |
|-----|--------|------|
| `files.getAll` | GET | `{BASE_URL}/files` |
| `files.getByProjectId` | GET | `{BASE_URL}/files/project/{id}` |
| `files.upload` | POST | `{BASE_URL}/files/upload` |
| `files.delete` | DELETE | `{BASE_URL}/files/{id}` |
| `files.search` | GET | `{BASE_URL}/files/search?q=` |
| `files.download` | GET | `{BASE_URL}/files/{id}/download` |
| `files.preview` | GET | `{BASE_URL}/files/{id}/preview` |
| `files.update` | PUT | `{BASE_URL}/files/edit/{id}` |
| `files.getFileTypes` | GET | `{BASE_URL}/type-files` |
| `files.createFileTypes` | POST | `{BASE_URL}/type-files/new` |

### Clients
| Key | Method | Path |
|-----|--------|------|
| `clients.getAll` | GET | `{BASE_URL}/client` |
| `clients.search` | GET | `{BASE_URL}/client/search` |
| `clients.create` | POST | `{BASE_URL}/client/new` |
| `clients.update` | PUT | `{BASE_URL}/client/edit/{id}` |
| `clients.withProjects` | GET | `{BASE_URL}/client/with-projects` |

### Projects
| Key | Method | Path |
|-----|--------|------|
| `projects.getAll` | GET | `{BASE_URL}/project` |
| `projects.search` | GET | `{BASE_URL}/project/search` |
| `projects.create` | POST | `{BASE_URL}/project/new` |
| `projects.update` | PUT | `{BASE_URL}/project/edit/{id}` |

---

## Frontend ↔ Backend Alignment

All frontend routes now match the backend. No mismatches.
