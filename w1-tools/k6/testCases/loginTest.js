import { generateTestObjects, generateRandomPassword, isEqual, isExists, testPostJson, generateRandomEmail, assert } from "../helper.js";

const loginNegativePayloads = generateTestObjects({
    email: { type: "string", notNull: true, isEmail: true },
    password: { type: "string", minLength: 5, maxLength: 15, notNull: true }
}, {
    email: generateRandomEmail(),
    password: generateRandomPassword()
})


const TEST_NAME = "(login test)"
/**
 * Test login API
 * @param {Config} config 
 * @param {Object} tags 
 * @param {User} user
 * @returns {User | null}
 */
export function TestLogin(config, user, tags = {}) {
    let res, currentTest
    // eslint-disable-next-line no-undef
    let route = __ENV.BASE_URL + "/v1/user/login"
    const currentFeature = TEST_NAME + " | post login"
    const positivePayload = {
        email: user.email,
        password: user.password
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
        loginNegativePayloads.forEach(payload => {
            res = testPostJson(route, payload, {}, tags)
            assert(res, currentFeature, config, {
                [`${currentTest} invalid payload should return 400`]: (r) => r.status === 400,
            }, payload)
        });

        // Negative case, user is not found
        currentTest = "non existing user"
        res = testPostJson(route, {
            email: generateRandomEmail(),
            password: generateRandomPassword()
        }, {}, tags)
        assert(res, currentFeature, config, {
            [`${currentTest} non existing user should return 404`]: (r) => r.status === 404,
        }, positivePayload)

        // Negative case, user password is wrong
        currentTest = "invalid password"
        res = testPostJson(route, {
            email: user.email,
            password: generateRandomPassword()
        }, {}, tags)
        assert(res, currentFeature, config, {
            [`${currentTest} invalid password user should return 400`]: (r) => r.status === 400,
        }, positivePayload)
    }

    // Positive case, login user
    currentTest = "login user"
    res = testPostJson(route, positivePayload, {}, tags)
    const positivePayloadPassAssertTest = assert(res, currentFeature, config, {
        [`${currentTest} valid payload should return 200`]: (r) => r.status === 200,
        [`${currentTest} valid payload should have user name`]: (r) => isEqual(r, 'data.name', user.name),
        [`${currentTest} valid payload should have user email`]: (r) => isEqual(r, 'data.email', positivePayload.email),
        [`${currentTest} valid payload should have user accessToken`]: (r) => isExists(r, 'data.accessToken'),
    }, positivePayload)

    if (!positivePayloadPassAssertTest) return null

    return Object.assign(positivePayload, {
        accessToken: res.json().data.accessToken,
    })
}
