import { Routes } from "@angular/router";
import path from "path";
import { AuthLayout } from '../../layouts'
import { App } from "../../app";
import { AppRoutes } from "../../core/constants.routes";

export const authRoutes: Routes = [
    {
        path: '',
        component: AuthLayout, 
        //grouping all auth related routes under the same layout
        children:[
            {
                path: AppRoutes.auth.login,
                
            }
        ]
    },
    
]