
import { fail } from "k6";
import { testGetAssert, testPostJsonAssert } from "../helpers/request.js";
import { isItUserValid, isNurseUserValid } from "../types/user.js";
import { generateRandomGender } from "../types/medicalRecord.js";

const { generateTestObjects, generateRandomNumber, generateRandomDescription, combine, generateRandomName, generateRandomDate, generateRandomImageUrl } = require("../helpers/generator.js");

const medicalRecordNegativePayload = (positivePayload) => generateTestObjects({
    identityNumber: { type: "number", notNull: true, min: 1000000000000000 },
    symptoms: { type: 'string', notNull: true, minLength: 1, maxLength: 2000 },
    medications: { type: 'string', notNull: true, minLength: 1, maxLength: 2000 },
}, positivePayload)
/**
 * 
 * @param {import("../config.js").Config} config 
 * @param {Object} tags 
 * @param {import("../types/user.js").ItUser} userIt
 * @param {import("../types/user.js").NurseUser}userNurse 
 */
export function TestMedicalRecordPost(config, userIt, userNurse, tags) {
    const currentRoute = `${config.BASE_URL}/v1/medical/record`
    const currentFeature = "medical record post"
    if (!isItUserValid(userIt)) {
        fail(`${currentFeature} Invalid userIt object`)
    }
    if (!isNurseUserValid(userNurse)) {
        fail(`${currentFeature} Invalid userNurse object`)
    }

    const headers = {
        Authorization: `Bearer ${userIt.accessToken}`
    }

    if (!config.LOAD_TEST) {
        for (let index = 0; index < 3; index++) {
            testPostJsonAssert(currentFeature, "add medical patient for record", `${config.BASE_URL}/v1/medical/patient`, {
                identityNumber: generateRandomNumber(1000000000000000, 9007199254740991),
                phoneNumber: `+62${generateRandomNumber(1000000000, 9999999999)}`,
                name: generateRandomName(),
                gender: generateRandomGender(),
                birthDate: generateRandomDate(new Date(1980, 1, 1), new Date(2024, 12, 31)),
                identityCardScanImg: generateRandomImageUrl()
            }, headers, {
                ['should return 201']: (res) => res.status === 201,
            }, config, tags);
        }
    }

    const paginationRes = testGetAssert(currentFeature, "get all patient", `${config.BASE_URL}/v1/medical/patient`, { limit: config.LOAD_TEST ? 10 : 1000 }, headers, {
        ['should return 200']: (res) => res.status === 200,
    }, config, tags);

    if (!paginationRes.isSuccess) {
        fail(currentFeature, "get all patient", paginationRes.res, config, tags)
    }

    /** @type {import("../types/medicalRecord.js").Patient[]} */
    const patients = paginationRes.res.json().data
    const patientToTry = patients[generateRandomNumber(0, patients.length - 1)]

    const medicalRecordPositivePayload = {
        identityNumber: patientToTry.identityNumber,
        symptoms: generateRandomDescription(2000),
        medications: generateRandomDescription(2000),
    }

    if (!config.POSITIVE_CASE) {
        testPostJsonAssert(currentFeature, "no header", currentRoute, {}, {}, {
            ['should return 401']: (res) => res.status === 401,
        }, config, tags);

        testPostJsonAssert(currentFeature, "no payload", currentRoute, {}, headers, {
            ['should return 400']: (res) => res.status === 400,
        }, config, tags);

        medicalRecordNegativePayload(medicalRecordPositivePayload).forEach((payload) => {
            testPostJsonAssert(currentFeature, "invalid payload", currentRoute, payload, headers, {
                ['should return 400']: (res) => res.status === 400,
            }, config, tags);
        });

        testPostJsonAssert(currentFeature, "not exists identity number", currentRoute, combine(medicalRecordPositivePayload, {
            identityNumber: generateRandomNumber(1000000000000000, 9007199254740991)
        }), headers, {
            ['should return 404']: (res) => res.status === 404,
        }, config, tags);
    }

    testPostJsonAssert(currentFeature, "add medical record", currentRoute, medicalRecordPositivePayload, headers, {
        ['should return 201']: (res) => res.status === 201,
    }, config, tags);

    return null
}