import { IsAdmin } from "../../entity/admin.js";
import { generateRandomMerchantCategory } from "../../entity/merchant.js";
import { isExists } from "../../helpers/assertion.js";
import { combine, generateRandomImageUrl, generateRandomName, generateTestObjects } from "../../helpers/generator.js";
import { testPostJsonAssert } from "../../helpers/request.js";

/**
 * @param {import("../../entity/admin.js").Admin} user
 * @param {import("../../entity/config").Config} config 
 * @param {import("../../entity/merchant").RawMerchant} merchantToAdd
 * @returns {import("../../entity/merchant").Merchant | null}
 */
export function MerchantPostTest(user, merchantToAdd, config, tags) {
    if (!IsAdmin(user)) {
        return;
    }

    const featureName = "Merchant Post";
    const route = config.BASE_URL + "/admin/merchants";

    /** @type {import("../../entity/merchant").Merchant} */
    const positivePayload = {
        name: generateRandomName(),
        merchantCategory: generateRandomMerchantCategory(),
        imageUrl: generateRandomImageUrl(),
        location: {
            lat: merchantToAdd.location.lat,
            long: merchantToAdd.location.long
        }
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
            merchantCategory: {
                type: "string", notNull: true, enum: [
                    "SmallRestaurant", "MediumRestaurant", "LargeRestaurant", "MerchandiseRestaurant", "BoothKiosk", "ConvenienceStore"]
            },
            imageUrl: { type: "string", notNull: true, isUrl: true },
            location: {
                type: 'object', properties: {
                    lat: { type: "number", notNull: true },
                    long: { type: "number", notNull: true }
                }, notNull: true
            }
        }, positivePayload)
        testObjects.forEach(payload => {
            testPostJsonAssert("invalid payload", featureName, route, payload, headers, {
                ['should return 400']: (res) => res.status === 400,
            }, config, tags);
        });
    }

    const res = testPostJsonAssert("valid payload", featureName, route, positivePayload, headers, {
        ['should return 201']: (v) => v.status === 201,
        ['should have merchantId']: (v) => isExists(v, 'merchantId')
    }, config, tags)

    if (res.isSuccess) {
        return combine(positivePayload, {
            pregeneratedId: merchantToAdd.pregeneratedId,
            merchantId: res.res.json().merchantId
        })
    }
    return null
}