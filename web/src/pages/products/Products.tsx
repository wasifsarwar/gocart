import { useState, useMemo } from 'react';

import ProductList from "../../components/products/ProductList";
import ProductSearch from '../../components/products/ProductSearch';
import ProductSort from "../../components/products/ProductSort";
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
    }, [products, sortBy, searchTerm, selectedCategory]);

    const handleClear = () => {
        setSearchTerm('');
        setSortBy('name-asc');
        setSelectedCategory('');
    }

    const getCategoryColor = (category: string) => {
        const colors: { [key: string]: string } = {
            'Electronics': '#3b82f6',
            'Clothing': '#10b981',
            'Books': '#f59e0b',
            'Home': '#8b5cf6',
            'Sports': '#ef4444',
            'Beauty': '#ec4899',
            'Food': '#f97316',
            'default': '#6b7280'
        };
        return colors[category] || colors['default'];
    };

    return (
        <div className="products-page page-container">
            <header className="hero-section">
                <div className="brand-container">
                    <img src="/assets/gopher_beer.gif" alt="GoCart Gopher" className="gopher-logo" />
                    <h1>GoCart Products</h1>
                </div>
                <p className="tagline">Browse our complete product catalog</p>
                <p className="subtitle">Search, filter, and discover products</p>
            </header>
            <section className="products-content">
                <div className="products-controls">
                    <ProductSearch onSearch={setSearchTerm} value={searchTerm} placeHolder="Search products" />
                    <ProductSort onSort={setSortBy} currentSort={sortBy} />
                </div>

                {categories.length > 0 && (
                    <div className="category-filters">
                        <h4>Filter by Category:</h4>
                        <div className="category-badges">
                            <button
                                className={`category-filter-btn ${!selectedCategory ? 'active' : ''}`}
                                onClick={() => setSelectedCategory('')}
                            >
                                All Categories
                            </button>
                            {categories.map(category => (
                                <button
                                    key={category}
                                    className={`category-filter-btn ${selectedCategory === category ? 'active' : ''}`}
                                    style={{
                                        backgroundColor: selectedCategory === category ? getCategoryColor(category) : 'transparent',
                                        borderColor: getCategoryColor(category),
                                        color: selectedCategory === category ? 'white' : getCategoryColor(category)
                                    }}
                                    onClick={() => setSelectedCategory(category)}
                                >
                                    {category}
                                </button>
                            ))}
                        </div>
                    </div>
                )}

                <div className="results-meta">
                    <span aria-live="polite">{filteredProducts.length} results</span>
                    {(searchTerm !== '' || sortBy !== 'name-asc' || selectedCategory !== '') && (
                        <button onClick={handleClear}>
                            Clear All Filters
                        </button>
                    )}
                </div>

                {error && (
                    <div role="alert" className="alert alert-error">
                        <span>{error}</span>
                        <button onClick={refetch}>Retry</button>
                    </div>
                )}

                <div className="table-container">
                    <ProductList products={filteredProducts} loading={loading} />
                </div>
            </section>

        </div >
    );
};

export default Products;