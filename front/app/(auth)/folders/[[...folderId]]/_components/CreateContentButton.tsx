"use client";

import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { FileText, Folder, Plus } from "lucide-react";

export default function CreateContentButton({
  onCreateNoteClick,
  onCreateFolderClick
}: {
  onCreateNoteClick: () => void;
  onCreateFolderClick: () => void;
}) {
  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button
          className="fixed bottom-6 right-6 h-12 w-12 rounded-full p-0 shadow-lg hover:shadow-xl transition-shadow"
          aria-label="Add new item"
        >
          <Plus className="h-7 w-7" />
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end" side="top">
        <DropdownMenuItem onClick={onCreateNoteClick}>
          <FileText className="h-4 w-4" />
          Note
        </DropdownMenuItem>
        <DropdownMenuItem onClick={onCreateFolderClick}>
          <Folder className="h-4 w-4" />
          Folder
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}