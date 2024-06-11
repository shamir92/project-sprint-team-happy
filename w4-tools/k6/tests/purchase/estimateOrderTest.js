import { IsAdmin } from "../../entity/admin.js";
import { generateRandomMerchantItemCategory } from "../../entity/merchantItem.js";
import { IsUser } from "../../entity/user.js";
import { isEqual, isEqualWith, isExists } from "../../helpers/assertion.js";
import { generateRandomImageUrl, generateRandomName, generateRandomNumber, generateTestObjects } from "../../helpers/generator.js";
import { testPostJsonAssert, testGetAssert } from "../../helpers/request.js";

/**
 * @param {import("../../entity/user.js").User} user
 * @param {import("../../entity/admin.js").Admin} admin 
 * @param {import("../../entity/config.js").Config} config 
 * @param {import("../../entity/merchant.js").RouteZone} zone1
 * @param {import("../../entity/merchant.js").RouteZone} zone2
 * 
 * @returns {string | null}
 */
export function EstimateOrderTest(user, admin, zone1, zone2, config, tags) {
    if (!IsUser(user) && !IsAdmin(admin)) {
        return;
    }

    const featureName = "User Estimate Order Test";
    const route = config.BASE_URL + `/users/estimate`;
    const headers = {
        Authorization: `Bearer ${user.token}`
    }

    const selectedZone1Route = zone1.routes[generateRandomNumber(0, zone1.routes.length - 1)]
    const selectedZone2Route = zone2.routes[generateRandomNumber(0, zone2.routes.length - 1)]
    const zone1StartingPoint = selectedZone1Route.startingPoint
    const zone2StartingPoint = selectedZone2Route.startingPoint

    /** @type {Record<string, import("../../entity/merchantItem.js").MerchantItem[]>} */
    const merchantItemMap = {}
    Object.values(selectedZone1Route.generatedRoutes)
        .concat(Object.values(selectedZone2Route.generatedRoutes))
        .forEach(merchant => {
            const res = testGetAssert("get merchant items", featureName,
                config.BASE_URL + `/admin/merchants/${merchant.merchantId}/items`, {},
                { Authorization: `Bearer ${admin.token}` },
                {
                    ['should return 200']: (v) => v.status === 200,
                    ['should have itemId']: (v) => isExists(v, 'data[].itemId'),
                },
                config, tags)
            if (res.isSuccess) {
                if (res.res.json().data.length == 0) {
                    const merchantItemToAdd = {
                        name: generateRandomName(),
                        productCategory: generateRandomMerchantItemCategory(),
                        imageUrl: generateRandomImageUrl(),
                        price: generateRandomNumber(1, 1000000),
                    }
                    const addRes = testPostJsonAssert("add merchant items if not exists",
                        featureName,
                        config.BASE_URL + `/admin/merchants/${merchant.merchantId}/items`,
                        merchantItemToAdd,
                        { Authorization: `Bearer ${admin.token}` },
                        {
                            ['should return 201']: (v) => v.status === 201,
                            ['should have itemId']: (v) => isExists(v, 'itemId'),
                        },
                        config, tags)
                    if (addRes.isSuccess) {
                        merchantItemToAdd.itemId = addRes.res.json().itemId
                    }
                    merchantItemMap[merchant.merchantId] = [merchantItemToAdd]
                } else {
                    merchantItemMap[merchant.merchantId] = res.res.json().data
                }
            }
        });

    let positivePayloadTotalPrice = 0
    const positivePayload = {
        userLocation: { lat: zone1StartingPoint.lat, long: zone1StartingPoint.long },
        orders: [],
    }
    const negativePayloadWithFalseStartingPoint = {
        userLocation: { lat: zone1StartingPoint.lat, long: zone1StartingPoint.long },
        orders: [],
    }
    const negativePayloadWithTrueStartingPoint = {
        userLocation: { lat: zone1StartingPoint.lat, long: zone1StartingPoint.long },
        orders: [],
    }
    const negativePayloadWithAllFalseItemId = {
        userLocation: { lat: zone1StartingPoint.lat, long: zone1StartingPoint.long },
        orders: [],
    }
    const negativePayloadWithAllFalseMerchantId = {
        userLocation: { lat: zone1StartingPoint.lat, long: zone1StartingPoint.long },
        orders: [],
    }
    const negativePayloadWithFarUserLocation = {
        userLocation: { lat: zone2StartingPoint.lat, long: zone2StartingPoint.long },
        orders: [],
    }
    for (let i = 0; i < Object.keys(selectedZone1Route.generatedRoutes).length; i++) {
        const currentMerchant = selectedZone1Route.generatedRoutes[`${i}`]
        const items = merchantItemMap[currentMerchant.merchantId]
        const itemsToAdd = []
        items.forEach(item => {
            const qty = generateRandomNumber(1, 5)
            positivePayloadTotalPrice += item.price * qty
            itemsToAdd.push({
                itemId: item.itemId,
                quantity: qty
            })
        });
        positivePayload.orders.push({
            merchantId: currentMerchant.merchantId,
            isStartingPoint: i == 1,
            items: itemsToAdd
        })
        if (!config.POSITIVE_CASE) {
            negativePayloadWithFalseStartingPoint.orders.push({
                merchantId: currentMerchant.merchantId,
                isStartingPoint: false,
                items: merchantItemMap[currentMerchant.merchantId].map(e => ({ itemId: e.itemId, quantity: 1 }))
            })
            negativePayloadWithTrueStartingPoint.orders.push({
                merchantId: currentMerchant.merchantId,
                isStartingPoint: true,
                items: merchantItemMap[currentMerchant.merchantId].map(e => ({ itemId: e.itemId, quantity: 1 }))
            })
            negativePayloadWithAllFalseItemId.orders.push({
                merchantId: currentMerchant.merchantId,
                isStartingPoint: i == 0,
                items: [{ itemId: "invalid", quantity: 1 }]
            })
            negativePayloadWithAllFalseMerchantId.orders.push({
                merchantId: "invalid",
                isStartingPoint: i == 0,
                items: []
            })
            negativePayloadWithFarUserLocation.orders.push({
                merchantId: currentMerchant.merchantId,
                isStartingPoint: i == 0,
                items: merchantItemMap[currentMerchant.merchantId].map(e => ({ itemId: e.itemId, quantity: 1 }))
            })
        }
    }

    if (!config.POSITIVE_CASE) {
        testPostJsonAssert(
            "empty auth",
            featureName,
            route, {}, {}, {
            ['should return 401']: (v) => v.status === 401
        }, config, tags)
        testPostJsonAssert(
            "empty body",
            featureName,
            route, {}, headers, {
            ['should return 400']: (v) => v.status === 400
        }, config, tags)

        const testObjects = generateTestObjects({
            userLocation: {
                type: 'object',
                properties: {
                    lat: { type: 'number' },
                    long: { type: 'number' }
                }
            },
            orders: {
                type: 'array',
                items: {
                    type: 'object',
                    properties: {
                        merchantId: { type: 'string' },
                        isStartingPoint: { type: 'boolean' },
                        items: {
                            type: 'array',
                            items: {
                                type: 'object',
                                properties: {
                                    itemId: { type: 'string' },
                                    quantity: { type: 'number' }
                                }
                            }
                        }
                    }
                }
            }
        }, positivePayload)
        testObjects.forEach(payload => {
            testPostJsonAssert("invalid payload", featureName, route, payload, headers, {
                ['should return 400']: (res) => res.status === 400,
            }, config, tags);
        });

        testPostJsonAssert("all false starting point", featureName, route, negativePayloadWithFalseStartingPoint, headers, {
            ['should return 400']: (res) => res.status === 400,
        }, config, tags);
        testPostJsonAssert("all true starting point", featureName, route, negativePayloadWithTrueStartingPoint, headers, {
            ['should return 400']: (res) => res.status === 400,
        }, config, tags);
        testPostJsonAssert("all false item id", featureName, route, negativePayloadWithAllFalseItemId, headers, {
            ['should return 404']: (res) => res.status === 404,
        }, config, tags);
        testPostJsonAssert("all false merchant id", featureName, route, negativePayloadWithAllFalseMerchantId, headers, {
            ['should return 404']: (res) => res.status === 404,
        }, config, tags);
        testPostJsonAssert("with far user location", featureName, route, negativePayloadWithFarUserLocation, headers, {
            ['should return 400']: (res) => res.status === 400,
        }, config, tags);
    }

    const res = testPostJsonAssert("valid payload", featureName, route, positivePayload, headers, {
        ['should return 200']: (v) => v.status === 200,
        ['should have totalPrice and equal to calculated total']: (v) => isEqual(v, 'totalPrice', positivePayloadTotalPrice),
        ['should have calculatedEstimateId']: (v) => isExists(v, 'calculatedEstimateId'),
        ['should have estimatedDeliveryTimeInMinutes and not far from precalculated result']:
            (v) => isEqualWith(v, 'estimatedDeliveryTimeInMinutes', (a) => a >= selectedZone1Route.totalDuration - 5 && a <= selectedZone1Route.totalDuration + 6)
    }, config, tags)

    if (res.isSuccess) {
        return res.res.json().calculatedEstimateId
    }
    return null
}