import { screen, render, waitFor } from "@testing-library/react"
import { CreateFolderDialog } from "./CreateFolderDialog"
import { queryClient, QueryClientProvider } from "@/lib/react-query-testing";
import userEvent from "@testing-library/user-event";
import nock from "nock";
import { Toaster } from "@/components/ui/sonner";

const createFolderMock = {
  "id": 1,
  "parent_id": null,
  "name": "folder name",
  "created_at": "2026-01-27T17:06:22.00707+04:00",
  "updated_at": "2026-01-27T17:06:22.00707+04:00"
}

const errorFolderMock = {
  error: "folder with this name already exists in this folder"
}

describe("Create folder dialog", () => {
  beforeEach(() => {
    nock.cleanAll();
    queryClient.clear();
  });

  it("should display name", async () => {
    render(
      <QueryClientProvider>
        <CreateFolderDialog
          folderId={undefined}
          folderQueryKey={["folders", { folderId: "root" }]}
          open={true}
          onClose={() => { }}
        />
      </QueryClientProvider>
    );

    expect(screen.getByText("Name")).toBeVisible();
  });

  it("should display validation message if trying to create folder with empty name", async () => {
    render(
      <QueryClientProvider>
        <CreateFolderDialog
          folderId={undefined}
          folderQueryKey={["folders", { folderId: "root" }]}
          open={true}
          onClose={() => { }}
        />
      </QueryClientProvider>
    );

    const createFolderButton = screen.getByText("Create Folder");
    await userEvent.click(createFolderButton);

    const errorMessage = screen.getByRole("alert");
    expect(errorMessage).toBeVisible();
    expect(errorMessage).toHaveTextContent("Input a folder name");
  });

  it("should display success toast and close the dialog when a folder is created successfully", async () => {
    nock("http://localhost")
      .post(
        "/api/folders/new",
        {
          name: "folder name",
          parent_id: undefined
        }
      )
      .reply(201, createFolderMock);
    const onCloseMock = jest.fn();

    render(
      <QueryClientProvider>
        <CreateFolderDialog
          folderId={undefined}
          folderQueryKey={["folders", { folderId: 1 }]}
          open={true}
          onClose={onCloseMock}
        />
        <Toaster />
      </QueryClientProvider>
    );

    const input = screen.getByLabelText("Name");
    await userEvent.type(input, "folder name");

    const createFolderButton = screen.getByText("Create Folder");
    await userEvent.click(createFolderButton);

    await waitFor(() => {
      expect(onCloseMock).toHaveBeenCalled();
    });
    expect(screen.getByText("Success creating folder")).toBeVisible();
  });

  it("should display error toast and close the dialog when a folder is failed to be created", async () => {
    nock("http://localhost")
      .post(
        "/api/folders/new",
        {
          name: "folder name",
          parent_id: undefined
        }
      )
      .reply(409, errorFolderMock);
    const onCloseMock = jest.fn();

    render(
      <QueryClientProvider>
        <CreateFolderDialog
          folderId={undefined}
          folderQueryKey={["folders", { folderId: 1 }]}
          open={true}
          onClose={onCloseMock}
        />
        <Toaster />
      </QueryClientProvider>
    );

    const input = screen.getByLabelText("Name");
    await userEvent.type(input, "folder name");

    const createFolderButton = screen.getByText("Create Folder");
    await userEvent.click(createFolderButton);

    await waitFor(() => {
      expect(onCloseMock).toHaveBeenCalled();
    });
    expect(screen.getByText("Error creating folder")).toBeVisible();
    expect(screen.getByText(errorFolderMock.error)).toBeVisible();
  });

  it("should display spinner during request execution", async () => {
    nock("http://localhost")
      .post(
        "/api/folders/new",
        {
          name: "folder name",
          parent_id: undefined
        }
      )
      .delay(100)
      .reply(201, createFolderMock);
    const onCloseMock = jest.fn();

    render(
      <QueryClientProvider>
        <CreateFolderDialog
          folderId={undefined}
          folderQueryKey={["folders", { folderId: 1 }]}
          open={true}
          onClose={onCloseMock}
        />
        <Toaster />
      </QueryClientProvider>
    );

    const input = screen.getByLabelText("Name");
    await userEvent.type(input, "folder name");

    const createFolderButton = screen.getByText("Create Folder");
    expect(createFolderButton).toHaveAttribute("aria-busy", "false");
    await userEvent.click(createFolderButton);

    expect(createFolderButton).toHaveAttribute("aria-busy", "true");
    expect(screen.getByLabelText("Loading")).toBeVisible();
  });
})

