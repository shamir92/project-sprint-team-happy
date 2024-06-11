import { IsUser } from "../../entity/user.js";
import { isExists } from "../../helpers/assertion.js";
import { generateTestObjects } from "../../helpers/generator.js";
import { testPostJsonAssert } from "../../helpers/request.js";

/**
 * @param {import("../../entity/user.js").User} user
 * @param {import("../../entity/config.js").Config} config 
 * @param {string} estimateId
 */
export function OrderTest(user, estimateId, config, tags) {
    if (!IsUser(user)) {
        return;
    }

    const featureName = "User Order Test";
    const route = config.BASE_URL + "/users/orders";

    const positivePayload = {
        calculatedEstimateId: estimateId
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
            calculatedEstimateId: { type: "string", notNull: true },
        }, positivePayload)
        testObjects.forEach(payload => {
            testPostJsonAssert("invalid payload", featureName, route, payload, headers, {
                ['should return 400']: (res) => res.status === 400,
            }, config, tags);
        });
    }

    const res = testPostJsonAssert("valid payload", featureName, route, positivePayload, headers, {
        ['should return 201']: (v) => v.status === 201,
        ['should have orderId']: (v) => isExists(v, 'orderId')
    }, config, tags)

    if (res.isSuccess) {
        return res.res.json().orderId
    }
    return null
}