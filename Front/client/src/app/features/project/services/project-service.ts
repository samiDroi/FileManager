import { inject, Injectable, signal } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { ApiRoutes } from '../../../core/routes/apiUrls';
import { ProjectData, ProjectRequest } from '../../../core/domain/dtos';

@Injectable({
  providedIn: 'root',
})
export class projectService {
  private http = inject(HttpClient)
  public state = signal({
    projects: new Map<number, ProjectData>(),
    loading: false
  })

  getAll(){
    this.state.set({...this.state(),loading: true});
    this.http.get< ProjectData []>(ApiRoutes.projects.getAll).subscribe({
      next: (getProjects) =>{
        const newMap = new Map(this.state().projects);
        getProjects.forEach(p => newMap.set(p.id,p))

        this.state.update(s => ({
          ...s,
          projects: newMap,
          loading: false
        }))
      },
      error: (err) => {
        console.error(err)
        this.state.update(s => ({ ...s, loading: false }));
      }
    })
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
      },
      error: (err) => {
        console.error(err);
        this.state.update(s => ({ ...s, loading: false }));
      }
    })
  }

  searchProjects(query: string): void {
    if (!query.trim()) {
      this.getAll();
      return;
    }

    const params = new HttpParams().set('q', query);

    this.http.get<ProjectData[]>(ApiRoutes.projects.search, { params }).subscribe({
      next: (projects) => {
        const newMap = new Map<number, ProjectData>();
        if (projects) {
          projects.forEach(p => newMap.set(p.id, p));
        }
        this.state.update(s => ({
          ...s,
          projects: newMap,
          loading: false
        }));
      },
      error: (err) => {
        console.error('Error al buscar proyectos:', err);
        this.state.update(s => ({ ...s, loading: false }));
      }
    });
  }
}