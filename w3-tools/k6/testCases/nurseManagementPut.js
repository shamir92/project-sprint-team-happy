import { fail } from "k6";
import { testGetAssert, testPutJsonAssert } from "../helpers/request.js";
import { isItUserValid } from "../types/user.js";
import { isEqual, isExists } from "../helpers/assertion.js";
const { generateTestObjects, generateRandomName, generateRandomNumber, combine } = require("../helpers/generator.js");

const nurseManagemenetNegativePayload = (positivePayload) => generateTestObjects({
    nip: { type: "number", notNull: true, min: 1000000000000, max: 999999999999999 },
    name: { type: "string", notNull: true, minLength: 5, maxLength: 50 },
}, positivePayload)
/**
 * 
 * @param {import("../config.js").Config} config 
 * @param {Object} tags 
 * @param {number} newNip
 * @param {import("../types/user.js").ItUser} user
 */
export function TestNurseManagementPut(config, user, newNip, tags) {
    const currentFeature = "nurse management put"
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
    if (!getItRes.isSuccess) {
        fail(currentFeature, "get all users with it role", getNurseRes.res, config, tags)
    }

    /** @type {import("../types/user.js").NurseUser[]}*/
    const nurses = getNurseRes.res.json().data
    const nurseToTry = nurses[nurses.length - 1]
    const nurseToEdit = nurses[generateRandomNumber(0, nurses.length / 2 - 1)]
    const originalRoute = `${config.BASE_URL}/v1/user/nurse/`
    const currentRoute = `${originalRoute}${nurseToEdit.userId}`

    /** @type {import("../types/user.js").ItUser[]} */
    const itUsers = getItRes.res.json().data
    const itUserToTry = itUsers[generateRandomNumber(0, itUsers.length - 1)]

    const nurseManagementPositivePayload = {
        name: generateRandomName(),
        nip: newNip
    }

    if (!config.POSITIVE_CASE) {
        testPutJsonAssert(currentFeature, "no header", currentRoute, {}, {}, {
            ['should return 401']: (res) => res.status === 401,
        }, config, tags);

        testPutJsonAssert(currentFeature, "no path value", originalRoute, {}, headers, {
            ['should return 404']: (res) => res.status === 404,
        }, config, tags);

        testPutJsonAssert(currentFeature, "no payload", currentRoute, {}, headers, {
            ['should return 400']: (res) => res.status === 400,
        }, config, tags);

        nurseManagemenetNegativePayload(nurseManagementPositivePayload).forEach((payload) => {
            testPutJsonAssert(currentFeature, "invalid payload", currentRoute, payload, headers, {
                ['should return 400']: (res) => res.status === 400,
            }, config, tags);
        });

        testPutJsonAssert(currentFeature, "editing an it", `${originalRoute}${itUserToTry.userId}`,
            nurseManagementPositivePayload, headers, {
            ['should return 404']: (res) => res.status === 404,
        }, config, tags);

        testPutJsonAssert(currentFeature, "edit with the existing nip", currentRoute,
            combine(nurseManagementPositivePayload, { nip: nurseToTry.nip }), headers, {
            ['should return 409']: (res) => res.status === 409,
        }, config, tags);
    }

    testPutJsonAssert(currentFeature, "editing nurse", currentRoute, nurseManagementPositivePayload, headers, {
        ['should return 200']: (res) => res.status === 200,
    }, config, tags);

    testGetAssert(currentFeature, "get user by userId after edit", `${config.BASE_URL}/v1/user`, { userId: nurseToEdit.userId }, headers, {
        ['should return 200']: (res) => res.status === 200,
        ['should all have a userId']: (res) => isExists(res, "data[].userId"),
        ['should have the same nip after edit']: (res) => isEqual(res, "data[].nip", nurseManagementPositivePayload.nip),
        ['should have the same name after edit']: (res) => isEqual(res, "data[].name", nurseManagementPositivePayload.name),
    }, config, tags);

    return null
}