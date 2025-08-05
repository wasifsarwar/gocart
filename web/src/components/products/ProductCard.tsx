import Product from "../../types/product";

interface ProductCardProps {
    product: Product
}

const ProductCard = ({ product }: ProductCardProps) => {
    return (
        <tr className="product-card">
            <td>{product.name}</td>
            <td className="price">${product.price.toFixed(2)}</td>
            <td className="category">{product.category}</td>
            <td className="description">{product.description}</td>

        </tr>
    );
}

export default ProductCard;