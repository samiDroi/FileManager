export interface ClientRequest {
  name: string;
  contact_phone?: string;
  email?: string;
}

export interface ClientData {
  id: number;
  name: string;
  contact_phone?: string;
  email?: string;
  count_projects?: number;
  created_at: string; // ISO date string
}