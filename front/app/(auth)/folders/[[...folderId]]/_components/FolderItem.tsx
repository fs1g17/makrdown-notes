import { Folder } from "@/types/folders";
import { Folder as FolderIcon } from "lucide-react";

export function FolderItem({ folder, onClick }: { folder: Folder; onClick: () => void }) {
  return (
    <div
      onClick={onClick}
      className="group flex cursor-pointer items-center gap-3 rounded-lg border border-border bg-card p-4 transition-colors hover:bg-accent"
    >
      <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-blue-100 dark:bg-blue-900/30">
        <FolderIcon className="h-5 w-5 text-blue-600 dark:text-blue-400" />
      </div>
      <div className="flex-1 overflow-hidden">
        <p className="truncate font-medium">{folder.name}</p>
      </div>
    </div>
  );
}

