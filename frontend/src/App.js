import React, { useState, useEffect } from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { getUserId } from './utils/storage';
import Header from './components/Header';
import Home from './components/Home';
import MovieDetail from './components/MovieDetail';
import SearchResults from './components/SearchResults';
import Watchlist from './components/Watchlist';
import GenreMovies from './components/GenreMovies';
import LoadingSpinner from './components/LoadingSpinner';

function App() {
  const [isLoading, setIsLoading] = useState(true);
  const [userId, setUserId] = useState(null);

  useEffect(() => {
    // Initialize user ID
    const id = getUserId();
    setUserId(id);
    setIsLoading(false);
  }, []);

  if (isLoading) {
    return <LoadingSpinner />;
  }

  return (
    <Router>
      <div className="min-h-screen bg-gray-50">
        <Header />
        <main className="container mx-auto px-4 py-8">
          <Routes>
            <Route path="/" element={<Home />} />
            <Route path="/movie/:id" element={<MovieDetail />} />
            <Route path="/search" element={<SearchResults />} />
            <Route path="/watchlist" element={<Watchlist />} />
            <Route path="/genre/:id" element={<GenreMovies />} />
          </Routes>
        </main>
      </div>
    </Router>
  );
}

export default App;
