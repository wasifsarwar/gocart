import React from 'react';
import './App.css';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';

import ProductList from './components/ProductList';
import useProducts from './hooks/useProducts';
import Home from './pages/Home';
import Products from './pages/Products';

function App() {

  const { products, loading, error } = useProducts();
  return (
    <Router>
      <div className="App">
        <Routes>
          <Route path='/' element={<Home />} />
          <Route path='/products' element={<Products />} />
        </Routes>
      </div>
    </Router >
  );
}

export default App;
