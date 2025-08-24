import Product from "../../types/product"
import ProductCard from "./ProductCard"
import './ProductList.css';
interface ProductListProps {
    products: Product[];
    loading: boolean;
}

const ProductList = ({ products, loading }: ProductListProps) => {

    const skeletonRows = Array.from({ length: 5 });

    return (
        <table className="data-table">
            <thead>
                <tr>
                    <th>Name</th>
                    <th>Price</th>
                    <th>Category</th>
                    <th>Description</th>
                </tr>
            </thead>
            <tbody>
                {loading ? (
                    skeletonRows.map((_, i) => (
                        <tr key={i}>
                            <td colSpan={4}>
                                <div className="skeleton-row" />
                            </td>
                        </tr>
                    ))
                ) : products.length === 0 ? (
                    <tr>
                        <td colSpan={4} style={{
                            textAlign: 'center', padding: '16px'
                        }}>
                            No Products Available
                        </td>
                    </tr>
                ) : (
                    products.map(product => (
                        <ProductCard key={product.productID} product={product} />
                    ))
                )}
            </tbody>
        </table>
    );
}

export default ProductList;