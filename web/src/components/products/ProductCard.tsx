import React, { useState } from "react";
import Product from "../../types/product";

interface ProductCardProps {
    product: Product
}

const usdFormatter = new Intl.NumberFormat('en-us', {
    style: 'currency',
    currency: 'USD'
});

const ProductCard = ({ product }: ProductCardProps) => {
    const [isHovered, setIsHovered] = useState(false);
    const [isExpanded, setIsExpanded] = useState(false);

    const handleRowClick = () => {
        setIsExpanded(!isExpanded);
    };

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
        <tr 
            className={`product-card ${isHovered ? 'hovered' : ''} ${isExpanded ? 'expanded' : ''}`}
            onMouseEnter={() => setIsHovered(true)}
            onMouseLeave={() => setIsHovered(false)}
            onClick={handleRowClick}
        >
            <td className="name">
                <div className="name-container">
                    <span className="product-name">{product.name}</span>
                    <span className="expand-indicator">
                        {isExpanded ? '▼' : '▶'}
                    </span>
                </div>
                {isExpanded && (
                    <div className="expanded-details">
                        <p><strong>Full Description:</strong></p>
                        <p>{product.description}</p>
                        <div className="product-meta">
                            <span className="product-id">ID: {product.productID}</span>
                        </div>
                    </div>
                )}
            </td>
            <td className="price">
                <span className="price-amount">{usdFormatter.format(product.price)}</span>
            </td>
            <td className="category">
                <span 
                    className="category-badge" 
                    style={{ backgroundColor: getCategoryColor(product.category) }}
                >
                    {product.category}
                </span>
            </td>
            <td className="description">
                <span className="description-text">
                    {product.description.length > 60 && !isExpanded
                        ? `${product.description.substring(0, 60)}...`
                        : product.description
                    }
                </span>
            </td>
        </tr>
    );
}

export default ProductCard;