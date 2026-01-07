import React, { useMemo } from 'react';
import ProductCard from './ProductCard';
import useRecentlyViewedProducts from '../../hooks/useRecentlyViewedProducts';
import './RecentlyViewedProducts.css';

interface RecentlyViewedProductsProps {
    limit?: number;
    excludeProductID?: string;
}

const RecentlyViewedProducts: React.FC<RecentlyViewedProductsProps> = ({
    limit,
    excludeProductID
}) => {
    const { recentlyViewed, clearRecentlyViewed, maxItems } = useRecentlyViewedProducts();

    const productsToShow = useMemo(() => {
        const filtered = excludeProductID
            ? recentlyViewed.filter((p) => p.productID !== excludeProductID)
            : recentlyViewed;

        return filtered.slice(0, limit ?? maxItems);
    }, [recentlyViewed, excludeProductID, limit, maxItems]);

    if (productsToShow.length === 0) return null;

    return (
        <section className="recently-viewed-section" aria-label="Recently viewed products">
            <div className="recently-viewed-header">
                <h2 className="recently-viewed-title">Recently Viewed</h2>
                <button className="recently-viewed-clear" onClick={clearRecentlyViewed} type="button">
                    Clear
                </button>
            </div>
            <div className="product-grid recently-viewed-grid">
                {productsToShow.map((product) => (
                    <ProductCard key={`recent-${product.productID}`} product={product} />
                ))}
            </div>
        </section>
    );
};

export default RecentlyViewedProducts;

