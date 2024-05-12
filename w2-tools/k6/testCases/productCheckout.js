import { fail } from "k6";
import { isEqualWith, isExists, isOrdered, isTotalDataInRange, isValidDate } from "../helpers/assertion.js";
import { clone, generateRandomName, generateRandomNumber, generateTestObjects } from "../helpers/generator.js";
import { testGetAssert, testPostJsonAssert } from "../helpers/request.js";
import { generateInternationalCallingCode, isUserValid } from "../types/user.js";
import { generateProduct } from "../types/product.js";

const registerCustomerNegativePayloads = (positivePayload) => generateTestObjects({
    phoneNumber: { type: "string", notNull: true, minLength: 10, maxLength: 16 },
    name: { type: "string", notNull: true, minLength: 5, maxLength: 50 },
}, positivePayload)

/**
 * 
 * @returns {import("../types/user.js").UserCustomer}
 */
function generateNewCustomer() {
    return {
        name: generateRandomName(),
        phoneNumber: `+${generateInternationalCallingCode()}${generateRandomNumber(1000000, 99999999)}`,
    }
}

/**
 * 
 * @param {import("../config.js").Config} config 
 * @param {Object} tags 
 * @param {import("../types/user.js").User} user
 * @returns {import("../types/user.js").UserCustomer | null}
 */
export function TestCustomerRegister(user, config, tags) {
    const currentRoute = `${config.BASE_URL}/v1/customer/register`
    const currentFeature = "Customer Register"
    /** @type {import("../types/user.js").UserCustomer} */
    const registerCustomerPositivePayload = generateNewCustomer()

    const headers = {
        Authorization: `Bearer ${user.accessToken}`
    }

    /** @type {import("../helpers/request.js").RequestAssertResponse} */
    let res;

    if (!config.POSITIVE_CASE) {
        testPostJsonAssert(currentFeature, "empty headers", currentRoute, {}, {}, {
            ['should return 401']: (res) => res.status === 401,
        }, config, tags);
        testPostJsonAssert(currentFeature, "invalid authorization header", currentRoute, {}, { Authorization: `Bearer ${headers.Authorization}a`, }, {
            ['should return 401']: (res) => res.status === 401,
        }, config, tags);
        registerCustomerNegativePayloads(registerCustomerPositivePayload).forEach((payload) => {
            testPostJsonAssert(currentFeature, "invalid payload", currentRoute, payload, headers, {
                ['should return 400']: (res) => res.status === 400,
            }, config, tags);
        });
    }

    res = testPostJsonAssert(currentFeature, "register with correct payload", currentRoute, registerCustomerPositivePayload, headers, {
        ['should return 201']: (res) => res.status === 201,
        ['should return have a userId']: (res) => isExists(res, "data.userId"),
        ['should return have a phoneNumber']: (res) => isExists(res, "data.phoneNumber"),
        ['should return have a name']: (res) => isExists(res, "data.name"),
    }, config, tags);

    if (!config.POSITIVE_CASE) {
        testPostJsonAssert(currentFeature, "register with existing phone number", currentRoute, registerCustomerPositivePayload, headers, {
            ['should return 409']: (res) => res.status === 409,
        }, config, tags);
    }
    if (res.isSuccess) {
        return Object.assign(registerCustomerPositivePayload, { userId: res.res.json().data.userId })
    }
    return null
}

/**
 * 
 * @param {import("../config.js").Config} config 
 * @param {Object} tags
 * @param {import("../types/user.js").User} user 
 */
