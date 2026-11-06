// export interface Permission {
//   name?: string;
//   action?: string;
//   resource?: string;
// }

// export interface Role {
//   name: string;
//   permissions: Array<Permission | string>;
// }

// export interface User {
//   isAdmin?: boolean;
//   roles?: Role[];
//   permissions?: Array<Permission | string>;
// }

// function permissionName(permission: Permission | string): string {
//   if (typeof permission === 'string') {
//     return permission.trim().toLowerCase();
//   }

//   if (permission.name) {
//     return permission.name.trim().toLowerCase();
//   }

//   if (permission.resource && permission.action) {
//     return `${permission.resource}.${permission.action}`.trim().toLowerCase();
//   }

//   return '';
// }

// function matchesPermission(granted: string, requested: string): boolean {
//   if (granted === requested || granted === '*') {
//     return true;
//   }

//   const [resource, action] = requested.split('.');
//   return granted === `${resource}.*` || granted === `*.${action}`;
// }

// export function can(user: User | null | undefined, action: string, resource: string): boolean {
//   const normalizedAction = action.trim().toLowerCase();
//   const normalizedResource = resource.trim().toLowerCase();
//   if (!user || !normalizedAction || !normalizedResource) {
//     return false;
//   }

//   const isAdmin = user.isAdmin || user.roles?.some(role =>
//     ['admin', 'administrador'].includes(role.name.trim().toLowerCase())
//   );
//   if (isAdmin) {
//     return true;
//   }

//   const requested = `${normalizedResource}.${normalizedAction}`;
//   const grantedPermissions = [
//     ...(user.permissions ?? []),
//     ...(user.roles ?? []).flatMap(role => role.permissions),
//   ].map(permissionName);

//   return grantedPermissions.some(permission => matchesPermission(permission, requested));
// }
export interface Permission {
  name?: string;
  action?: string;
  resource?: string;
}

export interface Role {
  name: string;
  permissions: Array<Permission | string>;
}

export interface User {
  isAdmin?: boolean;
  roles?: Role[];
  permissions?: Array<Permission | string>;
  // Opcional: Caché de permisos ya procesados para no recalcularlos en cada render
  _flatPermissions?: string[];
}

function matchesPermission(granted: string, requested: string): boolean {
  if (granted === requested || granted === '*' || granted === '*.*') {
    return true;
  }

  // Soporta recursos con puntos, cortando por el último punto
  const lastDotIndex = requested.lastIndexOf('.');
  const resource = lastDotIndex !== -1 ? requested.substring(0, lastDotIndex) : requested;
  const action = lastDotIndex !== -1 ? requested.substring(lastDotIndex + 1) : '*';

  return granted === `${resource}.*` || granted === `*.${action}`;
}

// Agregamos un parámetro opcional 'data' para validar ACL (ej. IDs de archivos)
export function can(
  user: User | null | undefined,
  action: string,
  resource: string,
  data?: { ownerId?: number, resourceId?: number }
): boolean {
  const normalizedAction = action.trim().toLowerCase();
  const normalizedResource = resource.trim().toLowerCase();

  if (!user || !normalizedAction || !normalizedResource) return false;

  // 1. Regla de negocio: Si el usuario es el dueño, siempre puede (opcional, ajustado a tu tabla)
  if (data?.ownerId && data.ownerId === (user as any).id) return true;

  // 2. Super Admin Bypass
  const isAdmin = user.isAdmin || user.roles?.some(role =>
    ['admin', 'administrador'].includes(role.name.trim().toLowerCase())
  );
  if (isAdmin) return true;

  const requested = `${normalizedResource}.${normalizedAction}`;

  // 3. Evaluar permisos usando caché si existe, o calculando si no
  let grantedPermissions = user._flatPermissions;

  if (!grantedPermissions) {
    grantedPermissions = [
      ...(user.permissions ?? []),
      ...(user.roles ?? []).flatMap(role => role.permissions),
    ].map(p => {
      if (typeof p === 'string') return p.trim().toLowerCase();
      if (p.name) return p.name.trim().toLowerCase();
      if (p.resource && p.action) return `${p.resource}.${p.action}`.trim().toLowerCase();
      return '';
    }).filter(Boolean); // Elimina strings vacíos

    // Aquí idealmente guardarías grantedPermissions en user._flatPermissions para la próxima vez
  }

  return grantedPermissions.some(permission => matchesPermission(permission, requested));
}
