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
  DialogContent,
  DialogHeader,
  DialogTitle
} from "@/components/ui/dialog";
import { toast } from "sonner"
import { Loader2 } from "lucide-react";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import clientFetch from "@/lib/client-side-fetching";
import { Controller, useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { parseErrorMessage } from "@/lib/utils";

const schema = z.object({
  name: z
    .string()
    .trim()
    .min(1, { message: "Input a folder name" }),
});

export function CreateFolderDialog({
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
      name: "",
    }
  });

  const { mutate: createFolder, isPending } = useMutation({
    mutationFn: ({
      name
    }: {
      name: string;
    }) =>
      clientFetch.post("/api/folders/new", {
        name,
        parent_id: folderId
      }),
    onError: (e) => {
      const errorMessage = parseErrorMessage(e);
      toast.error("Error creating folder", {
        description: errorMessage
      });
    },
    onSuccess: () => {
      toast.success("Success creating folder", {
        description: "Folder created successfully"
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
    name
  }: {
    name: string;
  }) => {
    createFolder({ name });
  }

  return (
    <Dialog open={open} onOpenChange={onClose}>
      <DialogContent aria-describedby={undefined}>
        <DialogHeader>
          <DialogTitle>Create a folder</DialogTitle>
        </DialogHeader>
        <form id="create-folder-form" onSubmit={form.handleSubmit(handleSubmit)}>
          <FieldGroup>
            <Controller
              name="name"
              control={form.control}
              render={({ field, fieldState }) => (
                <Field data-invalid={fieldState.invalid}>
                  <FieldLabel htmlFor="create-folder-name">Name</FieldLabel>
                  <Input
                    id="create-folder-name"
                    placeholder="Enter folder name"
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
                {isPending ? <Loader2 aria-label="Loading" className="h-8 w-8 animate-spin text-muted-foreground" /> : "Create Folder"}
              </Button>
            </Field>
          </FieldGroup>
        </form>
      </DialogContent>
    </Dialog>
  );
}

