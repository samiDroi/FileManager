# CattFileManager

Aplicación web para la gestión documental de clientes y proyectos. Centraliza el registro de clientes, sus proyectos y los archivos asociados, con autenticación de usuarios, control de acceso basado en roles (RBAC), grupos de trabajo y permisos por archivo.

El sistema está dividido en dos aplicaciones:

- **Backend**: API REST en Go con Echo, PostgreSQL y migraciones gestionadas con Goose.
- **Frontend**: SPA desarrollada con Angular 21, Angular Material y soporte de SSR.

## Características

### Autenticación y sesiones

- Registro e inicio de sesión con contraseñas hasheadas mediante bcrypt.
- Access token JWT con vigencia de 15 minutos.
- Refresh token opaco de 7 días almacenado en cookie `HttpOnly`, `Secure` y `SameSite=Strict`.
- Rotación de refresh tokens con detección de reutilización y revocación por familia de sesión.
- Cierre de sesión con revocación de la sesión asociada.
- Recuperación de contraseña por correo electrónico con enlace firmado de 15 minutos.
- Limitación de peticiones en los endpoints de autenticación.
- Limpieza periódica de sesiones expiradas.
- El primer usuario registrado se convierte automáticamente en administrador cuando aún no existe ningún administrador asignado.

### Roles y permisos (RBAC)

- Roles globales con permisos asociados.
- Rol `Administrador` creado por migración con todos los permisos del sistema.
- Verificación de permisos en el backend mediante middleware (`RequireAdmin`, `RequirePermission`).
- Guards y servicio de permisos en el frontend para ocultar o bloquear acciones según el usuario.
- Permisos disponibles para archivos, usuarios, clientes, proyectos, tipos de archivo, accesos y gestión de roles.

### Grupos de trabajo

- Creación, edición y eliminación de grupos.
- Asignación de usuarios a grupos con rol interno: `admin`, `editor` o `lector`.
- Gestión de archivos compartidos con un grupo.
- Acceso a grupos y archivos según el rol del usuario dentro del grupo.

### Clientes y proyectos

- CRUD de clientes con búsqueda.
- CRUD de proyectos con búsqueda.
- Consulta de clientes junto con sus proyectos.
- Navegación desde el cliente hacia sus proyectos y archivos.

### Gestión de archivos

- Carga de archivos con metadatos: título, año, proyecto y tipo de archivo.
- Listado general, listado por proyecto y búsqueda.
- Vista previa y descarga de archivos.
- Edición y eliminación de archivos.
- Administración de tipos de archivo.
- Permisos por archivo para usuarios y grupos con tipo de acceso: `leer`, `editar` o `descargar`.
- Dialogo de administración de accesos por archivo.

### Usuarios y perfil

- Administración de usuarios por parte del administrador.
- Creación y edición de usuarios con asignación de rol.
- Consulta de usuarios asignables para compartir archivos.
- Perfil de usuario y actualización de datos.

### Panel de administración

Rutas protegidas para administradores:

- `admin-dashboard`
- `admin-users`
- `admin-files`
- `admin-clients`
- `admin-projects`
- `admin-groups`
- `admin-roles`

El backend también expone estadísticas por tipo de archivo en `GET /stats/file-types` para usuarios administradores.

## Tecnologías

| Capa | Tecnologías |
|------|-------------|
| Backend | Go 1.25+, Echo v4, pgx v5, PostgreSQL, Goose, JWT, bcrypt |
| Frontend | Angular 21, Angular Material, Angular CDK, Angular SSR, TypeScript, Vitest, SweetAlert2, Font Awesome |
| Herramientas | Go Modules, pnpm, Prettier |

## Estructura del repositorio

