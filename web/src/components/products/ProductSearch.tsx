import React, { useState, useEffect } from "react";
import './ProductSearch.css';

interface ProductSearchProps {
    onSearch: (searchTerm: string) => void;
    placeHolder?: string;
    value?: string;
}

const ProductSearch = ({ onSearch, placeHolder, value }: ProductSearchProps) => {
    const [isFocused, setIsFocused] = useState(false);
    const [isTyping, setIsTyping] = useState(false);

    const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setIsTyping(true);
        onSearch(e.target.value);
        
        // Clear typing indicator after a delay
        setTimeout(() => setIsTyping(false), 500);
    }

    const handleFocus = () => setIsFocused(true);
    const handleBlur = () => setIsFocused(false);

    return (
        <div className={`product-search ${isFocused ? 'focused' : ''} ${isTyping ? 'typing' : ''}`}>
            <div className="search-container">
                <span className="search-icon">🔍</span>
                <input
                    type="text"
                    placeholder={placeHolder}
                    value={value ?? ''}
                    onChange={handleInputChange}
                    onFocus={handleFocus}
                    onBlur={handleBlur}
                    className="search-input"
                />
                {value && (
                    <button 
                        className="clear-button"
                        onClick={() => onSearch('')}
                        type="button"
                    >
                        ✕
                    </button>
                )}
            </div>
            {isTyping && <div className="typing-indicator">Searching...</div>}
        </div>
    );
};

export default ProductSearch;