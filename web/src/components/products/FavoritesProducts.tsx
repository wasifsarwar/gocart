import React, { useMemo } from 'react';
import Product from '../../types/product';
import ProductCard from './ProductCard';
import { useFavorites } from '../../context/FavoritesContext';
import './FavoritesProducts.css';

interface FavoritesProductsProps {
    products: Product[];
    limit?: number;
}

const FavoritesProducts: React.FC<FavoritesProductsProps> = ({ products, limit = 8 }) => {
    const { favoriteIds, clearFavorites } = useFavorites();

    const favoriteProducts = useMemo(() => {
        if (favoriteIds.length === 0 || products.length === 0) return [];
        const byId = new Map(products.map((p) => [p.productID, p] as const));
        return favoriteIds.map((id) => byId.get(id)).filter(Boolean) as Product[];
    }, [favoriteIds, products]);

    const productsToShow = useMemo(() => favoriteProducts.slice(0, limit), [favoriteProducts, limit]);

    if (productsToShow.length === 0) return null;

    return (
        <section className="favorites-section" aria-label="Favorite products">
            <div className="favorites-header">
                <h2 className="favorites-title">Favorites</h2>
                <button className="favorites-clear" onClick={clearFavorites} type="button">
                    Clear
                </button>
            </div>
            <div className="product-grid favorites-grid">
                {productsToShow.map((product) => (
                    <ProductCard key={`fav-${product.productID}`} product={product} />
                ))}
            </div>
        </section>
    );
};

export default FavoritesProducts;

