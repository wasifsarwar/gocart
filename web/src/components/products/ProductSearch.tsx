import React, { useState } from "react";
import { FaSearch, FaTimes } from "react-icons/fa";
import { IconType } from "react-icons";
import './ProductSearch.css';

interface ProductSearchProps {
    onSearch: (searchTerm: string) => void;
    placeHolder?: string;
    value?: string;
}

// Wrapper to fix TS2786 error with React 19 types
const Icon = ({ icon: IconComponent, className }: { icon: IconType; className?: string }) => {
    const Component = IconComponent as any;
    return <Component className={className} />;
};

const ProductSearch = ({ onSearch, placeHolder, value }: ProductSearchProps) => {
    const [isFocused, setIsFocused] = useState(false);

    const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        onSearch(e.target.value);
    }

    const handleFocus = () => setIsFocused(true);
    const handleBlur = () => setIsFocused(false);

    return (
        <div className={`product-search ${isFocused ? 'focused' : ''}`}>
            <div className="search-container">
                <span className="search-icon">
                    <Icon icon={FaSearch} />
                </span>
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
                        aria-label="Clear search"
                    >
                        <Icon icon={FaTimes} />
                    </button>
                )}
            </div>
        </div>
    );
};

export default ProductSearch;
