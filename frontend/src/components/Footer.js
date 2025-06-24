import React from "react";

const Footer = () => (
  <footer className="fixed bottom-0 left-0 w-full bg-gray-900 text-gray-300 py-4 border-t border-gray-800 z-50">
    <div className="container mx-auto px-4 flex flex-col md:flex-row items-center justify-between text-center md:text-left">
      <div className="text-sm mb-2 md:mb-0">
        © {new Date().getFullYear()} Movie & Shows Discovery. All rights reserved.
      </div>
      <div className="text-xs flex flex-col md:flex-row md:items-center gap-2 md:gap-4">
        <span>
          Powered by <a href="https://www.themoviedb.org/" target="_blank" rel="noopener noreferrer" className="underline hover:text-white">TMDB</a>
        </span>
        <span className="hidden md:inline">|</span>
        <span>
          Discover trending movies, search, and manage your watchlist with ease.
        </span>
      </div>
    </div>
  </footer>
);

export default Footer; 