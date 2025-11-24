import React, { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';
import { orderService, Order } from '../../services/orderService';
import { productService, ApiProduct } from '../../services/productService';
import { FaBox } from 'react-icons/fa';
import { IconType } from 'react-icons';
import './OrderHistory.css';

// Wrapper to fix TS2786 error with React 19 types
const Icon = ({ icon: IconComponent, className }: { icon: IconType; className?: string }) => {
    const Component = IconComponent as any;
    return <Component className={className} />;
};

const OrderHistory = () => {
    const { user } = useAuth();
    const [orders, setOrders] = useState<Order[]>([]);
    const [products, setProducts] = useState<Record<string, ApiProduct>>({});
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    const usdFormatter = new Intl.NumberFormat('en-us', {
        style: 'currency',
        currency: 'USD'
    });

    const dateFormatter = new Intl.DateTimeFormat('en-US', {
        year: 'numeric',
        month: 'long',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
    });

    useEffect(() => {
        const fetchOrdersAndProducts = async () => {
            if (!user) return;

            try {
                setLoading(true);
                const userOrders = await orderService.getOrdersByUserId(user.user_id);

                // Sort orders by date descending
                userOrders.sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime());
                setOrders(userOrders);

                // Collect all unique product IDs
                const productIds = new Set<string>();
                userOrders.forEach(order => {
                    order.items?.forEach((item: any) => {
                        if (item.product_id) productIds.add(item.product_id);
                    });
                });

                // Fetch product details
                const productMap: Record<string, ApiProduct> = {};
                await Promise.all(
                    Array.from(productIds).map(async (id) => {
                        try {
                            const product = await productService.getProductById(id);
                            productMap[id] = product;
                        } catch (err) {
                            console.error(`Failed to fetch product ${id}`, err);
                        }
                    })
                );
                setProducts(productMap);

            } catch (err) {
                setError('Failed to load order history. Please try again later.');
                console.error(err);
            } finally {
                setLoading(false);
            }
        };

        fetchOrdersAndProducts();
    }, [user]);

    if (loading) {
        return (
            <div className="order-history-page page-container">
                <div className="loading-spinner">Loading orders...</div>
            </div>
        );
    }

    if (error) {
        return (
            <div className="order-history-page page-container">
                <div className="error-message">{error}</div>
            </div>
        );
    }

    return (
        <div className="order-history-page page-container">
            <div className="orders-header">
                <h1>My Orders</h1>
                <p>Track and manage your recent purchases</p>
            </div>

            {orders.length === 0 ? (
                <div className="empty-orders">
                    <Icon icon={FaBox} className="empty-icon" />
                    <h2>No orders yet</h2>
                    <p>Looks like you haven't made any purchases yet.</p>
                    <Link to="/products" className="start-shopping-btn">
                        Start Shopping
                    </Link>
                </div>
            ) : (
                <div className="orders-list">
                    {orders.map(order => (
                        <div key={order.order_id} className="order-card">
                            <div className="order-header">
                                <div className="header-info-group">
                                    <div className="header-item">
                                        <span className="label">ORDER PLACED</span>
                                        <span className="value">{dateFormatter.format(new Date(order.created_at))}</span>
                                    </div>
                                    <div className="header-item">
                                        <span className="label">TOTAL</span>
                                        <span className="value">{usdFormatter.format(order.total_amount)}</span>
                                    </div>
                                    <div className="header-item">
                                        <span className="label">ORDER #</span>
                                        <span className="value">{order.friendly_id || order.order_id.slice(0, 8)}</span>
                                    </div>
                                </div>
                                <div className={`order-status status-${order.status?.toLowerCase() || 'pending'}`}>
                                    {order.status || 'Pending'}
                                </div>
                            </div>

                            <div className="order-content">
                                <div className="order-items">
                                    {order.items?.map((item: any, index: number) => {
                                        const product = products[item.product_id];
                                        return (
                                            <div key={index} className="order-item-row">
                                                <div className="item-image-placeholder">
                                                    {/* Image placeholder */}
                                                </div>
                                                <div className="item-details">
                                                    <span className="item-name">
                                                        {product ? product.name : `Product ID: ${item.product_id.slice(0, 8)}...`}
                                                    </span>
                                                    {product && <span className="item-category">{product.category}</span>}
                                                    <div className="item-meta-mobile">
                                                        <span>Qty: {item.quantity}</span>
                                                    </div>
                                                </div>
                                                <div className="item-price-qty">
                                                    <span className="item-price">{usdFormatter.format(item.price)}</span>
                                                    <span className="item-qty-badge">Qty: {item.quantity}</span>
                                                </div>
                                            </div>
                                        );
                                    })}
                                </div>

                                {(order.shipping_address || order.city) && (
                                    <div className="order-sidebar">
                                        <div className="sidebar-section">
                                            <h4>Shipping Address</h4>
                                            <div className="address-block">
                                                <p>{order.shipping_address}</p>
                                                <p>{order.city}, {order.zip_code}</p>
                                                <p>{order.country}</p>
                                            </div>
                                        </div>
                                    </div>
                                )}
                            </div>
                        </div>
                    ))}
                </div>
            )}
        </div>
    );
};

export default OrderHistory;

