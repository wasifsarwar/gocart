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
            <h1>GoCart Users</h1>
            <div className='users-controls' >
                <UserSort onSort={setSortBy} currentSort={sortBy} />
            </div>
            {loading && <p>Loading Users..</p>}
            {error && <p style={{ color: 'red' }}>{error}</p>}
            {!loading && !error && <UserList users={sortedUsers} />}
        </div>
    )
}

export default Users;