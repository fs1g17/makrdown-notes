import { screen, render } from "@testing-library/react"
import { CreateNoteDialog } from "./CreateNoteDialog"
import { queryClient, QueryClientProvider } from "@/lib/react-query-testing";

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
})