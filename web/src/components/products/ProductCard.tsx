import React from "react";
import { Link } from "react-router-dom";
import Product from "../../types/product";
import { FaShoppingBag } from "react-icons/fa";
import { IconType } from "react-icons";
import { useCart } from "../../context/CartContext";

interface ProductCardProps {
    product: Product
}

const usdFormatter = new Intl.NumberFormat('en-us', {
    style: 'currency',
    currency: 'USD'
});

// Wrapper to fix TS2786 error with React 19 types
const Icon = ({ icon: IconComponent, className }: { icon: IconType; className?: string }) => {
    const Component = IconComponent as any;
    return <Component className={className} />;
};

const ProductCard = ({ product }: ProductCardProps) => {
    const { addToCart } = useCart();

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

    const colorClass = getCategoryColor(product.category);

    const handleAddToCart = () => {
        addToCart(product);
        // Optional: Add visual feedback here (toast notification, etc.)
    };

    return (
        <div className="product-card">
            <Link to={`/products/${product.productID}`} className="product-link">
                <div className={`product-image-placeholder ${colorClass}`}>
                    <Icon icon={FaShoppingBag} />
                </div>
            </Link>
            <div className="product-content">
                <div className="product-header">
                    <span className={`category-badge ${colorClass}`}>
                        {product.category}
                    </span>
                    <span className="product-price">{usdFormatter.format(product.price)}</span>
                </div>
                <Link to={`/products/${product.productID}`} className="product-name-link">
                    <h3 className="product-name">{product.name}</h3>
                </Link>
                <p className="product-description">{product.description}</p>
                <button className="add-to-cart-btn" onClick={handleAddToCart}>
                    Add to Cart
                </button>
            </div>
        </div>
    );
}

export default ProductCard;
