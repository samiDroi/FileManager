-- +goose Up
CREATE TABLE IF NOT EXISTS public."clients"
(
    id serial PRIMARY KEY,
    name character varying(300) NOT NULL,
    created_by integer NOT NULL, -- Cambiado de serial a integer
    updated_by integer NOT NULL, -- Cambiado de serial a integer
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    updated_at timestamp with time zone NOT NULL DEFAULT now()
);

COMMENT ON TABLE public."clients" IS 'Save the client''s projects registered';

CREATE TABLE IF NOT EXISTS public."file_types"
(
    id serial PRIMARY KEY,
    name character varying(60) NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    updated_at timestamp with time zone NOT NULL DEFAULT now(),
    created_by integer NOT NULL,
    updated_by integer NOT NULL
);

CREATE TABLE IF NOT EXISTS public."user"
(
    id serial PRIMARY KEY,
    username character varying(300) NOT NULL,
    password text NOT NULL,
    email character varying(160) NOT NULL,
    phone character varying(120)[],
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    updated_at timestamp with time zone NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS public."projects"
(
    id bigserial PRIMARY KEY,
    client_id integer NOT NULL, -- Cambiado a integer
    name character varying(150) NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    updated_at timestamp with time zone NOT NULL DEFAULT now(),
    created_by integer NOT NULL,
    updated_by integer NOT NULL
);

CREATE TABLE IF NOT EXISTS public."contact_info"
(
    id serial PRIMARY KEY,
    client_id integer NOT NULL UNIQUE, -- Cambiado a integer
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
    type_id integer NOT NULL UNIQUE, -- Cambiado a integer
    tittle character varying(200) NOT NULL,
    year integer NOT NULL,
    path character varying(300) NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    updated_at timestamp with time zone NOT NULL DEFAULT now(),
    created_by integer NOT NULL,
    updated_by integer NOT NULL,
    project_id bigint NOT NULL -- Cambiado de bigserial a bigint
);

-- LLAVES FORÁNEAS (ALTER TABLES)
ALTER TABLE public."clients" ADD FOREIGN KEY (created_by) REFERENCES public."user" (id);
ALTER TABLE public."clients" ADD FOREIGN KEY (updated_by) REFERENCES public."user" (id);

ALTER TABLE public."file_types" ADD FOREIGN KEY (created_by) REFERENCES public."user" (id);
ALTER TABLE public."file_types" ADD FOREIGN KEY (updated_by) REFERENCES public."user" (id);

ALTER TABLE public."projects" ADD FOREIGN KEY (client_id) REFERENCES public."clients" (id);
ALTER TABLE public."projects" ADD FOREIGN KEY (created_by) REFERENCES public."user" (id);
ALTER TABLE public."projects" ADD FOREIGN KEY (updated_by) REFERENCES public."user" (id);

ALTER TABLE public."contact_info" ADD FOREIGN KEY (client_id) REFERENCES public."clients" (id);
ALTER TABLE public."contact_info" ADD FOREIGN KEY (created_by) REFERENCES public."user" (id);
ALTER TABLE public."contact_info" ADD FOREIGN KEY (updated_by) REFERENCES public."user" (id);

ALTER TABLE public."files" ADD FOREIGN KEY (type_id) REFERENCES public."file_types" (id);
ALTER TABLE public."files" ADD FOREIGN KEY (created_by) REFERENCES public."user" (id);
ALTER TABLE public."files" ADD FOREIGN KEY (updated_by) REFERENCES public."user" (id);
ALTER TABLE public."files" ADD FOREIGN KEY (project_id) REFERENCES public."projects" (id);


-- +goose Down
-- El orden inverso es vital por las dependencias de las llaves foráneas
DROP TABLE IF EXISTS public."files" CASCADE;
DROP TABLE IF EXISTS public."contact_Info" CASCADE;
DROP TABLE IF EXISTS public."projects" CASCADE;
DROP TABLE IF EXISTS public."file_Types" CASCADE;
DROP TABLE IF EXISTS public."clients" CASCADE;
DROP TABLE IF EXISTS public."user" CASCADE;
