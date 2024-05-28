/* eslint-disable no-undef */
/**
 * Configuration object for the EniQiloStoreTestCases.
 * @typedef {Object} Config
 * @property {string} BASE_URL - The base URL for the test cases.
 * @property {boolean} DEBUG_ALL - Flag indicating whether to enable debug mode for all test cases.
 * @property {boolean} POSITIVE_CASE - Flag indicating whether to run only positive test cases.
 * @property {boolean} LOAD_TEST - Flag indicating whether to run load test.
 */

/**
 * Configuration settings 
 * @type {Config}
 */
const config = {
    BASE_URL: __ENV.BASE_URL ? __ENV.BASE_URL : "http://localhost:8080",
    DEBUG_ALL: __ENV.DEBUG_ALL ? true : false,
    POSITIVE_CASE: __ENV.ONLY_POSITIVE ? true : false,
    LOAD_TEST: __ENV.LOAD_TEST ? true : false
}

module.exports = {
    config
}