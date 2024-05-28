import http from 'k6/http';
import { fail } from 'k6';
import { isItUserValid, isNurseUserValid } from '../types/user.js';
import { generateRandomNumber } from '../helpers/generator.js';
import { assert, isExists } from '../helpers/assertion.js';

// Prepare the payload using the file to be uploaded
var positivePayload = {
    // eslint-disable-next-line no-undef
    file: http.file(open('../figure/image15KB.jpg', 'b'), 'image1.jpg', 'image/jpeg'),
};

/**
 * 
 * @param {import("../config.js").Config} config 
 * @param {Object} tags 
 * @param {import("../types/user.js").ItUser} userIt
 * @param {import("../types/user.js").NurseUser}userNurse 
 */
export function TestUpload(config, userIt, userNurse, tags) {
    const currentRoute = `${config.BASE_URL}/v1/image`
    const currentFeature = "image post | "
    if (!isItUserValid(userIt)) {
        fail(`${currentFeature} Invalid userIt object`)
    }
    if (!isNurseUserValid(userNurse)) {
        fail(`${currentFeature} Invalid userNurse object`)
    }

    const headers = {
        Authorization: `Bearer ${generateRandomNumber(0, 1) ? userNurse.accessToken : userIt.accessToken}`,
        ['Accept']: 'multipart/form-data'
    }

    let res;
    if (!config.POSITIVE_CASE) {
        // Negative case, empty auth
        res = http.post(currentRoute, {}, {}, tags);
        assert(res, 'POST', {}, currentFeature, {
            [currentFeature + "post upload file empty auth should return 401"]: (v) => v.status === 401
        }, config)
        // Negative case, empty file 
        res = http.post(currentRoute, {}, { headers }, tags);
        assert(res, 'POST', {}, currentFeature, {
            [currentFeature + "post upload file empty file should return 400"]: (v) => v.status === 400
        }, config)
    }

    // Positive case, upload file
    res = http.post(currentRoute, positivePayload, { headers }, tags);
    assert(res, 'POST', positivePayload, currentFeature, {
        [currentFeature + "correct file should return 200"]: (v) => v.status === 200,
        [currentFeature + "correct file should have imageUrl"]: (v) => isExists(v, "data.imageUrl"),
    }, config)
}