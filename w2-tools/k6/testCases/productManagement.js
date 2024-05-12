/* eslint-disable no-loss-of-precision */
import { fail } from "k6";
import { MaxInt, generateRandomNumber, generateRandomName, generateTestObjects } from "../helpers/generator.js";
import { testDeleteAssert, testGetAssert, testPostJsonAssert, testPutJsonAssert } from "../helpers/request.js";
import { isUserValid } from "../types/user.js";
import { generateProduct, getRandomCategory } from "../types/product.js";
import { isEqualWith, isExists, isOrdered, isTotalDataInRange, isValidDate, isValidUrl } from "../helpers/assertion.js";

const productNegativePayload = (positivePayload) => generateTestObjects({
    name: { notNull: true, type: 'string', minLength: 1, maxLength: 30 },
    sku: { notNull: true, type: 'string', minLength: 1, maxLength: 30 },
    category: { notNull: true, type: 'string', enum: ["Clothing", "Accessories", "Footwear", "Beverages"] },
    imageUrl: { notNull: true, type: 'string', isUrl: true },
    notes: { notNull: true, type: 'string', minLength: 1, maxLength: 200 },
    price: { notNull: true, type: 'number', min: 1 },
    stock: { notNull: true, type: 'number', min: -1, max: 100000 },
    location: { notNull: true, type: 'string', minLength: 1, maxLength: 200 },
    isAvailable: { notNull: true, type: 'boolean' },
}, positivePayload)

/**
 * 
 * @param {import("../config.js").Config} config 
 * @param {Object} tags 
 * @param {import("../types/user.js").User} user
 * @returns {import("../types/product.js").Product}
 */
export function TestProductManagementPost(user, config, tags) {
    const currentRoute = `${config.BASE_URL}/v1/product`
    const currentFeature = "post product"

    if (!isUserValid(user)) {
        fail(`${currentFeature} Invalid user object`)
    }

    const headers = {
        Authorization: `Bearer ${user.accessToken}`
    }

    /** @type {import("../types/product.js").Product} */
    const productPositivePayload = generateProduct()

    /** @type {import("../helpers/request.js").RequestAssertResponse} */
    let res;

    if (!config.POSITIVE_CASE) {
        testPostJsonAssert(currentFeature, "empty headers", currentRoute, {}, {}, {
            ['should return 401']: (res) => res.status === 401,
        }, config, tags);
        testPostJsonAssert(currentFeature, "invalid authorization header", currentRoute, {}, { Authorization: `Bearer ${headers.Authorization}a`, }, {
            ['should return 401']: (res) => res.status === 401,
        }, config, tags);
        productNegativePayload(productPositivePayload).forEach((payload) => {
            testPostJsonAssert(currentFeature, "invalid payload", currentRoute, payload, headers, {
                ['should return 400']: (res) => res.status === 400,
            }, config, tags);
        });
    }

    res = testPostJsonAssert(currentFeature, "add product with correct payload", currentRoute, productPositivePayload, headers, {
        ['should return 201']: (res) => res.status === 201,
        ['should have id']: (res) => isExists(res, 'data.id'),
        ['should have createdAt and format should be date']: (res) => isEqualWith(res, 'data.createdAt', (v) => v.every(a => isValidDate(a))),
    }, config, tags);

    if (res.isSuccess) {
        return Object.assign(productPositivePayload, {
            id: res.res.json().data.id,
            createdAt: res.res.json().data.createdAt
        })
    }
    return null

}
/**
 * 
 * @param {import("../config.js").Config} config 
 * @param {Object} tags 
 * @param {import("../types/user.js").User} user
 * @returns {import("../types/product.js").Product[]}
 */
