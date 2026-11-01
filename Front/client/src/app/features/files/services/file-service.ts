import { computed, inject, Injectable, signal } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { ApiRoutes } from '../../../core/routes/apiUrls';
import { FileRequest, FileData, FileUpdateRequest } from '../../../core';
import { FileType } from '../../../core/domain/models';
import { Observable, catchError, EMPTY, map, tap } from 'rxjs';
import { FileTypeRequest } from '../../../core/domain/dtos/fileType';
import Swal from '../../../shared/utils/swal';
import { projectService } from '../../project/services';

@Injectable({
  providedIn: 'root',
})
export class FileService {
  private http = inject(HttpClient);
  private projectService = inject(projectService);
  private dataVersion = 0;
  private searchRequestId = 0;
  private projectRequestIds = new Map<number, number>();
  state = signal({
    filesByProject: new Map<number, FileData[]>(),
    fileTypes: [] as FileTypeRequest[],
    loading: false
  });

  allFiles = computed(() =>
    Array.from(this.state().filesByProject.values())
      .flat()
      .filter(file => file && file.id !== undefined)
  );

  private beginSearchRequest(): () => boolean {
    const version = this.dataVersion;
    const requestId = ++this.searchRequestId;

    return () => version === this.dataVersion && requestId === this.searchRequestId;
  }

  private beginProjectRequest(projectId: number): () => boolean {
    const version = this.dataVersion;
    const requestId = (this.projectRequestIds.get(projectId) ?? 0) + 1;
    this.projectRequestIds.set(projectId, requestId);

    return () => version === this.dataVersion && this.projectRequestIds.get(projectId) === requestId;
  }

  private invalidateLoads(): void {
    this.dataVersion++;
    this.searchRequestId++;
  }

  // getFileTypes(): void {
  //   this.http.get<FileTypeRequest[]>(ApiRoutes.files.getFileTypes()).subscribe({
  //     next: (types) => {
  //       this.state.update(s => ({ ...s, fileTypes: types }));
  //     },
  //     error: (err) => {
  //       console.error('Error al obtener tipos de archivo:', err);
  //     }
  //   });
  // }

  // createFileTypes(types: {name: string}[]): Observable<FileTypeRequest[]> {
  //   return this.http.post<FileTypeRequest[]>(ApiRoutes.files.createFileTypes(), { types }).pipe(
  //     tap(createdTypes => {
  //       this.state.update(s => ({
  //         ...s,
  //         fileTypes: [...s.fileTypes, ...createdTypes]
  //       }));
  //     }),
  //     catchError(error => {
  //       console.error('Error al crear tipos de archivo en bloque:', error);
  //       throw error;
  //     })
  //   );
  // }

