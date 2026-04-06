import { Routes } from '@angular/router';
import { AppRoutes } from './core/constants.routes';


export const routes: Routes = [
   {
        path: '',
        redirectTo: AppRoutes.auth.login, 
        pathMatch: 'full',
    },
    {
        path: '', 
        //import the auth routes.ts module lazily
        loadChildren: () => import('./auth/routes').then(m => m.AuthRoutes)
    }
]
