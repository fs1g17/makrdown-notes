import { render, screen } from "@testing-library/react";
import CreateContentButton from "./CreateContentButton"
import userEvent from "@testing-library/user-event";

describe("Create content button", () => {
  it("should display options to create folder and note when clicked", async () => {
    render(
      <CreateContentButton
        onCreateFolderClick={() => { }}
        onCreateNoteClick={() => { }}
      />
    );

    const button = screen.getByLabelText("Add new item");
    expect(button).toBeVisible();

    await userEvent.click(button);

    expect(screen.getByText("Note")).toBeVisible();
    expect(screen.getByText("Folder")).toBeVisible();
  });


})