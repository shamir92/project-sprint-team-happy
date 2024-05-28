import { fail } from "k6";
import { isExists } from "../helpers/assertion.js";
import { testPostJsonAssert } from "../helpers/request.js";
import { isItUserValid } from "../types/user.js";

const { generateTestObjects, generateRandomImageUrl, generateRandomName } = require("../helpers/generator.js");

const nurseManagemenetNegativePayload = (positivePayload) => generateTestObjects({
    nip: { type: "number", notNull: true, min: 1000000000000, max: 999999999999999 },
    name: { type: "string", notNull: true, minLength: 5, maxLength: 50 },
    identityCardScanImg: { type: "string", notNull: true, isUrl: true }
}, positivePayload)
/**
 * 
 * @param {import("../config.js").Config} config 
 * @param {Object} tags 
 * @param {import("../types/user.js").ItUser} user
 * @param {number} nip
 * @returns {import("../types/user.js").NurseUserWithoutAccess | null}
 */
export function TestNurseManagementPost(config, user, nip, tags) {
    const currentRoute = `${config.BASE_URL}/v1/user/nurse/register`
    const currentFeature = "nurse management post"
    if (!isItUserValid(user)) {
        fail(`${currentFeature} Invalid user object`)
    }

    /** @type {import("../helpers/request.js").RequestAssertResponse} */
    let res;
    const headers = {
        Authorization: `Bearer ${user.accessToken}`
    }

    const nurseManagementPositivePayload = {
        name: generateRandomName(),
        nip: nip,
        identityCardScanImg: generateRandomImageUrl()
    }

    if (!config.POSITIVE_CASE) {
        testPostJsonAssert(currentFeature, "no header", currentRoute, {}, {}, {
            ['should return 401']: (res) => res.status === 401,
        }, config, tags);

        testPostJsonAssert(currentFeature, "no payload", currentRoute, {}, headers, {
            ['should return 400']: (res) => res.status === 400,
        }, config, tags);

        nurseManagemenetNegativePayload(nurseManagementPositivePayload).forEach((payload) => {
            testPostJsonAssert(currentFeature, "invalid payload", currentRoute, payload, headers, {
                ['should return 400']: (res) => res.status === 400,
            }, config, tags);
        });
    }

    res = testPostJsonAssert(currentFeature, "register with correct payload", currentRoute, nurseManagementPositivePayload, headers, {
        ['should return 201']: (res) => res.status === 201,
        ['should return have a userId']: (res) => isExists(res, "data.userId"),
        ['should return have a nip']: (res) => isExists(res, "data.nip"),
        ['should return have a name']: (res) => isExists(res, "data.name"),
    }, config, tags);


    if (res.isSuccess) {
        return nurseManagementPositivePayload
    }
    fail(`${currentFeature} register nurse error: ${res.res.body}`)
}
