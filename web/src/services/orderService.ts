const API_URL = process.env.REACT_APP_API_URL || "http://localhost:8080";

export interface OrderItem {
    product_id: string;
    quantity: number;
    price: number;
}

export interface CreateOrderRequest {
    user_id: string;
    items: OrderItem[];
}

export interface Order {
    order_id: string;
    user_id: string;
    total_amount: number;
    status: string;
    created_at: string;
    items: any[]; // We can define a more specific type if needed for viewing orders
}

export const orderService = {
    async createOrder(orderData: CreateOrderRequest): Promise<Order> {
        const response = await fetch(`${API_URL}/orders`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(orderData)
        });

        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(`Failed to create order: ${errorText}`);
        }

        return response.json();
    },

    async getOrdersByUserId(userId: string): Promise<Order[]> {
        const response = await fetch(`${API_URL}/orders/user/${userId}`);

        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(`Failed to fetch orders: ${errorText}`);
        }

        return response.json();
    }
};

