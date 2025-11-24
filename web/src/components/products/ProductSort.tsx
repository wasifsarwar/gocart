import React from 'react';
import './ProductSort.css';

interface SortOption {
    value: string;
    label: string;
}

interface ProductSortProps {
    onSort: (sortValue: string) => void
    currentSort: string
}

const sortOptions: SortOption[] = [
    { value: 'name-asc', label: 'Name (A-Z)' },
    { value: 'name-desc', label: 'Name (Z-A)' },
    { value: 'price-asc', label: 'Price (Low to High)' },
    { value: 'price-desc', label: 'Price (High to Low)' },
    { value: 'category-asc', label: 'Category (A-Z)' },
];

const ProductSort = ({ onSort, currentSort }: ProductSortProps) => {
    return (
        <div className='product-sort'>
            <label htmlFor="sort-select">Sort by</label>
            <div className="select-wrapper">
                <select
                    id="sort-select"
                    value={currentSort}
                    onChange={(e) => onSort(e.target.value)}
                    className="sort-select"
                >
                    {sortOptions.map(option => (
                        <option key={option.value} value={option.value}>{option.label}</option>
                    ))}
                </select>
            </div>
        </div>
    );
};

export default ProductSort;
