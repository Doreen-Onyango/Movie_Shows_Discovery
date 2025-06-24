import React, { useState, useEffect } from 'react';
import { useSearchParams } from 'react-router-dom';
import { apiService } from '../services/api';
import { getLocalWatchlist } from '../utils/storage';
import MovieCard from './MovieCard';
import LoadingSpinner from './LoadingSpinner';

const SearchResults = () => {
  const [searchParams] = useSearchParams();
  const query = searchParams.get('q');
  
  const [movies, setMovies] = useState([]);
  const [watchlist, setWatchlist] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [currentPage, setCurrentPage] = useState(1);
  const [totalPages, setTotalPages] = useState(0);
  const [totalResults, setTotalResults] = useState(0);

  useEffect(() => {
    if (query) {
      setCurrentPage(1);
      searchMovies(query, 1);
    } else {
      // Clear results when query is empty
      setMovies([]);
      setTotalPages(0);
      setTotalResults(0);
      setError(null);
    }
  }, [query]);

  useEffect(() => {
    // Load watchlist from localStorage
    const localWatchlist = getLocalWatchlist();
    setWatchlist(localWatchlist);
  }, []);

  const searchMovies = async (searchQuery, page = 1) => {
    try {
      setLoading(true);
      setError(null);
      
      const response = await apiService.searchMovies(searchQuery, page);
      const { results: searchResults, meta } = response.data;
      
      const total_pages = meta?.total_pages || 0;
      const total_results = meta?.total_results || 0;
      
      if (page === 1) {
        setMovies(searchResults || []);
      } else {
        setMovies(prev => [...prev, ...(searchResults || [])]);
      }
      
      setTotalPages(total_pages);
      setTotalResults(total_results);
    } catch (err) {
      console.error('Error searching movies:', err);
      setError('Failed to search movies. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  const handleWatchlistChange = (newWatchlist) => {
    setWatchlist(newWatchlist);
  };

  if (!query) {
    return (
      <div className="text-center py-12">
        <h2 className="text-2xl font-semibold text-gray-900 mb-4">Search Movies</h2>
        <p className="text-gray-600">Enter a search term to find movies and TV shows.</p>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Search Header */}
      <div>
        {query && query.trim() && (
          <h1 className="text-3xl font-bold text-gray-900 mb-2">
            Search Results for "{query}"
          </h1>
        )}
        {totalResults > 0 && query && query.trim() && (
          <p className="text-gray-600">
            Found {totalResults} result{totalResults !== 1 ? 's' : ''}
          </p>
        )}
      </div>

      {/* Loading State */}
      {loading && movies.length === 0 && (
        <LoadingSpinner text="Searching movies..." />
      )}

      {/* Error State */}
      {error && (
        <div className="text-center py-8">
          <div className="text-red-600 text-lg mb-4">{error}</div>
          <button
            onClick={() => searchMovies(query, 1)}
            className="px-4 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700"
          >
            Try Again
          </button>
        </div>
      )}

      {/* Results */}
      {movies.length > 0 && (
        <>
          <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 gap-4">
            {movies.map((movie) => (
              <MovieCard
                key={movie.id}
                movie={movie}
                watchlist={watchlist}
                onWatchlistChange={handleWatchlistChange}
              />
            ))}
          </div>

          {/* Pagination Controls */}
          <div className="flex justify-center gap-4 mt-8">
            <button
              onClick={() => {
                if (currentPage > 1) {
                  const prevPage = currentPage - 1;
                  setCurrentPage(prevPage);
                  searchMovies(query, prevPage);
                }
              }}
              disabled={currentPage === 1 || loading}
              className="px-4 py-2 bg-primary-600 text-white rounded-lg disabled:opacity-50"
            >
              Back
            </button>
            <span>Page {currentPage} of {totalPages}</span>
            <button
              onClick={() => {
                if (currentPage < totalPages) {
                  const nextPage = currentPage + 1;
                  setCurrentPage(nextPage);
                  searchMovies(query, nextPage);
                }
              }}
              disabled={currentPage === totalPages || loading}
              className="px-4 py-2 bg-primary-600 text-white rounded-lg disabled:opacity-50"
            >
              Next
            </button>
          </div>
        </>
      )}

      {/* No Results */}
      {!loading && !error && movies.length === 0 && query && (
        <div className="text-center py-12">
          <div className="text-gray-500 text-lg mb-4">
            No movies found for "{query}"
          </div>
          <p className="text-gray-400">
            Try searching with different keywords or check your spelling.
          </p>
        </div>
      )}
    </div>
  );
};

export default SearchResults; 