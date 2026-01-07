const API_URL = process.env.REACT_APP_API_URL || "http://localhost:8080";

export function resolveImageUrl(imageUrl?: string): string | null {
    if (!imageUrl) return null;
    if (imageUrl.startsWith('http://') || imageUrl.startsWith('https://') || imageUrl.startsWith('data:')) {
        return imageUrl;
    }
    if (imageUrl.startsWith('/')) {
        return `${API_URL}${imageUrl}`;
    }
    return `${API_URL}/${imageUrl}`;
}

