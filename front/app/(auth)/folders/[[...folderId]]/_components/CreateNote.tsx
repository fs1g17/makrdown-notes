"use client";

import { useState } from "react";
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
import { Loader2, Plus } from "lucide-react";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import clientFetch from "@/lib/client-side-fetching";
import { Controller, useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { useMutation, useQueryClient } from "@tanstack/react-query";

const schema = z.object({
  title: z
    .string()
    .trim()
    .min(1, { message: "Input a note title" }),
});

export function CreateNote({ folderId, onClick, folderQueryKey }: { folderId: number | undefined, onClick?: () => void, folderQueryKey: unknown[] }) {
  const [open, setOpen] = useState(false);
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
      clientFetch.post("/api/notes/new", {
        title,
        note: "",
        folder_id: folderId
      }),
    onSettled: () => {
      queryClient.invalidateQueries({
        queryKey: folderQueryKey
      });
      setOpen(false);
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
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button
          onClick={onClick}
          className="fixed bottom-6 right-6 h-12 w-12 rounded-full p-0 shadow-lg hover:shadow-xl transition-shadow"
          aria-label="Add new item"
        >
          <Plus className="h-7 w-7" />
        </Button>
      </DialogTrigger>
      <DialogContent>
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
                {isPending ? <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" /> : "Create Note"}
              </Button>
            </Field>
          </FieldGroup>
        </form>
      </DialogContent>
    </Dialog>
  );
}
