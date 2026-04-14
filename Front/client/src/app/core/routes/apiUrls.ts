import { environment } from "../../../environments/environment.development";


// 1. URL Base general (ej. http://localhost:4100/api/v1)
const BASE_URL = environment.apiUrl;

// 2. Definimos los Prefijos de cada grupo (Tus "Route::prefix" de Laravel)
const PREFIX = {
    auth:  `${BASE_URL}`,
    users: `${BASE_URL}/users`,
    files: `${BASE_URL}/files`,
    clients: `${BASE_URL}/client`
};

// 3. Construimos el objeto exportable usando los prefijos
export const ApiRoutes = {
    auth: {
        login:        `${PREFIX.auth}/login`,
        register:     `${PREFIX.auth}/register`,
        refreshToken: `${PREFIX.auth}/refresh` 
    },
    
    users: {
        getAll: `${PREFIX.users}`,
        getById: (id: string | number) => `${PREFIX.users}/${id}`,
        update:  (id: string | number) => `${PREFIX.users}/${id}`,
        delete:  (id: string | number) => `${PREFIX.users}/${id}`
    },

    files: {
        upload:   `${PREFIX.files}/upload`,
        download: (fileId: string) => `${PREFIX.files}/${fileId}/download`
    },

    clients:{
        getAll: `${PREFIX.clients}`,
        search: `${PREFIX.clients}/search`,
        create: `${PREFIX.clients}/new`,
        getById: (id: string | number) => `${PREFIX.clients}/${id}`,
        update: (id: string | number) => `${PREFIX.clients}/edit/${id}`,
        delete: (id: string | number) => `${PREFIX.clients}/delete/${id}`
    }

};