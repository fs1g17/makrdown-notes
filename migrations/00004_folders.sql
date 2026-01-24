-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS folders (
  id BIGSERIAL PRIMARY KEY, 
  user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  parent_id BIGINT REFERENCES folders(id) ON DELETE CASCADE,
  name VARCHAR(255) NOT NULL, 
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (user_id, parent_id, name)
);

ALTER TABLE notes
  ADD COLUMN folder_id BIGINT NOT NULL REFERENCES folders(id) ON DELETE CASCADE;

ALTER TABLE notes 
  ADD CONSTRAINT notes_unique_title_per_folder UNIQUE (user_id, title, folder_id);

CREATE INDEX idx_folders_user_parent ON folders(user_id, parent_id);
CREATE INDEX idx_notes_user_folder ON notes(user_id, folder_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
