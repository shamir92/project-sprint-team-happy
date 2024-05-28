import { fail } from "k6";
import { testPostJsonAssert } from "../helpers/request.js";
import { isItUserValid, isNurseUserValid } from "../types/user.js";
import { generateRandomGender } from "../types/medicalRecord.js";

const { generateTestObjects, generateRandomImageUrl, generateRandomName, generateRandomNumber, generateRandomDate } = require("../helpers/generator.js");

const medicalPatientNegativePayload = (positivePayload) => generateTestObjects({
    identityNumber: { type: "number", notNull: true, min: 1000000000000000 },
    phoneNumber: { type: 'string', notNull: true, minLength: 10, maxLength: 15 },
    name: { type: 'string', notNull: true, minLength: 3, maxLength: 30 },
    gender: { type: 'string', notNull: true, enum: ['male', 'female'] },
    birthDate: { type: 'string', notNull: true },
    identityCardScanImg: { type: "string", notNull: true, isUrl: true }
}, positivePayload)
/**
 * 
 * @param {import("../config.js").Config} config 
 * @param {Object} tags 
 * @param {import("../types/user.js").ItUser} userIt
 * @param {import("../types/user.js").NurseUser} userNurse 
 * @returns {import("../types/medicalRecord.js").Patient | null}
 */
export function TestMedicalPatientPost(config, userIt, userNurse, tags) {
    const currentRoute = `${config.BASE_URL}/v1/medical/patient`
    const currentFeature = "medical patient post"
    if (!isItUserValid(userIt)) {
        fail(`${currentFeature} Invalid userIt object`)
    }
    if (!isNurseUserValid(userNurse)) {
        console.log("userNurse", userNurse)
        fail(`${currentFeature} Invalid userNurse object`)
    }

    const headers = {
        Authorization: `Bearer ${generateRandomNumber(0, 1) ? userNurse.accessToken : userIt.accessToken}`
    }

    const medicalPatientPositivePayload = {
        identityNumber: generateRandomNumber(1000000000000000, 9007199254740991),
        phoneNumber: `+62${generateRandomNumber(1000000000, 9999999999)}`,
        name: generateRandomName(),
        gender: generateRandomGender(),
        birthDate: generateRandomDate(new Date(1980, 1, 1), new Date(2024, 12, 31)),
        identityCardScanImg: generateRandomImageUrl()
    }

    if (!config.POSITIVE_CASE) {
        testPostJsonAssert(currentFeature, "no header", currentRoute, {}, {}, {
            ['should return 401']: (res) => res.status === 401,
        }, config, tags);

        testPostJsonAssert(currentFeature, "no payload", currentRoute, {}, headers, {
            ['should return 400']: (res) => res.status === 400,
        }, config, tags);

        medicalPatientNegativePayload(medicalPatientPositivePayload).forEach((payload) => {
            testPostJsonAssert(currentFeature, "invalid payload", currentRoute, payload, headers, {
                ['should return 400']: (res) => res.status === 400,
            }, config, tags);
        });
    }

    testPostJsonAssert(currentFeature, "add medical patient", currentRoute, medicalPatientPositivePayload, headers, {
        ['should return 201']: (res) => res.status === 201,
    }, config, tags);

    return null
}