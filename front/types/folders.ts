import { Note } from "./notes";

export type Folder = {
  id: number;
  user_id: number;
  parent_id: number | null;
  name: string;
  created_at: string;
  updated_at: string;
};

export type FolderContent = {
  folder_id: number;
  notes: Note[];
  folders: Folder[];
};

