import { HttpClient } from '@angular/common/http';
import { inject, Injectable, signal } from '@angular/core';
import { FileTypeRequest } from '../../../core/domain/dtos/fileType';
import { ApiRoutes } from '../../../core/routes';

@Injectable({
  providedIn: 'root',
})
export class TypeFiles {
  private http = inject(HttpClient);
   state = signal({fileTypes: [] as FileTypeRequest[]});

  getFileTypes(): void {
    this.http.get<FileTypeRequest[]>(ApiRoutes.files.getFileTypes).subscribe({
       next: (fileTypes) => {
          this.state.set({ fileTypes });
        },
      error: (err) => {
        console.error('Error al obtener tipos de archivo:', err);
      }
    });
  }
  createFileTypes(types: {name: string}[]): void {
    this.http.post<FileTypeRequest[]>(ApiRoutes.files.createFileTypes, { types }).subscribe({
      next: (createdTypes) => {
        //set all the previous data and add the new types created
        this.state.update(s => ({
          ...s,
          fileTypes: [...s.fileTypes, ...createdTypes]
        }));
      },
      error: (err) => {
        console.error('Error al crear tipos de archivo en bloque:', err);
      }
    });
  }

  
}
