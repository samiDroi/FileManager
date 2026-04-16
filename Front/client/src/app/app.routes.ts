import { Routes } from '@angular/router';
import { AppRoutes } from './core/routes/appRoutes';
import { AuthLayout, MainLayout } from './layouts';
import { authGuard, guestGuard } from './core/guards';


export const routes: Routes = [
   {
        path: '',
        redirectTo: AppRoutes.auth.login, 
        pathMatch: 'full',
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
                   //default route for auth, redirecting to login
                   path: '',
                   redirectTo: AppRoutes.auth.login,
                   pathMatch: 'full'
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
            }
        ]
       },
   
    {
        path: '**',
        loadComponent: () => import('./shared/components/not-found-component/not-found-component').then(m => m.NotFoundComponent)
    }
]
