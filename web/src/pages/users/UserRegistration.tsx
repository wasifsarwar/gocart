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
            <div className="registration-container">
                <div className="registration-header">
                    <h1>Create Account</h1>
                    <p>Join our community to start shopping</p>
                </div>

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
        </div>
    );
};

export default UserRegistration;
