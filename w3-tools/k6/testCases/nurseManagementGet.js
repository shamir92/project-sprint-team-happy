import { fail } from "k6";
import { isEqualWith, isExists, isOrdered, isTotalDataInRange, isValidDate } from "../helpers/assertion.js";
import { testGetAssert } from "../helpers/request.js";
import { isItUserValid } from "../types/user.js";
/**
 * 
 * @param {import("../config.js").Config} config 
 * @param {Object} tags 
 * @param {import("../types/user.js").ItUser} user
 */
export function TestNurseManagementGet(config, user, tags) {
    const currentRoute = `${config.BASE_URL}/v1/user`
    const currentFeature = "nurse management get"
    if (!isItUserValid(user)) {
        fail(`${currentFeature} Invalid user object`)
    }
    /** @type {import("../helpers/request.js").RequestAssertResponse} */
    const headers = {
        Authorization: `Bearer ${user.accessToken}`
    }

    if (!config.POSITIVE_CASE) {
        testGetAssert(currentFeature, "no header", currentRoute, {}, {}, {
            ['should return 401']: (res) => res.status === 401,
        }, config, tags);
    }

    testGetAssert(currentFeature, "get all users", currentRoute, {}, headers, {
        ['should return 200']: (res) => res.status === 200,
        ['should all have a userId']: (res) => isExists(res, "data[].userId"),
        ['should all have a nip']: (res) => isExists(res, "data[].nip"),
        ['should all have a name']: (res) => isExists(res, "data[].name"),
        ['should all have a createdAt']: (res) => isExists(res, "data[].createdAt"),
        ['should not have more than 5 result']: (res) => isTotalDataInRange(res, 'data[]', 1, 5),
        ['should have createdAt and format should be date']: (res) => isEqualWith(res, 'data[].createdAt', (v) => v.every(a => isValidDate(a))),
    }, config, tags);

    testGetAssert(currentFeature, "get all users with name", currentRoute, { name: 'a' }, headers, {
        ['should return 200']: (res) => res.status === 200,
        ['should all have a userId']: (res) => isExists(res, "data[].userId"),
        ['should all have a nip']: (res) => isExists(res, "data[].nip"),
        ['should all have a name']: (res) => isExists(res, "data[].name"),
        ['should all have a createdAt']: (res) => isExists(res, "data[].createdAt"),
        ['should not have more than 5 result']: (res) => isTotalDataInRange(res, 'data[]', 1, 5),
        ['should have createdAt and format should be date']: (res) => isEqualWith(res, 'data[].createdAt', (v) => v.every(a => isValidDate(a))),
        ['should have result with "a" in it']: (res) => isEqualWith(res, 'data[].name', (v) => v.every(a => a.toLowerCase().includes('a'))),
    }, config, tags);

    testGetAssert(currentFeature, "get all users with nip", currentRoute, { nip: `${user.nip}`.substring(0, 4) }, headers, {
        ['should return 200']: (res) => res.status === 200,
        ['should all have a userId']: (res) => isExists(res, "data[].userId"),
        ['should all have a nip']: (res) => isExists(res, "data[].nip"),
        ['should all have a name']: (res) => isExists(res, "data[].name"),
        ['should all have a createdAt']: (res) => isExists(res, "data[].createdAt"),
        ['should not have more than 5 result']: (res) => isTotalDataInRange(res, 'data[]', 1, 5),
        ['should have createdAt and format should be date']: (res) => isEqualWith(res, 'data[].createdAt', (v) => v.every(a => isValidDate(a))),
        ['should have result with four digit of searched nip in it']: (res) => isEqualWith(res, 'data[].nip', (v) => v.every(a => `${a}`.includes(`${user.nip}`.substring(0, 4)))),
    }, config, tags);

    testGetAssert(currentFeature, "get all users with it role", currentRoute, { role: 'it' }, headers, {
        ['should return 200']: (res) => res.status === 200,
        ['should all have a userId']: (res) => isExists(res, "data[].userId"),
        ['should all have a nip']: (res) => isExists(res, "data[].nip"),
        ['should all have a name']: (res) => isExists(res, "data[].name"),
        ['should all have a createdAt']: (res) => isExists(res, "data[].createdAt"),
        ['should have createdAt and format should be date']: (res) => isEqualWith(res, 'data[].createdAt', (v) => v.every(a => isValidDate(a))),
        ['should not have more than 5 result']: (res) => isTotalDataInRange(res, 'data[]', 1, 5),
        ['should have all nip that starts with 615']: (res) => isEqualWith(res, 'data[].nip', (v) => v.every(a => `${a}`.startsWith('615'))),
    }, config, tags);

    testGetAssert(currentFeature, "get all users with nurse role", currentRoute, { role: 'nurse' }, headers, {
        ['should return 200']: (res) => res.status === 200,
        ['should all have a userId']: (res) => isExists(res, "data[].userId"),
        ['should all have a nip']: (res) => isExists(res, "data[].nip"),
        ['should all have a name']: (res) => isExists(res, "data[].name"),
        ['should all have a createdAt']: (res) => isExists(res, "data[].createdAt"),
        ['should have createdAt and format should be date']: (res) => isEqualWith(res, 'data[].createdAt', (v) => v.every(a => isValidDate(a))),
        ['should not have more than 5 result']: (res) => isTotalDataInRange(res, 'data[]', 1, 5),
        ['should have all nip that starts with 303']: (res) => isEqualWith(res, 'data[].nip', (v) => v.every(a => `${a}`.startsWith('303'))),
    }, config, tags);

    testGetAssert(currentFeature, "get all users with createdAt asc", currentRoute, { createdAt: 'asc' }, headers, {
        ['should return 200']: (res) => res.status === 200,
        ['should all have a userId']: (res) => isExists(res, "data[].userId"),
        ['should all have a nip']: (res) => isExists(res, "data[].nip"),
        ['should all have a name']: (res) => isExists(res, "data[].name"),
        ['should all have a createdAt']: (res) => isExists(res, "data[].createdAt"),
        ['should have createdAt and format should be date']: (res) => isEqualWith(res, 'data[].createdAt', (v) => v.every(a => isValidDate(a))),
        ['should not have more than 5 result']: (res) => isTotalDataInRange(res, 'data[]', 1, 5),
        ['should have return ordered by oldest first']: (res) => isOrdered(res, 'data[].createdAt', 'asc', (v) => new Date(v)),
    }, config, tags);

    const paginationRes = testGetAssert(currentFeature, "get all users with limit", currentRoute, { limit: 2 }, headers, {
        ['should return 200']: (res) => res.status === 200,
        ['should all have a userId']: (res) => isExists(res, "data[].userId"),
        ['should all have a nip']: (res) => isExists(res, "data[].nip"),
        ['should all have a name']: (res) => isExists(res, "data[].name"),
        ['should all have a createdAt']: (res) => isExists(res, "data[].createdAt"),
        ['should have createdAt and format should be date']: (res) => isEqualWith(res, 'data[].createdAt', (v) => v.every(a => isValidDate(a))),
        ['should not have more than 2 result']: (res) => isTotalDataInRange(res, 'data[]', 1, 2),
    }, config, tags);

    if (paginationRes.isSuccess && !config.LOAD_TEST) {
        testGetAssert(currentFeature, "get all users with limit and offset", currentRoute, { limit: 2, offset: 2 }, headers, {
            ['should return 200']: (res) => res.status === 200,
            ['should all have a userId']: (res) => isExists(res, "data[].userId"),
            ['should all have a nip']: (res) => isExists(res, "data[].nip"),
            ['should all have a name']: (res) => isExists(res, "data[].name"),
            ['should all have a createdAt']: (res) => isExists(res, "data[].createdAt"),
            ['should have createdAt and format should be date']: (res) => isEqualWith(res, 'data[].createdAt', (v) => v.every(a => isValidDate(a))),
            ['should not have more than 2 result']: (res) => isTotalDataInRange(res, 'data[]', 1, 2),
            ['should have different data from offset 0']: (res) => {
                try {
                    return res.json().data.every(e => {
                        return paginationRes.res.json().data.every(a => a.userId !== e.userId)
                    })
                } catch (error) {
                    return false
                }
            },
        }, config, tags);
    }

    return null
}