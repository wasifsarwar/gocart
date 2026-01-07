import { useState, useMemo, useEffect } from 'react';

import ProductList from "../../components/products/ProductList";
import ProductSearch from '../../components/products/ProductSearch';
import ProductSort from "../../components/products/ProductSort";
import PriceRangeFilter from '../../components/products/PriceRangeFilter';
import ActiveFilterTags from '../../components/products/ActiveFilterTags';
import useProducts from "../../hooks/useProducts";
import { useFavorites } from '../../context/FavoritesContext';
import useRecentlyViewedProducts from '../../hooks/useRecentlyViewedProducts';
import ProductQuickViewModal from '../../components/products/ProductQuickViewModal';
import Product from '../../types/product';

import './Products.css'


const Products = () => {
    const { products, loading, error, refetch } = useProducts();
    const { favoriteIds } = useFavorites();
    const { recentlyViewed, clearRecentlyViewed } = useRecentlyViewedProducts();
    const [searchTerm, setSearchTerm] = useState('');
    const [sortBy, setSortBy] = useState('name-asc'); //default sort state value
    const [pageSize, setPageSize] = useState<5 | 10 | 15 | 25 | 50>(15);
    const [currentPage, setCurrentPage] = useState(1);
    const [activeTab, setActiveTab] = useState<'all' | 'favorites' | 'recent'>('all');
    const [quickViewProduct, setQuickViewProduct] = useState<Product | null>(null);

    const tabProducts = useMemo(() => {
        if (activeTab === 'all') return products;
        if (activeTab === 'recent') return recentlyViewed;
        if (products.length === 0 || favoriteIds.length === 0) return [];
        const byId = new Map(products.map((p) => [p.productID, p] as const));
        // Preserve favorites ordering (most recently favorited first)
        return favoriteIds.map((id) => byId.get(id)).filter(Boolean) as typeof products;
    }, [activeTab, products, favoriteIds, recentlyViewed]);

    // Get unique categories for filtering
    const categories = useMemo(() => {
        return Array.from(new Set(tabProducts.map(product => product.category))).sort();
    }, [tabProducts]);

    const [selectedCategory, setSelectedCategory] = useState<string>('');

    // Calculate min and max prices from products
    const { minPrice, maxPrice } = useMemo(() => {
        if (tabProducts.length === 0) return { minPrice: 0, maxPrice: 1000 };

        // Use reduce instead of spread operator to avoid stack overflow with large arrays
        let min = tabProducts[0].price;
        let max = tabProducts[0].price;

        for (let i = 1; i < tabProducts.length; i++) {
            const price = tabProducts[i].price;
            if (price < min) min = price;
            if (price > max) max = price;
        }

        return {
            minPrice: Math.floor(min),
            maxPrice: Math.ceil(max)
        };
    }, [tabProducts]);

    const [priceRange, setPriceRange] = useState<{ min: number; max: number }>(() => ({
        min: 0,
        max: 1000
    }));
    const [isPriceRangeDirty, setIsPriceRangeDirty] = useState(false);

    // Update price range when products change
    useEffect(() => {
        // If the user hasn't adjusted the slider, always track the full computed range.
        // This avoids a one-render mismatch (and flicker) when product prices exceed the initial defaults.
        if (!isPriceRangeDirty) {
            setPriceRange({ min: minPrice, max: maxPrice });
            return;
        }

        // If the user *has* adjusted it, clamp the stored range to the latest bounds.
        setPriceRange((prev) => {
            const nextMin = Math.max(minPrice, Math.min(prev.min, maxPrice));
            const nextMax = Math.min(maxPrice, Math.max(prev.max, minPrice));
            return { min: nextMin, max: nextMax };
        });
    }, [minPrice, maxPrice, isPriceRangeDirty]);

    const filteredProducts = useMemo(() => {
        let filtered = tabProducts.filter(
            product => product.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
                product.description.toLowerCase().includes(searchTerm.toLowerCase()) ||
                product.category.toLowerCase().includes(searchTerm.toLowerCase())
        );

        // Apply category filter
        if (selectedCategory) {
            filtered = filtered.filter(product => product.category === selectedCategory);
        }

        // Apply price range filter only if user has adjusted it.
        // (Do not infer "active" by comparing to min/max; priceRange state is synced in an effect.)
        if (isPriceRangeDirty) {
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
    }, [tabProducts, sortBy, searchTerm, selectedCategory, priceRange, isPriceRangeDirty]);

    // Reset pagination when filters change
    useEffect(() => {
        setCurrentPage(1);
    }, [searchTerm, sortBy, selectedCategory, isPriceRangeDirty, priceRange, pageSize, activeTab]);

    // Close modal when switching tabs
    useEffect(() => {
        setQuickViewProduct(null);
    }, [activeTab]);

    const handleClear = () => {
        setSearchTerm('');
        setSortBy('name-asc');
        setSelectedCategory('');
        setIsPriceRangeDirty(false);
        setPriceRange({ min: minPrice, max: maxPrice });
        setCurrentPage(1);
    }

    const handlePriceChange = (min: number, max: number) => {
        setIsPriceRangeDirty(true);
        setPriceRange({ min, max });
    }

    // Individual filter removal handlers
    const handleRemoveSearch = () => setSearchTerm('');
    const handleRemoveCategory = () => setSelectedCategory('');
    const handleRemovePriceRange = () => {
        setIsPriceRangeDirty(false);
        setPriceRange({ min: minPrice, max: maxPrice });
    };
    const handleRemoveSort = () => setSortBy('name-asc');

    const effectivePriceRange = isPriceRangeDirty ? priceRange : { min: minPrice, max: maxPrice };

    const totalPages = Math.max(1, Math.ceil(filteredProducts.length / pageSize));
    const safeCurrentPage = Math.min(currentPage, totalPages);
    const paginatedProducts = useMemo(() => {
        const start = (safeCurrentPage - 1) * pageSize;
        return filteredProducts.slice(start, start + pageSize);
    }, [filteredProducts, safeCurrentPage, pageSize]);

    return (
        <div className="products-page page-container">
            <div className="products-header">
                <div className="header-content">
                    <h1>Products</h1>
                    <p className="header-subtitle">Browse our complete catalog of quality items</p>
                </div>
            </div>

            <section className="products-content">
                <div className="products-tabs" role="tablist" aria-label="Products tabs">
                    <button
                        type="button"
                        role="tab"
                        aria-selected={activeTab === 'all'}
                        className={`products-tab ${activeTab === 'all' ? 'active' : ''}`}
                        onClick={() => setActiveTab('all')}
                    >
                        All products
                    </button>
                    <button
                        type="button"
                        role="tab"
                        aria-selected={activeTab === 'favorites'}
                        className={`products-tab ${activeTab === 'favorites' ? 'active' : ''}`}
                        onClick={() => setActiveTab('favorites')}
                    >
                        Favorites
                    </button>
                    <button
                        type="button"
                        role="tab"
                        aria-selected={activeTab === 'recent'}
                        className={`products-tab ${activeTab === 'recent' ? 'active' : ''}`}
                        onClick={() => setActiveTab('recent')}
                    >
                        Recently viewed
                    </button>
                </div>

                <div className="products-controls">
                    {activeTab === 'recent' && (
                        <div className="tab-toolbar">
                            <div className="tab-toolbar-title">Recently viewed</div>
                            <button
                                type="button"
                                className="tab-toolbar-action"
                                onClick={clearRecentlyViewed}
                                disabled={recentlyViewed.length === 0}
                            >
                                Clear
                            </button>
                        </div>
                    )}

                    <div className="search-sort-row">
                        <ProductSearch
                            onSearch={setSearchTerm}
                            value={searchTerm}
                            placeHolder="Search products..."
                            products={tabProducts}
                        />
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

                    {tabProducts.length > 0 && (
                        <PriceRangeFilter
                            minPrice={minPrice}
                            maxPrice={maxPrice}
                            currentMin={effectivePriceRange.min}
                            currentMax={effectivePriceRange.max}
                            onPriceChange={handlePriceChange}
                        />
                    )}
                </div>

                <ActiveFilterTags
                    searchTerm={searchTerm}
                    selectedCategory={selectedCategory}
                    isPriceFilterActive={isPriceRangeDirty}
                    priceRange={effectivePriceRange}
                    minPrice={minPrice}
                    maxPrice={maxPrice}
                    sortBy={sortBy}
                    onRemoveSearch={handleRemoveSearch}
                    onRemoveCategory={handleRemoveCategory}
                    onRemovePriceRange={handleRemovePriceRange}
                    onRemoveSort={handleRemoveSort}
                    onClearAll={handleClear}
                />

                {error && (
                    <div role="alert" className="alert alert-error">
                        <span>{error}</span>
                        <button onClick={refetch}>Retry</button>
                    </div>
                )}

                {activeTab === 'favorites' && !loading && tabProducts.length === 0 && (
                    <div className="empty-state favorites-empty">
                        <h3>No favorites yet</h3>
                        <p>Tap the heart icon on a product to save it here.</p>
                    </div>
                )}

                {activeTab === 'recent' && !loading && tabProducts.length === 0 && (
                    <div className="empty-state recent-empty">
                        <h3>No recently viewed products</h3>
                        <p>Open a product to start building your recently viewed list.</p>
                    </div>
                )}

                {!loading && filteredProducts.length > 0 && (
                    <div className="pager-controls pager-controls-top" aria-label="Product pagination">
                        <div className="pager">
                            <button
                                type="button"
                                className="pager-btn"
                                onClick={() => setCurrentPage(1)}
                                disabled={safeCurrentPage === 1}
                                aria-label="First page"
                            >
                                First
                            </button>
                            <button
                                type="button"
                                className="pager-btn"
                                onClick={() => setCurrentPage((p) => Math.max(1, p - 1))}
                                disabled={safeCurrentPage === 1}
                                aria-label="Previous page"
                            >
                                Prev
                            </button>

                            <span className="pager-status" aria-live="polite">
                                Page {safeCurrentPage} of {totalPages}
                            </span>

                            <button
                                type="button"
                                className="pager-btn"
                                onClick={() => setCurrentPage((p) => Math.min(totalPages, p + 1))}
                                disabled={safeCurrentPage === totalPages}
                                aria-label="Next page"
                            >
                                Next
                            </button>
                            <button
                                type="button"
                                className="pager-btn"
                                onClick={() => setCurrentPage(totalPages)}
                                disabled={safeCurrentPage === totalPages}
                                aria-label="Last page"
                            >
                                Last
                            </button>
                        </div>
                    </div>
                )}

                <div className="product-list-container">
                    <ProductList products={paginatedProducts} loading={loading} onQuickView={setQuickViewProduct} />
                </div>

                {!loading && filteredProducts.length > 0 && (
                    <div className="page-size-controls page-size-controls-bottom" aria-label="Items per page">
                        <div className="page-size page-size-centered">
                            <label htmlFor="pageSize" className="page-size-label">Show</label>
                            <select
                                id="pageSize"
                                value={pageSize}
                                onChange={(e) => setPageSize((Number(e.target.value) as 5 | 10 | 15 | 25 | 50) || 15)}
                                className="page-size-select"
                                aria-label="Items per page"
                            >
                                <option value={5}>5</option>
                                <option value={10}>10</option>
                                <option value={15}>15</option>
                                <option value={25}>25</option>
                                <option value={50}>50</option>
                            </select>
                            <span className="page-size-label">per page</span>
                        </div>
                    </div>
                )}
            </section>

            {quickViewProduct && (
                <ProductQuickViewModal
                    product={quickViewProduct}
                    onClose={() => setQuickViewProduct(null)}
                />
            )}

        </div >
    );
};

export default Products;
