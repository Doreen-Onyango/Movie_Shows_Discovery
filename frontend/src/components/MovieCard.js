import React, { useState } from 'react';
import { Link } from 'react-router-dom';
import { BookmarkIcon, BookmarkSlashIcon, StarIcon } from '@heroicons/react/24/outline';
import { BookmarkIcon as BookmarkSolidIcon } from '@heroicons/react/24/solid';
import { getPosterUrl, truncateText, generateStars, isInWatchlist } from '../utils/helpers';
import { addToLocalWatchlist, removeFromLocalWatchlist } from '../utils/storage';

const MovieCard = ({ movie, watchlist = [], onWatchlistChange }) => {
  const [imageError, setImageError] = useState(false);
  const [isAddingToWatchlist, setIsAddingToWatchlist] = useState(false);

  const isInUserWatchlist = isInWatchlist(movie.id, watchlist);

  const handleWatchlistToggle = async (e) => {
    e.preventDefault();
    e.stopPropagation();
    
    setIsAddingToWatchlist(true);
    
    try {
      if (isInUserWatchlist) {
        removeFromLocalWatchlist(movie.id);
        if (onWatchlistChange) {
          onWatchlistChange(watchlist.filter(item => item.id !== movie.id));
        }
      } else {
        addToLocalWatchlist(movie);
        if (onWatchlistChange) {
          onWatchlistChange([...watchlist, { ...movie, added_at: new Date().toISOString() }]);
        }
      }
    } catch (error) {
      console.error('Error updating watchlist:', error);
    } finally {
      setIsAddingToWatchlist(false);
    }
  };

  const handleImageError = () => {
    setImageError(true);
  };

  return (
    <div className="bg-white rounded-lg shadow-md overflow-hidden hover:shadow-lg transition-shadow duration-300 group">
      <Link to={`/media/${movie.media_type || 'movie'}/${movie.id}`} className="block">
        {/* Movie Poster */}
        <div className="relative aspect-[2/3] bg-gray-200">
          {!imageError ? (
            <img
              src={getPosterUrl(movie.poster_path)}
              alt={movie.title || movie.name}
              className="w-full h-full object-cover"
              onError={handleImageError}
            />
          ) : (
            <div className="w-full h-full flex items-center justify-center bg-gray-200">
              <span className="text-gray-400 text-sm">No Image</span>
            </div>
          )}
          {/* Type Badge */}
          <span className={`absolute top-2 left-2 px-2 py-1 text-xs rounded bg-primary-600 text-white font-semibold uppercase`}>{movie.media_type === 'tv' ? 'TV' : 'Movie'}</span>
          {/* Watchlist Button */}
          <button
            onClick={handleWatchlistToggle}
            disabled={isAddingToWatchlist}
            className="absolute top-2 right-2 p-1.5 bg-black bg-opacity-50 rounded-full text-white hover:bg-opacity-70 transition-all duration-200 group-hover:bg-opacity-70"
          >
            {isAddingToWatchlist ? (
              <div className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin"></div>
            ) : isInUserWatchlist ? (
              <BookmarkSolidIcon className="w-4 h-4 text-yellow-400" />
            ) : (
              <BookmarkIcon className="w-4 h-4" />
            )}
          </button>

          {/* Rating Badge */}
          {movie.vote_average && (
            <div className="absolute bottom-2 left-2 bg-black bg-opacity-75 text-white px-2 py-1 rounded text-xs flex items-center space-x-1">
              <StarIcon className="w-3 h-3 text-yellow-400" />
              <span>{movie.vote_average.toFixed(1)}</span>
            </div>
          )}
        </div>

        {/* Movie Info */}
        <div className="p-4">
          <h3 className="font-semibold text-gray-900 text-sm line-clamp-2 mb-1">
            {movie.title || movie.name}
          </h3>
          
          <div className="flex items-center justify-between text-xs text-gray-500 mb-2">
            <span>{(movie.release_date || movie.first_air_date) ? new Date(movie.release_date || movie.first_air_date).getFullYear() : 'Unknown'}</span>
            {movie.genre_ids && movie.genre_ids.length > 0 && (
              <span className="truncate ml-2">
                {movie.genre_ids.slice(0, 2).join(', ')}
              </span>
            )}
          </div>

          {movie.overview && (
            <p className="text-xs text-gray-600 line-clamp-2">
              {truncateText(movie.overview, 80)}
            </p>
          )}
        </div>
      </Link>
    </div>
  );
};

export default MovieCard; 