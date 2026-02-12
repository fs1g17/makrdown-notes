import nock from 'nock'
import NoteEditor from "./page";
import { Note } from '@/types/notes';
import { Toaster } from 'sonner';
import userEvent from '@testing-library/user-event';
import { useParams, useRouter } from 'next/navigation';
import { queryClient, QueryClientProvider } from "@/lib/react-query-testing";
import { render, screen, waitFor, waitForElementToBeRemoved } from '@testing-library/react';

const mockUseParams = useParams as jest.Mock
const mockUseRouter = useRouter as jest.Mock

const noteMock: Note = {
  "id": 1,
  "folder_id": 1,
  "title": "Test Note",
  "note": "hello world!",
  "created_at": "2026-01-25T19:00:35.0896+04:00",
  "updated_at": "2026-01-25T19:00:35.0896+04:00"
};

const errorMock = {
  "error": "note doesn't exist or you don't have access to it"
}

describe("/notes page", () => {
  beforeEach(() => {
    nock.cleanAll();
    queryClient.clear();
    mockUseRouter.mockReturnValue({ back: jest.fn() });
  });

  it("should display 'No note selected' when no noteId is provided", async () => {
    mockUseParams.mockReturnValue({ noteId: undefined });

    render(
      <QueryClientProvider>
        <NoteEditor />
      </QueryClientProvider>
    );

    expect(screen.getByText("No note selected")).toBeVisible();
  });

  it("should display loading spinner during fetch", async () => {
    mockUseParams.mockReturnValue({ noteId: ["1"] });
    nock("http://localhost").get("/api/notes/1").delay(100).reply(200, noteMock);

    render(
      <QueryClientProvider>
        <NoteEditor />
      </QueryClientProvider>
    );

    const loadingSpinner = screen.getByRole("status");
    expect(loadingSpinner).toBeVisible();
    await waitForElementToBeRemoved(() => screen.getByRole("status"));
    expect(screen.getByText("Test Note")).toBeVisible();
  });

  it("should display error when fetch fails", async () => {
    mockUseParams.mockReturnValue({ noteId: ["1"] });
    nock("http://localhost").get("/api/notes/1").reply(500, errorMock);

    render(
      <QueryClientProvider>
        <NoteEditor />
      </QueryClientProvider>
    );

    const errorMessage = await screen.findByText("Failed to load note");
    expect(errorMessage).toBeVisible();
  });

  it("should display note content when loaded", async () => {
    mockUseParams.mockReturnValue({ noteId: ["1"] });
    nock("http://localhost").get("/api/notes/1").reply(200, noteMock);

    render(
      <QueryClientProvider>
        <NoteEditor />
      </QueryClientProvider>
    );

    const noteTitle = await screen.findByText("Test Note");
    expect(noteTitle).toBeVisible();

    const textarea = screen.getByRole("textbox");
    expect(textarea).toHaveValue("hello world!");
  });

  it("should have save button disabled when content hasn't changed", async () => {
    mockUseParams.mockReturnValue({ noteId: ["1"] });
    nock("http://localhost").get("/api/notes/1").reply(200, noteMock);

    render(
      <QueryClientProvider>
        <NoteEditor />
      </QueryClientProvider>
    );

    await screen.findByText("Test Note");

    const saveButton = screen.getByRole("button", { name: /save/i });
    expect(saveButton).toBeDisabled();
  });

  it("should enable save button after editing content", async () => {
    mockUseParams.mockReturnValue({ noteId: ["1"] });
    nock("http://localhost").get("/api/notes/1").reply(200, noteMock);

    render(
      <QueryClientProvider>
        <NoteEditor />
      </QueryClientProvider>
    );

    await screen.findByText("Test Note");

    const textarea = screen.getByRole("textbox");
    await userEvent.type(textarea, " updated");

    const saveButton = screen.getByRole("button", { name: /save/i });
    expect(saveButton).toBeEnabled();
  });

  it("should save note successfully", async () => {
    mockUseParams.mockReturnValue({ noteId: ["1"] });
    nock("http://localhost").get("/api/notes/1").reply(200, noteMock);
    nock("http://localhost")
      .patch("/api/notes/1/save", { note: "hello world! updated" })
      .reply(200, { ...noteMock, note: "hello world! updated" });

    render(
      <QueryClientProvider>
        <NoteEditor />
      </QueryClientProvider>
    );

    await screen.findByText("Test Note");

    const textarea = screen.getByRole("textbox");
    await userEvent.type(textarea, " updated");

    const saveButton = screen.getByRole("button", { name: /save/i });
    await userEvent.click(saveButton);

    await waitFor(() => {
      expect(nock.isDone()).toBe(true);
    });
  });

  it("should show error toast when save fails", async () => {
    mockUseParams.mockReturnValue({ noteId: ["1"] });
    nock("http://localhost").get("/api/notes/1").reply(200, noteMock);
    nock("http://localhost")
      .patch("/api/notes/1/save", { note: "hello world! updated" })
      .reply(500, { error: "Failed to save" });

    render(
      <QueryClientProvider>
        <Toaster />
        <NoteEditor />
      </QueryClientProvider>
    );

    await screen.findByText("Test Note");

    const textarea = screen.getByRole("textbox");
    await userEvent.type(textarea, " updated");

    const saveButton = screen.getByRole("button", { name: /save/i });
    await userEvent.click(saveButton);

    const errorToast = await screen.findByText("Error saving note");
    expect(errorToast).toBeVisible();
  });

  it("should navigate back when go back button is clicked", async () => {
    const backMock = jest.fn();
    mockUseRouter.mockReturnValue({ back: backMock });
    mockUseParams.mockReturnValue({ noteId: ["1"] });
    nock("http://localhost").get("/api/notes/1").reply(200, noteMock);

    render(
      <QueryClientProvider>
        <NoteEditor />
      </QueryClientProvider>
    );

    await screen.findByText("Test Note");

    const backButtons = screen.getAllByRole("button");
    const backButton = backButtons[0]; // First button is the back button
    await userEvent.click(backButton);

    expect(backMock).toHaveBeenCalled();
  });

  it("should navigate back from error state", async () => {
    const backMock = jest.fn();
    mockUseRouter.mockReturnValue({ back: backMock });
    mockUseParams.mockReturnValue({ noteId: ["1"] });
    nock("http://localhost").get("/api/notes/1").reply(500, errorMock);

    render(
      <QueryClientProvider>
        <NoteEditor />
      </QueryClientProvider>
    );

    const goBackButton = await screen.findByRole("button", { name: /go back/i });
    await userEvent.click(goBackButton);

    expect(backMock).toHaveBeenCalled();
  });
});

