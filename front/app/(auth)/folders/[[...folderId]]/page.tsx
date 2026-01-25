"use client";

import clientFetch from "@/lib/client-side-fetching";
import { FolderContent } from "@/types/folders";
import { useQuery } from "@tanstack/react-query";
import { useParams } from "next/navigation";

async function getFolderContent(folderId: number | undefined): Promise<FolderContent> {
  const response = await clientFetch.get<FolderContent>(`/api/folders${folderId && `/${folderId}`}`);
  return response.data;
}

export default function Folders() {
  const params = useParams<{ folderId?: string[] }>();
  const folderId = params.folderId?.[0] ? Number(params.folderId[0]) : undefined;

  const { data, isPending, isError } = useQuery({
    queryKey: ["folders", { folderId: folderId ?? "root" }],
    queryFn: () => getFolderContent(folderId),
  })

  return (
    <div>folderId: {folderId}</div>
  )
}