export function TestCustomerGet(user, config, tags) {
    const currentRoute = `${config.BASE_URL}/v1/customer`
    const currentFeature = "get customer"

    /** @type {import("../helpers/request.js").RequestAssertResponse} */
    let res;
    const headers = {
        Authorization: `Bearer ${user.accessToken}`
    }

    if (!config.LOAD_TEST) {
        res = testGetAssert(currentFeature, "get customer", currentRoute, {}, headers, {
            ['should return 200']: (res) => res.status === 200,
            ['should return have a userId']: (res) => isExists(res, "data[].userId"),
            ['should return have a phoneNumber']: (res) => isExists(res, "data[].phoneNumber"),
            ['should return have a name']: (res) => isExists(res, "data[].name"),
        }, config, tags);
    } else {
        // limit search by name to prevent too many data
        let nameToSearch = generateRandomName()
        nameToSearch = nameToSearch.substring(0, nameToSearch.length - 6)
        res = testGetAssert(currentFeature, "get customer", currentRoute, {
            name: nameToSearch
        }, headers, {
            ['should return 200']: (res) => res.status === 200,
            ['should return have a userId']: (res) => isExists(res, "data[].userId"),
            ['should return have a phoneNumber']: (res) => isExists(res, "data[].phoneNumber"),
            ['should return have a name']: (res) => isExists(res, "data[].name"),
        }, config, tags);
    }

    if (!config.POSITIVE_CASE) {
        const phoneNumberToAdd = `+${generateInternationalCallingCode()}${generateRandomNumber(1000000, 99999999)}`
        let phoneNumberToSearch = phoneNumberToAdd.substring(0, phoneNumberToAdd.length - 3)
        phoneNumberToSearch = phoneNumberToSearch.substring(1)
        let nameToSearch = generateRandomName()
        if (!nameToSearch.includes("c")) {
            nameToSearch = nameToSearch + "c"
        }

        const postFeatureRoute = `${config.BASE_URL}/v1/customer/register`
        const postFeatureHeaders = {
            Authorization: `Bearer ${user.accessToken}`
        }

        testPostJsonAssert(currentFeature, 'add customer with searched phone number', postFeatureRoute, {
            name: generateRandomName(),
            phoneNumber: phoneNumberToAdd,
        }, postFeatureHeaders, {
            ['should return 201']: (res) => res.status === 201,
        }, config, tags)
        testPostJsonAssert(currentFeature, 'add customer with searched name', postFeatureRoute, {
            name: nameToSearch,
            phoneNumber: `+${generateInternationalCallingCode()}${generateRandomNumber(1000000, 99999999)}`,
        }, postFeatureHeaders, {
            ['should return 201']: (res) => res.status === 201,
        }, config, tags)

        testGetAssert(currentFeature, 'get customer with an "c" in the name', currentRoute, {
            name: "c"
        }, headers, {
            ['should return 200']: (res) => res.status === 200,
            ['should have an "c" in the result']: (res) => isEqualWith(res, 'data[].name', (v) => v.every(a => a.toLowerCase().includes("c"))),
        }, config, tags);
        testGetAssert(currentFeature, 'get customer with phone number', currentRoute, {
            phoneNumber: phoneNumberToSearch
        }, headers, {
            ['should return 200']: (res) => res.status === 200,
            ['should have a phone number that is searched']: (res) => isEqualWith(res, 'data[].phoneNumber', (v) => v.every(a => a.includes(`+${phoneNumberToSearch}`))),
        }, config, tags);
    }


    if (res.isSuccess) {
        return res.res.json().data
    }
    return null

}

const customerCheckoutNegativePayloads = (positivePayload) => generateTestObjects({
    customerId: { type: "string", notNull: true },
    productDetails: {
        type: "array", notNull: true, items: {
            type: "object",
            properties: {
                productId: { type: "string", notNull: true },
                quantity: { type: "number", notNull: true, min: 1 }
            }
        }
    },
    paid: { type: "number", notNull: true, min: 1 },
    change: { type: "number", notNull: true, min: 0 },
}, positivePayload)
/** 
 * @typedef {Object} ProductDetailsCheckout
 * @property {string} productId - The id of the customer
 * @property {number} quantity - The quantity of the product
 * @property {number} originalStock - The original stock of the product
 */

