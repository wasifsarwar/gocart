import React, { createContext, useCallback, useContext, useEffect, useMemo, useState, ReactNode } from 'react';

const STORAGE_KEY = 'gocart.favorites.v1';

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

function safeParseIds(raw: string | null): string[] {
    if (!raw) return [];
    try {
        const parsed = JSON.parse(raw);
        if (!Array.isArray(parsed)) return [];
        return parsed.map(String).filter(Boolean);
    } catch {
        return [];
    }
}

interface FavoritesContextType {
    favoriteIds: string[];
    isFavorite: (productId: string) => boolean;
    toggleFavorite: (productId: string) => void;
    clearFavorites: () => void;
}

const FavoritesContext = createContext<FavoritesContextType | undefined>(undefined);

export const FavoritesProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
    const [favoriteIds, setFavoriteIds] = useState<string[]>(() => safeParseIds(safeGetItem(STORAGE_KEY)));

    useEffect(() => {
        safeSetItem(STORAGE_KEY, JSON.stringify(favoriteIds));
    }, [favoriteIds]);

    useEffect(() => {
        const onStorage = (e: StorageEvent) => {
            if (e.key === STORAGE_KEY) {
                setFavoriteIds(safeParseIds(e.newValue));
            }
        };
        window.addEventListener('storage', onStorage);
        return () => window.removeEventListener('storage', onStorage);
    }, []);

    const favoriteIdSet = useMemo(() => new Set(favoriteIds), [favoriteIds]);

    const isFavorite = useCallback((productId: string) => favoriteIdSet.has(productId), [favoriteIdSet]);

    const toggleFavorite = useCallback((productId: string) => {
        setFavoriteIds((prev) => {
            // Keep a stable "most recently favorited first" ordering.
            if (prev.includes(productId)) {
                return prev.filter((id) => id !== productId);
            }
            return [productId, ...prev];
        });
    }, []);

    const clearFavorites = useCallback(() => {
        setFavoriteIds([]);
    }, []);

    return (
        <FavoritesContext.Provider value={{ favoriteIds, isFavorite, toggleFavorite, clearFavorites }}>
            {children}
        </FavoritesContext.Provider>
    );
};

export const useFavorites = () => {
    const context = useContext(FavoritesContext);
    if (context === undefined) {
        throw new Error('useFavorites must be used within a FavoritesProvider');
    }
    return context;
};

