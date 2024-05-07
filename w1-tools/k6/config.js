/* eslint-disable no-undef */
module.exports = {
    VERBOSE: __ENV.VERBOSE ? true : false,
    DEBUG_ALL: __ENV.DEBUG_ALL ? true : false,
    POSITIVE_CASE: __ENV.ONLY_POSITIVE ? true : false,
    GACHA: __ENV.GACHA ? true : false,
    LOAD_TEST: __ENV.LOAD_TEST ? true : false
}