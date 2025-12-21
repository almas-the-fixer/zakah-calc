# ZakahCalc

A Go + React tool to calculate Zakah using live gold/silver prices and optional currency conversion. Includes a Fiber API with Swagger docs and a Vite/React frontend.

## Features
- Calculate net assets (gold, silver, cash, business assets, liabilities) and Zakah payable at 2.5% when above the *silver-based* nisab.
- Live gold/silver prices from apised.com; supports USD plus common currencies via Frankfurter exchange API.
- Simple React UI with Tailwind styles, Axios call to the API, and result cards.
- Swagger docs served by the backend at `/swagger/*`.
- Dockerized dev setup for backend and frontend; optional Postgres service scaffolded but currently commented out.

## Tech Stack
- Backend: Go (Fiber), Swagger, godotenv
- Frontend: React 19, Vite, TypeScript, Tailwind
- Tooling: Docker Compose

## Prerequisites
- Docker & Docker Compose (for the simplest start)
- Or: Go 1.21+ (module file lists 1.25.5; use a modern stable Go toolchain)
- Node.js 18+ / npm for the frontend
- An apised.com API key for metal prices: `APISED_SECRET_KEY`

## Quick Start with Docker
1) Create a `.env` file in the repo root (used by the backend container):
   ```env
   APISED_SECRET_KEY=your_apised_key
   # DB_* variables are scaffolded but unused unless you enable Postgres
   ```
2) Build and run:
   ```bash
   docker compose up --build
   ```
3) Services:
   - Backend API: http://localhost:8080
   - Swagger UI: http://localhost:8080/swagger/index.html
   - Frontend (served by Nginx image built from Vite output): http://localhost:5173

> Note: The Postgres service is commented out in docker-compose.yml. Uncomment and configure if you later persist data.

## Running Backend Locally (without Docker)
1) Install dependencies (Go modules download automatically on build).
2) Create `.env` in repo root or `backend/.env`:
   ```env
   APISED_SECRET_KEY=your_apised_key
   ```
3) Start the server from `backend/`:
   ```bash
   go run ./cmd/main.go
   ```
4) Visit http://localhost:8080 and Swagger at http://localhost:8080/swagger/index.html

### Backend API
- `POST /calculate-zakah`
  - Body (JSON):
    ```json
    {
      "currency": "USD",
      "gold_grams": 10,
      "silver_grams": 50,
      "cash": 1000,
      "business_assets": 0,
      "liabilities": 200
    }
    ```
  - Response example:
    ```json
    {
      "total_assets": 1234.56,
      "nisab_threshold": 500.00,
      "zakah_payable": 30.86,
      "currency": "USD",
      "local_currency": "USD",
      "message": "Zakah is applicable"
    }
    ```

## Running Frontend Locally (without Docker)
1) From `frontend/`:
   ```bash
   npm install
   npm run dev
   ```
2) Open the Vite dev server URL (default http://localhost:5173). The form calls the backend at `http://localhost:8080/calculate-zakah`.

## Configuration Notes
- Environment variables read by the backend:
  - `APISED_SECRET_KEY` (required for metal prices)
  - `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_PORT`, `DB_HOST` are currently unused; they are placeholders for a future Postgres setup.
- CORS is enabled globally in the Fiber app for local frontend access.
- Currency conversion uses Frankfurter; if the currency is `USD` or empty, rate defaults to 1.0.

## Project Structure
- Backend Go entry: `backend/cmd/main.go`
- HTTP handlers: `backend/handlers/ZakahCalcHandler.go`
- Request/response types: `backend/types/types.go`
- Frontend entry: `frontend/src/main.tsx`, `frontend/src/App.tsx`
- Zakah form UI: `frontend/src/components/ZakahForm.tsx`

## Testing & Linting
- Frontend lint: `npm run lint` (from `frontend/`)
- Frontend build: `npm run build`
- Backend: no automated tests yet; run `go test ./...` after adding tests.

## Future Ideas
- Enable and persist calculations to Postgres.
- Add user auth or rate limiting for the API.
- Expand currency list and add offline fallbacks for price feeds.
- Add unit/integration tests for handler math and API calls.
