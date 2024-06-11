
/**
 * @typedef { "Beverage" | "Food" | "Snack" | "Condiments" | "Additons" } MerchantItemCategory
 */

/**
 * @typedef {Object} MerchantItem
 * @property {string} itemId 
 * @property {string} name
 * @property {MerchantItemCategory} productCategory 
 * @property {price} price 
 * @property {string} imageUrl
 */

/**
 * 
 * @returns {MerchantItemCategory}
 */
export function generateRandomMerchantItemCategory() {
    const categories = ["Beverage", "Food", "Snack", "Condiments", "Additions"];
    return categories[Math.floor(Math.random() * categories.length)];
}
