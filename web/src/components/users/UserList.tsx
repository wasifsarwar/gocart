import User from "../../types/user";
import UserCard from "./UserCard";

interface UserListProps {
    users: User[]
}

const UserList = ({ users }: UserListProps) => {
    if (users.length === 0) {
        return <div>No users available</div>
    }

    return (
        <table className="data-table">
            <thead>
                <th>Full Name</th>
                <th>Email</th>
                <th>Phone</th>
            </thead>
            <tbody>
                {users.map(user => (
                    <UserCard key={user.userId} user={user} />
                ))}
            </tbody>
        </table>

    );
}

export default UserList;