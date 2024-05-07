
export function generateRandomCatBreed() {
    const catBreeds = ['Persian', 'Maine Coon', 'Siamese', 'Ragdoll', 'Bengal', 'Sphynx', 'British Shorthair', 'Abyssinian', 'Scottish Fold', 'Birman'];
    const randomIndex = Math.floor(Math.random() * catBreeds.length);
    return catBreeds[randomIndex];
}

/**
 * Generates a random cat gender.
 * @param {CatSex} [oppositeSex] The opposite sex based on the provided CatSex value.
 * @returns {CatSex} A random cat gender or the opposite sex if provided.
 */
export function generateRandomCatGender(oppositeSex) {
    const catGenders = ['male', 'female'];
    const randomIndex = Math.floor(Math.random() * catGenders.length);
    const randomGender = catGenders[randomIndex];

    if (oppositeSex === 'male' || oppositeSex === 'female') {
        return oppositeSex === 'male' ? 'female' : 'male';
    }

    return randomGender;
}

