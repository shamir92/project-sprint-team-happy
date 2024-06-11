import { IsUser } from "../../entity/user.js";
import { isEqual, isEqualWith, isExists, isTotalDataInRange, isValidDate } from "../../helpers/assertion.js";
import { combine } from "../../helpers/generator.js";
import { testGetAssert } from "../../helpers/request.js";

/**
 * @param {import("../../entity/user.js").User} user
 * @param {import("../../entity/merchant.js").MerchantNearestRecord} nearestRecord
 * @param {import("../../entity/config.js").Config} config 
 */
export function GetNearbyMerchantTest(user, nearestRecord, config, tags) {
    if (!IsUser(user)) {
        return;
    }

    const featureName = "User Get Nearby Merchant";
    const route = config.BASE_URL + `/merchants/nearby/${nearestRecord.startingPoint.lat},${nearestRecord.startingPoint.long}`;

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
        testGetAssert(
            "invalid lat long",
            featureName,
            `${config.BASE_URL}/merchants/nearby/a,b`, {}, headers, {
            ['should return 400']: (v) => v.status === 400
        }, config, tags)
    }


    if (!config.LOAD_TEST) {
        console.log("five nearest location according to the test case")
        console.log("starting point: ", nearestRecord.startingPoint)
    }
    const nearestMerchantIdToVerify = []
    for (let i = 0; i < 5; i++) {
        const currentNearestMerchant = nearestRecord.merchants[`${i}`].merchantId
        nearestMerchantIdToVerify.push(currentNearestMerchant)
        if (!config.LOAD_TEST) {
            console.log(`nearest merchant ${i}:`, nearestRecord.merchants[`${i}`].merchantId)
            console.log("location:", nearestRecord.merchants[`${i}`].location,)
        }
    }

    const positiveTestCases = {
        ['should return 200']: (v) => v.status === 200,
        ['should have the correct total data based on pagination']: (v) => isTotalDataInRange(v, 'data[]', 1, 5),
        ['should have merchant.merchantId']: (v) => isExists(v, 'data[].merchant.merchantId'),
        ['should have merchant.name']: (v) => isExists(v, 'data[].merchant.name'),
        ['should have merchant.merchantCategory']: (v) => isExists(v, 'data[].merchant.merchantCategory'),
        ['should have merchant.imageUrl']: (v) => isExists(v, 'data[].merchant.imageUrl'),
        ['should have merchant.location.lat']: (v) => isExists(v, 'data[].merchant.location.lat'),
        ['should have merchant.location.long']: (v) => isExists(v, 'data[].merchant.location.long'),
        ['should have merchant.createdAt with correct format']: (v) => isEqualWith(v, 'data[].merchant.createdAt', (a) => a.every(b => isValidDate(b))),
        ['should have items[].itemId']: (v) => isExists(v, 'data[].items[].itemId'),
        ['should have items[].name']: (v) => isExists(v, 'data[].items[].name'),
        ['should have items[].productCategory']: (v) => isExists(v, 'data[].items[].productCategory'),
        ['should have items[].price']: (v) => isExists(v, 'data[].items[].price'),
        ['should have items[].imageUrl']: (v) => isExists(v, 'data[].items[].imageUrl'),
        ['should have items[].createdAt with correct format']: (v) => isEqualWith(v, 'data[].items[].createdAt', (a) => a.every(b => isValidDate(b))),
        ['should have meta.limit']: (v) => isExists(v, 'meta.limit'),
        ['should have meta.offset']: (v) => isExists(v, 'meta.offset'),
        ['should have meta.total']: (v) => isExists(v, 'meta.total'),
    }


    testGetAssert("no param", featureName, route, {}, headers, combine(positiveTestCases, {
        ['should have the correct nearest merchant']: (v) => isEqualWith(v, 'data[].merchant.merchantId', (e => e.every(a => nearestMerchantIdToVerify.includes(a))))
    }), config, tags)
    console.log(nearestMerchantIdToVerify)

    testGetAssert("with name=a param", featureName, route, { name: "a" }, headers, combine(positiveTestCases, {
        ['should have name with "a" in it']: (v) => {
            const hasMerchantName = isExists(v, 'data[].merchant.name')
            if (hasMerchantName) {
                v.json().data.forEach(e => {
                    if (!e.merchant.name.toLowerCase().includes('a')) {
                        if (!e.items.every(a => a.name.toLowerCase().includes('a'))) {
                            return false
                        }
                    }
                })
            } else {
                return false
            }
            return true
        }
    },), config, tags)

    testGetAssert("with merchantCategory=BoothKiosk param", featureName, route, { merchantCategory: "BoothKiosk" }, headers, combine(positiveTestCases, {
        ['should have "BoothKiosk" category in it']: (v) => isEqual(v, 'data[].merchant.merchantCategory', "BoothKiosk")
    }), config, tags)
}