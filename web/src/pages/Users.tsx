import UserList from "../components/users/UserList";
import useUsers from "../hooks/useUsers";
import Navigation
    from "../components/Navigation";
const Users = () => {
    const { users, loading, error } = useUsers();
    return (
        <div className="users-page">
            <Navigation title="GoCart Users" />
            {loading && <p>Loading Users..</p>}
            {error && <p style={{ color: 'red' }}>{error}</p>}
            {!loading && !error && <UserList users={users} />}
        </div>
    )
}

export default Users;