/**
 * @typedef {Object} ProductCheckout
 * @property {number} paid - The amount of money paid
 * @property {number} change - The amount of money change
 * @property {string} customerId - The id of the customer
 * @property {ProductDetailsCheckout[]} productDetails - The details of the product
 */

/**
 * 
 * @param {import("../types/user.js").User} user 
 * @param {import("../config.js").Config} config 
 * @param {Object} tags 
 * @returns {ProductCheckout[]}
 */
export function TestCustomerCheckout(user, config, tags) {
    const currentRoute = `${config.BASE_URL}/v1/product/checkout`
    const currentFeature = "customer checkout"

    if (!isUserValid(user)) {
        fail(`${currentFeature} Invalid user object`)
    }

    const headers = {
        Authorization: `Bearer ${user.accessToken}`
    }

    let res;

    /** @type {ProductDetailsCheckout[]} */
    let productsToBuy = []
    /** @type {ProductDetailsCheckout[]} */
    let productsToBuyButQuantityIsNotEnough = []
    let totalPrice = 0
    /**
     * @returns {ProductCheckout}
     */
    function composeProductToBuy() {
        productsToBuy = []
        productsToBuyButQuantityIsNotEnough = []
        totalPrice = 0
        /** @type {import("../helpers/request.js").RequestAssertResponse} */
        let res = testGetAssert(currentFeature, "get customer", `${config.BASE_URL}/v1/customer`, {}, headers, {
            ['should return 200']: (res) => res.status === 200,
        }, config, tags);
        if (!res.isSuccess) {
            fail(`${currentFeature}  Failed to get customer`)
        }
        /** @type {import("../types/user.js").UserCustomer[]} */
        const customers = res.res.json().data
        const customerToPay = customers[generateRandomNumber(0, customers.length - 1)]

        res = testGetAssert(currentFeature, "get product", `${config.BASE_URL}/v1/product/customer`, {
            inStock: true
        }, headers, {
            ['should return 200']: (res) => res.status === 200,
        }, config, tags);
        if (!res.isSuccess) {
            fail(`${currentFeature}  Failed to get product`)
        }
        /** @type {import("../types/product.js").Product[]} */
        const products = res.res.json().data
        let productsIndexToBuy = []
        for (let i = 0; i < generateRandomNumber(1, 4); i++) {
            productsIndexToBuy.push(generateRandomNumber(0, products.length - 1))
        }
        productsIndexToBuy = [...new Set(productsIndexToBuy)]
        productsIndexToBuy.forEach(i => {
            const product = products[i]
            const quantity = generateRandomNumber(1, product.stock)
            totalPrice += product.price * quantity
            productsToBuy.push({
                productId: product.id,
                quantity,
                originalStock: product.stock
            })
            productsToBuyButQuantityIsNotEnough.push({
                productId: product.id,
                quantity: quantity + 100000,
                originalStock: product.stock
            })
        });
        return {
            customerId: customerToPay.userId,
            productDetails: productsToBuy,
            paid: totalPrice,
            change: 0
        }
    }

    let customerCheckoutPositivePayload = composeProductToBuy()

    if (!config.POSITIVE_CASE) {
        testPostJsonAssert(currentFeature, "empty headers", currentRoute, {}, {}, {
            ['should return 401']: (res) => res.status === 401,
        }, config, tags);
        testPostJsonAssert(currentFeature, "invalid authorization header", currentRoute, {}, { Authorization: `Bearer ${headers.Authorization}a`, }, {
            ['should return 401']: (res) => res.status === 401,
        }, config, tags);
        customerCheckoutNegativePayloads(customerCheckoutPositivePayload).forEach((payload) => {
            testPostJsonAssert(currentFeature, "invalid payload", currentRoute, payload, headers, {
                ['should return 400']: (res) => res.status === 400,
            }, config, tags);
        });
        testPostJsonAssert(currentFeature, "productId is not found", currentRoute, Object.assign(clone(customerCheckoutPositivePayload), {
            productDetails: [{
                productId: "notfound",
                quantity: 1
            }]
        }), headers, {
            ['should return 404']: (res) => res.status === 404,
        }, config, tags);

        testPostJsonAssert(currentFeature, "paid is not enough", currentRoute, Object.assign(clone(customerCheckoutPositivePayload), {
            paid: totalPrice - 1
        }), headers, {
            ['should return 400']: (res) => res.status === 400,
        }, config, tags);

        testPostJsonAssert(currentFeature, "change is not right", currentRoute, Object.assign(clone(customerCheckoutPositivePayload), {
            paid: totalPrice + 10,
            change: 0
        }), headers, {
            ['should return 400']: (res) => res.status === 400,
        }, config, tags);

        testPostJsonAssert(currentFeature, "one of product ids is not enough", currentRoute, Object.assign(clone(customerCheckoutPositivePayload), {
            productDetails: productsToBuyButQuantityIsNotEnough
        }), headers, {
            ['should return 400']: (res) => res.status === 400,
        }, config, tags);


        const productIsAvailabeFalseToAdd = Object.assign(generateProduct(), {
            isAvailable: false,
        })
        res = testPostJsonAssert(currentFeature, 'add product with searched category', `${config.BASE_URL}/v1/product`, productIsAvailabeFalseToAdd, headers, {
            ['should return 201']: (res) => res.status === 201,
        }, config, tags)

        if (!res.isSuccess) {
            fail(`${currentFeature}  Failed to add product with isAvailable == false`)
        }
        const productIdIsAvailableFalse = res.res.json().data.id
        /** @type {ProductDetailsCheckout[]} */
        const productToBuyButOneItemIsAvailableFalse = [...productsToBuy, {
            productId: productIdIsAvailableFalse,
            quantity: productIsAvailabeFalseToAdd.stock - 1
        }]
        testPostJsonAssert(currentFeature, "one of product isAvailable == false", currentRoute, productToBuyButOneItemIsAvailableFalse, headers, {
            ['should return 400']: (res) => res.status === 400,
        }, config, tags);
    }

    if (!config.LOAD_TEST) {
        for (let i = 0; i < 10; i++) {
            let payload
            // only create new product to buy after the first iteration
            if (i === 0) {
                payload = clone(customerCheckoutPositivePayload)
            } else {
                payload = composeProductToBuy()
            }

            res = testPostJsonAssert(currentFeature, "checkout with correct payload loop " + i, currentRoute, payload, headers, {
                ['should return 200']: (res) => res.status === 200,
            }, config, tags);

            if (res.isSuccess && !config.POSITIVE_CASE && i === 0) {
                // only check product that already been checkouted after the first iteration
                productsToBuy.forEach(product => {
                    res = testGetAssert(currentFeature, "get product that already been checkouted", `${config.BASE_URL}/v1/product`, {
                        id: product.productId
                    }, headers, {
                        ['should return 200']: (res) => res.status === 200,
                        ['quantity should be less than previous get product']: (res) => isEqualWith(res, 'data[].stock', (v) => v.every(a => a < product.originalStock)),
                    }, config, tags);
                });
            }
        }
    } else {
        // run in normal behaviour if in load test mode
        testPostJsonAssert(currentFeature, "checkout with correct payload", currentRoute, clone(customerCheckoutPositivePayload), headers, {
            ['should return something']: (res) => res.status,
        }, config, tags);
    }


    return customerCheckoutPositivePayload
}

