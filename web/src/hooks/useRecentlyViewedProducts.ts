import { useCallback, useEffect, useMemo, useState } from 'react';
import Product from '../types/product';

const STORAGE_KEY = 'gocart.recentlyViewedProducts.v1';
const MAX_ITEMS = 8;

function safeGetItem(key: string): string | null {
    try {
        return localStorage.getItem(key);
    } catch {
        return null;
    }
}

function safeSetItem(key: string, value: string) {
    try {
        localStorage.setItem(key, value);
    } catch {
        // Ignore storage write failures (e.g. disabled storage, quota exceeded)
    }
}

function safeRemoveItem(key: string) {
    try {
        localStorage.removeItem(key);
    } catch {
        // Ignore storage removal failures
    }
}

function safeParseProducts(raw: string | null): Product[] {
    if (!raw) return [];
    try {
        const parsed = JSON.parse(raw);
        if (!Array.isArray(parsed)) return [];

        // Basic shape validation + coercion
        return parsed
            .filter((p) => p && typeof p === 'object')
            .map((p: any) => ({
                productID: String(p.productID),
                name: String(p.name ?? ''),
                description: String(p.description ?? ''),
                price: Number(p.price ?? 0),
                category: String(p.category ?? '')
            }))
            .filter((p) => p.productID && p.name);
    } catch {
        return [];
    }
}

function readRecentlyViewed(): Product[] {
    return safeParseProducts(safeGetItem(STORAGE_KEY));
}

function writeRecentlyViewed(products: Product[]) {
    safeSetItem(STORAGE_KEY, JSON.stringify(products));
}

function addToRecentlyViewed(product: Product): Product[] {
    const current = readRecentlyViewed();
    const deduped = current.filter((p) => p.productID !== product.productID);
    const next = [product, ...deduped].slice(0, MAX_ITEMS);
    writeRecentlyViewed(next);
    return next;
}

export default function useRecentlyViewedProducts() {
    const [recentlyViewed, setRecentlyViewed] = useState<Product[]>(() => readRecentlyViewed());

    useEffect(() => {
        const onStorage = (e: StorageEvent) => {
            if (e.key === STORAGE_KEY) {
                setRecentlyViewed(readRecentlyViewed());
            }
        };
        window.addEventListener('storage', onStorage);
        return () => window.removeEventListener('storage', onStorage);
    }, []);

    const addViewedProduct = useCallback((product: Product) => {
        setRecentlyViewed(addToRecentlyViewed(product));
    }, []);

    const clearRecentlyViewed = useCallback(() => {
        safeRemoveItem(STORAGE_KEY);
        setRecentlyViewed([]);
    }, []);

    const maxItems = useMemo(() => MAX_ITEMS, []);

    return {
        recentlyViewed,
        addViewedProduct,
        clearRecentlyViewed,
        maxItems
    };
}

