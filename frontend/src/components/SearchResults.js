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
  const [mediaType, setMediaType] = useState('all');

  useEffect(() => {
    if (query) {
      setCurrentPage(1);
      searchMovies(query, 1, mediaType);
    } else {
      // Clear results when query is empty
      setMovies([]);
      setTotalPages(0);
      setTotalResults(0);
      setError(null);
    }
  }, [query, mediaType]);

  useEffect(() => {
    // Load watchlist from localStorage
    const localWatchlist = getLocalWatchlist();
    setWatchlist(localWatchlist);
  }, []);

  const searchMovies = async (searchQuery, page = 1, type = mediaType) => {
    try {
      setLoading(true);
      setError(null);
      
      const response = await apiService.searchMovies(searchQuery, page, type);
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

  const handleTypeChange = (type) => {
    setMediaType(type);
    setCurrentPage(1);
  };

  if (!query) {
    return (
      <div className="text-center py-12">
        <h2 className="text-2xl font-semibold text-gray-900 mb-4">Search Movies & TV Shows</h2>
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

      {/* Type Filter */}
      <div className="flex gap-4 mb-4">
        <button
          className={`px-4 py-2 rounded-lg ${mediaType === 'all' ? 'bg-primary-600 text-white' : 'bg-gray-200 dark:bg-gray-700 text-gray-800 dark:text-gray-200'}`}
          onClick={() => handleTypeChange('all')}
        >
          All
        </button>
        <button
          className={`px-4 py-2 rounded-lg ${mediaType === 'movie' ? 'bg-primary-600 text-white' : 'bg-gray-200 dark:bg-gray-700 text-gray-800 dark:text-gray-200'}`}
          onClick={() => handleTypeChange('movie')}
        >
          Movies
        </button>
        <button
          className={`px-4 py-2 rounded-lg ${mediaType === 'tv' ? 'bg-primary-600 text-white' : 'bg-gray-200 dark:bg-gray-700 text-gray-800 dark:text-gray-200'}`}
          onClick={() => handleTypeChange('tv')}
        >
          TV Shows
        </button>
      </div>

      {/* Loading State */}
      {loading && movies.length === 0 && (
        <LoadingSpinner text="Searching..." />
      )}

      {/* Error State */}
      {error && (
        <div className="text-center py-8">
          <div className="text-red-600 text-lg mb-4">{error}</div>
          <button
            onClick={() => searchMovies(query, 1, mediaType)}
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
              <div key={movie.id} className="relative">
                <MovieCard
                  movie={movie}
                  watchlist={watchlist}
                  onWatchlistChange={handleWatchlistChange}
                />
                <span className={`absolute top-2 right-2 px-2 py-1 text-xs rounded bg-primary-600 text-white font-semibold uppercase`}>{movie.media_type === 'tv' ? 'TV' : 'Movie'}</span>
              </div>
            ))}
          </div>

          {/* Pagination Controls */}
          <div className="flex justify-center gap-4 mt-8">
            <button
              onClick={() => {
                if (currentPage > 1) {
                  const prevPage = currentPage - 1;
                  setCurrentPage(prevPage);
                  searchMovies(query, prevPage, mediaType);
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
                  searchMovies(query, nextPage, mediaType);
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
            No results found for "{query}"
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