import React from 'react';
import './PriceRangeFilter.css';

interface PriceRangeFilterProps {
    minPrice: number;
    maxPrice: number;
    currentMin: number;
    currentMax: number;
    onPriceChange: (min: number, max: number) => void;
}

const PriceRangeFilter: React.FC<PriceRangeFilterProps> = ({
    minPrice,
    maxPrice,
    currentMin,
    currentMax,
    onPriceChange
}) => {
    const usdFormatter = new Intl.NumberFormat('en-us', {
        style: 'currency',
        currency: 'USD',
        maximumFractionDigits: 0
    });

    const handleMinChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const newMin = Number(e.target.value);
        if (newMin <= currentMax) {
            onPriceChange(newMin, currentMax);
        }
    };

    const handleMaxChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const newMax = Number(e.target.value);
        if (newMax >= currentMin) {
            onPriceChange(currentMin, newMax);
        }
    };

    return (
        <div className="price-range-filter">
            <label className="price-filter-label">Price Range</label>
            <div className="price-display">
                <span>{usdFormatter.format(currentMin)}</span>
                <span className="separator">-</span>
                <span>{usdFormatter.format(currentMax)}</span>
            </div>
            <div className="price-sliders">
                <input
                    type="range"
                    min={minPrice}
                    max={maxPrice}
                    value={currentMin}
                    onChange={handleMinChange}
                    className="price-slider price-slider-min"
                    aria-label="Minimum price"
                />
                <input
                    type="range"
                    min={minPrice}
                    max={maxPrice}
                    value={currentMax}
                    onChange={handleMaxChange}
                    className="price-slider price-slider-max"
                    aria-label="Maximum price"
                />
            </div>
            <div className="price-labels">
                <span className="price-label-min">{usdFormatter.format(minPrice)}</span>
                <span className="price-label-max">{usdFormatter.format(maxPrice)}</span>
            </div>
        </div>
    );
};

export default PriceRangeFilter;
