import React from "react";
import { Link } from "react-router-dom";
import { FaShoppingBag, FaUsers, FaUserPlus, FaArrowRight } from "react-icons/fa";
import { IconType } from "react-icons";
import './Home.css'

// Wrapper to fix TS2786 error with React 19 types
const Icon = ({ icon: IconComponent, className }: { icon: IconType; className?: string }) => {
    const Component = IconComponent as any;
    return <Component className={className} />;
};

const Home = () => {
    return (
        <div className="home-page">
            <header className="hero-section">
                <div className="hero-content">
                    <div className="hero-text">
                        <h1>
                            The Future of <br />
                            <span className="text-highlight">Go-Commerce</span>
                        </h1>
                        <p className="tagline">
                            Experience the speed and reliability of a modern e-commerce platform built with Go and React.
                        </p>
                        <div className="hero-actions">
                            <Link to="/products" className="btn-primary">
                                Browse Products <Icon icon={FaArrowRight} />
                            </Link>
                            <Link to="/register" className="btn-secondary">
                                Sign Up Now
                            </Link>
                        </div>
                    </div>
                    <div className="hero-image">
                        <img src="/assets/gopher_beer.gif" alt="Gocart Gopher" className="gopher-mascot" />
                        <div className="blob-bg"></div>
                    </div>
                </div>
            </header>

            <section className="features-section">
                <h2 className="section-title">Everything you need</h2>
                <div className="features-grid">
                    <Link to="/products" className="feature-card">
                        <div className="icon-wrapper blue">
                            <Icon icon={FaShoppingBag} />
                        </div>
                        <h3>Product Catalog</h3>
                        <p>Browse our extensive collection of tech and lifestyle products.</p>
                        <span className="card-link">Shop Now &rarr;</span>
                    </Link>

                    <Link to="/users" className="feature-card">
                        <div className="icon-wrapper purple">
                            <Icon icon={FaUsers} />
                        </div>
                        <h3>User Management</h3>
                        <p>Administer user accounts, profiles, and permissions seamlessly.</p>
                        <span className="card-link">Manage Users &rarr;</span>
                    </Link>

                    <Link to="/register" className="feature-card">
                        <div className="icon-wrapper green">
                            <Icon icon={FaUserPlus} />
                        </div>
                        <h3>Easy Registration</h3>
                        <p>Onboard new customers in seconds with our streamlined flow.</p>
                        <span className="card-link">Register &rarr;</span>
                    </Link>
                </div>
            </section>
        </div>
    );
};

export default Home;
