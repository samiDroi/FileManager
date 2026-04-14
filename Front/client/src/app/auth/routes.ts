import { Routes } from "@angular/router";
import { AuthLayout } from "../layouts";
import { AppRoutes } from "../core/routes/appRoutes";

export const AuthRoutes: Routes = [
    {
        //path is '', because these routes will be loaded under the main path defined in app.routes.ts
        //so the routes are like /login, /register, etc, but they are grouped under the same layout
        path: '',
        component: AuthLayout, 
        //grouping all auth related routes under the same layout
        children:[
            {
                path: AppRoutes.auth.login,
                loadComponent: () => import('./components/login-component/login-component').then(m => m.LoginComponent),
            },
            {
                path: AppRoutes.auth.register,
                loadComponent: () => import('./components/register-component/register-component').then(m => m.RegisterComponent),
            },
            {
                //default route for auth, redirecting to login
                path: '',
                redirectTo: AppRoutes.auth.login,
                pathMatch: 'full'
            }
        ]
    }
];