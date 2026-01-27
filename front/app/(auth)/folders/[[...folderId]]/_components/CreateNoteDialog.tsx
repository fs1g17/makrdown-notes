"use client";

import z from "zod";
import {
  Field,
  FieldError,
  FieldGroup,
  FieldLabel,
} from "@/components/ui/field";
import {
  Dialog,
  DialogTrigger,
  DialogContent,
  DialogHeader,
  DialogTitle
} from "@/components/ui/dialog";
import { toast } from "sonner"
import { useState } from "react";
import { Loader2, Plus } from "lucide-react";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import clientFetch from "@/lib/client-side-fetching";
import { Controller, useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { AxiosError } from "axios";
import { parseErrorMessage } from "@/lib/utils";
import { CreateNoteResponse } from "@/types/notes";

const schema = z.object({
  title: z
    .string()
    .trim()
    .min(1, { message: "Input a note title" }),
});

export function CreateNoteDialog({
  folderId,
  folderQueryKey,
  open,
  onClose
}: {
  folderId: number | undefined,
  folderQueryKey: unknown[],
  open: boolean;
  onClose: () => void;
}) {
  const queryClient = useQueryClient();
  const form = useForm<z.infer<typeof schema>>({
    resolver: zodResolver(schema),
    defaultValues: {
      title: "",
    }
  });

  const { mutate: createNote, isPending } = useMutation({
    mutationFn: ({
      title
    }: {
      title: string;
    }) =>
      clientFetch.post<CreateNoteResponse>("/api/notes/new", {
        title,
        note: "",
        folder_id: folderId
      }),
    onError: (e) => {
      const errorMessage = parseErrorMessage(e);
      toast.error("Error creating note", {
        description: errorMessage
      });
    },
    onSuccess: () => {
      toast.success("Success creating note", {
        description: "Note created successfully"
      });
    },
    onSettled: () => {
      queryClient.invalidateQueries({
        queryKey: folderQueryKey
      });
      onClose();
    },
  })

  const handleSubmit = ({
    title
  }: {
    title: string;
  }) => {
    createNote({ title });
  }

  return (
    <Dialog open={open} onOpenChange={onClose}>
      <DialogContent aria-describedby={undefined}>
        <DialogHeader>
          <DialogTitle>Create a note</DialogTitle>
        </DialogHeader>
        <form id="create-note-form" onSubmit={form.handleSubmit(handleSubmit)}>
          <FieldGroup>
            <Controller
              name="title"
              control={form.control}
              render={({ field, fieldState }) => (
                <Field data-invalid={fieldState.invalid}>
                  <FieldLabel htmlFor="create-note-title">Title</FieldLabel>
                  <Input
                    id="create-note-title"
                    placeholder="Enter note title"
                    aria-invalid={fieldState.invalid}
                    {...field}
                  />
                  {fieldState.invalid && (
                    <FieldError errors={[fieldState.error]} />
                  )}
                </Field>
              )}
            />

            <Field>
              <Button type="submit" className="w-full" aria-busy={isPending} disabled={isPending}>
                {isPending ? <Loader2 aria-label="Loading" className="h-8 w-8 animate-spin text-muted-foreground" /> : "Create Note"}
              </Button>
            </Field>
          </FieldGroup>
        </form>
      </DialogContent>
    </Dialog>
  );
}
