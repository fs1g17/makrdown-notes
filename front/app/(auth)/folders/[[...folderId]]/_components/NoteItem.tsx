import { Note } from "@/types/notes";
import { FileText } from "lucide-react";

export function NoteItem({ note, onClick }: { note: Note; onClick: () => void }) {
  return (
    <div
      onClick={onClick}
      className="group flex cursor-pointer items-center gap-3 rounded-lg border border-border bg-card p-4 transition-colors hover:bg-accent"
    >
      <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-gray-100 dark:bg-gray-800">
        <FileText className="h-5 w-5 text-gray-600 dark:text-gray-400" />
      </div>
      <div className="flex-1 overflow-hidden">
        <p className="truncate font-medium">{note.title}</p>
      </div>
    </div>
  );
}

