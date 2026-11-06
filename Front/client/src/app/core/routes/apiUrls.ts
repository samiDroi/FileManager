import { environment } from "../../../environments/environment.development";


// 1. URL Base general (ej. http://localhost:4100/api/v1)
const BASE_URL = environment.apiUrl;

// 2. Definimos los Prefijos de cada grupo (Tus "Route::prefix" de Laravel)
const PREFIX = {
    auth:  `${BASE_URL}`,
    users: `${BASE_URL}/users`,
    files: `${BASE_URL}/files`,
    typeFiles: `${BASE_URL}/type-files`,
    clients: `${BASE_URL}/client`,
    projects: `${BASE_URL}/project`,
    groups: `${BASE_URL}/groups`,
    roles: `${BASE_URL}/roles`,
    permissions: `${BASE_URL}/permissions`
};

// 3. Construimos el objeto exportable usando los prefijos
export const ApiRoutes = {
    auth: {
        login:        `${PREFIX.auth}/login`,
        register:     `${PREFIX.auth}/register`,
        logout:       `${PREFIX.auth}/logout`,
        refreshToken: `${PREFIX.auth}/refresh`,
        forgotPassword: `${PREFIX.auth}/forgot-password`,
        recoveryPassword:  `${PREFIX.auth}/recovery-password`
    },
    
    users: {
        getAll: `${PREFIX.users}`,
        assignable: `${PREFIX.users}/assignable`,
        create: `${PREFIX.users}`,
        getById: (id: string | number) => `${PREFIX.users}/${id}`,
        update:  (id: string | number) => `${PREFIX.users}/edit/${id}`,
        delete:  (id: string | number) => `${PREFIX.users}/${id}`,
        search:  (query: string) => `${PREFIX.users}/search?query=${encodeURIComponent(query)}`
    },

    files: {
        upload:   `${PREFIX.files}/upload`,
        preview:  (fileId: number) => `${PREFIX.files}/${fileId}/preview`,
        download: (fileId: number) => `${PREFIX.files}/${fileId}/download`,
        delete:   (fileId: number) => `${PREFIX.files}/${fileId}`,
        getAll:   `${PREFIX.files}`,
        getFileTypes: `${PREFIX.typeFiles}`,
        createFileTypes: `${PREFIX.typeFiles}/new`,
        getByProjectId: (projectId: number) => `${PREFIX.files}/project/${projectId}`,
        update:   (fileId: number) => `${PREFIX.files}/edit/${fileId}`,
        search:   (query: string) => `${PREFIX.files}/search?q=${encodeURIComponent(query)}`,
        access:   (fileId: number) => `${PREFIX.files}/${fileId}/access`,
        grantUserAccess: (fileId: number) => `${PREFIX.files}/${fileId}/access/users`,
        revokeUserAccess: (fileId: number, userId: number) => `${PREFIX.files}/${fileId}/access/users/${userId}`,
        grantGroupAccess: (fileId: number) => `${PREFIX.files}/${fileId}/access/groups`,
        revokeGroupAccess: (fileId: number, groupId: number) => `${PREFIX.files}/${fileId}/access/groups/${groupId}`
    },

    clients:{
        getAll: `${PREFIX.clients}`,
        search: `${PREFIX.clients}/search`,
        create: `${PREFIX.clients}/new`,
        withProjects: `${PREFIX.clients}/with-projects`,
        update: (id: string | number) => `${PREFIX.clients}/edit/${id}`,
        delete: (id: string | number) => `${PREFIX.clients}/${id}`
    },

    projects: {
        getAll: `${PREFIX.projects}`,
        search: `${PREFIX.projects}/search`,
        create: `${PREFIX.projects}/new`,
        update: (id: string | number) => `${PREFIX.projects}/edit/${id}`,
        delete: (id: string | number) => `${PREFIX.projects}/${id}`
    },

    groups: {
        getAll: `${PREFIX.groups}`,
        create: `${PREFIX.groups}/new`,
        update: (id: string | number) => `${PREFIX.groups}/edit/${id}`,
        delete: (id: string | number) => `${PREFIX.groups}/${id}`,
        getUsers: (id: string | number) => `${PREFIX.groups}/${id}/users`,
        getAvailableUsers: (id: string | number) => `${PREFIX.groups}/${id}/available-users`,
        addUser: (id: string | number) => `${PREFIX.groups}/${id}/users`,
        removeUser: (id: string | number, userId: string | number) => `${PREFIX.groups}/${id}/users/${userId}`,
        getFiles: (id: string | number) => `${PREFIX.groups}/${id}/files`,
        addFile: (id: string | number) => `${PREFIX.groups}/${id}/files`,
        revokeFile: (id: string | number, fileId: string | number) => `${PREFIX.groups}/${id}/files/${fileId}`
    },

    roles: {
        getAll: `${PREFIX.roles}`,
        create: `${PREFIX.roles}/new`,
        update: (id: string | number) => `${PREFIX.roles}/edit/${id}`,
        delete: (id: string | number) => `${PREFIX.roles}/${id}`
    },

    permissions: {
        getAll: `${PREFIX.permissions}`,
        me: `${PREFIX.permissions}/me`
    }

};