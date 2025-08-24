import Product from "../../types/product";

interface ProductCardProps {
    product: Product
}

const usdFormatter = new Intl.NumberFormat('en-us', {
    style: 'currency',
    currency: 'USD'
});

const ProductCard = ({ product }: ProductCardProps) => {
    return (
        <tr className="product-card">
            <td>{product.name}</td>
            <td className="price">{usdFormatter.format(product.price)}</td>
            <td className="category">{product.category}</td>
            <td className="description">{product.description}</td>

        </tr>
    );
}

export default ProductCard;