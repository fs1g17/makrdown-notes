import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { FolderItem } from './FolderItem';
import { Folder } from '@/types/folders';

const mockFolder: Folder = {
  id: 0,
  user_id: 0,
  parent_id: null,
  name: 'test',
  created_at: '2026-01-25 12:59:45.059176+00',
  updated_at: '2026-01-25 12:59:45.059176+00'
}

describe("Folder items", () => {
  it("should display folder name", async () => {
    render(
      <FolderItem
        folder={mockFolder}
        onClick={() => { }}
      />
    );

    const element = screen.getByText(mockFolder.name);
    expect(element).toBeVisible();
  });

  it("should execute the callback on click", async () => {
    const mockCallback = jest.fn()
    render(
      <FolderItem
        folder={mockFolder}
        onClick={mockCallback}
      />
    );

    await userEvent.click(screen.getByText(mockFolder.name));
    expect(mockCallback).toHaveBeenCalled();
  });
})