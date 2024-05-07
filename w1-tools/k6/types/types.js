/**
 * Represents a match between two cats.
 * @typedef {Object} CatMatch
 * @property {string} id - The id of the match.
 * @property {User} issuedBy - The user who issued the match.
 * @property {Cat} matchCatDetail - The details of the cat being matched.
 * @property {Cat} userCatDetail - The details of the user's cat.
 * @property {string} message - Additional message for the match.
 * @property {string} createdAt - The timestamp when the match was created (in ISO 8601 format).
 */
/**
 * @typedef {Object} Config
 * @property {boolean} VERBOSE - Determines if verbose logging is enabled.
 * @property {boolean} DEBUG_ALL - Determines if debug mode is enabled for all cases.
 * @property {boolean} POSITIVE_CASE - Determines if only positive test cases should be executed.
 * @property {boolean} GACHA - Determines if data preparation should be executed.
 * @property {boolean} LOAD_TEST 
 */
/**
 * @typedef {Object} User
 * @property {string} email - User email.
 * @property {string} name - User name.
 * @property {string} password - User password.
 * @property {string} accessToken - User access token.
 */

/**
 * Enum representing different cat breeds.
 * @typedef {('Persian' | 'Maine Coon' | 'Siamese' | 'Ragdoll' | 'Bengal' | 'Sphynx' | 'British Shorthair' | 'Abyssinian' | 'Scottish Fold' | 'Birman')} CatBreed
 */

/**
 * Enum representing different cat sexes.
 * @typedef {('male' | 'female')} CatSex
 */

/**
 * Represents a cat.
 * @typedef {Object} Cat
 * @property {string} id - The id of the cat.
 * @property {string} name - The name of the cat.
 * @property {CatBreed} race - The breed of the cat.
 * @property {number} ageInMonth - The age of the cat in months.
 * @property {CatSex} sex - The gender of the cat.
 * @property {string} description - The description of the cat.
 * @property {boolean} hasMatched - Is the cat already h
 * @property {string[]} imageUrls - The URLs of the cat's images.
 */

/**
 * Generates a random cat breed.
 * @returns {CatBreed} A random cat breed.
 */

