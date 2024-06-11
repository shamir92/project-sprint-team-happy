/**
 * @typedef {"SmallRestaurant" | "MediumRestaurant" | "LargeRestaurant" | "MerchandiseRestaurant" | "BoothKiosk" | "ConvenienceStore" } MerchantCategory
 */

/**
 * @typedef {Object} Location
 * @property {int} lat 
 * @property {int} long 
 */

/**
 * Represents a Merchant.
 * @typedef {Object} RawMerchant
 * @property {string} pregeneratedId - The pregenerated identifier of the merchant.
 * @property {LocationPoint} location - The location point of the merchant.
 */

/**
 * Represents a Merchant.
 * @typedef {Object} Merchant
 * @property {string} merchantId - The unique identifier of the merchant.
 * @property {string} pregeneratedId - The pregenerated identifier of the merchant.
 * @property {string} name
 * @property {string} imageUrl
 * @property {MerchantCategory} merchantCategory
 * @property {LocationPoint} location - The location point of the merchant.
 */

/**
 * Represents a MerchantNearestRecord.
 * @typedef {Object} MerchantNearestRecord
 * @property {LocationPoint} startingPoint - The starting location point.
 * @property {Object.<string, Merchant>} merchants - A map where the key is an int64 and the value is a Merchant object.
 */

/**
 * Represents a collection of MerchantNearestRecords.
 * @typedef {Object} MerchantNearestRecordZone
 * @property {MerchantNearestRecord[]} records - An array of MerchantNearestRecord objects.
 */

/**
 * Represents a message containing all merchant nearest records.
 * @typedef {Object} AllMerchantNearestRecord
 * @property {MerchantNearestRecordZone[]} zones - An array of MerchantNearestRecordZone objects.
 */

/**
 * Represents a LocationPoint.
 * @typedef {Object} LocationPoint
 * @property {number} lat - The latitude of the location.
 * @property {number} long - The longitude of the location.
 */

/**
 * Represents generated routes.
 * @typedef {Object} GeneratedRoutes
 * @property {Object.<string, Merchant>} generatedRoutes - A map where the key is an int64 and the value is a Merchant object.
 * @property {LocationPoint} startingPoint - The starting location point.
 * @property {number} totalDistance - The total distance of the route.
 * @property {number} totalDuration - The total duration of the route.
 * @property {number} startingIndex - The starting index.
 */

/**
 * Represents a collection of generated routes.
 * @typedef {Object} AllGeneratedRoutes
 * @property {RouteZone[]} zone - An array of RouteZone objects.
 */

/**
 * Represents a RouteZone.
 * @typedef {Object} RouteZone
 * @property {GeneratedRoutes[]} routes - An array of GeneratedRoutes objects.
 */

/**
 * Represents pregenerated merchants.
 * @typedef {Object} PregeneratedMerchant
 * @property {RawMerchant[]} merchant - An array of Merchant objects.
 */

/**
 * Represents an AssignMerchant message.
 * @typedef {Object} AssignMerchant
 * @property {string} pregeneratedId - The pregenerated identifier of the merchant.
 * @property {string} merchantId - The unique identifier of the merchant.
 */

/**
 * 
 * @returns {MerchantCategory}
 */
export function generateRandomMerchantCategory() {
    const categories = ["SmallRestaurant", "MediumRestaurant", "LargeRestaurant", "MerchandiseRestaurant", "BoothKiosk", "ConvenienceStore"];
    return categories[Math.floor(Math.random() * categories.length)];
}

/**
 * 
 * @param {Object} obj 
 * @returns {boolean}
 */
export function isMerchant(obj) {
    return obj.merchantId
        && obj.name
        && obj.imageUrl
        && obj.merchantCategory
        && obj.location;
}