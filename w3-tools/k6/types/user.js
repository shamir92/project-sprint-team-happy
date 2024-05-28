/**
 * @typedef {Object} ItUser
 * @property {string} userId - The user id of the user.
 * @property {string} name - The name of the user.
 * @property {number} nip - The nip of the user.
 * @property {string} password - The password of the user.
 * @property {string} accessToken - The accessToken of the user.
 */
/**
 * @typedef {Object} NurseUserWithoutAccess
 * @property {string} userId - The user id of the user.
 * @property {string} name - The name of the user.
 * @property {number} nip - The nip of the user.
 */
/**
 * @typedef {Object} NurseUserWithoutLogin
 * @property {string} userId - The user id of the user.
 * @property {string} name - The name of the user.
 * @property {number} nip - The nip of the user.
 * @property {string} password - The password of the user.
 */
/**
 * @typedef {Object} NurseUser 
 * @property {string} userId - The user id of the user.
 * @property {string} name - The name of the user.
 * @property {number} nip - The nip of the user.
 * @property {string} password - The password of the user.
 * @property {string} accessToken - The accessToken of the user.
 */
import { generateRandomNumber } from "../helpers/generator.js";


/**
 * validate validate it user object
 * @param {ItUser} usr 
 * @returns {Boolean}
 */
export function isItUserValid(usr) {
    try {
        if (
            typeof usr === 'object' &&
            typeof usr.nip === 'number' && `${usr.nip}`.startsWith('615') &&
            typeof usr.userId === 'string' &&
            typeof usr.password === 'string' &&
            typeof usr.accessToken === 'string'
        ) {
            return true;
        }
        return false;
    } catch (error) {
        return false;
    }
}
export function isUserValidLogin(usr) {
    try {
        if (
            typeof usr === 'object' &&
            typeof usr.nip === 'number' && `${usr.nip}`.startsWith('615') &&
            typeof usr.password === 'string'
        ) {
            return true;
        }
        return false;
    } catch (error) {
        return false;
    }
}
/**
/**
 * validate validate nurse user object
 * @param {NurseUser} usr 
 * @returns {Boolean}
 */
export function isNurseUserValid(usr) {
    try {
        if (
            typeof usr === 'object' &&
            typeof usr.nip === 'number' && `${usr.nip}`.startsWith('303') &&
            typeof usr.userId === 'string' &&
            typeof usr.password === 'string' &&
            typeof usr.accessToken === 'string'
        ) {
            return true;
        }
        return false;
    } catch (error) {
        return false;
    }
}

export function generateItUserNip() {
    return parseInt(`615${generateNipSuffix()}`);
}

export function generateNurseUserNip() {
    return parseInt(`303${generateNipSuffix()}`);
}

export function generateNipSuffix() {
    const gender = generateRandomNumber(1, 2);
    const year = generateRandomNumber(2000, 2024);
    let month = generateRandomNumber(1, 12);
    if (`${month}`.length < 2) {
        month = `${month}`.padStart(2, '0');
    }
    let identifier = generateRandomNumber(0, 99999);
    if (`${identifier}`.length < 3) {
        identifier = `${identifier}`.padStart(3, '0');
    }
    return `${gender}${year}${month}${identifier}`;
}