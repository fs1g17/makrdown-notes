import { Note } from '@/types/notes';
import { render, screen } from '@testing-library/react';
import { NoteItem } from './NoteItem';
import userEvent from '@testing-library/user-event';

const mockNote: Note = {
  id: 0,
  folder_id: 0,
  title: 'title',
  note: 'title',
  created_at: '2026-01-25 12:59:45.059176+00',
  updated_at: '2026-01-25 12:59:45.059176+00'
}

describe("Note items", () => {
  it("should display note title", async () => {
    render(<NoteItem note={mockNote} onClick={() => { }} />);

    const element = screen.getByText(mockNote.title);
    expect(element).toBeVisible();
  });

  it("should execute the callback on click", async () => {
    const mockCallback = jest.fn()
    render(
      <NoteItem
        note={mockNote}
        onClick={mockCallback}
      />
    );

    await userEvent.click(screen.getByText(mockNote.title));
    expect(mockCallback).toHaveBeenCalled();
  });
})