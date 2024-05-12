/* eslint-disable no-loss-of-precision */
import { isValidUrl } from "../helpers/assertion.js";
import { generateRandomNumber, generateRandomDescription, generateRandomImageUrl, generateRandomName, MaxInt } from "../helpers/generator.js";

/**
 * Represents a product.
 * @typedef {Object} Product
 * @property {string} id - The id of the product
 * @property {string} name - The name of the product. Must not be null and must have a minimum length of 1 and a maximum length of 30.
 * @property {string} sku - The SKU (stock keeping unit) of the product. Must not be null and must have a minimum length of 1 and a maximum length of 30.
 * @property {"Clothing"|"Accessories"|"Footwear"|"Beverages"} category - The category of the product. Must not be null and must be one of the following: "Clothing", "Accessories", "Footwear", "Beverages".
 * @property {string} imageUrl - The URL of the product's image. Must not be null and must be a valid URL.
 * @property {string} notes - Additional notes about the product. Must not be null and must have a minimum length of 1 and a maximum length of 200.
 * @property {number} price - The price of the product. Must not be null and must be a minimum of 1.
 * @property {number} stock - The stock quantity of the product. Must not be null and must be a minimum of 0 and a maximum of 100000.
 * @property {string} location - The location of the product. Must not be null and must have a minimum length of 1 and a maximum length of 200.
 * @property {boolean} isAvailable - Indicates whether the product is available. Must not be null.
 */


/**
 * Generates a random category from a predefined list.
 * @returns {"Clothing"|"Accessories"|"Footwear"|"Beverages"} The randomly selected category.
 */
export function getRandomCategory() {
    const categories = ["Clothing", "Accessories", "Footwear", "Beverages"];
    const randomIndex = Math.floor(Math.random() * categories.length);
    return categories[randomIndex];
}

export function generateProduct() {
    return {
        name: generateRandomName(),
        sku: `${generateRandomNumber(10000000000, MaxInt)}`,
        category: getRandomCategory(),
        imageUrl: generateRandomImageUrl(),
        notes: generateRandomDescription(200),
        isAvailable: true,
        stock: generateRandomNumber(0, 100),
        price: generateRandomNumber(1, 1000),
        location: generateRandomDescription()
    }
}


/**
 * validates product
 * @param {Product} product 
 * @returns {Boolean}
 */
export function validateProduct(product) {
    if (!product) {
        return false;
    }

    if (!product.id || typeof product.id !== 'string') {
        return false;
    }

    if (!product.name || typeof product.name !== 'string' || product.name.length < 1 || product.name.length > 30) {
        return false;
    }

    if (!product.sku || typeof product.sku !== 'string' || product.sku.length < 1 || product.sku.length > 30) {
        return false;
    }

    if (!product.category || typeof product.category !== 'string' || !['Clothing', 'Accessories', 'Footwear', 'Beverages'].includes(product.category)) {
        return false;
    }

    if (!product.imageUrl || typeof product.imageUrl !== 'string' || !isValidUrl(product.imageUrl)) {
        return false;
    }

    if (!product.notes || typeof product.notes !== 'string' || product.notes.length < 1 || product.notes.length > 200) {
        return false;
    }

    if (!product.price || typeof product.price !== 'number' || product.price < 1) {
        return false;
    }

    if (!product.stock || typeof product.stock !== 'number' || product.stock < 0 || product.stock > 100000) {
        return false;
    }

    if (!product.location || typeof product.location !== 'string' || product.location.length < 1 || product.location.length > 200) {
        return false;
    }

    if (typeof product.isAvailable !== 'boolean') {
        return false;
    }

    return true;
}

