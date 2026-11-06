-- Legacy commented schema (inactive)
-- CREATE TABLE IF NOT EXISTS public."clients"
-- (
--     id serial PRIMARY KEY,
--     name character varying(300) NOT NULL,
--     created_by integer NOT NULL, -- Cambiado de serial a integer
--     updated_by integer NOT NULL, -- Cambiado de serial a integer
--     created_at timestamp with time zone NOT NULL DEFAULT now(),
--     updated_at timestamp with time zone NOT NULL DEFAULT now()
-- );

-- COMMENT ON TABLE public."clients" IS 'Save the client''s projects registered';

-- CREATE TABLE IF NOT EXISTS public."file_types"
-- (
--     id serial PRIMARY KEY,
--     name character varying(60) NOT NULL,
--     created_at timestamp with time zone NOT NULL DEFAULT now(),
--     updated_at timestamp with time zone NOT NULL DEFAULT now(),
--     created_by integer NOT NULL,
--     updated_by integer NOT NULL
-- );

-- CREATE TABLE IF NOT EXISTS public."user"
-- (
--     id serial PRIMARY KEY,
--     username character varying(300) NOT NULL,
--     password text NOT NULL,
--     email character varying(160) NOT NULL,
--     phone character varying(120)[],
--     created_at timestamp with time zone NOT NULL DEFAULT now(),
--     updated_at timestamp with time zone NOT NULL DEFAULT now()
-- );

-- CREATE TABLE IF NOT EXISTS public."projects"
-- (
--     id bigserial PRIMARY KEY,
--     client_id integer NOT NULL, -- Cambiado a integer
--     name character varying(150) NOT NULL,
--     created_at timestamp with time zone NOT NULL DEFAULT now(),
--     updated_at timestamp with time zone NOT NULL DEFAULT now(),
--     created_by integer NOT NULL,
--     updated_by integer NOT NULL
-- );

-- CREATE TABLE IF NOT EXISTS public."contact_info"
-- (
--     id serial PRIMARY KEY,
--     client_id integer NOT NULL UNIQUE, -- Cambiado a integer
--     contact_phone character varying(60),
--     email character varying(150),
--     created_at timestamp with time zone NOT NULL DEFAULT now(),
--     updated_at timestamp with time zone NOT NULL DEFAULT now(),
--     created_by integer NOT NULL,
--     updated_by integer NOT NULL
-- );

-- CREATE TABLE IF NOT EXISTS public."files"
-- (
--     id bigserial PRIMARY KEY,
--     type_id integer NOT NULL, -- Cambiado a integer
--     tittle character varying(200) NOT NULL,
--     year integer NOT NULL,
--     path character varying(300) NOT NULL,
--     created_at timestamp with time zone NOT NULL DEFAULT now(),
--     updated_at timestamp with time zone NOT NULL DEFAULT now(),
--     created_by integer NOT NULL,
--     updated_by integer NOT NULL,
--     project_id bigint NOT NULL -- Cambiado de bigserial a bigint
-- );

-- -- LLAVES FORÁNEAS (ALTER TABLES)
-- ALTER TABLE public."clients" ADD FOREIGN KEY (created_by) REFERENCES public."user" (id);
-- ALTER TABLE public."clients" ADD FOREIGN KEY (updated_by) REFERENCES public."user" (id);

-- ALTER TABLE public."file_types" ADD FOREIGN KEY (created_by) REFERENCES public."user" (id);
-- ALTER TABLE public."file_types" ADD FOREIGN KEY (updated_by) REFERENCES public."user" (id);

-- ALTER TABLE public."projects" ADD FOREIGN KEY (client_id) REFERENCES public."clients" (id);
-- ALTER TABLE public."projects" ADD FOREIGN KEY (created_by) REFERENCES public."user" (id);
-- ALTER TABLE public."projects" ADD FOREIGN KEY (updated_by) REFERENCES public."user" (id);

-- ALTER TABLE public."contact_info" ADD FOREIGN KEY (client_id) REFERENCES public."clients" (id);
-- ALTER TABLE public."contact_info" ADD FOREIGN KEY (created_by) REFERENCES public."user" (id);
-- ALTER TABLE public."contact_info" ADD FOREIGN KEY (updated_by) REFERENCES public."user" (id);

-- ALTER TABLE public."files" ADD FOREIGN KEY (type_id) REFERENCES public."file_types" (id);
-- ALTER TABLE public."files" ADD FOREIGN KEY (created_by) REFERENCES public."user" (id);
-- ALTER TABLE public."files" ADD FOREIGN KEY (updated_by) REFERENCES public."user" (id);
-- ALTER TABLE public."files" ADD FOREIGN KEY (project_id) REFERENCES public."projects" (id);


