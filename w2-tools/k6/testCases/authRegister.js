import { isExists } from "../helpers/assertion.js";
import { generateRandomName, generateRandomNumber, generateRandomPassword, generateTestObjects } from "../helpers/generator.js";
import { testPostJsonAssert } from "../helpers/request.js";
import { generateInternationalCallingCode } from "../types/user.js";

const registerNegativePayloads = (positivePayload) => generateTestObjects({
    phoneNumber: { type: "string", notNull: true, minLength: 10, maxLength: 16 },
    name: { type: "string", notNull: true, minLength: 5, maxLength: 50 },
    password: { type: "string", minLength: 5, maxLength: 15, notNull: true }
}, positivePayload)
/**
 * 
 * @param {import("../config.js").Config} config 
 * @param {Object} tags 
 * @returns {import("../types/user.js").User | null}
 */
export function TestRegister(config, tags) {
    const currentRoute = `${config.BASE_URL}/v1/staff/register`
    const currentFeature = "Register"
    /** @type {import("../types/user.js").User} */
    const registerPositivePayload = {
        name: generateRandomName(),
        phoneNumber: `+${generateInternationalCallingCode()}${generateRandomNumber(1000000, 99999999)}`,
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
        ['should return have a phoneNumber']: (res) => isExists(res, "data.phoneNumber"),
        ['should return have a name']: (res) => isExists(res, "data.name"),
        ['should return have an accessToken']: (res) => isExists(res, "data.accessToken"),
    }, config, tags);

    if (!config.POSITIVE_CASE) {
        testPostJsonAssert(currentFeature, "register with existing phone number", currentRoute, registerPositivePayload, {}, {
            ['should return 409']: (res) => res.status === 409,
        }, config, tags);
    }
    if (res.isSuccess) {
        return Object.assign(registerPositivePayload, { accessToken: res.res.json().data.accessToken })
    }
    return null
}