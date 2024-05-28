import { fail } from "k6";
import { testDeleteAssert, testGetAssert } from "../helpers/request.js";
import { isTotalDataInRange } from "../helpers/assertion.js";
import { isItUserValid } from "../types/user.js";
const { generateRandomNumber } = require("../helpers/generator.js");

/**
 * 
 * @param {import("../config.js").Config} config 
 * @param {Object} tags 
 * @param {import("../types/user.js").ItUser} user
 */
export function TestNurseManagementDelete(config, user, tags) {
    const currentFeature = "nurse management delete"
    if (!isItUserValid(user)) {
        fail(`${currentFeature} Invalid user object`)
    }
    const headers = {
        Authorization: `Bearer ${user.accessToken}`
    }

    const getNurseRes = testGetAssert(currentFeature, "get all users with nurse role", `${config.BASE_URL}/v1/user`, { role: 'nurse' }, headers, {
        ['should return 200']: (res) => res.status === 200,
    }, config, tags);
    if (!getNurseRes.isSuccess) {
        fail(currentFeature, "get all users with nurse role", getNurseRes.res, config, tags)
    }

    const getItRes = testGetAssert(currentFeature, "get all users with it role", `${config.BASE_URL}/v1/user`, { role: 'it' }, headers, {
        ['should return 200']: (res) => res.status === 200,
    }, config, tags);
    if (!getNurseRes.isSuccess) {
        fail(currentFeature, "get all users with nurse role", getNurseRes.res, config, tags)
    }

    /** @type {import("../types/user.js").NurseUser[]}*/
    const nurses = getNurseRes.res.json().data
    const nurseToRemove = nurses[generateRandomNumber(0, nurses.length - 1)]
    const originalRoute = `${config.BASE_URL}/v1/user/nurse/`
    const currentRoute = `${originalRoute}${nurseToRemove.userId}`

    /** @type {import("../types/user.js").ItUser[]} */
    const itUsers = getItRes.res.json().data
    const itUserToTry = itUsers[generateRandomNumber(0, itUsers.length - 1)]

    if (!config.POSITIVE_CASE) {
        testDeleteAssert(currentFeature, "no header", currentRoute, {}, {}, {
            ['should return 401']: (res) => res.status === 401,
        }, config, tags);

        testDeleteAssert(currentFeature, "no path value", originalRoute, {}, headers, {
            ['should return 404']: (res) => res.status === 404,
        }, config, tags);

        testDeleteAssert(currentFeature, "deleting an it", `${originalRoute}${itUserToTry.userId}`,
            {}, headers, {
            ['should return 404']: (res) => res.status === 404,
        }, config, tags);
    }

    testDeleteAssert(currentFeature, "deleting nurse", currentRoute, {}, headers, {
        ['should return 200']: (res) => res.status === 200,
    }, config, tags);

    testGetAssert(currentFeature, "get user by userId after edit", `${config.BASE_URL}/v1/user`, { userId: nurseToRemove.userId }, headers, {
        ['should return 200']: (res) => res.status === 200,
        ['should have no data']: (res) => isTotalDataInRange(res, "data[]", 0, 0),
    }, config, tags);

    return null
}