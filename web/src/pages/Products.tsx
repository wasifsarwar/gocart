import ProductList from "../components/ProductList";
import useProducts from "../hooks/useProducts";

const Products = () => {
    const { products, loading, error } = useProducts();
    return (
        <div className="products-page">
            <h1>GoCart Products</h1>
            {loading && <p>Loading Products</p>}
            {error && <p style={{ color: 'red' }} >{error}</p>}
            {!loading && !error && <ProductList products={products} />}
        </div >
    );
};

export default Products;