import React from 'react';
import './App.css';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';

import ProductList from './components/products/ProductList';
import useProducts from './hooks/useProducts';
import Home from './pages/Home';
import Products from './pages/Products';
import Users from './pages/Users';

function App() {

  const { products, loading, error } = useProducts();
  return (
    <Router>
      <div className="App">
        <Routes>
          <Route path='/' element={<Home />} />
          <Route path='/products' element={<Products />} />
          <Route path='/users' element={<Users />} />
        </Routes>
      </div>
    </Router >
  );
}

export default App;