export function TestProductManagementGet(user, config, tags) {
    const currentRoute = `${config.BASE_URL}/v1/product`
    const currentFeature = "get product"

    if (!isUserValid(user)) {
        fail(`${currentFeature} Invalid user object`)
    }

    const headers = {
        Authorization: `Bearer ${user.accessToken}`
    }

    /** @type {import("../helpers/request.js").RequestAssertResponse} */
    let res;

    if (!config.POSITIVE_CASE) {
        testGetAssert(currentFeature, "empty headers", currentRoute, {}, {}, {
            ['should return 401']: (res) => res.status === 401,
        }, config, tags);
        testGetAssert(currentFeature, "invalid authorization header", currentRoute, {}, { Authorization: `Bearer ${headers.Authorization}a`, }, {
            ['should return 401']: (res) => res.status === 401,
        }, config, tags);
    }

    testPostJsonAssert(currentFeature, 'add product for get product', currentRoute, generateProduct(), headers, {
        ['should return 201']: (res) => res.status === 201,
    }, config, tags)

    res = testGetAssert(currentFeature, "get product", currentRoute, {}, headers, {
        ['should return 200']: (res) => res.status === 200,
        ['should have id']: (res) => isExists(res, 'data[].id'),
        ['should have name']: (res) => isExists(res, 'data[].name'),
        ['should have sku']: (res) => isExists(res, 'data[].sku'),
        ['should have category']: (res) => isExists(res, 'data[].category'),
        ['should have imageUrl']: (res) => isEqualWith(res, 'data[].imageUrl', (v) => v.every(a => isValidUrl(a))),
        ['should have stock']: (res) => isExists(res, 'data[].stock'),
        ['should have notes']: (res) => isExists(res, 'data[].notes'),
        ['should have price']: (res) => isExists(res, 'data[].price'),
        ['should have location']: (res) => isExists(res, 'data[].location'),
        ['should have isAvailable']: (res) => isExists(res, 'data[].isAvailable'),
        ['should have createdAt']: (res) => isEqualWith(res, 'data[].createdAt', (v) => v.every(a => isValidDate(a))),
        ['should have return ordered by newest first']: (res) => isOrdered(res, 'data[].createdAt', 'desc', (v) => new Date(v)),
        ['should have no more than 5 data as default']: (res) => isTotalDataInRange(res, 'data[]', 0, 5),
    }, config, tags);

    if (!config.POSITIVE_CASE) {
        const categoryToSearch = getRandomCategory()
        const skuToSearch = `${generateRandomNumber(10000000000, MaxInt)}`

        let nameToAdd = generateRandomName()
        if (!nameToAdd.includes("a")) {
            nameToAdd = nameToAdd + "a"
        }

        testPostJsonAssert(currentFeature, 'add product with searched category', currentRoute, Object.assign(generateProduct(), {
            name: nameToAdd,
            category: categoryToSearch
        }), headers, {
            ['should return 201']: (res) => res.status === 201,
        }, config, tags)
        testPostJsonAssert(currentFeature, 'add product with seached sku', currentRoute, Object.assign(generateProduct(), {
            sku: skuToSearch
        }), headers, {
            ['should return 201']: (res) => res.status === 201,
        }, config, tags)
        testPostJsonAssert(currentFeature, 'add product with isAvalable == false', currentRoute, Object.assign(generateProduct(), {
            isAvailable: false,
        }), headers, {
            ['should return 201']: (res) => res.status === 201,
        }, config, tags)

        testPostJsonAssert(currentFeature, 'add product with isAvalable == true', currentRoute, Object.assign(generateProduct(), {
            isAvailable: true,
        }), headers, {
            ['should return 201']: (res) => res.status === 201,
        }, config, tags)



        testGetAssert(currentFeature, 'get product with an "a" in the name', currentRoute, {
            name: "a"
        }, headers, {
            ['should return 200']: (res) => res.status === 200,
            ['should have an "a" in the result']: (res) => isEqualWith(res, 'data[].name', (v) => v.every(a => a.toLowerCase().includes("a"))),
        }, config, tags);
        testGetAssert(currentFeature, 'get product filtered by category', currentRoute, {
            category: categoryToSearch
        }, headers, {
            ['should return 200']: (res) => res.status === 200,
            ['should have the category that is searced']: (res) => isEqualWith(res, 'data[].category', (v) => v.every(a => a === categoryToSearch)),
        }, config, tags);
        testGetAssert(currentFeature, 'get product filtered by sku', currentRoute, {
            sku: skuToSearch
        }, headers, {
            ['should return 200']: (res) => res.status === 200,
            ['should have the sku that is searced']: (res) => isEqualWith(res, 'data[].sku', (v) => v.every(a => a === skuToSearch)),
        }, config, tags);
        testGetAssert(currentFeature, 'get product filtered by isAvailable true', currentRoute, {
            isAvailable: true
        }, headers, {
            ['should return 200']: (res) => res.status === 200,
            ['should have isAvailable true']: (res) => isEqualWith(res, 'data[].isAvailable', (v) => v.every(a => a === true)),
        }, config, tags);
        testGetAssert(currentFeature, 'get product filtered by isAvailable false', currentRoute, {
            isAvailable: false
        }, headers, {
            ['should return 200']: (res) => res.status === 200,
            ['should have isAvailable false']: (res) => isEqualWith(res, 'data[].isAvailable', (v) => v.every(a => a === false)),
        }, config, tags);
        testGetAssert(currentFeature, 'get product createdAt asc', currentRoute, {
            createdAt: 'asc'
        }, headers, {
            ['should return 200']: (res) => res.status === 200,
            ['should have return ordered by oldest first']: (res) => isOrdered(res, 'data[].createdAt', 'asc', (v) => new Date(v)),
        }, config, tags);
        // TODO: add inStock search
        const paginationRes = testGetAssert(currentFeature, 'get product filtered by pagination', currentRoute, {
            limit: 2,
            offset: 0
        }, headers, {
            ['should return 200']: (res) => res.status === 200,
            ['should have no more than 2 data as default']: (res) => isTotalDataInRange(res, 'data[]', 1, 2)
        }, config, tags);
        if (paginationRes.isSuccess) {
            testGetAssert(currentFeature, 'get product filtered by pagination and offset', currentRoute, {
                limit: 2,
                offset: 2
            }, headers, {
                ['should return 200']: (res) => res.status === 200,
                ['should have no more than 2 data as default']: (res) => isTotalDataInRange(res, 'data[]', 1, 2),
                ['should have a different data than offset 0']: (res) => {
                    try {
                        return res.json().data.every(e => {
                            return paginationRes.res.json().data.every(a => a.id !== e.id)
                        })
                    } catch (error) {
                        return false
                    }
                },
            }, config, tags);
        }
    }


    if (res.isSuccess) {
        return res.res.json().data
    }
    return null

}


