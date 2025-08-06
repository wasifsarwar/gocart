import React from "react";

interface ProductSearchProps {
    onSearch: (searchTerm: string) => void
    placeHolder?: string
}

const ProductSearch = ({ onSearch, placeHolder }: ProductSearchProps) => {
    const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        onSearch(e.target.value)
    }

    return (
        <div className="product-search">
            <input
                type="text"
                placeholder={placeHolder}
                onChange={handleInputChange}
                className="search-input"
            />
        </div>
    );
};