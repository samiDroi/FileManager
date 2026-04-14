import { inject, Injectable, PLATFORM_ID } from '@angular/core';
import { CheckPlatformUtils } from '../../utils';

export enum LocalKeys{
  token = "token",
  userData = "userData"
}

@Injectable({
  providedIn: 'root',
})


export class LocalManagerService {
    checkPlatform = inject(CheckPlatformUtils);

    getElement(key: LocalKeys): string | null {
      if(this.checkPlatform.checkIfServer()) {
        return null;
      }
      return localStorage.getItem(key as string);
    }

    setElement(key: LocalKeys, value: string): void {
      if(this.checkPlatform.checkIfServer()) {
        return;
      }
      localStorage.setItem(key as string, value);
    }

    clearStorage(): void {
      if(this.checkPlatform.checkIfServer()) {
        return;
      }
      localStorage.clear();
    }
}
