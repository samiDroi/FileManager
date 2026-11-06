export interface PermissionData {
  id: number;
  name: string;
  description: string;
}

export interface RoleData {
  id: number;
  name: string;
  description: string;
  permissions: PermissionData[];
  created_at?: string;
  updated_at?: string;
}

export interface RoleRequest {
  name: string;
  description: string;
  permission_ids: number[];
}

export interface GroupData {
  id: number;
  name: string;
  description: string;
  file_count: number;
  user_count: number;
  role_in_group?: string;
  can_manage?: boolean;
  created_at?: string;
}

export interface GroupRequest {
  name: string;
  description: string;
}

export interface GroupUserData {
  id: number;
  username: string;
  email: string;
  role_in_group: string;
  added_at?: string;
}

export interface AddGroupUserRequest {
  user_id: number;
  role_in_group: string;
}

export interface AvailableUserData {
  id: number;
  username: string;
  email: string;
}

export interface AddGroupFileRequest {
  file_id: number;
  access_type: string;
}

export interface GroupFileData {
  id: number;
  tittle: string;
  project_name: string;
  file_type: string;
  access_type: string;
  granted_at?: string;
}

export interface UserPermissionsData {
  user_id: number;
  is_admin: boolean;
  roles: string[];
  permissions: string[];
}
