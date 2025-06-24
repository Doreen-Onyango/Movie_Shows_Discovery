import React, { useState, useEffect } from 'react';
import { getLocalWatchlist, removeFromLocalWatchlist, markAsWatchedInLocalWatchlist } from '../utils/storage';
import MovieCard from './MovieCard';

const Watchlist = () => {
  const [watchlist, setWatchlist] = useState([]);

  useEffect(() => {
    setWatchlist(getLocalWatchlist());
  }, []);

  const handleRemove = (movieId) => {
    const updated = removeFromLocalWatchlist(movieId);
    setWatchlist(updated);
  };

  const handleMarkAsWatched = (movieId) => {
    const updated = markAsWatchedInLocalWatchlist(movieId);
    setWatchlist(updated);
  };

  return (
    <div className="space-y-8">
      <h1 className="text-3xl font-bold text-gray-900 mb-6">My Watchlist</h1>
      {watchlist.length === 0 ? (
        <div className="text-center py-12">
          <p className="text-gray-500 text-lg mb-4">Your watchlist is empty.</p>
          <p className="text-gray-400">Browse and add movies to your watchlist to see them here.</p>
        </div>
      ) : (
        <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 gap-4">
          {watchlist.map((movie) => (
            <div key={movie.id} className="relative group bg-white dark:bg-gray-800 rounded-lg shadow-md overflow-hidden hover:shadow-lg transition-shadow duration-300">
              {/* Watched badge */}
              {movie.watched && (
                <span className="absolute top-2 left-2 bg-green-600 text-white text-xs font-semibold px-2 py-1 rounded z-10">Watched</span>
              )}
              <MovieCard
                movie={movie}
                watchlist={watchlist}
                onWatchlistChange={setWatchlist}
              />
              {/* Action buttons */}
              <div className="absolute top-2 right-2 flex flex-col space-y-2 opacity-0 group-hover:opacity-100 transition-opacity z-10">
                <button
                  onClick={() => handleRemove(movie.id)}
                  className="px-2 py-1 bg-red-600 text-white text-xs rounded hover:bg-red-700 shadow"
                  title="Remove from Watchlist"
                >
                  Remove
                </button>
                {!movie.watched && (
                  <button
                    onClick={() => handleMarkAsWatched(movie.id)}
                    className="px-2 py-1 bg-blue-600 text-white text-xs rounded hover:bg-blue-700 shadow"
                    title="Mark as Watched"
                  >
                    Mark as Watched
                  </button>
                )}
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

export default Watchlist; 