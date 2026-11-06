import { DestroyRef, inject, Injectable, signal } from '@angular/core';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { HttpClient, HttpParams } from '@angular/common/http';
import { catchError, EMPTY, map, Observable, Subject, switchMap, tap } from 'rxjs';
import { ApiRoutes } from '../../../core/routes/apiUrls';
import { ProjectData, ProjectRequest } from '../../../core/domain/dtos';
import Swal from '../../../shared/utils/swal';

@Injectable({
  providedIn: 'root',
})
export class projectService {
  private http = inject(HttpClient)
  private destroyRef = inject(DestroyRef);
  private readonly projectListRequests = new Subject<string | null>();
  public state = signal({
    projects: new Map<number, ProjectData>(),
    loading: false
  })

  constructor() {
    this.projectListRequests.pipe(
      tap(() => this.state.update(s => ({ ...s, loading: true }))),
      switchMap(query => this.requestProjects(query)),
      takeUntilDestroyed(this.destroyRef)
    ).subscribe(projects => {
      this.state.update(s => ({
        ...s,
        projects,
        loading: false
      }));
    });
  }

  private requestProjects(query: string | null): Observable<Map<number, ProjectData>> {
    const request$ = query === null
      ? this.http.get<ProjectData[] | null>(ApiRoutes.projects.getAll)
      : this.http.get<ProjectData[]>(ApiRoutes.projects.search, {
        params: new HttpParams().set('q', query)
      });

    return request$.pipe(
      map(projects => {
        const projectsMap = new Map<number, ProjectData>();
        (projects ?? []).forEach(project => projectsMap.set(project.id, project));
        return projectsMap;
      }),
      catchError(error => {
        console.error('Error al obtener proyectos:', error);
        this.state.update(s => ({ ...s, loading: false }));
        return EMPTY;
      })
    );
  }

  getAll(): void {
    this.projectListRequests.next(null);
  }

  onCreate(newProject: ProjectRequest){
    this.http.post<ProjectData>(ApiRoutes.projects.create,newProject).subscribe({
      next: (createdProject) => {
        this.state.update(s => {
          const createdMap = new Map(s.projects); 
          createdMap.set(createdProject.id, createdProject);
          return {
            ...s,
            projects: createdMap,
            loading: false
          };
        });
        this.getAll();
        Swal.fire({
          icon: 'success',
          title: 'Proyecto creado',
          text: 'El proyecto se ha creado correctamente',
          timer: 2000,
          showConfirmButton: false
        });
      },
      error: (err) => {
        console.error(err);
        this.state.update(s => ({ ...s, loading: false }));
      }
    })
  }

  onUpdate(id: number, updatedProject: ProjectRequest){
    this.http.put<ProjectData>(ApiRoutes.projects.update(id),updatedProject).subscribe({
      next: (updated) => {
        this.state.update(s => {
          const updatedMap = new Map(s.projects);
          updatedMap.set(updated.id, updated);
          return {
            ...s,
            projects: updatedMap,
            loading: false
          };
        });
        Swal.fire({
          icon: 'success',
          title: 'Proyecto actualizado',
          text: 'El proyecto se ha actualizado correctamente',
          timer: 2000,
          showConfirmButton: false
        });
      },
      error: (err) => {
        console.error(err);
        this.state.update(s => ({ ...s, loading: false }));
      }
    })
  }

  async onDelete(project: ProjectData): Promise<void> {
    const result = await Swal.fire({
      title: '¿Eliminar proyecto?',
      text: `Se eliminará el proyecto "${project.name}" y sus archivos dejarán de mostrarse.`,
      icon: 'warning',
      showCancelButton: true,
      confirmButtonText: 'Sí, eliminar',
      cancelButtonText: 'Cancelar',
      confirmButtonColor: '#d33',
      cancelButtonColor: '#3085d6',
    });

    if (!result.isConfirmed) {
      return;
    }

    this.http.delete<void>(ApiRoutes.projects.delete(project.id)).subscribe({
      next: () => {
        this.state.update(current => {
          const projects = new Map(current.projects);
          projects.delete(project.id);
          return { ...current, projects };
        });
        Swal.fire({
          icon: 'success',
          title: 'Proyecto eliminado',
          text: 'El proyecto se ha eliminado correctamente',
          timer: 2000,
          showConfirmButton: false,
        });
      },
      error: error => {
        console.error('Error al eliminar el proyecto:', error);
        Swal.fire('Error', 'No se pudo eliminar el proyecto', 'error');
      },
    });
  }

  searchProjects(query: string): void {
    if (!query.trim()) {
      this.getAll();
      return;
    }

    this.projectListRequests.next(query);
  }
}