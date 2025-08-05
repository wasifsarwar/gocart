import React, { useState, JSX } from 'react';
import './App.css';
import Product from './types/product'
import ProductList from './components/ProductList';

function App(): JSX.Element {

  const mockProducts: Product[] = [{
    productID: "111",
    name: "gaming laptop",
    description: "high performance laptop for gaming and productivity",
    price: 1299.99,
    category: "Electronics"
  },
  {
    productID: "2",
    name: "Wireless Headphones",
    description: "Noise-cancelling bluetooth headphones",
    price: 199.99,
    category: "Audio"
  },
  {
    productID: "3",
    name: "Mechanical Keyboard",
    description: "RGB backlit mechanical gaming keyboard",
    price: 149.99,
    category: "Accessories"
  }]

  return (
    <div className="App">
      <header className="App-header">
        <h1>GoCart Products</h1>
        <ProductList products={mockProducts} />
      </header>
    </div>
  );
}

export default App;
