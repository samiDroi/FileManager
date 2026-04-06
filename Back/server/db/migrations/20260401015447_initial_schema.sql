-- +goose Up
CREATE TABLE IF NOT EXISTS public."Client"
(
    id serial PRIMARY KEY,
    name character varying(300) NOT NULL,
    created_by integer NOT NULL, -- Cambiado de serial a integer
    updated_by integer NOT NULL, -- Cambiado de serial a integer
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    updated_at timestamp with time zone NOT NULL DEFAULT now()
);

COMMENT ON TABLE public."Client" IS 'Save the client''s proyects registered';

CREATE TABLE IF NOT EXISTS public."File_Tyipe"
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

CREATE TABLE IF NOT EXISTS public."Proyect"
(
    id bigserial PRIMARY KEY,
    client_id integer NOT NULL, -- Cambiado a integer
    name character varying(150) NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    updated_at timestamp with time zone NOT NULL DEFAULT now(),
    created_by integer NOT NULL,
    updated_by integer NOT NULL
);

CREATE TABLE IF NOT EXISTS public."Contact_Info"
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

CREATE TABLE IF NOT EXISTS public."File"
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
    proyect_id bigint NOT NULL -- Cambiado de bigserial a bigint
);

-- LLAVES FORÁNEAS (ALTER TABLES)
ALTER TABLE public."Client" ADD FOREIGN KEY (created_by) REFERENCES public."user" (id);
ALTER TABLE public."Client" ADD FOREIGN KEY (updated_by) REFERENCES public."user" (id);

ALTER TABLE public."File_Tyipe" ADD FOREIGN KEY (created_by) REFERENCES public."user" (id);
ALTER TABLE public."File_Tyipe" ADD FOREIGN KEY (updated_by) REFERENCES public."user" (id);

ALTER TABLE public."Proyect" ADD FOREIGN KEY (client_id) REFERENCES public."Client" (id);
ALTER TABLE public."Proyect" ADD FOREIGN KEY (created_by) REFERENCES public."user" (id);
ALTER TABLE public."Proyect" ADD FOREIGN KEY (updated_by) REFERENCES public."user" (id);

ALTER TABLE public."Contact_Info" ADD FOREIGN KEY (client_id) REFERENCES public."Client" (id);
ALTER TABLE public."Contact_Info" ADD FOREIGN KEY (created_by) REFERENCES public."user" (id);
ALTER TABLE public."Contact_Info" ADD FOREIGN KEY (updated_by) REFERENCES public."user" (id);

ALTER TABLE public."File" ADD FOREIGN KEY (type_id) REFERENCES public."File_Tyipe" (id);
ALTER TABLE public."File" ADD FOREIGN KEY (created_by) REFERENCES public."user" (id);
ALTER TABLE public."File" ADD FOREIGN KEY (updated_by) REFERENCES public."user" (id);
ALTER TABLE public."File" ADD FOREIGN KEY (proyect_id) REFERENCES public."Proyect" (id);


-- +goose Down
-- El orden inverso es vital por las dependencias de las llaves foráneas
DROP TABLE IF EXISTS public."File";
DROP TABLE IF EXISTS public."Contact_Info";
DROP TABLE IF EXISTS public."Proyect";
DROP TABLE IF EXISTS public."File_Tyipe";
DROP TABLE IF EXISTS public."Client";
DROP TABLE IF EXISTS public."user";
