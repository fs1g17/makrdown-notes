import nock from 'nock'
import Folders from "./page";
import { FolderContent } from '@/types/folders';
import userEvent from '@testing-library/user-event';
import { useParams, useRouter } from 'next/navigation';
import { queryClient, QueryClientProvider } from "@/lib/react-query-testing";
import { render, screen, waitForElementToBeRemoved } from '@testing-library/react';

const mockUseParams = useParams as jest.Mock
const mockUseRouter = useRouter as jest.Mock

const rootFolderContentMock: FolderContent = {
  "folder_id": 1,
  "notes": [
    {
      "id": 1,
      "folder_id": 1,
      "title": "test",
      "note": "hello world!",
      "created_at": "2026-01-25T19:00:35.0896+04:00",
      "updated_at": "2026-01-25T19:00:35.0896+04:00"
    }
  ],
  "folders": [
    {
      "id": 2,
      "user_id": 1,
      "parent_id": 1,
      "name": "subfolder",
      "created_at": "2026-01-25T19:01:04.921502+04:00",
      "updated_at": "2026-01-25T19:01:04.921502+04:00"
    }
  ]
};

const subFolderContentMock: FolderContent = {
  "folder_id": 2,
  "notes": [],
  "folders": []
};

const errorMock = {
  "error": "folder doesn't exist or you don't have access to it"
}

describe("/folders page", () => {
  beforeEach(() => {
    nock.cleanAll();
    queryClient.clear();
  });

  it("should display the folder and notes content when present", async () => {
    mockUseParams.mockReturnValue({ folderId: undefined });
    nock("http://localhost").get("/api/folders").reply(200, rootFolderContentMock);

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

  it("should display 'no files yet' in empty folder", async () => {
    mockUseParams.mockReturnValue({ folderId: ["2"] });
    nock("http://localhost").get("/api/folders/2").reply(200, subFolderContentMock);

    render(
      <QueryClientProvider>
        <Folders />
      </QueryClientProvider>
    );

    const emptyElement = await screen.findByText("No files yet");
    expect(emptyElement).toBeVisible();
  });

  it("should display loading spinner during fetch", async () => {
    mockUseParams.mockReturnValue({ folderId: ["1"] });
    nock("http://localhost").get("/api/folders/1").delay(100).reply(200, rootFolderContentMock);

    render(
      <QueryClientProvider>
        <Folders />
      </QueryClientProvider>
    );

    const loadingSpinner = screen.getByRole("status");
    expect(loadingSpinner).toBeVisible();
    await waitForElementToBeRemoved(() => screen.getByRole("status"));
    expect(screen.getByText("Folders")).toBeVisible();
  });

  it("should display error", async () => {
    mockUseParams.mockReturnValue({ folderId: ["2"] });
    nock("http://localhost").get("/api/folders/2").delay(100).reply(500, errorMock);

    render(
      <QueryClientProvider>
        <Folders />
      </QueryClientProvider>
    );

    const errorMessage = await screen.findByText("Failed to load folders");
    expect(errorMessage).toBeVisible();
  })

  it("should navigate to subfolder when clicked", async () => {
    mockUseParams.mockReturnValue({ folderId: undefined });
    const pushMock = jest.fn();
    mockUseRouter.mockReturnValue({ push: pushMock });
    nock("http://localhost").get("/api/folders").reply(200, rootFolderContentMock);

    render(
      <QueryClientProvider>
        <Folders />
      </QueryClientProvider>
    );

    const subfolder = await screen.findByText(rootFolderContentMock.folders[0].name);
    expect(subfolder).toBeVisible();

    await userEvent.click(subfolder)
    expect(pushMock).toHaveBeenCalledWith(`/folders/${rootFolderContentMock.folders[0].id}`);
  });

  it("should navigate to the note when clicked", async () => {
    mockUseParams.mockReturnValue({ folderId: undefined });
    const pushMock = jest.fn();
    mockUseRouter.mockReturnValue({ push: pushMock });
    nock("http://localhost").get("/api/folders").reply(200, rootFolderContentMock);

    render(
      <QueryClientProvider>
        <Folders />
      </QueryClientProvider>
    );

    const note = await screen.findByText(rootFolderContentMock.notes[0].title);
    expect(note).toBeVisible();

    await userEvent.click(note);
    expect(pushMock).toHaveBeenCalledWith(`/notes/${rootFolderContentMock.notes[0].id}`);
  });
});
