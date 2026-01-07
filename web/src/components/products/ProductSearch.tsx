import React, { useState, useMemo, useRef, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { FaSearch, FaTimes, FaShoppingBag } from "react-icons/fa";
import { IconType } from "react-icons";
import Product from "../../types/product";
import './ProductSearch.css';

interface ProductSearchProps {
    onSearch: (searchTerm: string) => void;
    placeHolder?: string;
    value?: string;
    products?: Product[];
}

interface SearchSuggestion {
    product: Product;
    matchType: 'name' | 'description' | 'category';
    matchText: string;
}

// Wrapper to fix TS2786 error with React 19 types
const Icon = ({ icon: IconComponent, className }: { icon: IconType; className?: string }) => {
    const Component = IconComponent as any;
    return <Component className={className} />;
};

const ProductSearch = ({ onSearch, placeHolder, value, products = [] }: ProductSearchProps) => {
    const [isFocused, setIsFocused] = useState(false);
    const [selectedIndex, setSelectedIndex] = useState(-1);
    const [showSuggestions, setShowSuggestions] = useState(false);
    const [localValue, setLocalValue] = useState(value || '');
    const inputRef = useRef<HTMLInputElement>(null);
    const dropdownRef = useRef<HTMLDivElement>(null);
    const navigate = useNavigate();

    // Sync local value with prop value when it changes externally (like when clearing)
    useEffect(() => {
        setLocalValue(value || '');
    }, [value]);

    // Generate suggestions based on local search term (not the filtered term)
    const suggestions = useMemo((): SearchSuggestion[] => {
        if (!localValue || localValue.trim().length < 2 || products.length === 0) {
            return [];
        }

        const searchLower = localValue.toLowerCase().trim();
        const results: SearchSuggestion[] = [];
        const seen = new Set<string>();

        for (const product of products) {
            if (seen.has(product.productID)) continue;

            // Check name match
            if (product.name.toLowerCase().includes(searchLower)) {
                results.push({
                    product,
                    matchType: 'name',
                    matchText: product.name
                });
                seen.add(product.productID);
                continue;
            }

            // Check category match
            if (product.category.toLowerCase().includes(searchLower)) {
                results.push({
                    product,
                    matchType: 'category',
                    matchText: product.category
                });
                seen.add(product.productID);
                continue;
            }

            // Check description match
            if (product.description.toLowerCase().includes(searchLower)) {
                results.push({
                    product,
                    matchType: 'description',
                    matchText: product.description
                });
                seen.add(product.productID);
            }
        }

        return results.slice(0, 8); // Limit to 8 suggestions
    }, [localValue, products]);

    // Reset selected index when suggestions change
    useEffect(() => {
        setSelectedIndex(-1);
    }, [suggestions]);

    const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const newValue = e.target.value;
        setLocalValue(newValue);
        setShowSuggestions(true);
        // Don't call onSearch here - only update local state for autocomplete
    };

    const handleFocus = () => {
        setIsFocused(true);
        setShowSuggestions(true);
    };

    const handleBlur = () => {
        setIsFocused(false);
        // Delay hiding suggestions to allow click events to fire
        setTimeout(() => setShowSuggestions(false), 200);
    };

    const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
        switch (e.key) {
            case 'ArrowDown':
                if (showSuggestions && suggestions.length > 0) {
                    e.preventDefault();
                    setSelectedIndex(prev =>
                        prev < suggestions.length - 1 ? prev + 1 : prev
                    );
                }
                break;
            case 'ArrowUp':
                if (showSuggestions && suggestions.length > 0) {
                    e.preventDefault();
                    setSelectedIndex(prev => prev > 0 ? prev - 1 : -1);
                }
                break;
            case 'Enter':
                e.preventDefault();
                if (showSuggestions && suggestions.length > 0 && selectedIndex >= 0 && selectedIndex < suggestions.length) {
                    // Navigate to selected suggestion
                    const selected = suggestions[selectedIndex];
                    navigate(`/products/${selected.product.productID}`);
                    setShowSuggestions(false);
                    inputRef.current?.blur();
                } else {
                    // Apply search filter to main list
                    onSearch(localValue);
                    setShowSuggestions(false);
                    inputRef.current?.blur();
                }
                break;
            case 'Escape':
                setShowSuggestions(false);
                inputRef.current?.blur();
                break;
        }
    };

    const handleSuggestionClick = (productId: string) => {
        navigate(`/products/${productId}`);
        setShowSuggestions(false);
        setLocalValue('');
        onSearch('');
        inputRef.current?.blur();
    };

    const handleClear = () => {
        setLocalValue('');
        onSearch('');
        setShowSuggestions(false);
    };

    const highlightMatch = (text: string, searchTerm: string): React.ReactNode => {
        if (!searchTerm) return text;

        const searchLower = searchTerm.toLowerCase();
        const textLower = text.toLowerCase();
        const index = textLower.indexOf(searchLower);

        if (index === -1) return text;

        const before = text.substring(0, index);
        const match = text.substring(index, index + searchTerm.length);
        const after = text.substring(index + searchTerm.length);

        return (
            <>
                {before}
                <mark className="search-highlight">{match}</mark>
                {after}
            </>
        );
    };

    const getCategoryColor = (category: string): string => {
        const colors: { [key: string]: string } = {
            'Electronics': 'blue',
            'Clothing': 'green',
            'Books': 'yellow',
            'Home': 'purple',
            'Sports': 'red',
            'Beauty': 'pink',
            'Food': 'orange',
            'default': 'gray'
        };
        return colors[category] || colors['default'];
    };

    const showDropdown = showSuggestions && suggestions.length > 0 && isFocused;

    return (
        <div className={`product-search ${isFocused ? 'focused' : ''}`}>
            <div className="search-container">
                <span className="search-icon">
                    <Icon icon={FaSearch} />
                </span>
                <input
                    ref={inputRef}
                    type="text"
                    placeholder={placeHolder}
                    value={localValue}
                    onChange={handleInputChange}
                    onFocus={handleFocus}
                    onBlur={handleBlur}
                    onKeyDown={handleKeyDown}
                    className="search-input"
                    autoComplete="off"
                    role="combobox"
                    aria-autocomplete="list"
                    aria-expanded={showDropdown}
                    aria-controls="search-suggestions"
                    aria-activedescendant={selectedIndex >= 0 ? `suggestion-${selectedIndex}` : undefined}
                />
                {localValue && (
                    <button
                        className="clear-button"
                        onClick={handleClear}
                        type="button"
                        aria-label="Clear search"
                    >
                        <Icon icon={FaTimes} />
                    </button>
                )}
            </div>

            {showDropdown && (
                <div
                    ref={dropdownRef}
                    id="search-suggestions"
                    className="search-suggestions"
                    role="listbox"
                >
                    {suggestions.map((suggestion, index) => {
                        const colorClass = getCategoryColor(suggestion.product.category);
                        const isSelected = index === selectedIndex;

                        return (
                            <div
                                key={`${suggestion.product.productID}-${index}`}
                                id={`suggestion-${index}`}
                                className={`suggestion-item ${isSelected ? 'selected' : ''}`}
                                onMouseDown={(e) => {
                                    e.preventDefault(); // Prevent input blur
                                    handleSuggestionClick(suggestion.product.productID);
                                }}
                                role="option"
                                aria-selected={isSelected}
                            >
                                <div className={`suggestion-icon ${colorClass}`}>
                                    <Icon icon={FaShoppingBag} />
                                </div>
                                <div className="suggestion-content">
                                    <div className="suggestion-name">
                                        {highlightMatch(suggestion.product.name, localValue)}
                                    </div>
                                    <div className="suggestion-meta">
                                        <span className={`suggestion-category ${colorClass}`}>
                                            {suggestion.product.category}
                                        </span>
                                        {suggestion.matchType === 'description' && (
                                            <span className="suggestion-match-type">
                                                in description
                                            </span>
                                        )}
                                    </div>
                                </div>
                                <div className="suggestion-price">
                                    ${suggestion.product.price.toFixed(2)}
                                </div>
                            </div>
                        );
                    })}
                </div>
            )}
        </div>
    );
};

export default ProductSearch;
