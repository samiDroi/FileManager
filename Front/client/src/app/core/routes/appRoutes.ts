export const AppRoutes = {
    //module auth routes
    auth: {
        register: 'register',
        login: 'login',
        forgotPassword: 'forgot-password',
        resetPassword: 'reset-password',
        recoveryPassword: 'recovery-password',
        notFound: '404'
    },
    features:{
        dashboard: 'overview',
        users: 'users',
        files: 'files',
        clients: 'clients',
        projects: 'projects',
        groups: 'groups'
    },
    admin:{
      dashboard: 'admin-dashboard',
      users: 'admin-users',
      files: 'admin-files',
      clients: 'admin-clients',
      projects: 'admin-projects',
      groups: 'admin-groups',
      roles: 'admin-roles'
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
