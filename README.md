# Social Network

A social network application with an antique gold aesthetic.

## Tech Stack
- **Frontend**: Next.js
- **Backend**: Go
- **Database**: SQLite

## Configuration
Before running, you must set the correct API URL in `frontend/src/config.mjs`:

- **For Docker**: Set `export const BASE_URL = buildURL;`
- **For Manual Setup**: Set `export const BASE_URL = devURL;`


## Run with Docker
```bash
docker-compose up --build
```

**What's inside?**
- **ns-backend**: The Go server (API & Database).
- **ns-frontend**: The Next.js web application.


## Manual Setup
```bash
cd backend
go run .
```
3. Run Frontend:
```bash
cd frontend
npm install
npm run dev
```

The application will be available at `http://localhost:3000`.
