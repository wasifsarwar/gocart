import React from "react";
import User from "../../types/user";

interface UserCardProps {
    user: User;
    index: number; // Used for deterministic color
}

const UserCard = ({ user, index }: UserCardProps) => {
    // Generate initials
    const getInitials = (first: string, last: string) => {
        return `${first.charAt(0)}${last.charAt(0)}`.toUpperCase();
    };

    // Deterministic color based on index
    const getColorClass = (idx: number) => {
        const colors = ['avatar-0', 'avatar-1', 'avatar-2', 'avatar-3', 'avatar-4', 'avatar-5'];
        return colors[idx % colors.length];
    };

    // Format phone number: +1-(AreaCode) Number
    const formatPhone = (phone: string) => {
        // Clean the input first (remove existing formatting chars)
        const cleaned = ('' + phone).replace(/\D/g, '');
        
        // Check if we have enough digits (assuming US-style 10 digits)
        // If it's already formatted or has different length, return as is or adapt logic
        const match = cleaned.match(/^(\d{3})(\d{3})(\d{4})$/);
        
        if (match) {
            return `+1-${match[1]} ${match[2]}-${match[3]}`;
        }

        // Handle case where it might already start with 1
        const matchWithOne = cleaned.match(/^1(\d{3})(\d{3})(\d{4})$/);
        if (matchWithOne) {
            return `+1-${matchWithOne[1]} ${matchWithOne[2]}-${matchWithOne[3]}`;
        }

        // Fallback: if it doesn't match expected format, return original
        // or basic formatting if possible
        return phone;
    };

    return (
        <div className="user-row">
            <div className={`user-avatar ${getColorClass(index)}`}>
                {getInitials(user.firstName, user.lastName)}
            </div>
            
            <div className="user-main-info">
                <span className="user-name">{user.firstName} {user.lastName}</span>
                <span className="user-email">{user.email}</span>
            </div>

            <div className="user-phone">
                {formatPhone(user.phone)}
            </div>

            <div className="user-actions">
                <button className="action-btn">
                    View
                </button>
            </div>
        </div>
    );
}

export default UserCard;
