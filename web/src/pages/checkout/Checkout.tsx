import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useCart } from '../../context/CartContext';
import { FaTrash, FaMinus, FaPlus, FaArrowLeft } from 'react-icons/fa';
import { IconType } from 'react-icons';
import './Checkout.css';

// Wrapper to fix TS2786 error with React 19 types
const Icon = ({ icon: IconComponent, className }: { icon: IconType; className?: string }) => {
    const Component = IconComponent as any;
    return <Component className={className} />;
};

const Checkout = () => {
    const { items, updateQuantity, removeFromCart, cartTotal, clearCart } = useCart();
    const navigate = useNavigate();
    const [isSubmitting, setIsSubmitting] = useState(false);
    const [orderStatus, setOrderStatus] = useState<'idle' | 'success' | 'error'>('idle');

    const usdFormatter = new Intl.NumberFormat('en-us', {
        style: 'currency',
        currency: 'USD'
    });

    const handleCheckout = async () => {
        setIsSubmitting(true);
        try {
            // TODO: Implement actual API call to create order
            // For now, simulate API call
            await new Promise(resolve => setTimeout(resolve, 1500));
            
            setOrderStatus('success');
            clearCart();
            setTimeout(() => {
                navigate('/');
            }, 3000);
        } catch (error) {
            console.error('Checkout failed:', error);
            setOrderStatus('error');
        } finally {
            setIsSubmitting(false);
        }
    };

    if (items.length === 0 && orderStatus !== 'success') {
        return (
            <div className="checkout-page page-container">
                <div className="empty-cart">
                    <h2>Your cart is empty</h2>
                    <p>Looks like you haven't added any items to your cart yet.</p>
                    <Link to="/products" className="continue-shopping-btn">
                        Start Shopping
                    </Link>
                </div>
            </div>
        );
    }

    if (orderStatus === 'success') {
        return (
            <div className="checkout-page page-container">
                <div className="order-success">
                    <div className="success-icon">🎉</div>
                    <h2>Order Placed Successfully!</h2>
                    <p>Thank you for your purchase. You will receive a confirmation email shortly.</p>
                    <p>Redirecting you to home page...</p>
                </div>
            </div>
        );
    }

    return (
        <div className="checkout-page page-container">
            <div className="checkout-header">
                <Link to="/products" className="back-link">
                    <Icon icon={FaArrowLeft} /> Continue Shopping
                </Link>
                <h1>Shopping Cart</h1>
            </div>

            <div className="checkout-content">
                <div className="cart-items">
                    {items.map(item => (
                        <div key={item.productID} className="cart-item">
                            <div className="item-image-placeholder">
                                {/* Placeholder for product image */}
                            </div>
                            <div className="item-details">
                                <h3>{item.name}</h3>
                                <p className="item-price">{usdFormatter.format(item.price)}</p>
                                <p className="item-category">{item.category}</p>
                            </div>
                            <div className="item-quantity">
                                <button 
                                    onClick={() => updateQuantity(item.productID, item.quantity - 1)}
                                    className="quantity-btn"
                                    aria-label="Decrease quantity"
                                >
                                    <Icon icon={FaMinus} />
                                </button>
                                <span>{item.quantity}</span>
                                <button 
                                    onClick={() => updateQuantity(item.productID, item.quantity + 1)}
                                    className="quantity-btn"
                                    aria-label="Increase quantity"
                                >
                                    <Icon icon={FaPlus} />
                                </button>
                            </div>
                            <div className="item-total">
                                {usdFormatter.format(item.price * item.quantity)}
                            </div>
                            <button 
                                onClick={() => removeFromCart(item.productID)}
                                className="remove-btn"
                                aria-label="Remove item"
                            >
                                <Icon icon={FaTrash} />
                            </button>
                        </div>
                    ))}
                </div>

                <div className="order-summary">
                    <h2>Order Summary</h2>
                    <div className="summary-row">
                        <span>Subtotal</span>
                        <span>{usdFormatter.format(cartTotal)}</span>
                    </div>
                    <div className="summary-row">
                        <span>Shipping</span>
                        <span>Free</span>
                    </div>
                    <div className="summary-row total">
                        <span>Total</span>
                        <span>{usdFormatter.format(cartTotal)}</span>
                    </div>
                    <button 
                        className="checkout-btn"
                        onClick={handleCheckout}
                        disabled={isSubmitting}
                    >
                        {isSubmitting ? 'Processing...' : 'Place Order'}
                    </button>
                    {orderStatus === 'error' && (
                        <p className="error-message">Something went wrong. Please try again.</p>
                    )}
                </div>
            </div>
        </div>
    );
};

export default Checkout;

