-- +goose Up

-- System-level rows (permissions and the default Administrator role) do not
-- belong to a concrete user, so their audit columns must accept NULL.
ALTER TABLE public."roles" ALTER COLUMN created_by DROP NOT NULL;
ALTER TABLE public."roles" ALTER COLUMN updated_by DROP NOT NULL;
ALTER TABLE public."permissions" ALTER COLUMN created_by DROP NOT NULL;
ALTER TABLE public."permissions" ALTER COLUMN updated_by DROP NOT NULL;
ALTER TABLE public."roles_permissions" ALTER COLUMN created_by DROP NOT NULL;
ALTER TABLE public."roles_permissions" ALTER COLUMN updated_by DROP NOT NULL;

INSERT INTO public."permissions" (name, description, created_by, updated_by) VALUES
('archivos.crear', 'Permite subir nuevos archivos', NULL, NULL),
('archivos.editar', 'Permite modificar datos de archivos', NULL, NULL),
('archivos.eliminar', 'Permite borrar archivos del sistema', NULL, NULL),
('usuarios.crear', 'Permite registrar nuevos usuarios', NULL, NULL),
('usuarios.editar', 'Permite modificar datos de otros usuarios', NULL, NULL),
('roles.gestionar', 'Permite crear y modificar roles de sistema', NULL, NULL),
('accesos.otorgar', 'Permite dar permisos sobre archivos a otros usuarios', NULL, NULL),
('clientes.crear', 'Permite registrar nuevos clientes', NULL, NULL),
('clientes.editar', 'Permite modificar datos de clientes', NULL, NULL),
('clientes.eliminar', 'Permite borrar clientes del sistema', NULL, NULL),
('proyectos.crear', 'Permite registrar nuevos proyectos', NULL, NULL),
('proyectos.editar', 'Permite modificar datos de proyectos', NULL, NULL),
('proyectos.eliminar', 'Permite borrar proyectos del sistema', NULL, NULL)
ON CONFLICT (name) DO UPDATE
SET description = EXCLUDED.description;

INSERT INTO public."roles" (name, description, created_by, updated_by)
VALUES ('Administrador', 'Control total sobre todos los módulos del sistema', NULL, NULL)
ON CONFLICT (name) DO NOTHING;

INSERT INTO public."roles_permissions" (role_id, permission_id, created_by, updated_by)
SELECT r.id, p.id, NULL, NULL
FROM public."roles" r
CROSS JOIN public."permissions" p
WHERE r.name = 'Administrador'
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- +goose Down

DELETE FROM public."roles_permissions"
WHERE role_id IN (SELECT id FROM public."roles" WHERE name = 'Administrador');

DELETE FROM public."roles" WHERE name = 'Administrador';

DELETE FROM public."permissions" WHERE name IN (
    'archivos.crear',
    'archivos.editar',
    'archivos.eliminar',
    'usuarios.crear',
    'usuarios.editar',
    'roles.gestionar',
    'accesos.otorgar',
    'clientes.crear',
    'clientes.editar',
    'clientes.eliminar',
    'proyectos.crear',
    'proyectos.editar',
    'proyectos.eliminar'
);

ALTER TABLE public."roles" ALTER COLUMN created_by SET NOT NULL;
ALTER TABLE public."roles" ALTER COLUMN updated_by SET NOT NULL;
ALTER TABLE public."permissions" ALTER COLUMN created_by SET NOT NULL;
ALTER TABLE public."permissions" ALTER COLUMN updated_by SET NOT NULL;
ALTER TABLE public."roles_permissions" ALTER COLUMN created_by SET NOT NULL;
ALTER TABLE public."roles_permissions" ALTER COLUMN updated_by SET NOT NULL;
