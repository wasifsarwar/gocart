import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';

import useProducts from './hooks/useProducts';
import Home from './pages/Home';
import Products from './pages/Products';
import Users from './pages/Users';
import UserRegistration from './pages/UserRegistration'

import './styles/App.css';
import './styles/globals.css';



function App() {

  const { products, loading, error } = useProducts();
  return (
    <Router>
      <div className="App">
        <Routes>
          <Route path='/' element={<Home />} />
          <Route path='/products' element={<Products />} />
          <Route path='/users' element={<Users />} />
          <Route path='/register' element={<UserRegistration />} />
        </Routes>
      </div>
    </Router >
  );
}

export default App;
