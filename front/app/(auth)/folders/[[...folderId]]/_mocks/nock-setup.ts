import nock from 'nock';

const API_BASE = "http://localhost";

export function mockRootFolder() {
  return nock(API_BASE)
    .get("/api/folders")
    .reply(200, {
      "folder_id": 1,
      "notes": [
        {
          "id": 1,
          "folder_id": 1,
          "title": "test",
          "note": "hello world!",
          "created_at": "2026-01-25T19:00:35.0896+04:00",
          "updated_at": "2026-01-25T19:00:35.0896+04:00"
        }
      ],
      "folders": [
        {
          "id": 2,
          "user_id": 1,
          "parent_id": 1,
          "name": "subfolder",
          "created_at": "2026-01-25T19:01:04.921502+04:00",
          "updated_at": "2026-01-25T19:01:04.921502+04:00"
        }
      ]
    });
}

export function mockFolders1() {
  return nock(API_BASE)
    .get("/api/folders/1")
    .reply(200, {
      "folder_id": 1,
      "notes": [
        {
          "id": 1,
          "folder_id": 1,
          "title": "test",
          "note": "hello world!",
          "created_at": "2026-01-25T19:00:35.0896+04:00",
          "updated_at": "2026-01-25T19:00:35.0896+04:00"
        }
      ],
      "folders": [
        {
          "id": 2,
          "user_id": 1,
          "parent_id": 1,
          "name": "subfolder",
          "created_at": "2026-01-25T19:01:04.921502+04:00",
          "updated_at": "2026-01-25T19:01:04.921502+04:00"
        }
      ]
    });
}

export function mockFolders2() {
  return nock(API_BASE)
    .get("/api/folders/1")
    .reply(200, {
      "folder_id": 2,
      "notes": [],
      "folders": []
    });
}

export function mockFolders3() {
  return nock(API_BASE)
    .get("/api/folders/1")
    .reply(500, {
      "error": "folder doesn't exist or you don't have access to it"
    });
}
