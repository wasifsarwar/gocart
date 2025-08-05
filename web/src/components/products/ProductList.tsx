import Product from "../../types/product"
import ProductCard from "./ProductCard"
interface ProductListProps {
    products: Product[]
}

const ProductList = ({ products }: ProductListProps) => {

    if (products.length === 0) {
        return <div>No products available</div>;
    }
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
                {products.map(product => (
                    <ProductCard key={product.productID} product={product} />
                ))}
            </tbody>
        </table>
    );
}

export default ProductList;