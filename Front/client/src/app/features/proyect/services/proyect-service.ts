import { inject, Injectable, signal } from '@angular/core';
import { Proyect } from '../../../core';
import { HttpClient } from '@angular/common/http';
import { ApiRoutes } from '../../../core/routes/apiUrls';

@Injectable({
  providedIn: 'root',
})
export class ProyectService {
  private http = inject(HttpClient)
  //state will save the data, updating the data how is required at the moment
  private state = signal({
    proyects: new Map<number,Proyect>()
  })

  getAll(){
    this.http.get<Proyect | null>(ApiRoutes.)
  }

}
