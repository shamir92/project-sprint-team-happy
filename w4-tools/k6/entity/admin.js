/**
 * @typedef {Object} Admin
 * @property {string} token 
 * @property {string} username
 * @property {string} password
 * @property {string} email
 */

export function IsAdmin(obj) {
    return obj.token
        && obj.username
        && obj.password
        && obj.email;
}