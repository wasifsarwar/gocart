import React from 'react'
import { Link } from 'react-router-dom'

interface NavigationProps {
    title: string;
}

const Navigation = ({ title }: NavigationProps) => {
    return (
        <div className="navigation-container">
            <div className="nav-header">
                <Link to="/" className="back-link">
                    <button className="back-button-modern">
                        <span className="back-arrow">←</span>
                        Back
                    </button>
                </Link>
                <h1>{title}</h1>
            </div>
        </div>
    );
};

export default Navigation;