```text
FileManager/
├── Back/
│   ├── migration.readme.txt          # Guía detallada de migraciones
│   └── server/
│       ├── cmd/
│       │   ├── api/                  # Punto de entrada de la API
│       │   └── flags/api-dto/        # Inspector de DTOs y rutas
│       ├── db/migrations/            # Migraciones SQL de Goose
│       ├── internal/
│       │   ├── config/               # Carga de variables de entorno
│       │   ├── domain/               # Modelos y DTOs
│       │   ├── http/                 # Handlers, middlewares y tokens
│       │   ├── routes/               # Registro central de rutas
│       │   ├── services/             # Lógica de negocio
│       │   └── store/                # Pool de PostgreSQL y repositorios
│       ├── pkg/
│       ├── uploads/                  # Archivos subidos (ignorado por Git)
│       ├── .env.example
│       └── go.mod
├── Front/
│   └── client/
│       ├── src/app/
│       │   ├── auth/                 # Login, registro y recuperación
│       │   ├── core/                 # Guards, interceptores, DTOs y servicios
│       │   ├── features/             # Clientes, proyectos, archivos, usuarios, grupos, roles
│       │   ├── layouts/              # Layouts público, principal y admin
│       │   └── shared/               # Componentes reutilizables
│       ├── src/environments/         # URL de la API por entorno
│       ├── angular.json
│       └── package.json
├── API_ENDPOINTS.md                  # Referencia de endpoints
└── README.md
```

## Requisitos

- Go 1.25 o superior.
- PostgreSQL en ejecución.
- Node.js `^20.19.0`, `^22.12.0` o `>=24.0.0`.
- pnpm.
- Una cuenta SMTP para el envío de correos de recuperación (opcional para desarrollo).
- Goose CLI solo si se desea ejecutar comandos manuales de migración.

## Puesta en marcha

### 1. Crear la base de datos

```sql
CREATE DATABASE "FileManagerDB";
```

### 2. Configurar el backend

```powershell
cd Back/server
Copy-Item .env.example .env
```

Edita `Back/server/.env` con tus valores:

| Variable | Descripción | Ejemplo |
|----------|-------------|---------|
| `PORT` | Puerto de la API | `4100` |
| `DB_CONN` | Cadena de conexión a PostgreSQL | `postgres://usuario:password@localhost:5432/FileManagerDB?sslmode=disable` |
| `SECRET_KEY` | Clave para firmar JWT y tokens de recuperación | Cadena aleatoria larga |
| `AUTH_KEY` | Clave heredada de sesiones de Echo; no se usa en la autenticación actual | Opcional |
| `ALLOW_ORIGINS` | Orígenes permitidos por CORS, separados por coma | `http://localhost:4200` |
| `FRONTEND_URL` | URL usada en el enlace de recuperación de contraseña | `http://localhost:4200` |
| `COOKIE_SECURE` | Marca las cookies como `Secure`; usar `true` con HTTPS | `false` |
| `UPLOAD_PATH` | Directorio donde se guardan los archivos subidos | `./uploads` |
| `MAIL_MAILER` | Transporte de correo | `smtp` |
| `MAIL_HOST` | Host SMTP | `smtp.ejemplo.com` |
| `MAIL_PORT` | Puerto SMTP | `587` |
| `MAIL_USERNAME` | Usuario SMTP | `correo@ejemplo.com` |
| `MAIL_PASSWORD` | Contraseña o clave de aplicación SMTP | Clave secreta |
| `MAIL_ENCRYPTION` | Tipo de cifrado | `tls` |
| `MAIL_FROM_ADDRESS` | Remitente de los correos | `correo@ejemplo.com` |
| `MAIL_FROM_NAME` | Nombre del remitente | `CattFileManager` |

Genera `SECRET_KEY` con un valor aleatorio, por ejemplo:

```bash
openssl rand -base64 48
```

Nunca subas el archivo `.env` al repositorio. Solo `.env.example` debe versionarse.

### 3. Aplicar migraciones e iniciar la API

```powershell
cd Back/server
go mod download
go run ./cmd/api -migrate
```

La API quedará disponible en `http://localhost:4100`. En arranques posteriores, si no necesitas aplicar migraciones:

```powershell
go run ./cmd/api
```

