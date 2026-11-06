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
  can_edit: boolean;
  can_delete: boolean;
  can_download: boolean;
}

export interface FileAccessUserData {
  user_id: number;
  username: string;
  access_type: string;
  granted_by: number;
  granted_at: string;
}

export interface FileAccessGroupData {
  group_id: number;
  group_name: string;
  access_type: string;
  granted_by: number;
  granted_at: string;
}

export interface FileAccessData {
  users: FileAccessUserData[];
  groups: FileAccessGroupData[];
}

export interface GrantFileUserAccessRequest {
  user_id: number;
  access_type: string;
}

export interface GrantFileGroupAccessRequest {
  group_id: number;
  access_type: string;
}