import { useState, useMemo, useEffect } from 'react';

import ProductList from "../../components/products/ProductList";
import ProductSearch from '../../components/products/ProductSearch';
import ProductSort from "../../components/products/ProductSort";
import PriceRangeFilter from '../../components/products/PriceRangeFilter';
import ActiveFilterTags from '../../components/products/ActiveFilterTags';
import useProducts from "../../hooks/useProducts";

import './Products.css'


const Products = () => {
    const { products, loading, error, refetch } = useProducts();
    const [searchTerm, setSearchTerm] = useState('');
    const [sortBy, setSortBy] = useState('name-asc'); //default sort state value

    // Get unique categories for filtering
    const categories = useMemo(() => {
        return Array.from(new Set(products.map(product => product.category))).sort();
    }, [products]);

    const [selectedCategory, setSelectedCategory] = useState<string>('');

    // Calculate min and max prices from products
    const { minPrice, maxPrice } = useMemo(() => {
        if (products.length === 0) return { minPrice: 0, maxPrice: 1000 };

        // Use reduce instead of spread operator to avoid stack overflow with large arrays
        let min = products[0].price;
        let max = products[0].price;

        for (let i = 1; i < products.length; i++) {
            const price = products[i].price;
            if (price < min) min = price;
            if (price > max) max = price;
        }

        return {
            minPrice: Math.floor(min),
            maxPrice: Math.ceil(max)
        };
    }, [products]);

    const [priceRange, setPriceRange] = useState<{ min: number; max: number }>(() => ({
        min: 0,
        max: 1000
    }));

    // Update price range when products change
    useEffect(() => {
        setPriceRange({ min: minPrice, max: maxPrice });
    }, [minPrice, maxPrice]);

    const filteredProducts = useMemo(() => {
        let filtered = products.filter(
            product => product.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
                product.description.toLowerCase().includes(searchTerm.toLowerCase()) ||
                product.category.toLowerCase().includes(searchTerm.toLowerCase())
        );

        // Apply category filter
        if (selectedCategory) {
            filtered = filtered.filter(product => product.category === selectedCategory);
        }

        // Apply price range filter only if user has adjusted it
        // Skip filtering if priceRange matches the full range (not user-adjusted)
        const isPriceFilterActive = priceRange.min !== minPrice || priceRange.max !== maxPrice;
        if (isPriceFilterActive) {
            filtered = filtered.filter(product =>
                product.price >= priceRange.min && product.price <= priceRange.max
            );
        }

        return filtered.sort((a, b) => {
            switch (sortBy) {
                case 'name-asc':
                    return a.name.localeCompare(b.name);
                case 'name-desc':
                    return b.name.localeCompare(a.name);
                case 'price-asc':
                    return a.price - b.price;
                case 'price-desc':
                    return b.price - a.price;
                case 'category-asc':
                    return a.category.localeCompare(b.category)
                default:
                    return 0;
            }
        });
    }, [products, sortBy, searchTerm, selectedCategory, priceRange, minPrice, maxPrice]);

    const handleClear = () => {
        setSearchTerm('');
        setSortBy('name-asc');
        setSelectedCategory('');
        setPriceRange({ min: minPrice, max: maxPrice });
    }

    const handlePriceChange = (min: number, max: number) => {
        setPriceRange({ min, max });
    }

    // Individual filter removal handlers
    const handleRemoveSearch = () => setSearchTerm('');
    const handleRemoveCategory = () => setSelectedCategory('');
    const handleRemovePriceRange = () => setPriceRange({ min: minPrice, max: maxPrice });
    const handleRemoveSort = () => setSortBy('name-asc');

    return (
        <div className="products-page page-container">
            <div className="products-header">
                <div className="header-content">
                    <h1>Products</h1>
                    <p className="header-subtitle">Browse our complete catalog of quality items</p>
                </div>
            </div>

            <section className="products-content">
                <div className="products-controls">
                    <div className="search-sort-row">
                        <ProductSearch onSearch={setSearchTerm} value={searchTerm} placeHolder="Search products..." />
                        <ProductSort onSort={setSortBy} currentSort={sortBy} />
                    </div>

                    {categories.length > 0 && (
                        <div className="category-filters">
                            <span className="filter-label">Filter:</span>
                            <div className="category-badges">
                                <button
                                    className={`category-filter-btn ${!selectedCategory ? 'active' : ''}`}
                                    onClick={() => setSelectedCategory('')}
                                >
                                    All
                                </button>
                                {categories.map(category => (
                                    <button
                                        key={category}
                                        className={`category-filter-btn ${selectedCategory === category ? 'active' : ''}`}
                                        onClick={() => setSelectedCategory(category)}
                                    >
                                        {category}
                                    </button>
                                ))}
                            </div>
                        </div>
                    )}

                    {products.length > 0 && (
                        <PriceRangeFilter
                            minPrice={minPrice}
                            maxPrice={maxPrice}
                            currentMin={priceRange.min}
                            currentMax={priceRange.max}
                            onPriceChange={handlePriceChange}
                        />
                    )}
                </div>

                <ActiveFilterTags
                    searchTerm={searchTerm}
                    selectedCategory={selectedCategory}
                    priceRange={priceRange}
                    minPrice={minPrice}
                    maxPrice={maxPrice}
                    sortBy={sortBy}
                    onRemoveSearch={handleRemoveSearch}
                    onRemoveCategory={handleRemoveCategory}
                    onRemovePriceRange={handleRemovePriceRange}
                    onRemoveSort={handleRemoveSort}
                    onClearAll={handleClear}
                />

                <div className="results-meta">
                    <span aria-live="polite">{filteredProducts.length} results found</span>
                </div>

                {error && (
                    <div role="alert" className="alert alert-error">
                        <span>{error}</span>
                        <button onClick={refetch}>Retry</button>
                    </div>
                )}

                <div className="product-list-container">
                    <ProductList products={filteredProducts} loading={loading} />
                </div>
            </section>

        </div >
    );
};

export default Products;
