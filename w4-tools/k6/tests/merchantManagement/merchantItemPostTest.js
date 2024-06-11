import { IsAdmin } from "../../entity/admin.js";
import { isMerchant } from "../../entity/merchant.js";
import { generateRandomMerchantItemCategory } from "../../entity/merchantItem.js";
import { isExists } from "../../helpers/assertion.js";
import { combine, generateRandomImageUrl, generateRandomName, generateRandomNumber, generateTestObjects } from "../../helpers/generator.js";
import { testPostJsonAssert } from "../../helpers/request.js";

/**
 * @param {import("../../entity/admin.js").Admin} user
 * @param {import("../../entity/config").Config} config 
 * @param {import("../../entity/merchant").Merchant} merchant
 * @returns {import("../../entity/merchantItem").MerchantItem | null}
 */
export function MerchantItemPostTest(user, merchant, config, tags) {
    if (!IsAdmin(user)
        && !isMerchant(merchant)) {
        return;
    }

    const featureName = "Merchant Item Post";
    const route = config.BASE_URL + `/admin/merchants/${merchant.merchantId}/items`;

    /** @type {import("../../entity/merchantItem").MerchantItemCategory} */
    const positivePayload = {
        name: generateRandomName(),
        productCategory: generateRandomMerchantItemCategory(),
        imageUrl: generateRandomImageUrl(),
        price: generateRandomNumber(1, 1000000),
    }

    const headers = {
        Authorization: `Bearer ${user.token}`
    }

    if (!config.POSITIVE_CASE) {
        testPostJsonAssert(
            "empty auth",
            featureName,
            route, {}, {}, {
            ['should return 401']: (v) => v.status === 401
        }, config, tags)

        testPostJsonAssert(
            "empty body",
            featureName,
            route, {}, headers, {
            ['should return 400']: (v) => v.status === 400
        }, config, tags)

        const testObjects = generateTestObjects({
            name: { type: "string", notNull: true, minLength: 2, maxLength: 30 },
            productCategory: {
                type: "string", notNull: true, enum: ["Beverage", "Food", "Snack", "Condiments", "Additions"]
            },
            imageUrl: { type: "string", notNull: true, isUrl: true },
            price: { type: "number", notNull: true, min: 1 }
        }, positivePayload)
        testObjects.forEach(payload => {
            testPostJsonAssert("invalid payload", featureName, route, payload, headers, {
                ['should return 400']: (res) => res.status === 400,
            }, config, tags);
        });
    }

    const res = testPostJsonAssert("valid payload", featureName, route, positivePayload, headers, {
        ['should return 201']: (v) => v.status === 201,
        ['should have itemId']: (v) => isExists(v, 'itemId'),
    }, config, tags)

    if (res.isSuccess) {
        return combine(positivePayload, {
            itemId: res.res.json().itemId,
        })
    }
    return null
}
