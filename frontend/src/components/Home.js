import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { apiService } from '../services/api';
import { getLocalWatchlist } from '../utils/storage';
import MovieCard from './MovieCard';
import LoadingSpinner from './LoadingSpinner';
import GenreFilter from './GenreFilter';

const Home = () => {
  const [trendingMovies, setTrendingMovies] = useState([]);
  const [genres, setGenres] = useState([]);
  const [watchlist, setWatchlist] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [currentPage, setCurrentPage] = useState(1);
  const [hasMore, setHasMore] = useState(true);
  const [mediaType, setMediaType] = useState('all');

  useEffect(() => {
    loadInitialData(mediaType);
    // eslint-disable-next-line
  }, [mediaType]);

  const loadInitialData = async (type = mediaType) => {
    try {
      setLoading(true);
      const localWatchlist = getLocalWatchlist();
      setWatchlist(localWatchlist);
      const [trendingResponse, genresResponse] = await Promise.all([
        apiService.getTrendingMovies(1, type),
        apiService.getGenres()
      ]);
      setTrendingMovies(trendingResponse.data.trending || []);
      setGenres(genresResponse.data.genres || []);
      setHasMore(trendingResponse.data.total_pages > 1);
      setCurrentPage(1);
    } catch (err) {
      console.error('Error loading initial data:', err);
      setError('Failed to load movies. Please try again later.');
    } finally {
      setLoading(false);
    }
  };

  const loadMoreMovies = async () => {
    try {
      const nextPage = currentPage + 1;
      const response = await apiService.getTrendingMovies(nextPage, mediaType);
      const newMovies = response.data.trending || [];
      setTrendingMovies(prev => [...prev, ...newMovies]);
      setCurrentPage(nextPage);
      setHasMore(response.data.total_pages > nextPage);
    } catch (err) {
      console.error('Error loading more movies:', err);
    }
  };

  const handleWatchlistChange = (newWatchlist) => {
    setWatchlist(newWatchlist);
  };

  const handleTypeChange = (type) => {
    setMediaType(type);
    setCurrentPage(1);
  };

  if (loading) {
    return <LoadingSpinner text="Loading trending..." />;
  }

  if (error) {
    return (
      <div className="text-center py-12">
        <div className="text-red-600 text-lg mb-4">{error}</div>
        <button
          onClick={() => loadInitialData(mediaType)}
          className="px-4 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700"
        >
          Try Again
        </button>
      </div>
    );
  }

  return (
    <div className="space-y-8">
      {/* Hero Section */}
      <div className="text-center py-8">
        <h1 className="text-4xl font-bold text-gray-900 mb-4">
          Discover Amazing Movies & TV Shows
        </h1>
        <p className="text-xl text-gray-600 max-w-2xl mx-auto">
          Explore trending content, search for your favorites, and build your personal watchlist
        </p>
      </div>

      {/* Type Filter */}
      <div className="flex gap-4 mb-8 justify-center">
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

      {/* Genre Filter */}
      {genres.length > 0 && (
        <div className="mb-8">
          <h2 className="text-2xl font-semibold text-gray-900 mb-4">Browse by Genre</h2>
          <GenreFilter genres={genres} />
        </div>
      )}

      {/* Trending Movies */}
      <div>
        <div className="flex items-center justify-between mb-6">
          <h2 className="text-2xl font-semibold text-gray-900">Trending Now</h2>
          <Link
            to="/watchlist"
            className="text-primary-600 hover:text-primary-700 font-medium"
          >
            View Watchlist ({watchlist.length})
          </Link>
        </div>

        {trendingMovies.length > 0 ? (
          <>
            <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 gap-4">
              {trendingMovies.map((movie) => (
                <div key={movie.id} className="relative">
                  <MovieCard
                    movie={movie}
                    watchlist={watchlist}
                    onWatchlistChange={handleWatchlistChange}
                  />
                </div>
              ))}
            </div>

            {/* Load More Button */}
            {hasMore && (
              <div className="text-center mt-8">
                <button
                  onClick={loadMoreMovies}
                  className="px-6 py-3 bg-primary-600 text-white rounded-lg hover:bg-primary-700 transition-colors"
                >
                  Load More
                </button>
              </div>
            )}
          </>
        ) : (
          <div className="text-center py-12">
            <p className="text-gray-500 text-lg">No trending content available at the moment.</p>
          </div>
        )}
      </div>
    </div>
  );
};

export default Home; 