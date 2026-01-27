import { screen, render, waitFor } from "@testing-library/react"
import { CreateNoteDialog } from "./CreateNoteDialog"
import { queryClient, QueryClientProvider } from "@/lib/react-query-testing";
import userEvent from "@testing-library/user-event";
import nock from "nock";
import { CreateNoteResponse } from "@/types/notes";
import { Toaster } from "@/components/ui/sonner";

const createNoteMock: CreateNoteResponse = {
  "id": 1,
  "folder_id": 1,
  "title": "title",
  "note": "",
  "created_at": "2026-01-27T17:06:22.00707+04:00",
  "updated_at": "2026-01-27T17:06:22.00707+04:00"
}

const errorNoteMock = {
  error: "note with this title already exists in this folder"
}

describe("Create note dialog", () => {
  beforeEach(() => {
    nock.cleanAll();
    queryClient.clear();
  });

  it("should display title", async () => {
    render(
      <QueryClientProvider>
        <CreateNoteDialog
          folderId={undefined}
          folderQueryKey={["folders", { folderId: "root" }]}
          open={true}
          onClose={() => { }}
        />
      </QueryClientProvider>
    );

    expect(screen.getByText("Title")).toBeVisible();
  });

  it("should display validation message if trying to create note with empty title", async () => {
    render(
      <QueryClientProvider>
        <CreateNoteDialog
          folderId={undefined}
          folderQueryKey={["folders", { folderId: "root" }]}
          open={true}
          onClose={() => { }}
        />
      </QueryClientProvider>
    );

    const createNoteButton = screen.getByText("Create Note");
    await userEvent.click(createNoteButton);

    const errorMessage = screen.getByRole("alert");
    expect(errorMessage).toBeVisible();
    expect(errorMessage).toHaveTextContent("Input a note title");
  });

  it("should display success toast and close the dialog when a note is created successfully", async () => {
    nock("http://localhost")
      .post(
        "/api/notes/new",
        {
          title: "title",
          note: "",
          folder_id: undefined
        }
      )
      .reply(201, createNoteMock);
    const onCloseMock = jest.fn();

    render(
      <QueryClientProvider>
        <CreateNoteDialog
          folderId={undefined}
          folderQueryKey={["folders", { folderId: 1 }]}
          open={true}
          onClose={onCloseMock}
        />
        <Toaster />
      </QueryClientProvider>
    );

    const input = screen.getByLabelText("Title");
    await userEvent.type(input, "title");

    const createNoteButton = screen.getByText("Create Note");
    await userEvent.click(createNoteButton);

    await waitFor(() => {
      expect(onCloseMock).toHaveBeenCalled();
    });
    expect(screen.getByText("Success creating note")).toBeVisible();
  });

  it("should display error toast and close the dialog when a note is failed to be created", async () => {
    nock("http://localhost")
      .post(
        "/api/notes/new",
        {
          title: "title",
          note: "",
          folder_id: undefined
        }
      )
      .reply(409, errorNoteMock);
    const onCloseMock = jest.fn();

    render(
      <QueryClientProvider>
        <CreateNoteDialog
          folderId={undefined}
          folderQueryKey={["folders", { folderId: 1 }]}
          open={true}
          onClose={onCloseMock}
        />
        <Toaster />
      </QueryClientProvider>
    );

    const input = screen.getByLabelText("Title");
    await userEvent.type(input, "title");

    const createNoteButton = screen.getByText("Create Note");
    await userEvent.click(createNoteButton);

    await waitFor(() => {
      expect(onCloseMock).toHaveBeenCalled();
    });
    expect(screen.getByText("Error creating note")).toBeVisible();
    expect(screen.getByText(errorNoteMock.error)).toBeVisible();
  });

  it("should display spinner during request execution", async () => {
    nock("http://localhost")
      .post(
        "/api/notes/new",
        {
          title: "title",
          note: "",
          folder_id: undefined
        }
      )
      .delay(100)
      .reply(201, createNoteMock);
    const onCloseMock = jest.fn();

    render(
      <QueryClientProvider>
        <CreateNoteDialog
          folderId={undefined}
          folderQueryKey={["folders", { folderId: 1 }]}
          open={true}
          onClose={onCloseMock}
        />
        <Toaster />
      </QueryClientProvider>
    );

    const input = screen.getByLabelText("Title");
    await userEvent.type(input, "title");

    const createNoteButton = screen.getByText("Create Note");
    expect(createNoteButton).toHaveAttribute("aria-busy", "false");
    await userEvent.click(createNoteButton);

    expect(createNoteButton).toHaveAttribute("aria-busy", "true");
    expect(screen.getByLabelText("Loading")).toBeVisible();
  });
})