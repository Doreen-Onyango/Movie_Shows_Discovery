import axios from 'axios';

const API_BASE_URL = process.env.REACT_APP_API_BASE_URL || 'http://localhost:8080/api/v1';

// Create axios instance with default config
const api = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor to add user ID header
api.interceptors.request.use((config) => {
  const userId = localStorage.getItem('userId');
  if (userId) {
    config.headers['X-User-ID'] = userId;
  }
  return config;
});

// Response interceptor for error handling
api.interceptors.response.use(
  (response) => response,
  (error) => {
    console.error('API Error:', error);
    return Promise.reject(error);
  }
);

// API methods
export const apiService = {
  // Health check
  healthCheck: () => api.get('/health'),

  // Search movies
  searchMovies: (query, page = 1) => 
    api.get(`/movies/search?q=${encodeURIComponent(query)}&page=${page}&per_page=10`),

  // Get trending movies
  getTrendingMovies: (page = 1) => 
    api.get(`/trending?page=${page}`),

  // Get movie details
  getMovieDetails: (movieId) => 
    api.get(`/movies/${movieId}`),

  // Get similar movies
  getSimilarMovies: (movieId, page = 1) => 
    api.get(`/movies/${movieId}/similar?page=${page}`),

  // Get genres
  getGenres: (page = 1) => 
    api.get('/movies/genres'),

  // Get movies by genre
  getMoviesByGenre: (genreId, page = 1) => 
    api.get(`/genres/${genreId}/movies?page=${page}`),

  // Watchlist operations
  getWatchlist: () => 
    api.get('/watchlist'),

  addToWatchlist: (movieId) => 
    api.post('/watchlist', { movie_id: movieId }),

  removeFromWatchlist: (movieId) => 
    api.delete(`/watchlist/${movieId}`),
};

export default apiService; 