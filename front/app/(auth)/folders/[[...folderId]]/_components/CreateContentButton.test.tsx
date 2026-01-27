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

  it("should execute folder callback when folder is clicked", async () => {
    const folderCallback = jest.fn();

    render(
      <CreateContentButton
        onCreateFolderClick={folderCallback}
        onCreateNoteClick={() => { }}
      />
    );

    const button = screen.getByLabelText("Add new item");
    await userEvent.click(button);

    const folder = screen.getByText("Folder");
    await userEvent.click(folder);

    expect(folderCallback).toHaveBeenCalled();
  });

  it("should execute note callback when note is clicked", async () => {
    const noteCallback = jest.fn();

    render(
      <CreateContentButton
        onCreateFolderClick={() => { }}
        onCreateNoteClick={noteCallback}
      />
    );

    const button = screen.getByLabelText("Add new item");
    await userEvent.click(button);

    const note = screen.getByText("Note");
    await userEvent.click(note);

    expect(noteCallback).toHaveBeenCalled();
  });
})