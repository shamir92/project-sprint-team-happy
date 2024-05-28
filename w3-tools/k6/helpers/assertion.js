
import { check } from "k6";

/**
 * Asserts the response of a k6 request.
 * @param {import("k6/http").RefinedResponse<ResponseType | undefined>} k6response 
 * @param {string} httpMethod 
 * @param {Object} requestPayload
 * @param {string} featureName
 * @param {import("k6").Checkers} conditions 
 * @param {Config} config
 * @returns {Boolean}
 */
export function assert(k6response, httpMethod, requestPayload, featureName, conditions, config) {
    const checks = {}

    Object.keys(conditions).forEach(testMsg => {
        let condition = conditions[testMsg]
        let testName = featureName + " | " + testMsg
        if (config.DEBUG_ALL) {
            condition = () => {
                const res = conditions[testMsg](k6response)
                console.log(testName + " | assert result:", res);
                return res
            }
        }
        checks[testName] = condition
    });

    if (config.DEBUG_ALL) {
        console.log(featureName + " | request path:", httpMethod, k6response.url);
        console.log(featureName + " | request payload:", requestPayload);
        console.log(featureName + " | response code:", k6response.status);
        console.log(featureName + " | response payload:", k6response.body);
    }

    return check(k6response, checks)
}

/**
 * Checks whether k6 response has the data that the query asks
 * @param {import("k6/http").RefinedResponse<ResponseType | undefined>} v 
 * @param {string} query 
 * @returns {Boolean}
 */
export function isExists(v, query) {
    try {
        const obj = v.json()
        const res = traverseObject(obj, query)
        return !res.includes(undefined)
    } catch (e) {
        return false
    }
}
/**
 * validate ISO date string
 * @param {string} dateString 
 * @returns {Boolean}
 */
export function isValidDate(dateString) {
    const date = new Date(dateString);
    return !isNaN(date.getTime()); // getTime() returns NaN for 'Invalid Date'
}

/**
 * Helper function to check if a string is a valid URL.
 * @param {string} url - The URL to check.
 * @returns {boolean} - Returns true if the URL is valid, otherwise false.
 */
export function isValidUrl(url) {
    // Implement your URL validation logic here
    // This is just a placeholder implementation
    return url.startsWith('http://') || url.startsWith('https://');
}

/**
 * This is a callback function.
 * @callback EqualWithCallback
 * @param {Array} arr - The array parameter.
 * @returns {boolean} - The return value.
 */

/**
 * Checks if the value `v` is equal to the result of traversing the object `v.json()` using the provided `query`.
 * 
 * @param {import("k6/http").RefinedResponse} v - The value to be checked.
 * @param {string} query - The query used to traverse the object.
 * @param {EqualWithCallback} cb - The callback function to be called with the result of the traversal.
 * @returns {boolean} - Returns `true` if the value is equal to the result of the traversal, otherwise `false`.
 */
export function isEqualWith(v, query, cb) {
    try {
        const obj = v.json()
        const res = traverseObject(obj, query)
        return cb(res)
    } catch (e) {
        return false
    }
}

/**
 * Checks whether k6 response has the data that the query asks and matches it
 * @param {import("k6/http").RefinedResponse<ResponseType | undefined >} v 
 * @param {string} query 
 * @param {any} expected 
 * @returns {Boolean}
 */
export function isEqual(v, query, expected) {
    try {
        const obj = v.json()
        const res = traverseObject(obj, query)
        return res.includes(expected)
    } catch (e) {
        return false
    }
}

/**
 * This is a callback function.
 * @callback ConversionCallback
 * @param {any} item - The array parameter.
 * @returns {any} - The return value.
 */

/**
 * Checks if the values in an object's property are ordered in a specified manner.
 * @param {import("k6/http").RefinedResponse} v - The object to be checked.
 * @param {"asc"|"desc"} ordered - The order in which the values should be checked. Can be 'asc' for ascending order or 'desc' for descending order.
 * @param {string} key - The key of the property to be checked.
 * @param {ConversionCallback} conversion - The callback function to convert the values to the desired type.
 * @returns {boolean} - Returns true if the values are ordered as specified, false otherwise.
 */
export function isOrdered(v, key, ordered, conversion) {
    try {
        const obj = v.json();
        const res = traverseObject(obj, key).map(conversion);
        if (ordered === 'asc') {
            return res.every((val, i) => i === 0 || val >= res[i - 1]);
        } else {
            return res.every((val, i) => i === 0 || val <= res[i - 1]);
        }
    } catch (error) {
        return false;
    }
}

/**
 * Checks if the total data in a given object is within a specified range.
 * @param {import("k6/http").RefinedResponse} v - The response object.
 * @param {string} key - The key to traverse in the object.
 * @param {number} min - The minimum number of elements allowed in the range.
 * @param {number} max - The maximum number of elements allowed in the range.
 * @returns {boolean} - Returns true if the total data is within the specified range, false otherwise.
 */
export function isTotalDataInRange(v, key, min, max) {
    try {
        const obj = v.json();
        const res = traverseObject(obj, key);
        return res.length >= min && res.length <= max;
    } catch (error) {
        return false;
    }
}

function flatMap(arr, callback) {
    return arr.reduce((acc, item) => acc.concat(callback(item)), []);
}

/**
 * Traverses an object and retrieves the values based on the provided query.
 *
 * @param {Object} obj - The object to traverse.
 * @param {string} query - The query to specify the path of the values to retrieve.
 * @returns {Array} - An array of values retrieved from the object based on the query.
 */
function traverseObject(obj, query) {
    const keys = query.split('.');
    let result = [obj];

    for (const key of keys) {
        if (key.endsWith('[]')) {
            const arrayKey = key.slice(0, -2);
            result = flatMap(result, item => item[arrayKey]);
        } else {
            result = result.map(item => item[key]);
        }
    }

    return result;
}