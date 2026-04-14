import { HttpInterceptorFn } from "@angular/common/http";
import { inject } from "@angular/core";
import { LocalKeys, LocalManagerService } from "../services/local-manager-service";

export const authInterceptor: HttpInterceptorFn = (req, next) => {
    const localManagerService = inject(LocalManagerService);
    const token = localManagerService.getElement(LocalKeys.token);  
    let headers = req.headers;

    if(token){
        headers = headers.set('Authorization', `Bearer ${token}`);
    }

    const authReq = req.clone({ headers });
    return next(authReq);
}
