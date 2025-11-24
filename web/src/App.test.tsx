import React from 'react';
import { render, screen } from '@testing-library/react';
import App from './App';

test('renders GoCart app', () => {
  render(<App />);
  // Look for the brand name in the navbar or home page
  const brandElements = screen.getAllByText(/GoCart/i);
  expect(brandElements.length).toBeGreaterThan(0);
});
