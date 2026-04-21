import { features } from "process";
import { rootCertificates } from "tls";

export const AppRoutes = {
    //module auth routes
    auth: {
        register: 'register',
        login: 'login',
        notFound: '404'
    },
    features:{
        dashboard: 'overview',
        users: 'users',
        files: 'files',
        clients: 'clients',
        projects: 'projects'
    }
}

export const GlobalPaths = {

}
// export const AppRoutes = {
//     auth: {
//         register: 'register',
//         login: 'login',
//         notFound: '404'
//     },
//     features: {
//         dashboard: 'overview',
//         users: 'users',
//         files: 'files',
//     },
//     //register globalpaths to components
//     globalPath:{
//         //auth
//         login: () => `/${AppRoutes.auth.login}`,
//         register: () => `/${AppRoutes.auth.register}`,
//         //features
//         dashboard: () => `/${AppRoutes.features.dashboard}`,
//         users: () => `/${AppRoutes.features.users}`,
//         files: () => `/${AppRoutes.features.files}`,
//     }
// }