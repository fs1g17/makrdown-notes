import { screen, render } from "@testing-library/react"
import { CreateNoteDialog } from "./CreateNoteDialog"
import { queryClient, QueryClientProvider } from "@/lib/react-query-testing";
import userEvent from "@testing-library/user-event";

describe("Create note dialog", () => {
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
})