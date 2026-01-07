import React, { useState, useEffect } from 'react';
import { useParams, useNavigate, Link } from 'react-router-dom';
import { FaArrowLeft, FaShoppingBag, FaShoppingCart } from 'react-icons/fa';
import { IconType } from 'react-icons';
import { productService, ApiProduct } from '../../services/productService';
import { useCart } from '../../context/CartContext';
import useProducts from '../../hooks/useProducts';
import Product from '../../types/product';
import './ProductDetails.css';

// Wrapper to fix TS2786 error with React 19 types
const Icon = ({ icon: IconComponent, className }: { icon: IconType; className?: string }) => {
    const Component = IconComponent as any;
    return <Component className={className} />;
};

const ProductDetails = () => {
    const { id } = useParams<{ id: string }>();
    const navigate = useNavigate();
    const { addToCart } = useCart();
    const { products } = useProducts();

    const [product, setProduct] = useState<ApiProduct | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    const usdFormatter = new Intl.NumberFormat('en-us', {
        style: 'currency',
        currency: 'USD'
    });

    useEffect(() => {
        // Reset state synchronously when id changes to prevent showing stale content
        setLoading(true);
        setError(null);
        setProduct(null);

        // Use AbortController to cancel fetch if component unmounts or id changes
        const abortController = new AbortController();
        let isCancelled = false;

        const fetchProduct = async () => {
            if (!id) {
                setError('Product ID is missing');
                setLoading(false);
                return;
            }

            try {
                const data = await productService.getProductById(id);

                // Only update state if this fetch hasn't been cancelled
                if (!isCancelled) {
                    setProduct(data);
                    setError(null);
                }
            } catch (err) {
                if (!isCancelled) {
                    console.error('Failed to fetch product:', err);
                    setError('Failed to load product details. Please try again.');
                }
            } finally {
                if (!isCancelled) {
                    setLoading(false);
                }
            }
        };

        fetchProduct();

        // Cleanup function to cancel stale requests
        return () => {
            isCancelled = true;
            abortController.abort();
        };
    }, [id]);

    const handleAddToCart = () => {
        if (product) {
            const cartProduct: Product = {
                productID: product.product_id,
                name: product.name,
                description: product.description,
                price: product.price,
                category: product.category
            };
            addToCart(cartProduct);
        }
    };

    const getCategoryColor = (category: string) => {
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

    // Get related products (same category, exclude current)
    const relatedProducts = product
        ? products
            .filter(p => p.category === product.category && p.productID !== product.product_id)
            .slice(0, 4)
        : [];

    if (loading) {
        return (
            <div className="product-details-page page-container">
                <div className="loading-spinner">Loading product details...</div>
            </div>
        );
    }

    if (error || !product) {
        return (
            <div className="product-details-page page-container">
                <div className="error-container">
                    <h2>Product Not Found</h2>
                    <p>{error || 'The product you are looking for does not exist.'}</p>
                    <Link to="/products" className="back-to-products-btn">
                        <Icon icon={FaArrowLeft} /> Back to Products
                    </Link>
                </div>
            </div>
        );
    }

    const colorClass = getCategoryColor(product.category);

    return (
        <div className="product-details-page page-container">
            <div className="breadcrumb">
                <Link to="/">Home</Link>
                <span className="separator">/</span>
                <Link to="/products">Products</Link>
                <span className="separator">/</span>
                <span className="current">{product.name}</span>
            </div>

            <button onClick={() => navigate('/products')} className="back-link">
                <Icon icon={FaArrowLeft} /> Back to Products
            </button>

            <div className="product-details-container">
                <div className="product-image-section">
                    <div className={`product-image-large ${colorClass}`}>
                        <Icon icon={FaShoppingBag} className="product-icon-large" />
                    </div>
                </div>

                <div className="product-info-section">
                    <div className="product-header">
                        <span className={`category-badge-large ${colorClass}`}>
                            {product.category}
                        </span>
                        <h1 className="product-title">{product.name}</h1>
                        <div className="product-price-large">
                            {usdFormatter.format(product.price)}
                        </div>
                    </div>

                    <div className="product-description-section">
                        <h2>Description</h2>
                        <p className="product-description">{product.description}</p>
                    </div>

                    <div className="product-actions">
                        <button onClick={handleAddToCart} className="add-to-cart-btn-large">
                            <Icon icon={FaShoppingCart} />
                            Add to Cart
                        </button>
                    </div>
                </div>
            </div>

            {relatedProducts.length > 0 && (
                <div className="related-products-section">
                    <h2>Related Products</h2>
                    <div className="related-products-grid">
                        {relatedProducts.map(relatedProduct => (
                            <Link
                                key={relatedProduct.productID}
                                to={`/products/${relatedProduct.productID}`}
                                className="related-product-card"
                            >
                                <div className={`related-product-image ${getCategoryColor(relatedProduct.category)}`}>
                                    <Icon icon={FaShoppingBag} />
                                </div>
                                <div className="related-product-info">
                                    <h3>{relatedProduct.name}</h3>
                                    <p className="related-product-price">
                                        {usdFormatter.format(relatedProduct.price)}
                                    </p>
                                </div>
                            </Link>
                        ))}
                    </div>
                </div>
            )}
        </div>
    );
};

export default ProductDetails;
