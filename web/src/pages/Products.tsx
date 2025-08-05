import ProductList from "../components/products/ProductList";
import useProducts from "../hooks/useProducts";
import Navigation from "../components/navigation/Navigation";

const Products = () => {
    const { products, loading, error } = useProducts();
    return (
        <div className="products-page">
            <Navigation title="GoCart Products" />
            {loading && <p>Loading Products</p>}
            {error && <p style={{ color: 'red' }} >{error}</p>}
            {!loading && !error && <ProductList products={products} />}
        </div >
    );
};

export default Products;