import React from 'react'
import { Link } from 'react-router-dom'

interface NavigationProps {
    title: string;
}

const Navigation = ({ title }: NavigationProps) => {
    return (
        <div>
            <Link to="/">
                <button className="back-button">← Back to Home</button>
            </Link>
            <h1>{title}</h1>
        </div>
    );
};

export default Navigation;