import React from "react";
import './ProductSearch.css';

interface ProductSearchProps {
    onSearch: (searchTerm: string) => void;
    placeHolder?: string;
    value?: string;
}

const ProductSearch = ({ onSearch, placeHolder, value }: ProductSearchProps) => {
    const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        onSearch(e.target.value)
    }

    return (
        <div className="product-search">
            <input
                type="text"
                placeholder={placeHolder}
                value={value ?? ''}
                onChange={handleInputChange}
                className="search-input"
            />
        </div>
    );
};

export default ProductSearch;