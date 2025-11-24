import User from "../../types/user";
import UserCard from "./UserCard";
import './UserList.css';

interface UserListProps {
    users: User[]
}

const UserList = ({ users }: UserListProps) => {
    if (users.length === 0) {
        return (
            <div className="user-list-container">
                <div className="no-users">No users found.</div>
            </div>
        );
    }

    return (
        <div className="user-list-container">
            <div className="user-list-header">
                <span></span> {/* Avatar spacer */}
                <span>User</span>
                <span>Contact</span>
                <span></span> {/* Actions spacer */}
            </div>
            <div className="user-list-body">
                {users.map((user, index) => (
                    <UserCard key={user.userId} user={user} index={index} />
                ))}
            </div>
        </div>
    );
}

export default UserList;