/**
 * 
 * @param {import("../config.js").Config} config 
 * @param {Object} tags 
 * @param {import("../types/user.js").User} user
 * @returns {import("../types/product.js").Product}
 */
export function TestProductManagementPut(user, config, tags) {
    const currentRoute = `${config.BASE_URL}/v1/product`
    const currentFeature = "put product"

    if (!isUserValid(user)) {
        fail(`${currentFeature} Invalid user object`)
    }

    const headers = {
        Authorization: `Bearer ${user.accessToken}`
    }

    /** @type {import("../types/product.js").Product} */
    const productPositivePayload = generateProduct()

    /** @type {import("../helpers/request.js").RequestAssertResponse} */
    let res;

    res = testGetAssert(currentFeature, "get product", currentRoute, {}, headers, {
        ['should return 200']: (res) => res.status === 200,
    }, config, tags);

    if (!res.isSuccess) {
        fail(`${currentFeature} get product to edit failed`, res)
    }
    /** @type {import("../types/product.js").Product[]} */
    const products = res.res.json().data
    const productToEdit = products[generateRandomNumber(0, products.length - 1)]

    if (!config.POSITIVE_CASE) {
        testPutJsonAssert(currentFeature, "empty headers", `${currentRoute}/a`, {}, {}, {
            ['should return 401']: (res) => res.status === 401,
        }, config, tags);
        testPutJsonAssert(currentFeature, "invalid authorization header", `${currentRoute}/a`, {}, { Authorization: `Bearer ${headers.Authorization}a`, }, {
            ['should return 401']: (res) => res.status === 401,
        }, config, tags);
        productNegativePayload(productPositivePayload).forEach((payload) => {
            testPutJsonAssert(currentFeature, "invalid payload", `${currentRoute}/${productToEdit.id}`, payload, headers, {
                ['should return 400']: (res) => res.status === 400,
            }, config, tags);
        });
        testPutJsonAssert(currentFeature, "invalid id", `${currentRoute}/asd`, productPositivePayload, headers, {
            ['should return 404']: (res) => res.status === 404,
        }, config, tags);
    }


    res = testPutJsonAssert(currentFeature, "edit product with correct payload", `${currentRoute}/${productToEdit.id}`, productPositivePayload, headers, {
        ['should return 200']: (res) => res.status === 200,
    }, config, tags);

    if (!config.LOAD_TEST) {
        res = testGetAssert(currentFeature, "get product after edit", currentRoute, {
            id: productToEdit.id
        }, headers, {
            ['should return 200']: (res) => res.status === 200,
            ['should have edited data']: (res) => isEqualWith(res, 'data[]', (v) => {
                try {
                    Object.keys(productPositivePayload).forEach(key => {
                        if (v[0][key] !== productPositivePayload[key]) {
                            return false
                        }
                    })
                    return true
                } catch (error) {
                    return false
                }
            })
        }, config, tags);
    }

    return null
}

