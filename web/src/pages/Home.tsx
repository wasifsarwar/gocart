import React from "react";
import useProducts from "../hooks/useProducts";
import { Link } from "react-router-dom";

const Home = () => {
    const { products, loading, error } = useProducts();
    return (
        <div className="home-page">
            <h1>Welcome to GoCart</h1>
            <p>Your one-stop-ecommerce solution</p>
            <Link to="/products">
                <button className="view-products-btn">
                    View Products
                </button>
            </Link>
        </div>
    );
};


export default Home;