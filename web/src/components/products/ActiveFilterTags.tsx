import React from 'react';
import { FaTimes } from 'react-icons/fa';
import { IconType } from 'react-icons';
import './ActiveFilterTags.css';

// Wrapper to fix TS2786 error with React 19 types
const Icon = ({ icon: IconComponent, className }: { icon: IconType; className?: string }) => {
    const Component = IconComponent as any;
    return <Component className={className} />;
};

interface ActiveFilterTagsProps {
    searchTerm: string;
    selectedCategory: string;
    priceRange: { min: number; max: number };
    minPrice: number;
    maxPrice: number;
    sortBy: string;
    onRemoveSearch: () => void;
    onRemoveCategory: () => void;
    onRemovePriceRange: () => void;
    onRemoveSort: () => void;
    onClearAll: () => void;
}

const ActiveFilterTags: React.FC<ActiveFilterTagsProps> = ({
    searchTerm,
    selectedCategory,
    priceRange,
    minPrice,
    maxPrice,
    sortBy,
    onRemoveSearch,
    onRemoveCategory,
    onRemovePriceRange,
    onRemoveSort,
    onClearAll
}) => {
    const usdFormatter = new Intl.NumberFormat('en-us', {
        style: 'currency',
        currency: 'USD',
        maximumFractionDigits: 0
    });

    const hasSearchFilter = searchTerm !== '';
    const hasCategoryFilter = selectedCategory !== '';
    const hasPriceFilter = priceRange.min !== minPrice || priceRange.max !== maxPrice;
    const hasSortFilter = sortBy !== 'name-asc';

    const hasAnyFilter = hasSearchFilter || hasCategoryFilter || hasPriceFilter || hasSortFilter;

    if (!hasAnyFilter) {
        return null;
    }

    const getSortLabel = (sort: string) => {
        const sortLabels: { [key: string]: string } = {
            'name-asc': 'Name (A-Z)',
            'name-desc': 'Name (Z-A)',
            'price-asc': 'Price (Low to High)',
            'price-desc': 'Price (High to Low)',
            'category-asc': 'Category (A-Z)'
        };
        return sortLabels[sort] || sort;
    };

    return (
        <div className="active-filter-tags">
            <div className="filter-tags-header">
                <span className="filter-tags-label">Active Filters:</span>
                <button onClick={onClearAll} className="clear-all-btn">
                    Clear All
                </button>
            </div>
            <div className="filter-tags-container">
                {hasSearchFilter && (
                    <div className="filter-tag filter-tag-search">
                        <span className="filter-tag-label">Search:</span>
                        <span className="filter-tag-value">"{searchTerm}"</span>
                        <button onClick={onRemoveSearch} className="filter-tag-remove" aria-label="Remove search filter">
                            <Icon icon={FaTimes} />
                        </button>
                    </div>
                )}

                {hasCategoryFilter && (
                    <div className="filter-tag filter-tag-category">
                        <span className="filter-tag-label">Category:</span>
                        <span className="filter-tag-value">{selectedCategory}</span>
                        <button onClick={onRemoveCategory} className="filter-tag-remove" aria-label="Remove category filter">
                            <Icon icon={FaTimes} />
                        </button>
                    </div>
                )}

                {hasPriceFilter && (
                    <div className="filter-tag filter-tag-price">
                        <span className="filter-tag-label">Price:</span>
                        <span className="filter-tag-value">
                            {usdFormatter.format(priceRange.min)} - {usdFormatter.format(priceRange.max)}
                        </span>
                        <button onClick={onRemovePriceRange} className="filter-tag-remove" aria-label="Remove price filter">
                            <Icon icon={FaTimes} />
                        </button>
                    </div>
                )}

                {hasSortFilter && (
                    <div className="filter-tag filter-tag-sort">
                        <span className="filter-tag-label">Sort:</span>
                        <span className="filter-tag-value">{getSortLabel(sortBy)}</span>
                        <button onClick={onRemoveSort} className="filter-tag-remove" aria-label="Remove sort filter">
                            <Icon icon={FaTimes} />
                        </button>
                    </div>
                )}
            </div>
        </div>
    );
};

export default ActiveFilterTags;
