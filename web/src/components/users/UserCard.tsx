import User from "../../types/user"

interface UserCardProps {
    user: User
}

const UserCard = ({ user }: UserCardProps) => {
    return (
        <tr className="user-card">
            <td className="first-last-name">{user.firstName} {user.lastName}</td>
            <td className="email">{user.email}</td>
            <td className="phone">{user.phone}</td>
        </tr>
    );
}
export default UserCard;