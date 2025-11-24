import React from 'react';
import { Link, useLocation } from 'react-router-dom';
import { FaShoppingCart, FaSignOutAlt } from 'react-icons/fa';
import { IconType } from 'react-icons';
import { useCart } from '../../context/CartContext';
import { useAuth } from '../../context/AuthContext';
import './Navbar.css';

// Wrapper to fix TS2786 error with React 19 types
const Icon = ({ icon: IconComponent, className }: { icon: IconType; className?: string }) => {
    const Component = IconComponent as any;
    return <Component className={className} />;
};

const Navbar = () => {
    const location = useLocation();
    const { cartCount } = useCart();
    const { isAuthenticated, logout, user } = useAuth();

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
                    {!isAuthenticated ? (
                        <>
                            <Link to="/login" className={`nav-link ${isActive('/login')}`}>
                                Login
                            </Link>
                            <Link to="/register" className={`nav-link ${isActive('/register')}`}>
                                Register
                            </Link>
                        </>
                    ) : (
                        <div className="user-menu">
                            <Link to="/orders" className={`nav-link ${isActive('/orders')}`}>
                                My Orders
                            </Link>
                            <span className="user-greeting">Hi, {user?.first_name}</span>
                            <button onClick={logout} className="logout-btn" aria-label="Logout">
                                <Icon icon={FaSignOutAlt} />
                            </button>
                        </div>
                    )}
                    <Link to="/checkout" className="cart-icon-container">
                        <Icon icon={FaShoppingCart} className="cart-icon" />
                        {cartCount > 0 && <span className="cart-badge">{cartCount}</span>}
                    </Link>
                </div>
            </div>
        </nav>
    );
};

export default Navbar;
