import http from "k6/http";
import { assert } from "./assertion.js";
import { generateParamFromObj } from "./generator.js";

/**
 * @typedef {Object} RequestAssertResponse
 * @property {import("k6/http").RefinedResponse} res - k6 http response.
 * @property {Boolean} isSuccess - Whether the assertion was successful.
 */

/**
 * Sends a Get request with JSON data to the specified route.
 * @param {string} route - The route to send the request to.
 * @param {object} params - The params that will be parsed into the URL.
 * @param {string[]} options - Additional options for the request. 
 *                             Available options: `"noContentType"`
 * @returns {import("k6/http").RefinedResponse} - k6 http response.
 */
export function testGet(route, params, headersObj, tags = {}) {
    const queryParams = generateParamFromObj(params);
    const modifiedRoute = route + "?" + queryParams;

    return http.get(modifiedRoute, { headers: headersObj, tags });
}

/**
 * Sends a Get request with JSON data to the specified route and asserts the response.
 * @param {string} featureName 
 * @param {string} currentTestName 
 * @param {string} route 
 * @param {Object} params 
 * @param {Object} headersObj 
 * @param {import("k6").Checkers} expectedCase 
 * @param {import("../config.js").Config} config 
 * @param {Object} tags
 * @returns {RequestAssertResponse} - k6 http response. 
 */
export function testGetAssert(featureName, currentTestName, route, params, headersObj, expectedCase, config, tags) {
    const res = testGet(route, params, headersObj, tags);
    const isSuccess = assert(res, "GET", generateParamFromObj(params), `${featureName} | ${currentTestName}`, expectedCase, config);
    return {
        isSuccess, res
    }
}

/**
 * Sends a POST request with JSON data to the specified route.
 * @param {string} route - The route to send the request to.
 * @param {object} body - The JSON data to send in the request body.
 * @param {object} headersObj - External headers other than `Content-Type`
 * @param {string[]} options - Additional options for the request. 
 *                             Available options: `"noContentType"`, `"plainBody"`
 * @returns {import("k6/http").RefinedResponse} - k6 http response.
 */
export function testPostJson(route, body, headersObj, tags = {}, options = [],) {
    const headers = options.includes("noContentType") ? Object.assign({}, headersObj) : Object.assign({ "Content-Type": "application/json" }, headersObj)
    const parsedBody = options.includes("plainBody") ? body : JSON.stringify(body);
    return http.post(route, parsedBody, { headers, tags });
}

/**
 * Sends a POST request with JSON data to the specified route and asserts the response.
 * @param {string} featureName 
 * @param {string} currentTestName 
 * @param {string} route 
 * @param {Object} body 
 * @param {Object} headersObj 
 * @param {import("k6").Checkers} expectedCase 
 * @param {import("../config.js").Config} config
 * @param {Object} tags 
 * @returns {RequestAssertResponse} - k6 http response.
 */
export function testPostJsonAssert(featureName, currentTestName, route, body, headersObj, expectedCase, config, tags) {
    const res = testPostJson(route, body, headersObj, tags);
    const isSuccess = assert(res, "POST", body, `${featureName} | ${currentTestName}`, expectedCase, config);
    return {
        isSuccess, res
    }
}

/**
 * Sends a PUT request with JSON data to the specified route.
 * @param {string} route - The route to send the request to.
 * @param {object} body - The JSON data to send in the request body.
 * @param {string[]} options - Additional options for the request. 
 *                             Available options: `"noContentType"`, `"plainBody"`
 * @param {object} headersObj - External headers other than `Content-Type`
 * @returns {import("k6/http").RefinedResponse} - k6 http response.
 */
export function testPutJson(route, body, headersObj, tags = {}, options = []) {
    const headers = options.includes("noContentType") ? Object.assign({}, headersObj) : Object.assign({ "Content-Type": "application/json" }, headersObj)
    const parsedBody = options.includes("plainBody") ? body : JSON.stringify(body);

    return http.put(route, parsedBody, { headers, tags });
}

/**
 * Sends a PUT request with JSON data to the specified route and asserts the response.
 * @param {string} featureName 
 * @param {string} currentTestName
 * @param {string} route 
 * @param {Object} body 
 * @param {Object} headersObj 
 * @param {import("k6").Checkers} expectedCase 
 * @param {import("../config.js").Config} config
 * @param {Object} tags
 * @returns {RequestAssertResponse} - k6 http response. 
 */
export function testPutJsonAssert(featureName, currentTestName, route, body, headersObj, expectedCase, config, tags) {
    const res = testPutJson(route, body, headersObj, tags);
    const isSuccess = assert(res, "PUT", body, `${featureName} | ${currentTestName}`, expectedCase, config);
    return {
        isSuccess, res
    }
}

/**
 * Sends a DELETE request to the specified route.
 * @param {string} route - The route to send the request to.
 * @param {object} params - The params that will be parsed into the URL.
 * @param {object} headersObj - External headers other than `Content-Type`.
 * @returns {import("k6/http").RefinedResponse} - k6 http response.
 */
export function testDelete(route, params, headersObj, tags = {}) {
    const queryParams = Object.entries(params)
        .map(([key, value]) => `${encodeURIComponent(key)}=${encodeURIComponent(value)}`)
        .join('&');
    const modifiedRoute = route + "?" + queryParams;
    const headers = Object.assign({}, headersObj);

    return http.del(modifiedRoute, null, { headers, tags });
}

/**
 * Sends a DELETE request with JSON data to the specified route and asserts the response.
 * @param {string} featureName 
 * @param {string} currentTestName 
 * @param {string} params 
 * @param {Object} body 
 * @param {Object} headersObj 
 * @param {import("k6").Checkers} expectedCase 
 * @param {import("../config.js").Config} config
 * @param {Object} tags
 * @returns {RequestAssertResponse} - k6 http response. 
 */
export function testDeleteAssert(featureName, currentTestName, params, body, headersObj, expectedCase, config, tags) {
    const res = testDelete(params, body, headersObj, tags);
    const isSuccess = assert(res, "DELETE", body, `${featureName} | ${currentTestName}`, expectedCase, config);
    return {
        isSuccess, res
    }
}