import React, { useEffect } from 'react';
import { Link } from 'react-router-dom';
import { FaHeart, FaRegHeart, FaShoppingCart, FaTimes } from 'react-icons/fa';
import { IconType } from 'react-icons';
import Product from '../../types/product';
import { useCart } from '../../context/CartContext';
import { useFavorites } from '../../context/FavoritesContext';
import { resolveImageUrl } from '../../utils/resolveImageUrl';
import './ProductQuickViewModal.css';

// Wrapper to fix TS2786 error with React 19 types
const Icon = ({ icon: IconComponent, className }: { icon: IconType; className?: string }) => {
    const Component = IconComponent as any;
    return <Component className={className} />;
};

interface ProductQuickViewModalProps {
    product: Product;
    onClose: () => void;
}

const usdFormatter = new Intl.NumberFormat('en-us', {
    style: 'currency',
    currency: 'USD'
});

const ProductQuickViewModal: React.FC<ProductQuickViewModalProps> = ({ product, onClose }) => {
    const { addToCart } = useCart();
    const { isFavorite, toggleFavorite } = useFavorites();

    const favorited = isFavorite(product.productID);
    const imageSrc = resolveImageUrl(product.imageUrl);

    useEffect(() => {
        const prevOverflow = document.body.style.overflow;
        document.body.style.overflow = 'hidden';
        return () => {
            document.body.style.overflow = prevOverflow;
        };
    }, []);

    useEffect(() => {
        const onKeyDown = (e: KeyboardEvent) => {
            if (e.key === 'Escape') onClose();
        };
        window.addEventListener('keydown', onKeyDown);
        return () => window.removeEventListener('keydown', onKeyDown);
    }, [onClose]);

    return (
        <div className="qv-overlay" role="dialog" aria-modal="true" aria-label={`Quick view: ${product.name}`}>
            <button className="qv-backdrop" aria-label="Close quick view" onClick={onClose} type="button" />
            <div className="qv-modal">
                <div className="qv-header">
                    <h2 className="qv-title">{product.name}</h2>
                    <button className="qv-close" onClick={onClose} aria-label="Close" type="button">
                        <Icon icon={FaTimes} />
                    </button>
                </div>

                <div className="qv-body">
                    {imageSrc && (
                        <img className="qv-image" src={imageSrc} alt={product.name} loading="lazy" />
                    )}
                    <div className="qv-meta">
                        <span className="qv-category">{product.category}</span>
                        <span className="qv-price">{usdFormatter.format(product.price)}</span>
                    </div>
                    <p className="qv-description">{product.description}</p>
                </div>

                <div className="qv-actions">
                    <button className="qv-action qv-primary" onClick={() => addToCart(product)} type="button">
                        <Icon icon={FaShoppingCart} />
                        Add to cart
                    </button>
                    <button
                        className={`qv-action qv-secondary ${favorited ? 'favorited' : ''}`}
                        onClick={() => toggleFavorite(product.productID)}
                        type="button"
                    >
                        <Icon icon={favorited ? FaHeart : FaRegHeart} />
                        {favorited ? 'Favorited' : 'Favorite'}
                    </button>
                    <Link className="qv-link" to={`/products/${product.productID}`} onClick={onClose}>
                        View full details
                    </Link>
                </div>
            </div>
        </div>
    );
};

export default ProductQuickViewModal;