-- Legacy down migration (inactive)
-- -- El orden inverso es vital por las dependencias de las llaves foráneas
-- DROP TABLE IF EXISTS public."files" CASCADE;
-- DROP TABLE IF EXISTS public."contact_info" CASCADE;
-- DROP TABLE IF EXISTS public."projects" CASCADE;
-- DROP TABLE IF EXISTS public."file_types" CASCADE;
-- DROP TABLE IF EXISTS public."clients" CASCADE;
-- DROP TABLE IF EXISTS public."user" CASCADE;
-- +goose Up

-- 1. CREACIÓN DE TIPOS ENUM PARA PERMISOS Y ROLES DE GRUPO
CREATE TYPE file_access_enum AS ENUM ('leer', 'editar', 'descargar');
CREATE TYPE group_role_enum AS ENUM ('admin', 'editor', 'lector');

-- 2. TABLAS BASE (Tu estructura original corregida)
CREATE TABLE IF NOT EXISTS public."user"
(
    id serial PRIMARY KEY,
    username character varying(300) NOT NULL,
    password text NOT NULL,
    email character varying(160) NOT NULL,
    phone character varying(120)[],
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    updated_at timestamp with time zone NOT NULL DEFAULT now(),
    created_by integer NOT NULL,
    updated_by integer NOT NULL,
    is_deleted integer NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS public."clients"
(
    id serial PRIMARY KEY,
    name character varying(300) NOT NULL,
    created_by integer NOT NULL,
    updated_by integer NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    updated_at timestamp with time zone NOT NULL DEFAULT now(),
    is_deleted integer NOT NULL DEFAULT 0
);
COMMENT ON TABLE public."clients" IS 'Save the client''s projects registered';

CREATE TABLE IF NOT EXISTS public."file_types"
(
    id serial PRIMARY KEY,
    name character varying(60) NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    updated_at timestamp with time zone NOT NULL DEFAULT now(),
    created_by integer NOT NULL,
    updated_by integer NOT NULL,
    is_deleted integer NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS public."projects"
(
    id bigserial PRIMARY KEY,
    client_id integer NOT NULL,
    name character varying(150) NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    updated_at timestamp with time zone NOT NULL DEFAULT now(),
    created_by integer NOT NULL,
    updated_by integer NOT NULL,
    is_deleted integer NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS public."contact_info"
(
    id serial PRIMARY KEY,
    client_id integer NOT NULL UNIQUE,
    contact_phone character varying(60),
    email character varying(150),
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    updated_at timestamp with time zone NOT NULL DEFAULT now(),
    created_by integer NOT NULL,
    updated_by integer NOT NULL
);

CREATE TABLE IF NOT EXISTS public."files"
(
    id bigserial PRIMARY KEY,
    type_id integer NOT NULL,
    title character varying(200) NOT NULL,
    year integer NOT NULL,
    path character varying(300) NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    updated_at timestamp with time zone NOT NULL DEFAULT now(),
    created_by integer NOT NULL, -- Actúa como el Propietario/Owner del archivo
    updated_by integer NOT NULL,
    project_id bigint NOT NULL,
    is_deleted integer NOT NULL DEFAULT 0
);

-- 3. NUEVAS TABLAS DE SISTEMA (RBAC GLOBAL)
CREATE TABLE IF NOT EXISTS public."roles"
(
    id serial PRIMARY KEY,
    name text UNIQUE NOT NULL,
    description text,
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    updated_at timestamp with time zone NOT NULL DEFAULT now(),
    created_by integer NOT NULL,
    updated_by integer NOT NULL
);

CREATE TABLE IF NOT EXISTS public."permissions"
(
    id serial PRIMARY KEY,
    name text UNIQUE NOT NULL,
    description text,
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    updated_at timestamp with time zone NOT NULL DEFAULT now(),
    created_by integer NOT NULL,
    updated_by integer NOT NULL
);

CREATE TABLE IF NOT EXISTS public."groups"
(
    id serial PRIMARY KEY,
    name text UNIQUE NOT NULL,
    description text,
    created_by integer NOT NULL,
    updated_by integer NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    updated_at timestamp with time zone NOT NULL DEFAULT now()
);

-- 4. TABLAS INTERMEDIAS DE RELACIÓN Y ACCESO (ACL)
CREATE TABLE IF NOT EXISTS public."roles_permissions"
(
    role_id integer NOT NULL,
    permission_id integer NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    updated_at timestamp with time zone NOT NULL DEFAULT now(),
    created_by integer NOT NULL,
    updated_by integer NOT NULL,
    PRIMARY KEY (role_id, permission_id)
);

CREATE TABLE IF NOT EXISTS public."users_roles"
(
    user_id integer NOT NULL,
    role_id integer NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    updated_at timestamp with time zone NOT NULL DEFAULT now(),
    created_by integer NOT NULL,
    updated_by integer NOT NULL,
    PRIMARY KEY (user_id, role_id)
);

CREATE TABLE IF NOT EXISTS public."users_groups"
(
    user_id integer NOT NULL,
    group_id integer NOT NULL,
    role_in_group group_role_enum NOT NULL DEFAULT 'lector', -- Rol interno del grupo
    added_at timestamp with time zone NOT NULL DEFAULT now(),
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    updated_at timestamp with time zone NOT NULL DEFAULT now(),
    created_by integer NOT NULL,
    updated_by integer NOT NULL,
    PRIMARY KEY (user_id, group_id)
);

CREATE TABLE IF NOT EXISTS public."files_users"
(
    file_id bigint NOT NULL,
    user_id integer NOT NULL,
    access_type file_access_enum NOT NULL,
    granted_by integer NOT NULL,
    granted_at timestamp with time zone NOT NULL DEFAULT now(),
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    updated_at timestamp with time zone NOT NULL DEFAULT now(),
    created_by integer NOT NULL,
    updated_by integer NOT NULL,
    PRIMARY KEY (file_id, user_id)
);

CREATE TABLE IF NOT EXISTS public."files_groups"
(
    file_id bigint NOT NULL,
    group_id integer NOT NULL,
    access_type file_access_enum NOT NULL,
    granted_by integer NOT NULL,
    granted_at timestamp with time zone NOT NULL DEFAULT now(),
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    updated_at timestamp with time zone NOT NULL DEFAULT now(),
    created_by integer NOT NULL,
    updated_by integer NOT NULL,
    PRIMARY KEY (file_id, group_id)
);
CREATE TABLE IF NOT EXISTS public."sessions"
(
    id BIGSERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    token_hash VARCHAR(64) UNIQUE NOT NULL,
    is_revoked BOOLEAN NOT NULL DEFAULT FALSE,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    revoked_at TIMESTAMP WITH TIME ZONE,
    ip_address VARCHAR(45),
    user_agent TEXT,

    CONSTRAINT sessions_user_id_fkey
        FOREIGN KEY(user_id)
        REFERENCES public."user" (id) ON DELETE CASCADE
);

-- Índice para buscar todas las sesiones activas de un usuario (para la vista de "Dispositivos")
CREATE INDEX idx_sessions_user_id ON sessions(user_id);
-- 5. LLAVES FORÁNEAS DE TABLAS BASE ORIGINALES
ALTER TABLE public."clients" ADD FOREIGN KEY (created_by) REFERENCES public."user" (id);
ALTER TABLE public."clients" ADD FOREIGN KEY (updated_by) REFERENCES public."user" (id);

ALTER TABLE public."file_types" ADD FOREIGN KEY (created_by) REFERENCES public."user" (id);
ALTER TABLE public."file_types" ADD FOREIGN KEY (updated_by) REFERENCES public."user" (id);

ALTER TABLE public."projects" ADD FOREIGN KEY (client_id) REFERENCES public."clients" (id) ON DELETE CASCADE;
ALTER TABLE public."projects" ADD FOREIGN KEY (created_by) REFERENCES public."user" (id);
ALTER TABLE public."projects" ADD FOREIGN KEY (updated_by) REFERENCES public."user" (id);

ALTER TABLE public."contact_info" ADD FOREIGN KEY (client_id) REFERENCES public."clients" (id) ON DELETE CASCADE;
ALTER TABLE public."contact_info" ADD FOREIGN KEY (created_by) REFERENCES public."user" (id);
ALTER TABLE public."contact_info" ADD FOREIGN KEY (updated_by) REFERENCES public."user" (id);

ALTER TABLE public."files" ADD FOREIGN KEY (type_id) REFERENCES public."file_types" (id);
ALTER TABLE public."files" ADD FOREIGN KEY (created_by) REFERENCES public."user" (id);
ALTER TABLE public."files" ADD FOREIGN KEY (updated_by) REFERENCES public."user" (id);
ALTER TABLE public."files" ADD FOREIGN KEY (project_id) REFERENCES public."projects" (id) ON DELETE CASCADE;

-- 6. LLAVES FORÁNEAS DE NUEVAS TABLAS DE PERMISOS Y GRUPOS
ALTER TABLE public."roles" ADD FOREIGN KEY (created_by) REFERENCES public."user" (id);
ALTER TABLE public."roles" ADD FOREIGN KEY (updated_by) REFERENCES public."user" (id);

ALTER TABLE public."permissions" ADD FOREIGN KEY (created_by) REFERENCES public."user" (id);
ALTER TABLE public."permissions" ADD FOREIGN KEY (updated_by) REFERENCES public."user" (id);

ALTER TABLE public."groups" ADD FOREIGN KEY (created_by) REFERENCES public."user" (id);
ALTER TABLE public."groups" ADD FOREIGN KEY (updated_by) REFERENCES public."user" (id);

ALTER TABLE public."roles_permissions" ADD FOREIGN KEY (role_id) REFERENCES public."roles" (id) ON DELETE CASCADE;
ALTER TABLE public."roles_permissions" ADD FOREIGN KEY (permission_id) REFERENCES public."permissions" (id) ON DELETE CASCADE;
ALTER TABLE public."roles_permissions" ADD FOREIGN KEY (created_by) REFERENCES public."user" (id);
ALTER TABLE public."roles_permissions" ADD FOREIGN KEY (updated_by) REFERENCES public."user" (id);

ALTER TABLE public."users_roles" ADD FOREIGN KEY (user_id) REFERENCES public."user" (id) ON DELETE CASCADE;
ALTER TABLE public."users_roles" ADD FOREIGN KEY (role_id) REFERENCES public."roles" (id) ON DELETE CASCADE;
ALTER TABLE public."users_roles" ADD FOREIGN KEY (created_by) REFERENCES public."user" (id);
ALTER TABLE public."users_roles" ADD FOREIGN KEY (updated_by) REFERENCES public."user" (id);

ALTER TABLE public."users_groups" ADD FOREIGN KEY (user_id) REFERENCES public."user" (id) ON DELETE CASCADE;
ALTER TABLE public."users_groups" ADD FOREIGN KEY (group_id) REFERENCES public."groups" (id) ON DELETE CASCADE;
ALTER TABLE public."users_groups" ADD FOREIGN KEY (created_by) REFERENCES public."user" (id);
ALTER TABLE public."users_groups" ADD FOREIGN KEY (updated_by) REFERENCES public."user" (id);

ALTER TABLE public."files_users" ADD FOREIGN KEY (file_id) REFERENCES public."files" (id) ON DELETE CASCADE;
ALTER TABLE public."files_users" ADD FOREIGN KEY (user_id) REFERENCES public."user" (id) ON DELETE CASCADE;
ALTER TABLE public."files_users" ADD FOREIGN KEY (granted_by) REFERENCES public."user" (id);
ALTER TABLE public."files_users" ADD FOREIGN KEY (created_by) REFERENCES public."user" (id);
ALTER TABLE public."files_users" ADD FOREIGN KEY (updated_by) REFERENCES public."user" (id);

ALTER TABLE public."files_groups" ADD FOREIGN KEY (file_id) REFERENCES public."files" (id) ON DELETE CASCADE;
ALTER TABLE public."files_groups" ADD FOREIGN KEY (group_id) REFERENCES public."groups" (id) ON DELETE CASCADE;
ALTER TABLE public."files_groups" ADD FOREIGN KEY (granted_by) REFERENCES public."user" (id);
ALTER TABLE public."files_groups" ADD FOREIGN KEY (created_by) REFERENCES public."user" (id);
ALTER TABLE public."files_groups" ADD FOREIGN KEY (updated_by) REFERENCES public."user" (id);


-- +goose Down
-- El orden inverso es vital por las dependencias de las llaves foráneas y tipos

-- 1. Eliminar tablas intermedias y dependientes
DROP TABLE IF EXISTS public."files_groups" CASCADE;
DROP TABLE IF EXISTS public."files_users" CASCADE;
DROP TABLE IF EXISTS public."users_groups" CASCADE;
DROP TABLE IF EXISTS public."users_roles" CASCADE;
DROP TABLE IF EXISTS public."roles_permissions" CASCADE;

-- 2. Eliminar entidades de permisos y grupos
DROP TABLE IF EXISTS public."groups" CASCADE;
DROP TABLE IF EXISTS public."permissions" CASCADE;
DROP TABLE IF EXISTS public."roles" CASCADE;

-- 3. Eliminar tablas base
DROP TABLE IF EXISTS public."files" CASCADE;
DROP TABLE IF EXISTS public."contact_info" CASCADE;
DROP TABLE IF EXISTS public."projects" CASCADE;
DROP TABLE IF EXISTS public."file_types" CASCADE;
DROP TABLE IF EXISTS public."clients" CASCADE;
DROP TABLE IF EXISTS public."sessions" CASCADE;
DROP TABLE IF EXISTS public."user" CASCADE;

-- 4. Eliminar tipos ENUM
DROP TYPE IF EXISTS file_access_enum;
DROP TYPE IF EXISTS group_role_enum;