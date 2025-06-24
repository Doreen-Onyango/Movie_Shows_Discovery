// Debounce function for search input
export const debounce = (func, wait) => {
  let timeout;
  return function executedFunction(...args) {
    const later = () => {
      clearTimeout(timeout);
      func(...args);
    };
    clearTimeout(timeout);
    timeout = setTimeout(later, wait);
  };
};

// Format date
export const formatDate = (dateString) => {
  if (!dateString) return 'Unknown';
  try {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric'
    });
  } catch (error) {
    return 'Unknown';
  }
};

// Format runtime (minutes to hours and minutes)
export const formatRuntime = (minutes) => {
  if (!minutes) return 'Unknown';
  const hours = Math.floor(minutes / 60);
  const mins = minutes % 60;
  return hours > 0 ? `${hours}h ${mins}m` : `${mins}m`;
};

// Get poster URL with fallback
export const getPosterUrl = (posterPath, size = 'w500') => {
  if (!posterPath) {
    return '/placeholder-poster.jpg'; // You can add a placeholder image
  }
  return `https://image.tmdb.org/t/p/${size}${posterPath}`;
};

// Get backdrop URL
export const getBackdropUrl = (backdropPath, size = 'w1280') => {
  if (!backdropPath) {
    return '/placeholder-backdrop.jpg'; // You can add a placeholder image
  }
  return `https://image.tmdb.org/t/p/${size}${backdropPath}`;
};

// Truncate text
export const truncateText = (text, maxLength = 150) => {
  if (!text) return '';
  if (text.length <= maxLength) return text;
  return text.substring(0, maxLength) + '...';
};

// Calculate average rating
export const calculateAverageRating = (ratings) => {
  if (!ratings || ratings.length === 0) return 0;
  const sum = ratings.reduce((acc, rating) => acc + (rating.value || 0), 0);
  return Math.round((sum / ratings.length) * 10) / 10;
};

// Validate movie data
export const validateMovie = (movie) => {
  return movie && 
         movie.id && 
         movie.title && 
         typeof movie.title === 'string' &&
         movie.title.trim().length > 0;
};

// Generate star rating display
export const generateStars = (rating, maxRating = 5) => {
  const stars = [];
  const filledStars = Math.floor(rating);
  const hasHalfStar = rating % 1 >= 0.5;

  for (let i = 0; i < maxRating; i++) {
    if (i < filledStars) {
      stars.push('★'); // Filled star
    } else if (i === filledStars && hasHalfStar) {
      stars.push('☆'); // Half star (you can use a different character)
    } else {
      stars.push('☆'); // Empty star
    }
  }

  return stars.join('');
};

// Format currency
export const formatCurrency = (amount, currency = 'USD') => {
  if (!amount) return 'Unknown';
  try {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: currency
    }).format(amount);
  } catch (error) {
    return `$${amount}`;
  }
};

// Get genre name by ID
export const getGenreName = (genreId, genres) => {
  if (!genres || !genreId) return '';
  const genre = genres.find(g => g.id === genreId);
  return genre ? genre.name : '';
};

// Check if movie is in watchlist
export const isInWatchlist = (movieId, watchlist) => {
  return watchlist.some(movie => movie.id === movieId);
}; 