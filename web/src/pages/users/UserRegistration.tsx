import React, { useState } from "react";
import UserRegistrationForm from "../../components/users/UserRegistrationForm";

import './UserRegistration.css'

const UserRegistration = () => {
    const [message, setMessage] = useState<string>('');
    const [messageType, setMessageType] = useState<'success' | 'error' | ''>('');

    const handleSuccess = (successMessage: string) => {
        setMessage(successMessage);
        setMessageType('success');

        //clear message after 5 seconds
        setTimeout(() => {
            setMessage('');
            setMessageType('');
        }, 5000);
    };

    const handleError = (errorMessage: string) => {
        setMessage(errorMessage);
        setMessageType('error');

        // clear message after 5 seconds
        setTimeout(() => {
            setMessage('');
            setMessageType('');
        }, 5000);
    };

    return (
        <div className="user-registration-page page-container">
            <h1>Register New User</h1>
            {message && (
                <div className={`alert ${messageType === 'success' ? 'alert-success' : 'alert-error'}`}>
                    {message}
                </div>
            )}
            <UserRegistrationForm
                onError={handleError}
                onSuccess={handleSuccess}
            />
        </div>
    );
};

export default UserRegistration;