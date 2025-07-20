
```sql
+---------------+          1    N            +-----------------+
|    Order      | <------------------------> |   OrderItem     |
+---------------+                            +-----------------+
| OrderID (PK)  |                            | OrderItemID (PK)|
| UserID (FK) * |                            | OrderID (FK)   |
| Status        |                            | ProductID (FK)**|
| TotalAmount   |                            | Quantity       |
| CreatedAt     |                            | Price          |
| UpdatedAt     |                            +-----------------+
+---------------+

*  FK to User Service (users table: UserID)
** FK to Product Service (products table: ProductID)

Relationships:
- Order 1:N OrderItem (one order has many items)
- Integrates with: User Service (validate user) and Product Service (check stock/price
```