/**
 * 
 * @param {import("../types/user.js").User} user 
 * @param {ProductCheckout} productCheckoutToCheck 
 * @param {import("../config.js").Config} config 
 * @param {Object} tags 
 * @returns 
 */
export function TestCustomerCheckoutHistory(user, productCheckoutToCheck, config, tags) {
    const currentRoute = `${config.BASE_URL}/v1/product/checkout/history`
    const currentFeature = "get customer product checkout history"

    if (!isUserValid(user)) {
        fail(`${currentFeature} Invalid user object`)
    }
    if (!productCheckoutToCheck && !config.LOAD_TEST) {
        fail(`${currentFeature} Invalid productCheckoutToCheck object`)
    }

    const headers = {
        Authorization: `Bearer ${user.accessToken}`
    }

    /** @type {import("../helpers/request.js").RequestAssertResponse} */
    if (!config.POSITIVE_CASE) {
        testGetAssert(currentFeature, "empty headers", currentRoute, {}, {}, {
            ['should return 401']: (res) => res.status === 401,
        }, config, tags);
        testGetAssert(currentFeature, "invalid authorization header", currentRoute, {}, { Authorization: `Bearer ${headers.Authorization}a`, }, {
            ['should return 401']: (res) => res.status === 401,
        }, config, tags);
    }

    testGetAssert(currentFeature, "get customer checkout history", currentRoute, {}, headers, {
        ['should return 200']: (res) => res.status === 200,
        ['should have transactionId']: (res) => isExists(res, 'data[].transactionId'),
        ['should have customerId']: (res) => isExists(res, 'data[].customerId'),
        ['should have productDetails']: (res) => isExists(res, 'data[].productDetails'),
        ['should have productDetails.productId']: (res) => isExists(res, 'data[].productDetails[].productId'),
        ['should have productDetails.quantity']: (res) => isExists(res, 'data[].productDetails[].quantity'),
        ['should have paid']: (res) => isExists(res, 'data[].paid'),
        ['should have change']: (res) => isExists(res, 'data[].change'),
        ['should have createdAt']: (res) => isEqualWith(res, 'data[].createdAt', (v) => v.every(a => isValidDate(a))),
        ['should have return ordered by newest first']: (res) => isOrdered(res, 'data[].createdAt', 'desc', (v) => new Date(v)),
        ['should have no more than 5 data as default']: (res) => isTotalDataInRange(res, 'data[]', 1, 5),
    }, config, tags);

    if (!config.POSITIVE_CASE) {
        testGetAssert(currentFeature, 'get customer checkout history get createdAt oldest first', currentRoute, {
            createdAt: 'asc'
        }, headers, {
            ['should return 200']: (res) => res.status === 200,
            ['should have return ordered by oldest first']: (res) => isOrdered(res, 'data[].createdAt', 'asc', (v) => new Date(v)),
        }, config, tags);

        const paginationRes = testGetAssert(currentFeature, 'get customer checkout history filtered by pagination', currentRoute, {
            limit: 2,
            offset: 0
        }, headers, {
            ['should return 200']: (res) => res.status === 200,
            ['should have no more than 2 data as default']: (res) => isTotalDataInRange(res, 'data[]', 1, 2),
        }, config, tags);

        if (paginationRes.isSuccess) {
            testGetAssert(currentFeature, 'get customer checkout history filtered by pagination and offset', currentRoute, {
                limit: 2,
                offset: 2
            }, headers, {
                ['should return 200']: (res) => res.status === 200,
                ['should have no more than 2 data as default']: (res) => isTotalDataInRange(res, 'data[]', 1, 2),
                ['should have a different data than offset 0']: (res) => {
                    try {
                        return res.json().data.every(e => {
                            return paginationRes.res.json().data.every(a => a.transactionId !== e.transactionId)
                        })
                    } catch (error) {
                        return false
                    }
                },
            }, config, tags);
        }
        for (let index = 0; index < 990; index += 10) {
            const searchProductBought = testGetAssert(currentFeature, "get customer checkout history based on the product bought", currentRoute, {
                limit: 10,
                offset: index,
            }, headers, {
                ['should return 200']: (res) => res.status === 200,
            }, config, tags);
            if (searchProductBought.isSuccess) {
                if (searchProductBought.res.json().data.some(b =>
                    b.customerId === productCheckoutToCheck.customerId &&
                    b.paid === productCheckoutToCheck.paid &&
                    b.change === productCheckoutToCheck.change
                )) {
                    break;
                }
            }
        }
    }


    return null
}