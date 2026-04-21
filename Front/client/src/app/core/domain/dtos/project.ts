
export interface projectCreated{
    name: String,
    client_id: number
}
export interface ProjectResponse {
      id: number;
      name: string;
      client_id: number;
      client_name: string;
      count_files: number | null;
      created_at: string;
    }