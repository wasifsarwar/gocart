import React, { useState } from "react";
import { CreateUserRequest, userService } from "../../services/userService";
import './UserRegistrationForm.css'
interface UserRegistrationFormProps {
    onSuccess?: (message: string) => void;
    onError?: (error: string) => void;
}

const UserRegistrationForm = ({ onSuccess, onError }: UserRegistrationFormProps) => {

    const [formData, setFormData] = useState<CreateUserRequest>({
        first_name: '',
        last_name: '',
        email: '',
        phone: ''
    });

    const [isLoading, setIsLoading] = useState(false);
    const [errors, setErrors] = useState<Partial<CreateUserRequest>>({});

    const validateForm = (): boolean => {
        const newErrors: Partial<CreateUserRequest> = {};
        if (!formData.first_name.trim()) {
            newErrors.first_name = 'First Name is required'
        }
        if (!formData.last_name.trim()) {
            newErrors.last_name = 'Last Name is required'
        }
        if (!formData.email.trim()) {
            newErrors.email = 'Email address is required'
        } else if (!/\S+@\S+\.\S+/.test(formData.email)) {
            newErrors.email = 'Please enter a valid email address';
        }
        if (!formData.phone.trim()) {
            newErrors.phone = 'Phone number is required'
        }

        setErrors(newErrors);
        return Object.keys(newErrors).length === 0;
    };

    const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const { name, value } = e.target;
        setFormData(prev => ({
            ...prev,
            [name]: value
        }));


        //clear error for this field when user starts typing
        if (errors[name as keyof CreateUserRequest]) {
            setErrors(prev => ({
                ...prev,
                [name]: undefined
            }));
        }
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();

        if (!validateForm()) return;

        setIsLoading(true);

        try {
            const newUser = await userService.createUser(formData)

            //Reset form
            setFormData({
                first_name: '',
                last_name: '',
                email: '',
                phone: ''
            });

            onSuccess?.(`User ${newUser.first_name} ${newUser.last_name} created successfully!`);
        } catch (error) {
            const errorMessage = error instanceof Error ? error.message : 'Failed to create user';
            onError?.(errorMessage);
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <form onSubmit={handleSubmit} className="user-registration-form">
            <h2>Register New User</h2>
            <div className="form-row">
                <div className="form-group">
                    <label htmlFor="first_name">First Name</label>
                    <input
                        type="text"
                        id="first_name"
                        name="first_name"
                        value={formData.first_name}
                        onChange={handleInputChange}
                        className={errors.first_name ? 'error' : ''}
                        disabled={isLoading}
                    />
                    {errors.first_name && <span className="error-message">{errors.first_name}</span>}
                </div>

                <div className="form-group">
                    <label htmlFor="last_name">Last Name</label>
                    <input
                        type="text"
                        id="last_name"
                        name="last_name"
                        value={formData.last_name}
                        onChange={handleInputChange}
                        className={errors.last_name ? 'error' : ''}
                        disabled={isLoading}
                    />
                    {errors.last_name && <span className="error-message">{errors.last_name}</span>}
                </div>
            </div>
            <div className="form-group">
                <label htmlFor="email">Email</label>
                <input
                    type="email"
                    id="email"
                    name="email"
                    value={formData.email}
                    onChange={handleInputChange}
                    className={errors.email ? 'error' : ''}
                    disabled={isLoading}
                />
                {errors.email && <span className="error-message">{errors.email}</span>}

            </div>

            <div className="form-group">
                <label htmlFor="phone">Phone Number</label>
                <input
                    type="tel"
                    id="phone"
                    name="phone"
                    value={formData.phone}
                    onChange={handleInputChange}
                    className={errors.phone ? 'error' : ''}
                    disabled={isLoading}
                />
                {errors.phone && <span className="error-message">{errors.phone}</span>}
            </div>

            <div className="form-actions">
                <button
                    type="submit"
                    disabled={isLoading}
                    className="submit-button">
                    {isLoading ? 'Creating User..' : 'Create User'}
                </button>
            </div>

        </form>
    );
};

export default UserRegistrationForm;