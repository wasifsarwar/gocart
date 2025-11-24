import React, { useState } from 'react';
import { Link, useNavigate, useLocation } from 'react-router-dom';
import { toast } from 'react-hot-toast';
import { useCart } from '../../context/CartContext';
import { useAuth } from '../../context/AuthContext';
import { orderService } from '../../services/orderService';
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
    const { user, isAuthenticated } = useAuth();
    const navigate = useNavigate();
    const location = useLocation();

    const [step, setStep] = useState<'cart' | 'shipping'>('cart');
    const [shippingDetails, setShippingDetails] = useState({
        address: '',
        city: '',
        zipCode: '',
        country: ''
    });

    const [isSubmitting, setIsSubmitting] = useState(false);
    const [orderStatus, setOrderStatus] = useState<'idle' | 'success' | 'error'>('idle');
    const [errorMessage, setErrorMessage] = useState('');

    const usdFormatter = new Intl.NumberFormat('en-us', {
        style: 'currency',
        currency: 'USD'
    });

    const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const { name, value } = e.target;
        setShippingDetails(prev => ({
            ...prev,
            [name]: value
        }));
    };

    const handleProceedToShipping = () => {
        if (!isAuthenticated) {
            navigate('/login', { state: { from: location } });
            return;
        }
        setStep('shipping');
    };

    const handleBackToCart = () => {
        setStep('cart');
    };

    const handlePlaceOrder = async (e: React.FormEvent) => {
        e.preventDefault();

        if (!user) return;

        setIsSubmitting(true);
        setErrorMessage('');
        try {
            const orderItems = items.map(item => ({
                product_id: item.productID,
                quantity: item.quantity,
                price: item.price
            }));

            await orderService.createOrder({
                user_id: user.user_id,
                items: orderItems,
                shipping_address: shippingDetails.address,
                city: shippingDetails.city,
                zip_code: shippingDetails.zipCode,
                country: shippingDetails.country
            });

            setOrderStatus('success');
            clearCart();
            toast.success('Order placed successfully! Check your email for confirmation.');
            setTimeout(() => {
                navigate('/orders');
            }, 3000);
        } catch (error) {
            console.error('Checkout failed:', error);
            const msg = error instanceof Error ? error.message : 'Failed to place order';
            setErrorMessage(msg);
            toast.error(msg);
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
                {step === 'cart' ? (
                    <Link to="/products" className="back-link">
                        <Icon icon={FaArrowLeft} /> Continue Shopping
                    </Link>
                ) : (
                    <button onClick={handleBackToCart} className="back-link-btn">
                        <Icon icon={FaArrowLeft} /> Back to Cart
                    </button>
                )}
                <h1>{step === 'cart' ? 'Shopping Cart' : 'Shipping Details'}</h1>
            </div>

            <div className="checkout-content">
                {step === 'cart' ? (
                    <>
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
                                onClick={handleProceedToShipping}
                            >
                                {isAuthenticated ? 'Proceed to Checkout' : 'Login to Checkout'}
                            </button>
                        </div>
                    </>
                ) : (
                    <div className="shipping-form-container">
                        <form onSubmit={handlePlaceOrder} className="shipping-form">
                            <div className="form-group">
                                <label htmlFor="address">Address</label>
                                <input
                                    type="text"
                                    id="address"
                                    name="address"
                                    value={shippingDetails.address}
                                    onChange={handleInputChange}
                                    required
                                    placeholder="123 Main St"
                                />
                            </div>

                            <div className="form-row">
                                <div className="form-group">
                                    <label htmlFor="city">City</label>
                                    <input
                                        type="text"
                                        id="city"
                                        name="city"
                                        value={shippingDetails.city}
                                        onChange={handleInputChange}
                                        required
                                        placeholder="New York"
                                    />
                                </div>
                                <div className="form-group">
                                    <label htmlFor="zipCode">Zip Code</label>
                                    <input
                                        type="text"
                                        id="zipCode"
                                        name="zipCode"
                                        value={shippingDetails.zipCode}
                                        onChange={handleInputChange}
                                        required
                                        placeholder="10001"
                                    />
                                </div>
                            </div>

                            <div className="form-group">
                                <label htmlFor="country">Country</label>
                                <input
                                    type="text"
                                    id="country"
                                    name="country"
                                    value={shippingDetails.country}
                                    onChange={handleInputChange}
                                    required
                                    placeholder="USA"
                                />
                            </div>

                            <div className="order-summary-mini">
                                <h3>Order Total: {usdFormatter.format(cartTotal)}</h3>
                                <p className="shipping-note">Free Shipping</p>
                            </div>

                            <div className="form-actions">
                                <button
                                    type="submit"
                                    className="checkout-btn"
                                    disabled={isSubmitting}
                                >
                                    {isSubmitting ? 'Processing...' : 'Place Order'}
                                </button>
                            </div>

                            {orderStatus === 'error' && (
                                <p className="error-message">{errorMessage || 'Something went wrong. Please try again.'}</p>
                            )}
                        </form>
                    </div>
                )}
            </div>
        </div>
    );
};

export default Checkout;
