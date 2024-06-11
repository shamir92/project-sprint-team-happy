import { IsAdmin } from "../../entity/admin.js";
import { isEqual, isEqualWith, isExists, isOrdered, isTotalDataInRange, isValidDate } from "../../helpers/assertion.js";
import { combine, generateRandomImageUrl, generateRandomName } from "../../helpers/generator.js";
import { testPostJsonAssert, testGetAssert } from "../../helpers/request.js";

/**
 * @param {import("../../entity/admin.js").Admin} user
 * @param {import("../../entity/config").Config} config 
 * @param {import("../../entity/merchant").RawMerchant[]} rawMerchants
 * @returns {import("../../entity/merchant").Merchant[]}
 */
export function MerchantGetTest(user, rawMerchants, config, tags) {
    if (!IsAdmin(user)) {
        return;
    }

    const featureName = "Merchant Get";
    const route = config.BASE_URL + "/admin/merchants";

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
        ['should have the correct total data based on pagination']: (v) => isTotalDataInRange(v, 'data[]', 1, 5),
        ['should have merchantId']: (v) => isExists(v, 'data[].merchantId'),
        ['should have name']: (v) => isExists(v, 'data[].name'),
        ['should have imageUrl']: (v) => isExists(v, 'data[].imageUrl'),
        ['should have merchantCategory']: (v) => isExists(v, 'data[].merchantCategory'),
        ['should have createdAt with correct format']: (v) => isEqualWith(v, 'data[].createdAt', (a) => a.every(b => isValidDate(b))),
        ['should have return ordered correctly']: (res) => isOrdered(res, 'data[].createdAt', 'desc', (v) => new Date(v)),
        ['should have location.lat']: (v) => isExists(v, 'data[].location.lat'),
        ['should have location.long']: (v) => isExists(v, 'data[].location.long'),
        ['should have meta.limit']: (v) => isExists(v, 'meta.limit'),
        ['should have meta.offset']: (v) => isExists(v, 'meta.offset'),
        ['should have meta.total']: (v) => isExists(v, 'meta.total'),
    }

    /** @type {import("../../entity/merchant").Merchant[]} */
    const addedMerchants = []
    if (!config.LOAD_TEST) {
        rawMerchants.forEach(merchant => {
            const merchantToAdd = {
                name: generateRandomName() + "a",
                merchantCategory: "BoothKiosk",
                imageUrl: generateRandomImageUrl(),
                location: {
                    lat: merchant.location.lat,
                    long: merchant.location.long
                }
            }
            const postRes = testPostJsonAssert("add merchant for search", featureName, route, merchantToAdd, headers, {}, config, tags)
            if (postRes.isSuccess) {
                addedMerchants.push(combine(merchantToAdd, { merchantId: postRes.res.json().merchantId, pregeneratedId: merchant.pregeneratedId }))
            }
        });

    }

    testGetAssert("no param", featureName, route, {}, headers, positiveTestCases, config, tags)

    testGetAssert("with name=a param", featureName, route, { name: "a" }, headers, combine(positiveTestCases, {
        ['should have name with "a" in it']: (v) => isEqualWith(v, 'data[].name', (a) => a.every(b => b.toLowerCase().includes('a')))
    },), config, tags)

    testGetAssert("with merchantCategory=BoothKiosk param", featureName, route, { merchantCategory: "BoothKiosk" }, headers, combine(positiveTestCases, {
        ['should have "BoothKiosk" category in it']: (v) => isEqual(v, 'data[].merchantCategory', "BoothKiosk")
    }), config, tags)

    testGetAssert("with createdAt=asc param", featureName, route, { createdAt: "asc" }, headers, combine(positiveTestCases, {
        ['should have return ordered correctly']: (v) => isOrdered(v, 'data[].createdAt', "asc", (a) => new Date(a))
    }), config, tags)

    const paginationRes = testGetAssert("pagination", featureName, route, { limit: 2, offset: 0 }, headers, combine(positiveTestCases, {
        ['should have the correct total data based on pagination']: (v) => isTotalDataInRange(v, 'data[]', 1, 2),
    }), config, tags)
    if (!config.LOAD_TEST && paginationRes.isSuccess) {
        testGetAssert("pagination offset", featureName, route, { limit: 2, offset: 2 }, headers, combine(positiveTestCases, {
            ['should have the correct total data based on pagination']: (v) => isTotalDataInRange(v, 'data[]', 1, 2),
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
    return addedMerchants
}