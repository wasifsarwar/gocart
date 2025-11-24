import React, { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';
import { orderService, Order } from '../../services/orderService';
import { FaBox, FaCalendar, FaMoneyBillWave } from 'react-icons/fa';
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
        const fetchOrders = async () => {
            if (!user) return;

            try {
                setLoading(true);
                const userOrders = await orderService.getOrdersByUserId(user.user_id);
                setOrders(userOrders);
            } catch (err) {
                setError('Failed to load order history. Please try again later.');
                console.error(err);
            } finally {
                setLoading(false);
            }
        };

        fetchOrders();
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
                                <div className="order-id">
                                    <span className="label">Order ID:</span>
                                    <span className="value">#{order.order_id.slice(0, 8)}</span>
                                </div>
                                <div className="order-date">
                                    <Icon icon={FaCalendar} className="icon" />
                                    <span>{dateFormatter.format(new Date(order.created_at))}</span>
                                </div>
                                <div className={`order-status status-${order.status || 'pending'}`}>
                                    {order.status || 'Pending'}
                                </div>
                            </div>
                            
                            <div className="order-items">
                                {order.items?.map((item: any, index: number) => (
                                    <div key={index} className="order-item-row">
                                        <span className="item-name">Product ID: {item.product_id.slice(0, 8)}...</span>
                                        <span className="item-quantity">x{item.quantity}</span>
                                        <span className="item-price">{usdFormatter.format(item.price)}</span>
                                    </div>
                                ))}
                            </div>

                            <div className="order-footer">
                                <div className="order-total">
                                    <Icon icon={FaMoneyBillWave} className="icon" />
                                    <span className="label">Total:</span>
                                    <span className="value">{usdFormatter.format(order.total_amount)}</span>
                                </div>
                            </div>
                        </div>
                    ))}
                </div>
            )}
        </div>
    );
};

export default OrderHistory;

