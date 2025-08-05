const API_URL = "http://localhost:8080/";

export interface ApiProduct {
    product_id: string;
    name: string;
    description: string;
    price: number;
    category: string;
}

export const productService = {
    async getAllProducts(): Promise<ApiProduct[]> {
        const response = await fetch(`${API_URL}/products`);
        if (!response.ok) {
            throw new Error(`HTTP error: status: ${response.status}`)
        }
        return response.json();
    },

    async getProductById(id: string): Promise<ApiProduct> {
        const response = await fetch(`${API_URL}/products/${id}`)
        if (!response.ok) {
            throw new Error(`HTTP error: status ${response.status}`);
        }
        return response.json();
    }
};