# markdown-notes

This is a little project for storing notes in markdown format. https://roadmap.sh/projects/markdown-note-taking-app

TODO:

- write an endpoint for getting contents of a folder
- ensure users can only access their own folders & files
- verify ownership in the notes store also

### GCloud run 


- `cd back && docker build -t markdown-notes-backend .`

tagging:
- `docker tag markdown-notes-backend europe-west2-docker.pkg.dev/markdown-notes-487211/markdown-notes-backend/markdown-notes-backend:v1.0.0` 
pushing:
- `docker push europe-west2-docker.pkg.dev/markdown-notes-487211/markdown-notes-backend/markdown-notes-backend:v1.0.0`

