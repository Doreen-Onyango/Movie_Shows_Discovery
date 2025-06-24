import React, { useState, useEffect } from 'react';
import { useParams, Link } from 'react-router-dom';
import { BookmarkIcon, StarIcon, CalendarIcon, ClockIcon, CurrencyDollarIcon } from '@heroicons/react/24/outline';
import { BookmarkIcon as BookmarkSolidIcon } from '@heroicons/react/24/solid';
import { apiService } from '../services/api';
import { getLocalWatchlist, addToLocalWatchlist, removeFromLocalWatchlist } from '../utils/storage';
import { 
  getPosterUrl, 
  getBackdropUrl, 
  formatDate, 
  formatRuntime, 
  formatCurrency, 
  generateStars,
  calculateAverageRating,
  truncateText
} from '../utils/helpers';
import LoadingSpinner from './LoadingSpinner';
import MovieCard from './MovieCard';

const MovieDetail = () => {
  const { type, id } = useParams();
  const [movie, setMovie] = useState(null);
  const [similarMovies, setSimilarMovies] = useState([]);
  const [watchlist, setWatchlist] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [isAddingToWatchlist, setIsAddingToWatchlist] = useState(false);

  useEffect(() => {
    loadMovieData();
    loadWatchlist();
  }, [type, id]);

  const loadMovieData = async () => {
    try {
      setLoading(true);
      setError(null);

      const response = await apiService.getMediaDetails(type, id);
      let mediaData = response.data.data;
      if (!mediaData && response.data && response.data.success !== false) {
        mediaData = response.data;
      }
      if (!mediaData || response.data.success === false) {
        setError(type === 'tv' ? 'TV show not found.' : 'Movie not found.');
        setMovie(null);
        return;
      }
      setMovie(mediaData);
      setSimilarMovies(response.data.movies || []);
    } catch (err) {
      console.error('Error loading details:', err);
      setError('Failed to load details. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  const loadWatchlist = () => {
    const localWatchlist = getLocalWatchlist();
    setWatchlist(localWatchlist);
  };

  const handleWatchlistToggle = async () => {
    if (!movie) return;

    setIsAddingToWatchlist(true);
    try {
      const isInWatchlist = watchlist.some(item => item.id === movie.id);
      
      if (isInWatchlist) {
        removeFromLocalWatchlist(movie.id);
        setWatchlist(prev => prev.filter(item => item.id !== movie.id));
      } else {
        addToLocalWatchlist(movie);
        setWatchlist(prev => [...prev, { ...movie, added_at: new Date().toISOString() }]);
      }
    } catch (error) {
      console.error('Error updating watchlist:', error);
    } finally {
      setIsAddingToWatchlist(false);
    }
  };

  const isInUserWatchlist = watchlist.some(item => item.id === parseInt(id));

  if (loading) {
    return <LoadingSpinner text="Loading movie details..." />;
  }

  if (error) {
    return (
      <div className="text-center py-12">
        <div className="text-red-600 text-lg mb-4">{error}</div>
        <button
          onClick={loadMovieData}
          className="px-4 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700"
        >
          Try Again
        </button>
      </div>
    );
  }

  if (!movie) {
    return (
      <div className="text-center py-12">
        <h2 className="text-2xl font-semibold text-gray-900 mb-4">Movie Not Found</h2>
        <p className="text-gray-600">The movie you're looking for doesn't exist.</p>
      </div>
    );
  }

  const averageRating = calculateAverageRating(movie.ratings || []);

  return (
    <div className="space-y-8">
      {/* Trailer */}
      {movie.trailerKey && (
        <div className="my-6">
          <h2 className="text-2xl font-semibold text-gray-900 dark:text-gray-100 mb-3">Trailer</h2>
          <div className="w-full flex flex-col md:flex-row md:items-start gap-6">
            <iframe
              src={`https://www.youtube.com/embed/${movie.trailerKey}`}
              title="Movie Trailer"
              allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
              allowFullScreen
              className="w-full h-48 md:h-64 rounded-lg border border-gray-300 dark:border-gray-700"
              style={{ minWidth: 0 }}
            ></iframe>
          </div>
        </div>
      )}

      {/* Movie Header */}
      <div className="relative">
        {/* Backdrop Image */}
        <div className="relative h-96 bg-gray-900 rounded-lg overflow-hidden">
          <img
            src={getBackdropUrl(movie.backdrop_path)}
            alt={movie.title || movie.name}
            className="w-full h-full object-cover"
          />
          <div className="absolute inset-0 bg-gradient-to-t from-black/80 via-black/40 to-transparent"></div>
          
          {/* Movie Info Overlay */}
          <div className="absolute bottom-0 left-0 right-0 p-6">
            <div className="flex flex-col md:flex-row gap-6">
              {/* Poster */}
              <div className="flex-shrink-0">
                <img
                  src={getPosterUrl(movie.poster_path, 'w500')}
                  alt={movie.title || movie.name}
                  className="w-32 h-48 object-cover rounded-lg shadow-lg"
                />
              </div>
              
              {/* Movie Info */}
              <div className="flex-1 text-white">
                <h1 className="text-4xl font-bold mb-2">{movie.title || movie.name}</h1>
                
                <div className="flex items-center space-x-4 text-sm mb-3">
                  {movie.release_date && (
                    <div className="flex items-center space-x-1">
                      <CalendarIcon className="w-4 h-4" />
                      <span>{formatDate(movie.release_date)}</span>
                    </div>
                  )}
                  {movie.runtime && (
                    <div className="flex items-center space-x-1">
                      <ClockIcon className="w-4 h-4" />
                      <span>{formatRuntime(movie.runtime)}</span>
                    </div>
                  )}
                  {movie.budget && movie.budget > 0 && (
                    <div className="flex items-center space-x-1">
                      <CurrencyDollarIcon className="w-4 h-4" />
                      <span>{formatCurrency(movie.budget)}</span>
                    </div>
                  )}
                </div>

                {/* Genres */}
                {movie.genres && movie.genres.length > 0 && (
                  <div className="flex flex-wrap gap-2 mb-4">
                    {movie.genres.map((genre) => (
                      <Link
                        key={genre.id}
                        to={`/genre/${genre.id}`}
                        className="px-3 py-1 bg-white/20 text-white rounded-full text-sm hover:bg-white/30 transition-colors"
                      >
                        {genre.name}
                      </Link>
                    ))}
                  </div>
                )}

                {/* Watchlist Button */}
                <button
                  onClick={handleWatchlistToggle}
                  disabled={isAddingToWatchlist}
                  className="inline-flex items-center space-x-2 px-4 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700 disabled:opacity-50 transition-colors"
                >
                  {isAddingToWatchlist ? (
                    <div className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin"></div>
                  ) : isInUserWatchlist ? (
                    <BookmarkSolidIcon className="w-5 h-5 text-yellow-400" />
                  ) : (
                    <BookmarkIcon className="w-5 h-5" />
                  )}
                  <span>
                    {isInUserWatchlist ? 'Remove from Watchlist' : 'Add to Watchlist'}
                  </span>
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Movie Details */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        {/* Main Content */}
        <div className="lg:col-span-2 space-y-6">
          {/* Overview */}
          {movie.overview && (
            <div>
              <h2 className="text-2xl font-semibold text-gray-900 mb-3">Overview</h2>
              <p className="text-gray-700 leading-relaxed">{movie.overview}</p>
            </div>
          )}

          {/* Ratings */}
          {movie.ratings && movie.ratings.length > 0 && (
            <div>
              <h2 className="text-2xl font-semibold text-gray-900 mb-3">Ratings</h2>
              <div className="space-y-3">
                {movie.ratings.map((rating, index) => (
                  <div key={index} className="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
                    <div className="flex items-center space-x-2">
                      <span className="font-medium">{rating.source}</span>
                      <div className="flex items-center space-x-1">
                        <StarIcon className="w-4 h-4 text-yellow-400" />
                        <span>{rating.value}</span>
                      </div>
                    </div>
                    <div className="text-sm text-gray-500">
                      {rating.max_value && `${rating.value}/${rating.max_value}`}
                    </div>
                  </div>
                ))}
                {averageRating > 0 && (
                  <div className="flex items-center justify-between p-3 bg-primary-50 rounded-lg">
                    <div className="flex items-center space-x-2">
                      <span className="font-medium">Average Rating</span>
                      <div className="flex items-center space-x-1">
                        <StarIcon className="w-4 h-4 text-yellow-400" />
                        <span>{averageRating}</span>
                      </div>
                    </div>
                    <div className="text-sm text-primary-600">
                      {generateStars(averageRating)}
                    </div>
                  </div>
                )}
              </div>
            </div>
          )}

          {/* Cast */}
          {movie.cast && movie.cast.length > 0 && (
            <div>
              <h2 className="text-2xl font-semibold text-gray-900 mb-3">Cast</h2>
              <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-4">
                {movie.cast.slice(0, 8).map((person, index) => (
                  <div key={index} className="text-center">
                    <div className="w-16 h-16 bg-gray-200 rounded-full mx-auto mb-2 flex items-center justify-center">
                      <span className="text-gray-500 text-sm">{person.name?.charAt(0) || '?'}</span>
                    </div>
                    <p className="text-sm font-medium text-gray-900">{person.name}</p>
                    <p className="text-xs text-gray-500">{person.character}</p>
                  </div>
                ))}
              </div>
            </div>
          )}
        </div>

        {/* Sidebar */}
        <div className="space-y-6">
          {/* Quick Info */}
          <div className="bg-white p-4 rounded-lg shadow">
            <h3 className="text-lg font-semibold text-gray-900 mb-3">Quick Info</h3>
            <div className="space-y-2 text-sm">
              {movie.status && (
                <div className="flex justify-between">
                  <span className="text-gray-600">Status:</span>
                  <span className="font-medium">{movie.status}</span>
                </div>
              )}
              {movie.original_language && (
                <div className="flex justify-between">
                  <span className="text-gray-600">Language:</span>
                  <span className="font-medium">{movie.original_language.toUpperCase()}</span>
                </div>
              )}
              {movie.production_companies && movie.production_companies.length > 0 && (
                <div className="flex justify-between">
                  <span className="text-gray-600">Production:</span>
                  <span className="font-medium text-right">
                    {movie.production_companies[0].name}
                  </span>
                </div>
              )}
            </div>
          </div>

          {/* Vote Average */}
          {movie.vote_average && (
            <div className="bg-white p-4 rounded-lg shadow text-center">
              <div className="text-3xl font-bold text-primary-600 mb-1">
                {movie.vote_average.toFixed(1)}
              </div>
              <div className="text-sm text-gray-600 mb-2">User Score</div>
              <div className="flex justify-center text-yellow-400">
                {generateStars(movie.vote_average / 2)}
              </div>
            </div>
          )}
        </div>
      </div>

      {/* Similar Movies */}
      {similarMovies.length > 0 && (
        <div>
          <h2 className="text-2xl font-semibold text-gray-900 mb-4">Similar Movies</h2>
          <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 gap-4">
            {similarMovies.slice(0, 6).map((movie) => (
              <MovieCard
                key={movie.id}
                movie={movie}
                watchlist={watchlist}
                onWatchlistChange={setWatchlist}
              />
            ))}
          </div>
        </div>
      )}
    </div>
  );
};

export default MovieDetail; 