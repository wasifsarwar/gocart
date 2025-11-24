import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';

import Home from './pages/home/Home';
import Products from './pages/products/Products';
import Users from './pages/users/Users';
import UserRegistration from './pages/users/UserRegistration'
import Checkout from './pages/checkout/Checkout';
import Login from './pages/login/Login';
import OrderHistory from './pages/orders/OrderHistory';
import Navbar from './components/layout/Navbar';
import { CartProvider } from './context/CartContext';
import { AuthProvider } from './context/AuthContext';

import './styles/globals.css';

function App() {
  return (
    <AuthProvider>
      <CartProvider>
        <Router>
          <Navbar />
          <div className="App">
            <Routes>
              <Route path='/' element={<Home />} />
              <Route path='/products' element={<Products />} />
              <Route path='/users' element={<Users />} />
              <Route path='/register' element={<UserRegistration />} />
              <Route path='/login' element={<Login />} />
              <Route path='/orders' element={<OrderHistory />} />
              <Route path='/checkout' element={<Checkout />} />
            </Routes>
          </div>
        </Router >
      </CartProvider>
    </AuthProvider>
  );
}

export default App;
