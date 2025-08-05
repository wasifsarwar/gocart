import React from "react";
import { Link } from "react-router-dom";

const Home = () => {
    return (
        <div className="home-page">
            <h1>Welcome to GoCart</h1>
            <p>Your one-stop-ecommerce solution</p>
            <div className="button-container">
                <Link to="/products">
                    <button className="nav-button">
                        View Products
                    </button>
                </Link>
                <Link to="/users">
                    <button className="nav-button">
                        View Users
                    </button>
                </Link>
            </div>
        </div>
    );
};


export default Home;