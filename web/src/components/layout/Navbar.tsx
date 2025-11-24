import React from 'react';
import { Link, useLocation } from 'react-router-dom';
import './Navbar.css';

const Navbar = () => {
    const location = useLocation();

    const isActive = (path: string) => {
        return location.pathname === path ? 'active' : '';
    };

    return (
        <nav className="navbar">
            <div className="navbar-container">
                <Link to="/" className="navbar-brand">
                    <img src="/assets/gopher_beer.gif" alt="GoCart" className="navbar-logo" />
                    <span>GoCart</span>
                </Link>
                <div className="navbar-links">
                    <Link to="/" className={`nav-link ${isActive('/')}`}>
                        Home
                    </Link>
                    <Link to="/products" className={`nav-link ${isActive('/products')}`}>
                        Products
                    </Link>
                    <Link to="/users" className={`nav-link ${isActive('/users')}`}>
                        Users
                    </Link>
                    <Link to="/register" className={`nav-link ${isActive('/register')}`}>
                        Register
                    </Link>
                </div>
            </div>
        </nav>
    );
};

export default Navbar;
