import { useState, useMemo } from 'react';

import ProductList from "../../components/products/ProductList";
import ProductSearch from '../../components/products/ProductSearch';
import ProductSort from "../../components/products/ProductSort";
import useProducts from "../../hooks/useProducts";
import Navigation from "../../components/navigation/Navigation";

import './Products.css'


const Products = () => {
    const { products, loading, error } = useProducts();
    const [searchTerm, setSearchTerm] = useState('');
    const [sortBy, setSortBy] = useState('name-asc'); //default sort state value

    const filteredProducts = useMemo(() => {
        let filtered = products.filter(
            product => product.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
                product.description.toLowerCase().includes(searchTerm.toLowerCase()) ||
                product.category.toLowerCase().includes(searchTerm.toLowerCase())
        );

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
    }, [products, sortBy, searchTerm]);

    return (
        <div className="products-page">
            <Navigation title="GoCart Products" />
            <div className="products-controls" >
                <ProductSearch onSearch={setSearchTerm} placeHolder="🔍 Search products..." />
                <ProductSort onSort={setSortBy} currentSort={sortBy} />
            </div>
            {loading && <p>Loading Products</p>}
            {error && <p style={{ color: 'red' }} >{error}</p>}
            {!loading && !error && <ProductList products={filteredProducts} />}
        </div >
    );
};

export default Products;