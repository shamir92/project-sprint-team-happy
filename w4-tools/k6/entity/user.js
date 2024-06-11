/**
 * @typedef {Object} User
 * @property {string} token 
 * @property {string} username
 * @property {string} password
 * @property {string} email
 */

export function IsUser(obj) {
    return obj.token
        && obj.username
        && obj.password
        && obj.email;
}