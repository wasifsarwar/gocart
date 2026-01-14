import { useState, useEffect, useCallback } from 'react'
import Product from '../types/product';
import { productService, ApiProduct } from '../services/productService';

const useProducts = () => {
    const [products, setProducts] = useState<Product[]>([])
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    const convertApiProduct = (apiProduct: ApiProduct): Product => ({
        productID: apiProduct.product_id,
        name: apiProduct.name,
        description: apiProduct.description,
        price: apiProduct.price,
        category: apiProduct.category,
        imageUrl: apiProduct.image_url,

    });

    const fetchProducts = useCallback(async () => {
        try {
            setLoading(true);
            const apiProducts = await productService.getAllProducts();
            const convertedProducts = apiProducts.map(convertApiProduct);
            setProducts(convertedProducts);
            setError(null);
        } catch (err) {
            setError('Failed to fetch products. Please try again later.')
            console.error('Error fetching products: ', err)
        } finally {
            setLoading(false);
        }
    }, []);
    // effect to fetch on mount
    useEffect(() => {
        fetchProducts();
    }, [fetchProducts]);
    return {
        products,
        loading,
        error,
        refetch: fetchProducts
    };
};

export default useProducts;