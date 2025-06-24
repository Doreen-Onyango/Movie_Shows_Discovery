import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { apiService } from '../services/api';
import { getLocalWatchlist } from '../utils/storage';
import MovieCard from './MovieCard';
import LoadingSpinner from './LoadingSpinner';

const GenreMovies = () => {
  const { id } = useParams();
  const [movies, setMovies] = useState([]);
  const [genreName, setGenreName] = useState('');
  const [watchlist, setWatchlist] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [currentPage, setCurrentPage] = useState(1);
  const [totalPages, setTotalPages] = useState(0);

  useEffect(() => {
    fetchMoviesByGenre(1);
    // eslint-disable-next-line
  }, [id]);

  useEffect(() => {
    setWatchlist(getLocalWatchlist());
  }, []);

  const fetchMoviesByGenre = async (page = 1) => {
    try {
      setLoading(true);
      setError(null);
      const response = await apiService.getMoviesByGenre(id, page);
      setMovies(page === 1 ? response.data.movies : prev => [...prev, ...response.data.movies]);
      setGenreName(response.data.genre_name || '');
      setTotalPages(response.data.total_pages || 0);
      setCurrentPage(page);
    } catch (err) {
      setError('Failed to load movies for this genre.');
    } finally {
      setLoading(false);
    }
  };

  const loadMore = () => {
    if (currentPage < totalPages) {
      fetchMoviesByGenre(currentPage + 1);
    }
  };

  const handleWatchlistChange = (newWatchlist) => {
    setWatchlist(newWatchlist);
  };

  if (loading && movies.length === 0) {
    return <LoadingSpinner text="Loading movies by genre..." />;
  }

  if (error) {
    return <div className="text-center py-12 text-red-600">{error}</div>;
  }

  return (
    <div className="space-y-8">
      <h1 className="text-3xl font-bold text-gray-900 mb-6">
        {genreName ? `Genre: ${genreName}` : 'Movies by Genre'}
      </h1>
      {movies.length === 0 ? (
        <div className="text-center py-12 text-gray-500">No movies found for this genre.</div>
      ) : (
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
          {currentPage < totalPages && (
            <div className="text-center mt-8">
              <button
                onClick={loadMore}
                className="px-6 py-3 bg-primary-600 text-white rounded-lg hover:bg-primary-700 transition-colors"
              >
                Load More
              </button>
            </div>
          )}
        </>
      )}
    </div>
  );
};

export default GenreMovies; 