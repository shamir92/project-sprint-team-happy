import { isExists } from "../helpers/assertion.js";
import { combine, generateRandomName, generateRandomPassword, generateTestObjects } from "../helpers/generator.js";
import { testPostJsonAssert } from "../helpers/request.js";

const registerNegativePayloads = (positivePayload) => generateTestObjects({
    nip: { type: "number", notNull: true, min: 1000000000000, max: 999999999999999 },
    name: { type: "string", notNull: true, minLength: 5, maxLength: 50 },
    password: { type: "string", minLength: 5, maxLength: 33, notNull: true }
}, positivePayload)
/**
 * 
 * @param {import("../config.js").Config} config 
 * @param {Object} tags 
 * @returns {import("../types/user.js").ItUser | null}
 */
export function TestRegister(config, nip, tags) {
    const currentRoute = `${config.BASE_URL}/v1/user/it/register`
    const currentFeature = "register"
    /** @type {import("../types/user.js").ItUser} */
    const registerPositivePayload = {
        name: generateRandomName(),
        nip: nip,
        password: generateRandomPassword()
    }

    /** @type {import("../helpers/request.js").RequestAssertResponse} */
    let res;

    if (!config.POSITIVE_CASE) {
        registerNegativePayloads(registerPositivePayload).forEach((payload) => {
            testPostJsonAssert(currentFeature, "invalid payload", currentRoute, payload, {}, {
                ['should return 400']: (res) => res.status === 400,
            }, config, tags);
        });
    }

    res = testPostJsonAssert(currentFeature, "register with correct payload", currentRoute, registerPositivePayload, {}, {
        ['should return 201']: (res) => res.status === 201,
        ['should return have a userId']: (res) => isExists(res, "data.userId"),
        ['should return have a nip']: (res) => isExists(res, "data.nip"),
        ['should return have a name']: (res) => isExists(res, "data.name"),
        ['should return have an accessToken']: (res) => isExists(res, "data.accessToken"),
    }, config, tags);

    if (!config.POSITIVE_CASE) {
        testPostJsonAssert(currentFeature, "register with existing nip", currentRoute, registerPositivePayload, {}, {
            ['should return 409']: (res) => res.status === 409,
        }, config, tags);
    }
    if (res.isSuccess) {
        const user = res.res.json().data
        return combine(registerPositivePayload, user)
    }
    return null
}