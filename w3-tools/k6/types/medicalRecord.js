/**
 * @typedef {Object} Patient
 * @property {int} identityNumber
 * @property {string} phoneNumber
 * @property {string} name
 * @property {string} birthDate
 * @property {"male"|"female"} gender
 * @property {string} identityCardScanImg
 */

import { generateRandomNumber } from "../helpers/generator.js"


const genders = ['male', 'female']

export function generateRandomGender(versusGender) {
    if (versusGender) {
        return genders[genders.findIndex(v => v !== versusGender)]
    }
    return genders[generateRandomNumber(0, 1)]
}