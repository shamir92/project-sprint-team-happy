import { fail } from "k6";
import { combine, generateTestObjects } from "../helpers/generator.js";
import { testPostJsonAssert } from "../helpers/request.js";
import { isUserValidLogin } from "../types/user.js";
import { isExists } from "../helpers/assertion.js";

const registerNegativePayloads = (positivePayload) => generateTestObjects({
    nip: { type: "number", notNull: true },
    password: { type: "string", minLength: 5, maxLength: 33, notNull: true }
}, positivePayload)
/**
 * 
 * @param {import("../config.js").Config} config 
 * @param {Object} tags 
 * @param {import("../types/user.js").ItUser} user
 * @param {number} newNip
 * @returns {import("../types/user.js").ItUser | null}
 */
export function TestLogin(user, config, newNip, tags) {
    const currentRoute = `${config.BASE_URL}/v1/user/it/login`
    const currentFeature = "login"
    if (!isUserValidLogin(user)) {
        fail(`${currentFeature} Invalid user object`)
    }

    /** @type {import("../types/user.js").ItUser} */
    const registerPositivePayload = {
        name: user.name,
        nip: user.nip,
        password: user.password
    }

    /** @type {import("../helpers/request.js").RequestAssertResponse} */
    let res;

    if (!config.POSITIVE_CASE) {
        registerNegativePayloads(registerPositivePayload).forEach((payload) => {
            testPostJsonAssert(currentFeature, "invalid payload", currentRoute, payload, {}, {
                ['should return 400']: (res) => res.status === 400,
            }, config, tags);
        });
        testPostJsonAssert(currentFeature, "not existing nip", currentRoute,
            combine(registerPositivePayload, {
                nip: newNip
            }), {}, {
            ['should return 404']: (res) => res.status === 404,
        }, config, tags);
    }

    res = testPostJsonAssert(currentFeature, "login with correct payload", currentRoute, registerPositivePayload, {}, {
        ['should return 200']: (res) => res.status === 200,
        ['should return have a userId']: (res) => isExists(res, "data.userId"),
        ['should return have a nip']: (res) => isExists(res, "data.nip"),
        ['should return have a name']: (res) => isExists(res, "data.name"),
        ['should return have an accessToken']: (res) => isExists(res, "data.accessToken"),
    }, config, tags);

    if (res.isSuccess) {
        const user = res.res.json().data
        return combine(registerPositivePayload, user)
    }
    return null

}