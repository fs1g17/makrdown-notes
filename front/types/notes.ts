export type Note = {
  id: number;
  folder_id: number;
  title: string;
  note: string;
  created_at: string;
  updated_at: string;
};

export type CreateNoteResponse = {
  id: number;
  folder_id: number;
  title: string;
  note: string;
  created_at: string;
  updated_at: string;
}
