"use client";

import clientFetch from "@/lib/client-side-fetching";
import { Folder, FolderContent } from "@/types/folders";
import { Note } from "@/types/notes";
import { useQuery } from "@tanstack/react-query";
import { useParams, useRouter } from "next/navigation";
import { Folder as FolderIcon, FileText, Loader2 } from "lucide-react";

async function getFolderContent(folderId: number | undefined): Promise<FolderContent> {
  const response = await clientFetch.get<FolderContent>(`/api/folders${folderId ? `/${folderId}` : ""}`);
  return response.data;
}

function FolderItem({ folder, onClick }: { folder: Folder; onClick: () => void }) {
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

function NoteItem({ note, onClick }: { note: Note; onClick: () => void }) {
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

export default function Folders() {
  const router = useRouter();
  const params = useParams<{ folderId?: string[] }>();
  const folderId = params.folderId?.[0] ? Number(params.folderId[0]) : undefined;

  const { data, isPending, isError } = useQuery({
    queryKey: ["folders", { folderId: folderId ?? "root" }],
    queryFn: () => getFolderContent(folderId),
  });

  if (isPending) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    );
  }

  if (isError) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <p className="text-destructive">Failed to load folders</p>
      </div>
    );
  }

  const folders = data?.folders ?? [];
  const notes = data?.notes ?? [];
  const isEmpty = folders.length === 0 && notes.length === 0;

  const handleFolderClick = (folderId: number) => {
    router.push(`/folders/${folderId}`);
  };

  const handleNoteClick = (noteId: number) => {
    router.push(`/notes/${noteId}`);
  };

  return (
    <div className="min-h-screen p-6">
      <div className="mx-auto max-w-6xl">
        <h1 className="mb-6 text-2xl font-bold">My Files</h1>

        {isEmpty ? (
          <div className="flex flex-col items-center justify-center py-20 text-center">
            <FolderIcon className="mb-4 h-16 w-16 text-muted-foreground/50" />
            <p className="text-lg font-medium text-muted-foreground">No files yet</p>
            <p className="text-sm text-muted-foreground">
              Create a folder or note to get started
            </p>
          </div>
        ) : (
          <div className="space-y-6">
            {folders.length > 0 && (
              <section>
                <h2 className="mb-3 text-sm font-medium text-muted-foreground">
                  Folders
                </h2>
                <div className="grid grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
                  {folders.map((folder) => (
                    <FolderItem
                      key={folder.id}
                      folder={folder}
                      onClick={() => handleFolderClick(folder.id)}
                    />
                  ))}
                </div>
              </section>
            )}

            {notes.length > 0 && (
              <section>
                <h2 className="mb-3 text-sm font-medium text-muted-foreground">
                  Notes
                </h2>
                <div className="grid grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
                  {notes.map((note) => (
                    <NoteItem
                      key={note.id}
                      note={note}
                      onClick={() => handleNoteClick(note.id)}
                    />
                  ))}
                </div>
              </section>
            )}
          </div>
        )}
      </div>
    </div>
  );
}