  getByProjectId(projectId: number): void {
    const isCurrentRequest = this.beginProjectRequest(projectId);
    this.state.update(s => ({ ...s, loading: true }));
    
    this.http.get<FileData[]>(ApiRoutes.files.getByProjectId(projectId)).subscribe({
      next: (files) => {
        if (!isCurrentRequest()) return;

        this.state.update(s => {
          const updatedMap = new Map(s.filesByProject);
          updatedMap.set(projectId, files);
          return { ...s, filesByProject: updatedMap, loading: false };
        });
      },
      error: (err) => {
        console.error('Error al obtener archivos:', err);
        if (isCurrentRequest()) {
          this.state.update(s => ({ ...s, loading: false }));
        }
      }
    });
  
  }
  // deleteFile(file: FileData){
  //   this.http.delete(ApiRoutes.files.delete(file.id)).pipe(
  //         catchError(error => {
  //           console.error("error",error)
  //           this.state.update(s => ({...s, loading: false}))
  //           return EMPTY
  //         })
  //       ).subscribe({
  //         next: () => {
  //           this.state.update(s => {
  //             const updatedMap = new Map(s.filesByProject);
  //             const currentFiles = updatedMap.get(file.project_id) ?? [];
  //             updatedMap.set(file.project_id, currentFiles.filter(f => f.id !== file.id));
  //             return {
  //               ...s,
  //               filesByProject: updatedMap,
  //               loading: false
  //             };
  //           });
  //         }
  //       })
  // }
//create file
async deleteFile(file: FileData): Promise<void> {
    // const result = await Swal.fire({
    //   title: '¿Eliminar archivo?',
    //   text: `Se eliminará "${file.tittle}" permanentemente`,
    //   icon: 'warning',
    //   showCancelButton: true,
    //   confirmButtonText: 'Sí, eliminar',
    //   cancelButtonText: 'Cancelar',
    //   confirmButtonColor: '#d33',
    //   cancelButtonColor: '#3085d6',
    // });

    // if (!result.isConfirmed) return;

    this.http.delete(ApiRoutes.files.delete(file.id)).pipe(
      catchError(error => {
        console.error("error", error)
        this.state.update(s => ({...s, loading: false}))
        return EMPTY
      })
    ).subscribe({
      next: () => {
        this.invalidateLoads();
        this.state.update(s => {
          const updatedMap = new Map(s.filesByProject);
          const projectIdKey = file.project_id ?? 0;
          const currentFiles = updatedMap.get(projectIdKey) ?? [];
          updatedMap.set(projectIdKey, currentFiles.filter(f => f.id !== file.id));
          return { ...s, filesByProject: updatedMap, loading: false };
        });
         this.projectService.state.update(s => {
          const updatedMap = new Map(s.projects);
          const projectIdKey = file.project_id ?? 0;  
          const project = updatedMap.get(projectIdKey); // ← usar la key ya sanitizada

          if (project && project.count_files !== null) {
              updatedMap.set(projectIdKey, {
                  ...project,
                  //restar 1 archivo al ser borrado
                  count_files: (project.count_files ?? 0) - 1
              });
          }
          
        return { ...s, projects: updatedMap };
         });
        Swal.fire({
          icon: 'success',
          title: 'Archivo eliminado',
          text: 'El archivo se ha eliminado correctamente',
          timer: 2000,
          showConfirmButton: false
        });
        this
      }
    })
  }
uploadFile(fileRequest: FileRequest) {
    this.state.update(s => ({ ...s, loading: true }));

    const projectId = fileRequest.project_id;
    if (projectId === undefined) {
      this.state.update(s => ({ ...s, loading: false }));
      return;
    }

    const formData = new FormData();
    formData.append('file', fileRequest.file);
    formData.append('project_id', projectId.toString());
    formData.append('tittle', fileRequest.tittle);
    formData.append('year', fileRequest.year.toString());
    formData.append('type_id', fileRequest.type_id.toString());

    // 1. Condición vital: Solo agregamos el project_id si viene en la petición
    // Si viene del formulario global (sin proyecto), simplemente no se envía este dato al backend
    // if (fileRequest.project_id !== undefined && fileRequest.project_id !== null) {
    //   formData.append('project_id', fileRequest.project_id.toString());
    // }

    this.http.post<FileData>(ApiRoutes.files.upload, formData).pipe(
      catchError(error => {
        console.error("Error al subir archivo:", error);
        this.state.update(s => ({ ...s, loading: false }));
        return EMPTY;
      })
    ).subscribe({
      next: (newFile) => {
        this.invalidateLoads();
        this.state.update(s => {
          const fileType = s.fileTypes.find(ft => ft.id === fileRequest.type_id);
          const resolvedFile = fileType ? { ...newFile, file_type: fileType.name } : newFile;
          const updatedMap = new Map(s.filesByProject);
          const projectIdKey = fileRequest.project_id ?? 0;
          const currentFiles = updatedMap.get(projectIdKey) ?? [];
          updatedMap.set(projectIdKey, [...currentFiles, resolvedFile]);
          
          return {
            ...s,
            filesByProject: updatedMap,
            loading: false
          };
        })
          this.projectService.state.update(s => {
        const updatedMap = new Map(s.projects);
        const project = updatedMap.get(projectId);
        if (project) {
            updatedMap.set(projectId, {
                ...project,
                count_files: (project.count_files ?? 0) + 1
            });
        }
        return { ...s, projects: updatedMap };
    });
        Swal.fire({
          icon: 'success',
          title: 'Archivo subido',
          text: 'El archivo se ha subido correctamente',
          timer: 2000,
          showConfirmButton: false
        });
      }
    });
  }
  downloadFile(fileId: number): void {
    window.open(ApiRoutes.files.download(fileId), '_blank');
  }
  // onSearch(query: string) {
  //   const params = new HttpParams().set('q', query);
  //   this.http.get<FileData[]>(ApiRoutes.files.search, { params }).pipe(
  //     map(files => {
  //       const filesMap = new Map<number, FileData[]>();
  //       if (files && files.length > 0) {
  //         files.forEach(file => {
  //           const projectFiles = filesMap.get(file.project_id) || [];
  //           projectFiles.push(file);
  //           filesMap.set(file.project_id, projectFiles);
  //         });
  //       }
  //       return filesMap;
  //     }),
  //     catchError(error => {
  //       console.error('Error al buscar archivos:', error);
  //       return EMPTY;
  //     })
  //   ).subscribe(filesMap => {
  //     this.state.update(s => ({ ...s, filesByProject: filesMap }));
  //   });
  // }
  onSearch(query: string) {
    const isCurrentRequest = this.beginSearchRequest();
    this.state.update(s => ({ ...s, loading: true }));

    this.http.get<FileData[]>(ApiRoutes.files.search(query)).pipe(
      map(files => {
        const filesMap = new Map<number, FileData[]>();
        if (files && files.length > 0) {
          files.forEach(file => {
            
            // Si el archivo no tiene proyecto, lo agrupamos en la llave 0
            const projectIdKey = file.project_id ?? 0; 
            
            const projectFiles = filesMap.get(projectIdKey) || [];
            projectFiles.push(file);
            filesMap.set(projectIdKey, projectFiles);
            
          });
        }
        return filesMap;
      }),
      catchError(error => {
        console.error('Error al buscar archivos:', error);
        if (isCurrentRequest()) {
          this.state.update(s => ({ ...s, loading: false }));
        }
        return EMPTY;
      })
    ).subscribe(filesMap => {
      if (!isCurrentRequest()) return;

      this.state.update(s => ({ ...s, filesByProject: filesMap, loading: false }));
    });
  }
  editFile(fileId: number, data: FileUpdateRequest): void {
    const formData = new FormData();
    if (data.tittle !== undefined) formData.append('tittle', data.tittle);
    if (data.year !== undefined) formData.append('year', data.year.toString());
    if (data.type_id !== undefined) formData.append('type_id', data.type_id.toString());
    if (data.project_id !== undefined) formData.append('project_id', data.project_id.toString());

    this.http.put<FileData>(ApiRoutes.files.update(fileId), formData).pipe(
      catchError(error => {
        console.error("Error al editar archivo:", error);
        return EMPTY;
      })
    ).subscribe({
      next: (updatedFile) => {
        this.invalidateLoads();
        this.state.update(s => {
          const updatedMap = new Map(s.filesByProject);
          for (const [key, files] of updatedMap) {
            const idx = files.findIndex(f => f.id === fileId);
            if (idx !== -1) {
              const fileType = s.fileTypes.find(ft => ft.id === (data.type_id ?? files[idx].type_id));
              const resolvedFile = fileType ? { ...files[idx], ...updatedFile, file_type: fileType.name } : { ...files[idx], ...updatedFile };
              files[idx] = resolvedFile;
              updatedMap.set(key, [...files]);
              break;
            }
          }
          return { ...s, filesByProject: updatedMap };
        });
        Swal.fire({
          icon: 'success',
          title: 'Archivo actualizado',
          text: 'El archivo se ha actualizado correctamente',
          timer: 2000,
          showConfirmButton: false
        });
      }
    });
  }
}