/**
 * 
 * @param {import("../config.js").Config} config 
 * @param {Object} tags 
 * @param {import("../types/user.js").User} user
 * @returns {import("../types/product.js").Product}
 */
export function TestProductManagementDelete(user, config, tags) {
    const currentRoute = `${config.BASE_URL}/v1/product`
    const currentFeature = "delete product"

    if (!isUserValid(user)) {
        fail(`${currentFeature} Invalid user object`)
    }

    const headers = {
        Authorization: `Bearer ${user.accessToken}`
    }

    /** @type {import("../helpers/request.js").RequestAssertResponse} */
    let res;

    res = testGetAssert(currentFeature, "get product", currentRoute, {}, headers, {
        ['should return 200']: (res) => res.status === 200,
    }, config, tags);

    /** @type {import("../types/product.js").Product[]} */
    const products = res.res.json().data
    const productToDelete = products[generateRandomNumber(0, products.length - 1)]

    if (!config.POSITIVE_CASE) {
        testDeleteAssert(currentFeature, "empty headers", `${currentRoute}/a`, {}, {}, {
            ['should return 401']: (res) => res.status === 401,
        }, config, tags);
        testDeleteAssert(currentFeature, "invalid authorization header", `${currentRoute}/a`, {}, { Authorization: `Bearer ${headers.Authorization}a`, }, {
            ['should return 401']: (res) => res.status === 401,
        }, config, tags);
        testDeleteAssert(currentFeature, "invalid id", `${currentRoute}/asd`, {}, headers, {
            ['should return 404']: (res) => res.status === 404,
        }, config, tags);
    }


    res = testDeleteAssert(currentFeature, "delete product with correct payload", `${currentRoute}/${productToDelete.id}`, {}, headers, {
        ['should return 200']: (res) => res.status === 200,
    }, config, tags);

    res = testGetAssert(currentFeature, "get product after delete", currentRoute, {
        id: productToDelete.id
    }, headers, {
        ['should return 200']: (res) => res.status === 200,
        ['should not return anything']: (res) => {
            try {
                return res.json().data.length === 0
            } catch (error) {
                return false
            }
        }
    }, config, tags);

    return null
}