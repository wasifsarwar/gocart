import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { Toaster } from 'react-hot-toast';

import Home from './pages/home/Home';
import Products from './pages/products/Products';
import ProductDetails from './pages/products/ProductDetails';
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
            <Toaster
              position="bottom-center"
              toastOptions={{
                duration: 2000,
                style: {
                  background: '#1f2937',
                  color: '#fff',
                  fontSize: '0.875rem',
                  padding: '8px 12px',
                  borderRadius: '8px',
                  maxWidth: '400px',
                },
                success: {
                  style: {
                    background: '#1f2937',
                    borderLeft: '4px solid #10B981',
                  },
                  iconTheme: {
                    primary: '#10B981',
                    secondary: '#fff',
                  },
                },
                error: {
                  style: {
                    background: '#1f2937',
                    borderLeft: '4px solid #EF4444',
                  },
                  iconTheme: {
                    primary: '#EF4444',
                    secondary: '#fff',
                  },
                },
              }}
            />
            <Routes>
              <Route path='/' element={<Home />} />
              <Route path='/products' element={<Products />} />
              <Route path='/products/:id' element={<ProductDetails />} />
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
