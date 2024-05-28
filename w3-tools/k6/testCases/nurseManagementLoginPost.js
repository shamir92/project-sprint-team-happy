import { fail } from "k6";
import { combine, generateTestObjects } from "../helpers/generator.js";
import { testPostJsonAssert } from "../helpers/request.js";
import { isUserValidLogin } from "../types/user.js";

const nurseAccesstNegativePayload = (positivePayload) => generateTestObjects({
    nip: { type: "number", notNull: true },
    password: { type: "string", minLength: 5, maxLength: 33, notNull: true }
}, positivePayload)
/**
 * 
 * @param {import("../config.js").Config} config 
 * @param {Object} tags 
 * @param {import("../types/user.js").NurseUserWithoutLogin} userNurseToLogin
 * @returns {import("../types/user.js").NurseUser | null}
 */
export function TestNurseManagementLoginPost(config, userNurseToLogin, tags) {
    const currentFeature = "nurse management login"
    const currentRoute = `${config.BASE_URL}/v1/user/nurse/login`
    if (isUserValidLogin(userNurseToLogin)) {
        fail(`${currentFeature} Invalid user object`)
    }
    const positivePayload = {
        nip: userNurseToLogin.nip,
        password: userNurseToLogin.password
    }

    if (!config.POSITIVE_CASE) {
        nurseAccesstNegativePayload(positivePayload).forEach((payload) => {
            testPostJsonAssert(currentFeature, "invalid payload", currentRoute, payload, {}, {
                ['should return 400']: (res) => res.status === 400,
            }, config, tags);
        });
    }

    const res = testPostJsonAssert(currentFeature, "login with correct payload", currentRoute, positivePayload, {}, {
        'should return 200': (res) => res.status === 200
    }, config, tags);
    if (!res.isSuccess) {
        fail(`${currentFeature} login nurse error: ${res.res.body}`)
    }
    const user = res.res.json().data
    return combine(positivePayload, user)
}