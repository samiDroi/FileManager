import { Router, Routes } from '@angular/router';
import { AppRoutes } from './core/routes/appRoutes';
import { AdminLayout, AuthLayout, MainLayout } from './layouts';
import { authGuard, guestGuard } from './core/guards';
import { AuthStateService } from './auth/auth-state.service';
import { inject } from '@angular/core/primitives/di';


export const routes: Routes = [
   {
        path: '',
        redirectTo: AppRoutes.auth.login, 
        pathMatch: 'full',
    },
    {
        path: AppRoutes.auth.recoveryPassword,
        loadComponent: () => import('./auth/components/recovery-password-component/recovery-password-component').then(m => m.RecoveryPasswordComponent),
        canActivate: [guestGuard]
    },
    //lazy loading auth module
       {
           //path is '', because these routes will be loaded under the main path defined in app.routes.ts
           //so the routes are like /login, /register, etc, but they are grouped under the same layout
           path: '',
           component: AuthLayout, 
           //grouping all auth related routes under the same layout
            children:[
                {
                    path: AppRoutes.auth.login,
                    loadComponent: () => import('./auth/components/login-component/login-component').then(m => m.LoginComponent),
                    canActivate: [guestGuard]
                },
                {
                    path: AppRoutes.auth.register,
                    loadComponent: () => import('./auth/components/register-component/register-component').then(m => m.RegisterComponent),
                    canActivate: [guestGuard]
                },
                {
                    path: AppRoutes.auth.forgotPassword,
                    loadComponent: () => import('./auth/components/forgot-password-component/forgot-password-component').then(m => m.ForgotPasswordComponent),
                    canActivate: [guestGuard]
                },
                {
                    path: AppRoutes.auth.resetPassword,
                    loadComponent: () => import('./auth/components/reset-password-component/reset-password-component').then(m => m.ResetPasswordComponent),
                    canActivate: [guestGuard]
                },
               {
                    //default route for auth, redirecting to login
                    path: '',
                    redirectTo: AppRoutes.auth.login,
                    pathMatch: 'full'
                }
           ]
       },
       {
        path: '',
        component: AdminLayout,
        canActivate: [authGuard],
        children: [
            {
                path: AppRoutes.admin.dashboard,
                loadComponent: () => import('./features/dashboard/dashboard-component/dashboard-component').then(m => m.DashboardComponent)
            },
            {
                path: AppRoutes.admin.users,
                loadComponent: () => import('./features/users/components/user-component/user-component').then(m => m.UserComponent)
            },
            {
                path: AppRoutes.admin.files,
                loadComponent: () => import('./features/files/components/files-component/files-component').then(m => m.FilesComponent)
            },
            {
                path: AppRoutes.admin.clients,
                loadComponent: () => import('./features/clients/components/clients-component/clients-component').then(m => m.ClientsComponent)
            },
            {
                path: AppRoutes.admin.projects,
                loadComponent: () => import('./features/project/projects-component/projects-component').then(m => m.projectsComponent)
            }
        ]
       },
       //lazy loading dashboard module
       {
        path: AppRoutes.features.dashboard,
        component: MainLayout,
        canActivate: [authGuard], // Protegemos esta ruta con el guard de autenticación
        children:[
            {
                path: '',
                loadComponent: () => import('./features/dashboard/dashboard-component/dashboard-component').then(m => m.DashboardComponent)
            },
            {
                path: AppRoutes.features.clients,
                loadComponent: () => import('./features/clients/components/clients-component/clients-component').then(m => m.ClientsComponent)
            },
            {
                path:AppRoutes.features.projects,
                loadComponent: () => import('./features/project/projects-component/projects-component').then(m => m.projectsComponent)
            },
            {
                path: AppRoutes.features.users,
                loadComponent: () => import('./features/users/components/user-component/user-component').then(m => m.UserComponent)
            },
            {
                path: AppRoutes.features.files,
                loadComponent: () => import('./features/files/components/files-component/files-component').then(m => m.FilesComponent)
            },
            {
                path: 'client/:id',
                loadComponent: () => import('./features/project/dynamic-list-component/dynamic-list-component').then(m => m.DynamicListComponent)
            }
        ]
       },
   // app.routes.ts — ruta wildcard modificada
{
  path: '**',
  canActivate: [() => {
    const authState = inject(AuthStateService);
    const router = inject(Router);
    if (!authState.currentUser()) {
      return router.createUrlTree([AppRoutes.auth.login]);
    }
    return true;
  }],
  loadComponent: () => import('./shared/components/not-found-component/not-found-component')
    .then(m => m.NotFoundComponent)
}
    // {
    //     path: '**',
    //     loadComponent: () => import('./shared/components/not-found-component/not-found-component').then(m => m.NotFoundComponent)
    // }
]
