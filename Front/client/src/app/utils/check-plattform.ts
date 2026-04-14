import { isPlatformBrowser, isPlatformServer } from "@angular/common";
import { inject, Injectable, PLATFORM_ID } from "@angular/core";

@Injectable({
  providedIn: 'root',
})
export class CheckPlatformUtils {
    private plattformId = inject(PLATFORM_ID);

    checkIfBrowser(): boolean {
        return isPlatformBrowser(this.plattformId);
    }

    checkIfServer(): boolean {
        return isPlatformServer(this.plattformId);
    }
}