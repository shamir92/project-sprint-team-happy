import { IsAdmin } from "../../entity/admin.js";
import { isEqual, isEqualWith, isExists, isOrdered, isTotalDataInRange, isValidDate } from "../../helpers/assertion.js";
import { combine, generateRandomImageUrl, generateRandomName, generateRandomNumber } from "../../helpers/generator.js";
import { testPostJsonAssert, testGetAssert } from "../../helpers/request.js";

/**
 * @param {import("../../entity/admin.js").Admin} user
 * @param {import("../../entity/config").Config} config 
 * @param {import("../../entity/merchant").Merchant} merchant
 * @param {number} addItemCount
 * 
 * @returns {import("../../entity/merchantItem").MerchantItem[] | null}
 */
export function MerchantItemGetTest(user, merchant, addItemCount, config, tags) {
    if (!IsAdmin(user)) {
        return;
    }

    const featureName = "Merchant Item Get";
    const route = config.BASE_URL + `/admin/merchants/${merchant.merchantId}/items`;

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
        ['should have itemId']: (v) => isExists(v, 'data[].itemId'),
        ['should have name']: (v) => isExists(v, 'data[].name'),
        ['should have imageUrl']: (v) => isExists(v, 'data[].imageUrl'),
        ['should have productCategory']: (v) => isExists(v, 'data[].productCategory'),
        ['should have createdAt with correct format']: (v) => isEqualWith(v, 'data[].createdAt', (a) => a.every(b => isValidDate(b))),
        ['should have return ordered correctly']: (res) => isOrdered(res, 'data[].createdAt', 'desc', (v) => new Date(v)),
        ['should have meta.limit']: (v) => isExists(v, 'meta.limit'),
        ['should have meta.offset']: (v) => isExists(v, 'meta.offset'),
        ['should have meta.total']: (v) => isExists(v, 'meta.total'),
    }

    /** @type {import("../../entity/merchantItem").MerchantItem[]} */
    const addedMerchantItems = []
    if (!config.LOAD_TEST) {
        for (let i = 0; i < addItemCount; i++) {
            const merchantItemToAdd = {
                name: generateRandomName() + "a",
                productCategory: "Snack",
                imageUrl: generateRandomImageUrl(),
                price: generateRandomNumber(1000, 10000),
            }
            const postRes = testPostJsonAssert("add merchant item for search", featureName, route, merchantItemToAdd, headers, {}, config, tags)
            if (postRes.isSuccess) {
                addedMerchantItems.push(combine(merchantItemToAdd, { itemId: postRes.res.json().itemId }))
            }
        }
    }

    testGetAssert("no param", featureName, route, {}, headers, positiveTestCases, config, tags)

    testGetAssert("with name=a param", featureName, route, { name: "a" }, headers, combine(positiveTestCases, {
        ['should have name with "a" in it']: (v) => isEqualWith(v, 'data[].name', (a) => a.every(b => b.toLowerCase().includes('a')))
    },), config, tags)

    testGetAssert("with productCategory=Snack param", featureName, route, { productCategory: "Snack" }, headers, combine(positiveTestCases, {
        ['should have "Snack" category in it']: (v) => isEqual(v, 'data[].productCategory', "Snack")
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
                        return paginationRes.res.json().data.every(a => a.itemId !== e.itemId)
                    })
                } catch (err) {
                    return err
                }
            },
        }), config, tags)
    }

    return addedMerchantItems
}