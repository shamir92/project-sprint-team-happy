import { generateTestObjects, generateUniqueName, generateRandomPassword, isEqual, isExists, testPostJson, generateRandomEmail, assert } from "../helper.js";

const registerNegativePayloads = generateTestObjects({
    email: { type: "string", notNull: true, isEmail: true },
    name: { type: "string", minLength: 5, maxLength: 50, notNull: true },
    password: { type: "string", minLength: 5, maxLength: 15, notNull: true }
}, {
    email: generateRandomEmail(),
    name: generateUniqueName(),
    password: generateRandomPassword()
})


const TEST_NAME = "(register test)"
/**
 * Test registration API
 * @param {Config} config 
 * @param {Object} tags 
 * @returns {User | null}
 */
export function TestRegistration(config, tags = {}) {
    let res, currentTest
    // eslint-disable-next-line no-undef
    let route = __ENV.BASE_URL + "/v1/user/register"
    const currentFeature = TEST_NAME + " | post register"
    const positivePayload = {
        email: generateRandomEmail(),
        name: generateUniqueName(),
        password: generateRandomPassword()
    }

    if (!config.POSITIVE_CASE) {
        // Negative case, no body
        currentTest = "no payload"
        res = testPostJson(route, {}, {}, tags, ["noContentType"])
        assert(res, currentFeature, config, {
            [`${currentTest} should return 400`]: (r) => r.status === 400
        })


        // Negative case, invalid payload
        currentTest = "invalid payload"
        registerNegativePayloads.forEach(payload => {
            res = testPostJson(route, payload, {}, tags)
            assert(res, currentFeature, config, {
                [`${currentTest} should return 400`]: (r) => r.status === 400,
            }, payload)
        });
    }

    // Positive case, register new user
    currentTest = "register new user"
    res = testPostJson(route, positivePayload, {}, tags)
    const positivePayloadPassAssertTest = assert(res, currentFeature, config, {
        [`${currentTest} should return 201`]: (r) => r.status === 201,
        [`${currentTest} should have user name`]: (r) => isEqual(r, 'data.name', positivePayload.name),
        [`${currentTest} should have user email`]: (r) => isEqual(r, 'data.email', positivePayload.email),
        [`${currentTest} should have user accessToken`]: (r) => isExists(r, 'data.accessToken'),
    }, positivePayload)

    if (!config.POSITIVE_CASE) {
        // Negative case, email already exists
        currentTest = "email already exists"
        const exRes = testPostJson(route, positivePayload, {}, tags)
        assert(exRes, currentFeature, config, {
            [`${currentTest} should return 409`]: (r) => r.status === 409,
        }, positivePayload)
    }

    if (!positivePayloadPassAssertTest) return null

    return Object.assign(positivePayload, {
        accessToken: res.json().data.accessToken,
    })
}
