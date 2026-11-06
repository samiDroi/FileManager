# API Endpoints

## Base URL
- **Development**: `http://localhost:4100`
- **Production**: `https://api.cattfilemanager.com`

---

## Backend Endpoints (Go — Echo Framework)

### Auth Routes
**Group**: No prefix  
**Middleware**: `/logout` and `/refresh` use the refresh-token cookie
**Handler**: `internal/http/handlers/auth/authHandler.go`

| Method | Endpoint | Auth |
|--------|----------|------|
| POST | `/login` | No |
| POST | `/register` | No |
| POST | `/logout` | Refresh cookie |
| POST | `/refresh` | Refresh cookie |
| POST | `/forgot-password` | No |
| POST | `/recovery-password` | No |

`/login` and `/register` return `access_token` in the JSON body. The API sets a
7-day opaque refresh token only in an `HttpOnly`, `Secure`, `SameSite=Strict`
cookie. Protected routes require `Authorization: Bearer <access_token>`.
`/refresh` validates that cookie against the hashed token in `sessions` and
returns a new 15-minute access token in JSON. `/logout` revokes the associated
session and clears the refresh cookie. Login/register and refresh responses also
include the user's flat `roles`, `permissions`, and `is_admin` values for
conditional UI rendering. For example, a file edit route can enforce
`middlewares.RequirePermission(db, "archivos", "editar")`; an admin-only route
can enforce `middlewares.RequireAdmin(db)`. Both must follow
`middlewares.IsAuthenticated`.

`GET /roles` includes each role's `id`, `name`, and `description`.
`POST /users` accepts `username`, `email`, `password`, and `role_id`; its `201`
response includes the assigned role in `roles` with `id`, `name`, and
`description`, which the user-creation confirmation displays. `GET /users/:id`
also returns assigned roles with descriptions. `PUT /users/edit/:id` accepts an
optional `role_id` and returns the updated user's roles; changing a role also
requires the `roles.gestionar` permission.

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
**Middleware**: `IsAuthenticated` (all); creating and deleting users also require an administrator  
**Handler**: `internal/http/handlers/users/user.go`

| Method | Endpoint | Auth |
|--------|----------|------|
| GET | `/users` | Yes |
| GET | `/users/search` | Yes |
| GET | `/users/:id` | Yes |
| POST | `/users` | Admin |
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

### Role Routes
**Group Prefix**: `/roles` and `/permissions`  
**Middleware**: `IsAuthenticated` (all); every route also requires administrator or the `roles.gestionar` permission, except `/permissions/me`  
**Handler**: `internal/http/handlers/roles/role.go`

| Method | Endpoint | Auth |
|--------|----------|------|
| GET | `/permissions` | Yes |
| GET | `/permissions/me` | Yes |
| GET | `/roles` | Yes |
| POST | `/roles/new` | Yes |
| PUT | `/roles/edit/:id` | Yes |
| DELETE | `/roles/:id` | Yes |

> The `Administrador` role cannot be deleted.

---

### Group Routes
**Group Prefix**: `/groups`  
**Middleware**: `IsAuthenticated` (all); every route also requires administrator or the `roles.gestionar` permission  
**Handler**: `internal/http/handlers/groups/group.go`

| Method | Endpoint | Auth |
|--------|----------|------|
| GET | `/groups` | Yes |
| POST | `/groups/new` | Yes |
| PUT | `/groups/edit/:id` | Yes |
| DELETE | `/groups/:id` | Yes |
| GET | `/groups/:id/users` | Yes |
| POST | `/groups/:id/users` | Yes |
| DELETE | `/groups/:id/users/:userId` | Yes |
| GET | `/groups/:id/files` | Yes |
| DELETE | `/groups/:id/files/:fileId` | Yes |

---

## Total: 46 live endpoints

---

## Frontend API Routes (`apiUrls.ts`)

Requests use credentials for the refresh cookie and the auth interceptor adds the
access-token bearer header. Routes are defined in `src/app/core/routes/apiUrls.ts`:

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

### Roles and Permissions
| Key | Method | Path |
|-----|--------|------|
| `roles.getAll` | GET | `{BASE_URL}/roles` |
| `roles.create` | POST | `{BASE_URL}/roles/new` |
| `roles.update` | PUT | `{BASE_URL}/roles/edit/{id}` |
| `roles.delete` | DELETE | `{BASE_URL}/roles/{id}` |
| `permissions.getAll` | GET | `{BASE_URL}/permissions` |
| `permissions.me` | GET | `{BASE_URL}/permissions/me` |

### Groups
| Key | Method | Path |
|-----|--------|------|
| `groups.getAll` | GET | `{BASE_URL}/groups` |
| `groups.create` | POST | `{BASE_URL}/groups/new` |
| `groups.update` | PUT | `{BASE_URL}/groups/edit/{id}` |
| `groups.delete` | DELETE | `{BASE_URL}/groups/{id}` |
| `groups.getUsers` | GET | `{BASE_URL}/groups/{id}/users` |
| `groups.addUser` | POST | `{BASE_URL}/groups/{id}/users` |
| `groups.removeUser` | DELETE | `{BASE_URL}/groups/{id}/users/{userId}` |
| `groups.getFiles` | GET | `{BASE_URL}/groups/{id}/files` |
| `groups.revokeFile` | DELETE | `{BASE_URL}/groups/{id}/files/{fileId}` |

---

## Frontend ↔ Backend Alignment

All frontend routes now match the backend. No mismatches.
