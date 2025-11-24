import { useState, useMemo } from 'react'
import UserList from "../../components/users/UserList";
import useUsers from "../../hooks/useUsers";

import './Users.css'
import UserSort from '../../components/users/UserSort';

const Users = () => {
    const { users, loading, error } = useUsers();
    const [sortBy, setSortBy] = useState('name-asc');

    const sortedUsers = useMemo(() => {
        return [...users].sort((a, b) => {
            switch (sortBy) {
                case 'name-asc':
                    const fullNameA = `${a.firstName} ${a.lastName}`.toLowerCase();
                    const fullNameB = `${b.firstName} ${b.lastName}`.toLowerCase();
                    return fullNameA.localeCompare(fullNameB);
                case 'name-desc':
                    const fullNameA2 = `${a.firstName} ${a.lastName}`.toLowerCase();
                    const fullNameB2 = `${b.firstName} ${b.lastName}`.toLowerCase();
                    return fullNameB2.localeCompare(fullNameA2);
                case 'email-asc':
                    return a.email.toLowerCase().localeCompare(b.email.toLowerCase());
                default:
                    return 0;
            }
        })
    }, [users, sortBy]);

    return (
        <div className="users-page page-container">
            <div className="users-header-area">
                <h1>Users</h1>
                <UserSort onSort={setSortBy} currentSort={sortBy} />
            </div>

            {loading && (
                <div className="user-list-container">
                    <div className="user-list-header">
                        <span></span>
                        <span>User</span>
                        <span>Contact</span>
                        <span></span>
                    </div>
                    <div style={{ padding: '2rem', textAlign: 'center', color: '#64748b' }}>
                        Loading users...
                    </div>
                </div>
            )}
            {error && (
                <div className="alert alert-error">
                    <span>{error}</span>
                </div>
            )}
            {!loading && !error && <UserList users={sortedUsers} />}
        </div>
    )
}

export default Users;
