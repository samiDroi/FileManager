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
    projects: `${BASE_URL}/project`
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
        search:   (query: string) => `${PREFIX.files}/search?q=${encodeURIComponent(query)}`
    },

    clients:{
        getAll: `${PREFIX.clients}`,
        search: `${PREFIX.clients}/search`,
        create: `${PREFIX.clients}/new`,
        withProjects: `${PREFIX.clients}/with-projects`,
        update: (id: string | number) => `${PREFIX.clients}/edit/${id}`
    },

    projects: {
        getAll: `${PREFIX.projects}`,
        search: `${PREFIX.projects}/search`,
        create: `${PREFIX.projects}/new`,
        update: (id: string | number) => `${PREFIX.projects}/edit/${id}`
    }

};