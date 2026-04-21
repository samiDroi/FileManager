export interface ContactInfo{
  email?: string;
  contact_phone?: string;
}

export interface Client {
  id: number;
  name: string;
  contact_info?: ContactInfo;
  countprojects?: number;
  created_at?: string;
  updated_at?: string;
}
