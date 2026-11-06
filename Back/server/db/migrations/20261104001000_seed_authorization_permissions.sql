-- +goose Up

INSERT INTO public."permissions" (name, description, created_by, updated_by) VALUES
('usuarios.leer', 'Permite consultar el listado y detalle de usuarios', NULL, NULL),
('usuarios.eliminar', 'Permite eliminar usuarios del sistema', NULL, NULL),
('tipos_archivo.crear', 'Permite registrar nuevos tipos de archivo', NULL, NULL)
ON CONFLICT (name) DO UPDATE
SET description = EXCLUDED.description;

INSERT INTO public."roles_permissions" (role_id, permission_id, created_by, updated_by)
SELECT r.id, p.id, NULL, NULL
FROM public."roles" r
CROSS JOIN public."permissions" p
WHERE r.name = 'Administrador'
  AND p.name IN ('usuarios.leer', 'usuarios.eliminar', 'tipos_archivo.crear')
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- +goose Down

DELETE FROM public."roles_permissions"
WHERE permission_id IN (
    SELECT id FROM public."permissions"
    WHERE name IN ('usuarios.leer', 'usuarios.eliminar', 'tipos_archivo.crear')
);

DELETE FROM public."permissions"
WHERE name IN ('usuarios.leer', 'usuarios.eliminar', 'tipos_archivo.crear');
