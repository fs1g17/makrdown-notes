"use client";

import { useState, useEffect } from "react";
import { Note } from "@/types/notes";
import { useQuery, useMutation } from "@tanstack/react-query";
import clientFetch from "@/lib/client-side-fetching";
import { useParams, useRouter } from "next/navigation";
import { ArrowLeft, FileText, Loader2, Save } from "lucide-react";
import { Button } from "@/components/ui/button";
import { toast } from "sonner";
import { parseErrorMessage } from "@/lib/utils";

async function getNote(noteId: number | undefined): Promise<Note | undefined> {
  if (!noteId) return;
  const result = await clientFetch.get<Note>(`/api/notes/${noteId}`);
  return result.data;
}

export default function NoteEditor() {
  const router = useRouter();
  const params = useParams<{ noteId?: string[] }>();
  const noteId = params.noteId?.[0] ? Number(params.noteId[0]) : undefined;
  const [noteContent, setNoteContent] = useState("");

  const { data: note, isPending, isError, refetch } = useQuery({
    queryKey: ["notes", { noteId }],
    queryFn: () => getNote(noteId),
    enabled: !!noteId,
  });

  useEffect(() => {
    if (note) {
      setNoteContent(note.note);
    }
  }, [note]);

  const { mutate: saveNote, isPending: isSaving } = useMutation({
    mutationFn: (content: string) =>
      clientFetch.patch(`/api/notes/${noteId}/save`, { note: content }),
    onError: (e) => {
      const errorMessage = parseErrorMessage(e);
      toast.error("Error saving note", {
        description: errorMessage || "Your changes haven't been saved"
      });
    },
    onSuccess: () => {
      toast.success("Note saved", {
        description: "Your changes have been saved"
      });
      refetch();
    },
  });

  const handleSave = () => {
    saveNote(noteContent);
  };

  if (!noteId) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <div className="text-center">
          <FileText className="mx-auto mb-4 h-16 w-16 text-muted-foreground/50" />
          <p className="text-lg font-medium text-muted-foreground">No note selected</p>
        </div>
      </div>
    );
  }

  if (isPending) {
    return (
      <div role="status" aria-label="Loading note" className="flex min-h-screen items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
      </div>
    );
  }

  if (isError || !note) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <div className="text-center">
          <p className="text-lg font-medium text-destructive">Failed to load note</p>
          <Button
            variant="outline"
            className="mt-4"
            onClick={() => router.back()}
          >
            <ArrowLeft className="mr-2 h-4 w-4" />
            Go back
          </Button>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen p-6">
      <div className="mx-auto max-w-4xl">
        {/* Header */}
        <div className="mb-6 flex items-center gap-4">
          <Button
            variant="ghost"
            size="icon"
            onClick={() => router.back()}
          >
            <ArrowLeft className="h-5 w-5" />
          </Button>
          <div className="flex-1">
            <h1 className="text-2xl font-bold">{note.title}</h1>
            <p className="text-sm text-muted-foreground">
              Last updated: {new Date(note.updated_at).toLocaleDateString()}
            </p>
          </div>
          <Button
            onClick={handleSave}
            disabled={isSaving || noteContent === note.note}
          >
            {isSaving ? (
              <Loader2 className="mr-2 h-4 w-4 animate-spin" />
            ) : (
              <Save className="mr-2 h-4 w-4" />
            )}
            Save
          </Button>
        </div>

        {/* Editor */}
        <div className="rounded-lg border border-border bg-card">
          <textarea
            value={noteContent}
            onChange={(e) => setNoteContent(e.target.value)}
            className="min-h-[calc(100vh-250px)] w-full resize-none bg-transparent p-6 font-mono text-sm leading-relaxed focus:outline-none"
            placeholder="No content"
          />
        </div>
      </div>
    </div>
  );
}
