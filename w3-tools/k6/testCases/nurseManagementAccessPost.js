import { fail } from "k6";
import { combine, generateRandomNumber, generateRandomPassword, generateTestObjects } from "../helpers/generator.js";
import { testGetAssert, testPostJsonAssert } from "../helpers/request.js";
import { isItUserValid } from "../types/user.js";

const nurseAccesstNegativePayload = (positivePayload) => generateTestObjects({
    password: { type: "string", minLength: 5, maxLength: 33, notNull: true }
}, positivePayload)
/**
 * 
 * @param {import("../config.js").Config} config 
 * @param {Object} tags 
 * @param {import("../types/user.js").ItUser} user
 * @returns {import("../types/user.js").NurseUserWithoutLogin | null}
 */
export function TestNurseManagementAccessPost(config, userIt, tags) {
    if (!isItUserValid(userIt)) {
        fail(`${currentFeature} Invalid user object`)
    }
    const currentFeature = "nurse management access post"
    const headers = {
        Authorization: `Bearer ${userIt.accessToken}`
    }

    const getNurseRes = testGetAssert(currentFeature, "get all users with nurse role", `${config.BASE_URL}/v1/user`, { role: 'nurse' }, headers, {
        ['should return 200']: (res) => res.status === 200,
    }, config, tags);
    if (!getNurseRes.isSuccess) {
        fail(currentFeature, "get all users with nurse role", getNurseRes.res, config, tags)
    }

    const nurses = getNurseRes.res.json().data
    const nurseToTry = nurses[generateRandomNumber(0, nurses.length - 1)]

    const currentRoute = `${config.BASE_URL}/v1/user/nurse/${nurseToTry.userId}/access`
    const positivePayload = {
        password: generateRandomPassword()
    }

    if (!config.POSITIVE_CASE) {
        nurseAccesstNegativePayload(positivePayload).forEach((payload) => {
            testPostJsonAssert(currentFeature, "invalid payload", currentRoute, payload, headers, {
                ['should return 400']: (res) => res.status === 400,
            }, config, tags);
        });
    }

    const res = testPostJsonAssert(currentFeature, "give nurse access", currentRoute, positivePayload, headers, {
        'should return 200': (res) => res.status === 200
    }, config, tags);

    if (!res.isSuccess) {
        fail(`${currentFeature} give nurse access error: ${res.res.body}`)
    }

    return combine(nurseToTry, positivePayload)
}