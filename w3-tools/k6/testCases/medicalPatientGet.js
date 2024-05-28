import { fail } from "k6";
import { isEqualWith, isExists, isOrdered, isTotalDataInRange, isValidDate } from "../helpers/assertion.js";
import { testGetAssert, testPostJsonAssert } from "../helpers/request.js";
import { isItUserValid, isNurseUserValid } from "../types/user.js";
import { generateRandomDate, generateRandomImageUrl, generateRandomName, generateRandomNumber } from "../helpers/generator.js";
import { generateRandomGender } from "../types/medicalRecord.js";
/**
 * 
 * @param {import("../config.js").Config} config 
 * @param {Object} tags 
 * @param {import("../types/user.js").ItUser} userIt
 * @param {import("../types/user.js").NurseUser} userNurse 
 */
export function TestMedicalPatientGet(config, userIt, userNurse, tags) {
    const currentRoute = `${config.BASE_URL}/v1/medical/patient`
    const currentFeature = "medical patient get"
    if (!isItUserValid(userIt)) {
        fail(`${currentFeature} Invalid userIt object`)
    }
    if (!isNurseUserValid(userNurse)) {
        fail(`${currentFeature} Invalid userNurse object`)
    }
    const headers = {
        Authorization: `Bearer ${generateRandomNumber(0, 1) ? userNurse.accessToken : userIt.accessToken}`
    }

    if (!config.POSITIVE_CASE) {
        testGetAssert(currentFeature, "no header", currentRoute, {}, {}, {
            ['should return 401']: (res) => res.status === 401,
        }, config, tags);
    }
    let medicalPatientPositivePayload;
    const phoneNumberToSearch = `62${generateRandomNumber(1, 9)}`
    if (!config.LOAD_TEST) {
        for (let index = 0; index < 10; index++) {
            medicalPatientPositivePayload = {
                identityNumber: generateRandomNumber(1000000000000000, 9007199254740991),
                phoneNumber: `+${phoneNumberToSearch}${generateRandomNumber(100000000, 999999999)}`,
                name: generateRandomName(),
                gender: generateRandomGender(),
                birthDate: generateRandomDate(new Date(1980, 1, 1), new Date(2024, 12, 31)),
                identityCardScanImg: generateRandomImageUrl()
            }
            testPostJsonAssert(currentFeature, "add medical patient", currentRoute, medicalPatientPositivePayload, headers, {
                ['should return 201']: (res) => res.status === 201,
            }, config, tags);
        }
    }

    testGetAssert(currentFeature, "get all patient", currentRoute, {}, headers, {
        ['should return 200']: (res) => res.status === 200,
        ['should all have a identityNumber']: (res) => isExists(res, "data[].identityNumber"),
        ['should all have phoneNumber']: (res) => isExists(res, "data[].phoneNumber"),
        ['should all have a name']: (res) => isExists(res, "data[].name"),
        ['should have birthDate and format should be date']: (res) => isEqualWith(res, 'data[].birthDate', (v) => v.every(a => isValidDate(a))),
        ['should all have gender']: (res) => isExists(res, "data[].gender"),
        ['should all have a createdAt']: (res) => isExists(res, "data[].createdAt"),
        ['should have createdAt and format should be date']: (res) => isEqualWith(res, 'data[].createdAt', (v) => v.every(a => isValidDate(a))),
        ['should not have more than 5 result']: (res) => isTotalDataInRange(res, 'data[]', 1, 5),
    }, config, tags);

    testGetAssert(currentFeature, "get all patient with name", currentRoute, { name: 'a' }, headers, {
        ['should return 200']: (res) => res.status === 200,
        ['should all have a identityNumber']: (res) => isExists(res, "data[].identityNumber"),
        ['should all have phoneNumber']: (res) => isExists(res, "data[].phoneNumber"),
        ['should all have a name']: (res) => isExists(res, "data[].name"),
        ['should have birthDate and format should be date']: (res) => isEqualWith(res, 'data[].birthDate', (v) => v.every(a => isValidDate(a))),
        ['should all have gender']: (res) => isExists(res, "data[].gender"),
        ['should all have a createdAt']: (res) => isExists(res, "data[].createdAt"),
        ['should have createdAt and format should be date']: (res) => isEqualWith(res, 'data[].createdAt', (v) => v.every(a => isValidDate(a))),
        ['should not have more than 5 result']: (res) => isTotalDataInRange(res, 'data[]', 1, 5),
        ['should have result with "a" in it']: (res) => isEqualWith(res, 'data[].name', (v) => v.every(a => a.toLowerCase().includes('a'))),
    }, config, tags);

    if (!config.LOAD_TEST) {
        testGetAssert(currentFeature, "get all patient with identityNumber", currentRoute, { identityNumber: medicalPatientPositivePayload.identityNumber }, headers, {
            ['should return 200']: (res) => res.status === 200,
            ['should all have a identityNumber']: (res) => isExists(res, "data[].identityNumber"),
            ['should all have phoneNumber']: (res) => isExists(res, "data[].phoneNumber"),
            ['should all have a name']: (res) => isExists(res, "data[].name"),
            ['should have birthDate and format should be date']: (res) => isEqualWith(res, 'data[].birthDate', (v) => v.every(a => isValidDate(a))),
            ['should all have gender']: (res) => isExists(res, "data[].gender"),
            ['should all have a createdAt']: (res) => isExists(res, "data[].createdAt"),
            ['should have createdAt and format should be date']: (res) => isEqualWith(res, 'data[].createdAt', (v) => v.every(a => isValidDate(a))),
            ['should not have more than 5 result']: (res) => isTotalDataInRange(res, 'data[]', 1, 5),
            ['should have result with four digit of searched identityNumber in it']: (res) => isEqualWith(res, 'data[].identityNumber', (v) => v.every(a => `${a}`.includes(medicalPatientPositivePayload.identityNumber))),
        }, config, tags);
    }

    testGetAssert(currentFeature, "get all patient with phoneNumber", currentRoute, { phoneNumber: phoneNumberToSearch }, headers, {
        ['should return 200']: (res) => res.status === 200,
        ['should all have a identityNumber']: (res) => isExists(res, "data[].identityNumber"),
        ['should all have phoneNumber']: (res) => isExists(res, "data[].phoneNumber"),
        ['should all have a name']: (res) => isExists(res, "data[].name"),
        ['should have birthDate and format should be date']: (res) => isEqualWith(res, 'data[].birthDate', (v) => v.every(a => isValidDate(a))),
        ['should all have gender']: (res) => isExists(res, "data[].gender"),
        ['should all have a createdAt']: (res) => isExists(res, "data[].createdAt"),
        ['should have createdAt and format should be date']: (res) => isEqualWith(res, 'data[].createdAt', (v) => v.every(a => isValidDate(a))),
        ['should not have more than 5 result']: (res) => isTotalDataInRange(res, 'data[]', 1, 5),
        ['should have result with four digit of searched phoneNumber in it']: (res) => isEqualWith(res, 'data[].phoneNumber', (v) => v.every(a => `${a}`.includes(phoneNumberToSearch))),
    }, config, tags);

    testGetAssert(currentFeature, "get all patient with createdAt asc", currentRoute, { createdAt: 'asc' }, headers, {
        ['should return 200']: (res) => res.status === 200,
        ['should all have a identityNumber']: (res) => isExists(res, "data[].identityNumber"),
        ['should all have phoneNumber']: (res) => isExists(res, "data[].phoneNumber"),
        ['should all have a name']: (res) => isExists(res, "data[].name"),
        ['should have birthDate and format should be date']: (res) => isEqualWith(res, 'data[].birthDate', (v) => v.every(a => isValidDate(a))),
        ['should all have gender']: (res) => isExists(res, "data[].gender"),
        ['should all have a createdAt']: (res) => isExists(res, "data[].createdAt"),
        ['should have createdAt and format should be date']: (res) => isEqualWith(res, 'data[].createdAt', (v) => v.every(a => isValidDate(a))),
        ['should not have more than 5 result']: (res) => isTotalDataInRange(res, 'data[]', 1, 5),
        ['should have return ordered by oldest first']: (res) => isOrdered(res, 'data[].createdAt', 'asc', (v) => new Date(v)),
    }, config, tags);

    const paginationRes = testGetAssert(currentFeature, "get all patient with limit", currentRoute, { limit: 2 }, headers, {
        ['should return 200']: (res) => res.status === 200,
        ['should all have a identityNumber']: (res) => isExists(res, "data[].identityNumber"),
        ['should all have phoneNumber']: (res) => isExists(res, "data[].phoneNumber"),
        ['should all have a name']: (res) => isExists(res, "data[].name"),
        ['should have birthDate and format should be date']: (res) => isEqualWith(res, 'data[].birthDate', (v) => v.every(a => isValidDate(a))),
        ['should all have gender']: (res) => isExists(res, "data[].gender"),
        ['should all have a createdAt']: (res) => isExists(res, "data[].createdAt"),
        ['should have createdAt and format should be date']: (res) => isEqualWith(res, 'data[].createdAt', (v) => v.every(a => isValidDate(a))),
        ['should not have more than 2 result']: (res) => isTotalDataInRange(res, 'data[]', 1, 2),
    }, config, tags);

    if (paginationRes.isSuccess && !config.LOAD_TEST) {
        testGetAssert(currentFeature, "get all patient with limit and offset", currentRoute, { limit: 2, offset: 2 }, headers, {
            ['should return 200']: (res) => res.status === 200,
            ['should all have a identityNumber']: (res) => isExists(res, "data[].identityNumber"),
            ['should all have phoneNumber']: (res) => isExists(res, "data[].phoneNumber"),
            ['should all have a name']: (res) => isExists(res, "data[].name"),
            ['should have birthDate and format should be date']: (res) => isEqualWith(res, 'data[].birthDate', (v) => v.every(a => isValidDate(a))),
            ['should all have gender']: (res) => isExists(res, "data[].gender"),
            ['should all have a createdAt']: (res) => isExists(res, "data[].createdAt"),
            ['should have createdAt and format should be date']: (res) => isEqualWith(res, 'data[].createdAt', (v) => v.every(a => isValidDate(a))),
            ['should not have more than 2 result']: (res) => isTotalDataInRange(res, 'data[]', 1, 2),
            ['should have different data from offset 0']: (res) => {
                try {
                    return res.json().data.every(e => {
                        return paginationRes.res.json().data.every(a => a.identityNumber !== e.identityNumber)
                    })
                } catch (error) {
                    return false
                }
            },
        }, config, tags);
    }
}