Para compilar:

```powershell
go build -o api.exe ./cmd/api
.\api.exe -migrate
```

En Linux/macOS:

```bash
go build -o api ./cmd/api
./api -migrate
```

### 4. Configurar e iniciar el frontend

```powershell
cd Front/client
pnpm install
pnpm start
```

La aplicación estará disponible en `http://localhost:4200`.

La URL de la API se define en:

- `Front/client/src/environments/environment.development.ts` para desarrollo.
- `Front/client/src/environments/environment.ts` para producción.

### 5. Primer acceso

1. Abre `http://localhost:4200/register`.
2. Registra el primer usuario.
3. Ese usuario recibirá el rol `Administrador` si todavía no existe un administrador en el sistema.
4. Inicia sesión y accede al panel de administración.

## Autenticación en la API

`POST /login` y `POST /register` devuelven un `access_token` en el cuerpo de la respuesta. Las rutas protegidas requieren:

```http
Authorization: Bearer <access_token>
```

El refresh token se envía únicamente en la cookie `refresh_token`. `POST /refresh` renueva el access token y rota el refresh token; `POST /logout` revoca la sesión y elimina la cookie.

Flujo básico:

```bash
curl -X POST http://localhost:4100/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"tu-password"}'
```

La referencia completa de endpoints está en [API_ENDPOINTS.md](API_ENDPOINTS.md).

## Migraciones

Las migraciones se encuentran en `Back/server/db/migrations` y usan el formato de Goose (`-- +goose Up` / `-- +goose Down`).

Aplicar migraciones y arrancar la API:

```powershell
cd Back/server
go run ./cmd/api -migrate
```

Comandos manuales con Goose:

```bash
goose postgres "$DB_CONN" -dir=db/migrations status
goose postgres "$DB_CONN" -dir=db/migrations up
goose postgres "$DB_CONN" -dir=db/migrations down
goose postgres "$DB_CONN" -dir=db/migrations version
```

Crear una migración nueva:

```bash
cd Back/server
goose -dir=db/migrations create nombre_de_la_migracion sql
```

Para reiniciar todas las migraciones existen los scripts `Back/server/cmd/migrate-fresh.bat` (Windows) y `Back/server/cmd/migrate-fresh.sh` (Linux/macOS). Por defecto asumen `localhost`, puerto `5432`, usuario `postgres` y base `FileManagerDB`; ajústalos si tu configuración es diferente. Consulta [Back/migration.readme.txt](Back/migration.readme.txt) para más detalle.

## Pruebas

Backend:

```powershell
cd Back/server
go test ./...
```

Frontend:

```powershell
cd Front/client
pnpm test
```

## Utilidades

El binario del backend incluye un inspector de DTOs y rutas por feature:

```powershell
cd Back/server
go run ./cmd/api -all
go run ./cmd/api -feature files
go run ./cmd/api -dto -feature users
```

## Seguridad

- `.env` está excluido por `.gitignore`; versiona únicamente `.env.example`.
- Usa una clave `SECRET_KEY` aleatoria y específica por entorno.
- En producción, activa `COOKIE_SECURE=true` y sirve la aplicación mediante HTTPS.
- Restringe `ALLOW_ORIGINS` a los dominios reales del frontend.
- Utiliza un usuario de PostgreSQL dedicado con permisos mínimos en lugar del superusuario `postgres`.
- Los archivos subidos se guardan en `UPLOAD_PATH`; protege ese directorio y respáldalo.
- Si una credencial se expone, rótala de inmediato: cambiar `SECRET_KEY` invalida los tokens emitidos y obliga a los usuarios a iniciar sesión de nuevo.

## Documentación adicional

- [API_ENDPOINTS.md](API_ENDPOINTS.md): referencia de endpoints del backend y rutas del frontend.
- [Back/migration.readme.txt](Back/migration.readme.txt): guía de migraciones con Goose.
- [Front/client/README.md](Front/client/README.md): documentación base de Angular CLI.
