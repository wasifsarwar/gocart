import Product from "../../types/product"
import ProductCard from "./ProductCard"
import './ProductList.css';

interface ProductListProps {
    products: Product[];
    loading: boolean;
    onQuickView?: (product: Product) => void;
}

const ProductList = ({ products, loading, onQuickView }: ProductListProps) => {

    const skeletonCards = Array.from({ length: 8 });

    if (loading) {
        return (
            <div className="product-grid">
                {skeletonCards.map((_, i) => (
                    <div key={i} className="product-card skeleton-card">
                        <div className="skeleton-image" />
                        <div className="skeleton-content">
                            <div className="skeleton-line w-25" />
                            <div className="skeleton-line w-75" />
                            <div className="skeleton-line w-50" />
                            <div className="skeleton-button" />
                        </div>
                    </div>
                ))}
            </div>
        );
    }

    if (products.length === 0) {
        return (
            <div className="no-products">
                <div className="empty-state">
                    <h3>No Products Found</h3>
                    <p>Try adjusting your search or filter to find what you're looking for.</p>
                </div>
            </div>
        );
    }

    return (
        <div className="product-grid">
            {products.map(product => (
                <ProductCard key={product.productID} product={product} onQuickView={onQuickView} />
            ))}
        </div>
    );
}

export default ProductList;
