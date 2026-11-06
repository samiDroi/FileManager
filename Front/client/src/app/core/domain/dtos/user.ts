export interface UserData{
    id: number;
    username: string;
    email: string;
    createdBy: string;
    updatedBy: string;
    uploadedFiles?: number;
    roles?: AssignedRoleData[];
}

export interface AssignedRoleData {
    id: number;
    name: string;
    description: string;
}

export interface CreateUserRequest {
    username: string;
    email: string;
    password: string;
    role_id: number;
}

export interface UpdateUserRequest {
    username: string;
    email: string;
    role_id?: number;
}

export interface CreatedUserData extends UserData {
    roles: AssignedRoleData[];
}