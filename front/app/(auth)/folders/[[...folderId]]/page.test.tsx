import { Note } from "@/types/notes";
import { Folder } from '@/types/folders';
import userEvent from '@testing-library/user-event'
import { FolderItem } from './_components/FolderItem';
import { render, screen } from '@testing-library/react';
import { renderHook, waitFor } from "@testing-library/react";
import Folders from "./page";
import { QueryClientProvider } from "@/lib/react-query-testing";
import { useParams, useRouter } from 'next/navigation'
import { mockRootFolder } from './_mocks/nock-setup'
import nock from 'nock'

const mockUseParams = useParams as jest.Mock
const mockUseRouter = useRouter as jest.Mock

describe("/folders page", () => {
  beforeEach(() => {
    mockUseRouter.mockReturnValue({ push: jest.fn() });
    nock.cleanAll();
  });

  it("should display the root folder content", async () => {
    mockUseParams.mockReturnValue({ folderId: undefined });
    mockRootFolder();

    render(
      <QueryClientProvider>
        <Folders />
      </QueryClientProvider>
    );

    const foldersElement = await screen.findByText("Folders");
    expect(foldersElement).toBeVisible();

    const notesElement = await screen.findByText("Notes");
    expect(notesElement).toBeVisible();
  });
});
