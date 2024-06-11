import { IsAdmin } from "../../entity/admin.js";
import { IsUser } from "../../entity/user.js";
import { isEqual, isEqualWith, isExists, isTotalDataInRange, isValidDate } from "../../helpers/assertion.js";
import { combine } from "../../helpers/generator.js";
import { testGetAssert } from "../../helpers/request.js";

/**
 * @param {import("../../entity/user.js").User} user
 * @param {import("../../entity/config").Config} config 
 */

export function GetOrderTest(user, config, tags) {
    if (!IsUser(user)) {
        return;
    }

    const featureName = "User Get Order Test";
    const route = config.BASE_URL + "/users/orders";

    const headers = {
        Authorization: `Bearer ${user.token}`
    }

    if (!config.POSITIVE_CASE) {
        testGetAssert(
            "empty auth",
            featureName,
            route, {}, {}, {
            ['should return 401']: (v) => v.status === 401
        }, config, tags)
    }

    const positiveTestCases = {
        ['should return 200']: (v) => v.status === 200,
        ['should have the correct total data based on pagination']: (v) => isTotalDataInRange(v, '', 1, 5),
        ['should have merchant.merchantId']: (v) => isArrayExists(v, 'orders[].merchant.merchantId'),
        ['should have merchant.name']: (v) => isArrayExists(v, 'orders[].merchant.name'),
        ['should have merchant.imageUrl']: (v) => isArrayExists(v, 'orders[].merchant.imageUrl'),
        ['should have merchant.merchantCategory']: (v) => isArrayExists(v, 'orders[].merchant.merchantCategory'),
        ['should have merchant.createdAt']: (v) => isArrayEqualWith(v, 'orders[].merchant.createdAt', (a) => a.every(b => isValidDate(b))),
        ['should have merchant.location.lat']: (v) => isArrayExists(v, 'orders[].merchant.location.lat'),
        ['should have merchant.location.long']: (v) => isArrayExists(v, 'orders[].merchant.location.long'),
        ['should have items.itemId']: (v) => isArrayExists(v, 'orders[].items[].itemId'),
        ['should have items.name']: (v) => isArrayExists(v, 'orders[].items[].name'),
        ['should have items.imageUrl']: (v) => isArrayExists(v, 'orders[].items[].imageUrl'),
        ['should have items.productCategory']: (v) => isArrayExists(v, 'orders[].items[].productCategory'),
        ['should have items.createdAt']: (v) => isArrayEqualWith(v, 'orders[].items[].createdAt', (a) => a.every(b => isValidDate(b))),
    }

    testGetAssert("no param", featureName, route, {}, headers, positiveTestCases, config, tags)

    testGetAssert("with name=a param", featureName, route, { name: "a" }, headers, combine(positiveTestCases, {
        ['should have name with "a" in it']: (v) => {
            const hasMerchantName = isArrayExists(v, 'orders[].merchant.name')
            if (hasMerchantName) {
                v.json().forEach(orderDetail => {
                    orderDetail.orders.forEach(order => {
                        if (!order.merchant.name.toLowerCase().includes('a')) {
                            if (!order.items.every(a => a.name.toLowerCase().includes('a'))) {
                                return false
                            }
                        }
                    })
                })
            } else {
                return false
            }
            return true
        }
    },), config, tags)

    testGetAssert("with merchantCategory=BoothKiosk param", featureName, route, { merchantCategory: "BoothKiosk" }, headers, combine(positiveTestCases, {
        ['should have "BoothKiosk" category in it']: (v) => isArrayEqual(v, 'orders[].merchant.merchantCategory', "BoothKiosk")
    }), config, tags)

    const paginationRes = testGetAssert("pagination", featureName, route, { limit: 2, offset: 0 }, headers, combine(positiveTestCases, {
        ['should have the correct total data based on pagination']: (v) => isTotalDataInRange(v, '', 1, 2),
    }), config, tags)
    if (!config.LOAD_TEST && paginationRes.isSuccess) {
        testGetAssert("pagination offset", featureName, route, { limit: 2, offset: 2 }, headers, combine(positiveTestCases, {
            ['should have the correct total data based on pagination']: (v) => isTotalDataInRange(v, '', 1, 2),
            ['should have different data from offset 0']: (res) => {
                try {
                    return res.json().data.every(e => {
                        return paginationRes.res.json().data.every(a => a.merchantId !== e.merchantId)
                    })
                } catch (err) {
                    return err
                }
            },
        }), config, tags)
    }
}

function isArrayEqual(v, query, expected) {
    try {
        const a = v.json();
        return isEqual({
            json() {
                return { a: a }
            }
        }, `a[].${query}`, expected);
    } catch (error) {
        return false
    }
}

function isArrayEqualWith(v, query, expected) {
    try {
        const a = v.json();
        return isEqualWith({
            json() {
                return { a: a }
            }
        }, `a[].${query}`, expected);
    } catch (error) {
        return false
    }
}

function isArrayExists(v, query) {
    try {
        const a = v.json();
        return isExists({
            json() {
                return { a: a }
            }
        }, `a[].${query}`);
    } catch (error) {
        return false
    }
}