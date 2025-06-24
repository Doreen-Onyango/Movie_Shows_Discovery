// User ID management
export const getUserId = () => {
  let userId = localStorage.getItem('userId');
  if (!userId) {
    userId = `user_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
    localStorage.setItem('userId', userId);
  }
  return userId;
};

export const setUserId = (userId) => {
  localStorage.setItem('userId', userId);
};

// Watchlist management
export const getLocalWatchlist = () => {
  try {
    const watchlist = localStorage.getItem('watchlist');
    return watchlist ? JSON.parse(watchlist) : [];
  } catch (error) {
    console.error('Error parsing watchlist from localStorage:', error);
    return [];
  }
};

export const setLocalWatchlist = (watchlist) => {
  try {
    localStorage.setItem('watchlist', JSON.stringify(watchlist));
  } catch (error) {
    console.error('Error saving watchlist to localStorage:', error);
  }
};

export const addToLocalWatchlist = (movie) => {
  const watchlist = getLocalWatchlist();
  const existingIndex = watchlist.findIndex(item => item.id === movie.id);
  
  if (existingIndex === -1) {
    watchlist.push({ ...movie, added_at: new Date().toISOString() });
    setLocalWatchlist(watchlist);
  }
  
  return watchlist;
};

export const removeFromLocalWatchlist = (movieId) => {
  const watchlist = getLocalWatchlist();
  const filteredWatchlist = watchlist.filter(item => item.id !== movieId);
  setLocalWatchlist(filteredWatchlist);
  return filteredWatchlist;
};

export const markAsWatchedInLocalWatchlist = (movieId) => {
  const watchlist = getLocalWatchlist();
  const updated = watchlist.map(item =>
    item.id === movieId ? { ...item, watched: true } : item
  );
  setLocalWatchlist(updated);
  return updated;
};

// User preferences
export const getUserPreferences = () => {
  try {
    const preferences = localStorage.getItem('userPreferences');
    return preferences ? JSON.parse(preferences) : {
      theme: 'light',
      language: 'en',
      autoPlay: false,
      notifications: true
    };
  } catch (error) {
    console.error('Error parsing user preferences:', error);
    return {
      theme: 'light',
      language: 'en',
      autoPlay: false,
      notifications: true
    };
  }
};

export const setUserPreferences = (preferences) => {
  try {
    localStorage.setItem('userPreferences', JSON.stringify(preferences));
  } catch (error) {
    console.error('Error saving user preferences:', error);
  }
}; 