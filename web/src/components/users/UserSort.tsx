import React from 'react';

interface SortOption {
    value: string;
    label: string;
}

interface ProductSortProps {
    onSort: (sortValue: string) => void
    currentSort: string
}

const sortOptions: SortOption[] = [
    { value: 'name-asc', label: 'Name (A-Z)' },
    { value: 'name-desc', label: 'Name (Z-A)' },
    { value: 'email-asc', label: 'Email (A-Z)' }
];

const UserSort = ({ onSort, currentSort }: ProductSortProps) => {
    return (
        <div className='user-sort'>
            <label htmlFor='sort-select'>Sort By</label>
            <select
                id='sort-select'
                value={currentSort}
                onChange={(e) => onSort(e.target.value)}
                className='sort-select'
            >
                {sortOptions.map(option => (
                    <option key={option.value} value={option.value}>{option.label}</option>
                ))}
            </select>
        </div>
    );
};

export default UserSort;