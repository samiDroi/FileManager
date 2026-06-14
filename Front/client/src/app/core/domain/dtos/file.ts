//what we need to send to the back
export interface FileRequest {
  file: File;
  project_id?: number;
  tittle: string;
  year: number; // <- Corregido a number
  type_id: number;
}

export interface FileUpdateRequest {
  tittle?: string;
  year?: number;
  project_id?: number;
  type_id?: number;
}
//what we need to recieve from the back
export interface FileData {
  id: number;
  project_id?: number;
  tittle: string;
  path: string;
  year: number;
  type_id: number;
  file_type: string;
  created_at: string;
  project_name: string;
}