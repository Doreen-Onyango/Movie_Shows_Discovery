# Movie Shows Discovery Frontend

A modern React frontend for discovering movies, searching, browsing genres, managing a watchlist, and more. This app connects to the Movie Shows Discovery backend API.

---

## Features
- Debounced search for movies and TV shows
- Trending movies and genre browsing
- Movie detail pages with ratings
- Add/remove movies to your watchlist (localStorage)
- Pagination (10 movies per page, with Back/Next navigation)
- Responsive design with Tailwind CSS
- Recommendations and similar movies
- Loading and error states, fallback UI

---

## Prerequisites
- **Node.js** (v18 or higher recommended)
- **npm** (v9 or higher recommended)
- The backend API running (see backend/README.md)

---

## ⚡ Installation & Setup

1. **Clone the repository**
   ```bash
   git clone https://github.com/Doreen-Onyango/Movie_Shows_Discovery
   cd Movie_Shows_Discovery/frontend
   ```

2. **Install dependencies**
   ```bash
   npm install
   ```

3. **Configure environment variables**
   - By default, the frontend expects the backend API at `http://localhost:8080/api/v1`.
   - To use a different backend URL, create a `.env` file in the `frontend` directory:
     ```env
     REACT_APP_API_BASE_URL=http://localhost:8080/api/v1
     ```
   - If you use the default backend port and path, you can skip this step.

4. **Tailwind CSS**
   - Tailwind CSS is already configured via `tailwind.config.js` and `postcss.config.js`.
   - No extra setup is needed.

5. **Start the backend API**
   - Make sure the backend is running (see backend/README.md for setup and API keys).

6. **Run the frontend app**
   ```bash
   npm start
   ```
   - The app will open at [http://localhost:3000](http://localhost:3000)
   - The frontend will proxy API requests to the backend.

---

## Available Scripts
- `npm start` — Start the development server
- `npm run build` — Build for production


---

## How to Use the App
- **Search**: Use the search bar to find movies/TV shows. Results are paginated (10 per page). Use Back/Next to navigate.
- **Trending**: See trending movies on the home page.
- **Genres**: Browse by genre from the home page.
- **Movie Details**: Click a movie for full details, ratings, and similar movies.
- **Watchlist**: Add/remove movies to your watchlist (stored in your browser).
- **Recommendations**: Get personalized recommendations from your watchlist page.

---

## Troubleshooting
- **CORS errors**: Ensure the backend is running and accessible at the URL set in `REACT_APP_API_BASE_URL`.
- **API errors**: Make sure your backend `.env` has valid TMDB and OMDB API keys, and the backend is running.
- **Tailwind/PostCSS errors**: Ensure all dependencies are installed (`npm install`).
- **Blank page or network errors**: Check the browser console and network tab for failed API requests.

---

## Need Help?
If you run into issues, check the backend is running and reachable, and review the browser console for errors. For further help, open an issue.
