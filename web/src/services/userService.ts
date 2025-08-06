const API_URL = "http://localhost:8081"

export interface ApiUser {
    user_id: string;
    first_name: string;
    last_name: string;
    email: string;
    phone: string;
    created_at: string;
    updated_at: string;
}

export interface CreateUserRequest {
    first_name: string;
    last_name: string;
    email: string;
    phone: string;
}

export const userService = {
    async getAllUsers(): Promise<ApiUser[]> {
        const response = await fetch(`${API_URL}/users`)
        if (!response.ok) {
            throw new Error(`HTTP error: status: ${response.status}`)
        }
        return response.json();
    },

    async getUserById(id: string): Promise<ApiUser> {
        const response = await fetch(`${API_URL}/users/${id}`)
        if (!response.ok) {
            throw new Error(`HTTP error: status: ${response.status}`)
        }
        return response.json();
    },

    async createUser(userData: CreateUserRequest): Promise<ApiUser> {
        const response = await fetch(`${API_URL}/users/register`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(userData)
        });

        if (!response.ok) {
            const errorData = await response.text();
            throw new Error(`HTTP error: status: ${response.status}, message: ${errorData}`)
        }
        return response.json();

    }
}