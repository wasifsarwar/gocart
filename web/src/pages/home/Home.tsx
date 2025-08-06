import React from "react";
import { Link } from "react-router-dom";
import './Home.css'
const Home = () => {
    return (
        <div className="home-page">
            <header className="hero-section">
                <div className="brand-container">
                    <img src="/assets/gopher_beer.gif" alt="Gocart Gopher" className="gopher-logo" />
                    <h1>GoCart</h1>
                </div>
                <p className="tagline">Your Go-powered e-commerce solution</p>
                <p className="subtitle">Fast, reliable, and built for modern e-commerce</p>
            </header>

            <section className="features-grid">
                <Link to="/products" className="feature-card">
                    <span className="card-icon">🛒</span>
                    <h3>Browse Products</h3>
                    <p>Explore our product catalog</p>
                </Link>

                <Link to="/users" className="feature-card">
                    <span className="card-icon">👥</span>
                    <h3>Manage Users</h3>
                    <p>View and manage user accounts</p>
                </Link>

                <Link to="/register" className="feature-card">
                    <span className="card-icon">➕</span>
                    <h3>Add New User</h3>
                    <p>Register new customers</p>
                </Link>
            </section>
        </div>
    );
};


export default Home;