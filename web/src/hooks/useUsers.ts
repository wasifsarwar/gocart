import { useEffect, useState } from "react";
import User from "../types/user";
import { userService, ApiUser } from "../services/userService";

const useUsers = () => {
    const [users, setUsers] = useState<User[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null)

    const convertApiUser = (apiUser: ApiUser): User => ({
        userId: apiUser.user_id,
        firstName: apiUser.first_name,
        lastName: apiUser.last_name,
        email: apiUser.email,
        phone: apiUser.phone,
        createdAt: new Date(apiUser.created_at),
        updatedAt: new Date(apiUser.updated_at)

    });
    const fetchUsers = async () => {
        try {
            setLoading(true);
            const apiUsers = await userService.getAllUsers()
            const convertedUser = apiUsers.map(convertApiUser)
            setUsers(convertedUser);
            setError(null)
        } catch (err) {
            setError('Failed to fetch users. Please try again')
            console.log('Error fetching all users: ', err)
        } finally {
            setLoading(false);
        }
    }

    useEffect(() => {
        fetchUsers();
    }, [fetchUsers]);

    return {
        users,
        loading,
        error,
        refetch: fetchUsers
    };
};

export default useUsers;