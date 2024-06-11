import { combine, generateRandomEmail, generateRandomPassword, generateRandomUsername, generateTestObjects } from "../../helpers/generator.js";
import { testPostJsonAssert } from "../../helpers/request.js";

/**
 * @param {import("../../entity/config").Config} config 
 * @returns {null | import("../../entity/user").User}
 */
export function UserRegisterTest(config, tags) {
    const featureName = "User Register";
    const route = config.BASE_URL + "/users/register";

    /** @type {import("../../entity/user").User} */
    const positivePayload = {
        username: generateRandomUsername(),
        password: generateRandomPassword(),
        email: generateRandomEmail()
    }

    if (!config.POSITIVE_CASE) {
        testPostJsonAssert(
            "empty body",
            featureName,
            route, {}, {}, {
            ['should return 400']: (v) => v.status === 400
        }, config, tags)

        const testObjects = generateTestObjects({
            username: { type: "string", notNull: true, minLength: 5, maxLength: 30 },
            password: { type: "string", notNull: true, minLength: 5, maxLength: 30 },
            email: { type: "string", notNull: true, isEmail: true },
        }, positivePayload)
        testObjects.forEach(payload => {
            testPostJsonAssert("invalid payload", featureName, route, payload, {}, {
                ['should return 400']: (res) => res.status === 400,
            }, config, tags);
        });
    }

    const res = testPostJsonAssert("valid payload", featureName, route, positivePayload, {}, {
        ['should return 201']: (v) => v.status === 201
    }, config, tags)

    if (res.isSuccess) {
        testPostJsonAssert("register twice", featureName, route, positivePayload, {}, {
            ['should return 409']: (v) => v.status === 409
        }, config, tags)
    }

    if (res.isSuccess) {
        return combine(positivePayload, {
            token: res.res.json().token
        })
    }
}
