import { Injectable, signal } from '@angular/core';

export interface ProjectBreadcrumb {
  clientName: string;
  projectName: string;
}

@Injectable({
  providedIn: 'root',
})
export class TopbarTitleService {
  private readonly projectBreadcrumb = signal<ProjectBreadcrumb | null>(null);
  readonly projectContext = this.projectBreadcrumb.asReadonly();

  setProjectBreadcrumb(context: ProjectBreadcrumb): void {
    this.projectBreadcrumb.set(context);
  }

  clearProjectBreadcrumb(): void {
    this.projectBreadcrumb.set(null);
  }
}
