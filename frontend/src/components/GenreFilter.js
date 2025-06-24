import React from 'react';
import { Link } from 'react-router-dom';

const GenreFilter = ({ genres }) => {
  return (
    <div className="flex flex-wrap gap-2">
      {genres.map((genre) => (
        <Link
          key={genre.id}
          to={`/genre/${genre.id}`}
          className="px-4 py-2 bg-gray-100 text-gray-700 rounded-full hover:bg-primary-100 hover:text-primary-700 transition-colors text-sm font-medium"
        >
          {genre.name}
        </Link>
      ))}
    </div>
  );
};

export default GenreFilter; 