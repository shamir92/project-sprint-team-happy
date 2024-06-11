import grpc from "k6/net/grpc";
import { config } from "./k6/entity/config.js";
import { initMerchantNearestLocations, initPegeneratedTSPMerchants, initZonesWithPregeneratedMerchants, resetAll } from "./k6/repository/initCaddyRepository.js";
import { assignPregeneratedMerchant, getAllMerchantNearestLocations, getAllMerchantRoutes, getAllPregeneratedMerchants, getMerchantNearestLocations, getPregeneratedMerchant, getTwoZoneMerchantRoutes } from "./k6/repository/getterCaddyRepository.js";
import { AdminRegisterTest } from "./k6/tests/auth/adminRegisterTest.js";
import { fail, sleep } from "k6";
import { AdminLoginTest } from "./k6/tests/auth/adminLoginTest.js";
import { UserRegisterTest } from "./k6/tests/auth/userRegisterTest.js";
import { UserLoginTest } from "./k6/tests/auth/userLoginTest.js";
import { MerchantPostTest } from "./k6/tests/merchantManagement/merchantPostTest.js";
import { MerchantGetTest } from "./k6/tests/merchantManagement/merchantGetTest.js";
import { MerchantItemGetTest } from "./k6/tests/merchantManagement/merchantItemGetTest.js";
import { MerchantItemPostTest } from "./k6/tests/merchantManagement/merchantItemPostTest.js";
import { EstimateOrderTest } from "./k6/tests/purchase/estimateOrderTest.js";
import { GetNearbyMerchantTest } from "./k6/tests/purchase/getNearbyMerchantTest.js";
import { OrderTest } from "./k6/tests/purchase/orderTest.js";
import { GetOrderTest } from "./k6/tests/purchase/getOrderTest.js";
import { combine, generateRandomPassword } from "./k6/helpers/generator.js";

const client = new grpc.Client();
client.load([], 'contract.proto');
export const options = {
    stages: [
        { duration: '3s', target: 10 },
        { duration: '5s', target: 30 },
        { duration: '10s', target: 50 },
        { duration: '10s', target: 100 },
        { duration: '5s', target: 300 },
        { duration: '5s', target: 200 },
        { duration: '5s', target: 100 },
        { duration: '5s', target: 50 },
    ],
};
export function setup() {
    client.connect('127.0.0.1:50051', {
        plaintext: true
    });

    if (config.LOAD_TEST) {
        console.log(`⛑️ ARMING... ⛑️`)
        initZonesWithPregeneratedMerchants(client, {
            area: 1,
            gap: 10,
            numberOfMerchantsPerZone: 50,
            numberOfZones: 30
        });
        initMerchantNearestLocations(client, {
            generateCount: 300
        });
        initPegeneratedTSPMerchants(client, {
            generateCount: 300
        });
        console.log(`🚀 NUKE READY!`)
        console.log(`💥 FAYYYAAAAH! 🇯🇵`)
        console.log('Adding all pregenerated merchants... 😰')
        let admin = AdminRegisterTest(positiveConfig, { feature: "Admin Register" })
        checkRes(admin, "Admin Register failed")
        admin = AdminLoginTest(admin, positiveConfig, { feature: "Admin Login" })
        let loop = true
        while (loop) {
            const rawMerchant = getPregeneratedMerchant(client)
            if (rawMerchant.pregeneratedId == "") {
                loop = false
                continue
            }
            const merchantFromPost = MerchantPostTest(admin, rawMerchant, positiveConfig, { feature: "Merchant Post" })
            checkRes(merchantFromPost, "Merchant Post failed")
            assignPregeneratedMerchant(client, { merchantId: merchantFromPost.merchantId, pregeneratedId: rawMerchant.pregeneratedId })
        }
        console.log(`All pregenerated merchants added! 🥳`)
    } else {
        initZonesWithPregeneratedMerchants(client, {
            area: 1,
            gap: 10,
            numberOfMerchantsPerZone: 10,
            numberOfZones: 2
        });
        initMerchantNearestLocations(client, {
            generateCount: 2
        });
        initPegeneratedTSPMerchants(client, {
            generateCount: 2
        });
    }
    client.close();
}
export function teardown() {
    client.connect('127.0.0.1:50051', {
        plaintext: true
    });

    resetAll(client);
    client.close();

    console.log("Test ended, I recomend to remove all of your merchants before re-run the test to avoid miss estimation on nearest merchant location.")
}
const positiveConfig = combine(config, { POSITIVE_CASE: true })
export default function () {
    client.connect('127.0.0.1:50051', {
        plaintext: true,
        timeout: 60000
    });

    if (config.LOAD_TEST) {
        // Admin test
        let admin = AdminRegisterTest(positiveConfig, { feature: "Admin Register" })
        checkRes(admin, "Admin Register failed")
        admin = AdminLoginTest(admin, positiveConfig, { feature: "Admin Login" })
        checkRes(admin, "Admin Login failed")

        MerchantGetTest(admin, [], positiveConfig, { feature: "Merchant Get" })

        let user = UserRegisterTest(positiveConfig, { feature: "User Register" })
        checkRes(user, "User Register failed")
        user = UserLoginTest(user, positiveConfig, { feature: "User Login" })
        checkRes(user, "User Login failed")

        let allNearestMerchantLocations = getMerchantNearestLocations(client)
        checkRes(allNearestMerchantLocations, "Get all nearest merchant locations failed")
        GetNearbyMerchantTest(user, allNearestMerchantLocations, positiveConfig, { feature: "Get Nearby Merchant" })

        let allMerchantRoutes = getTwoZoneMerchantRoutes(client)
        checkRes(allMerchantRoutes, "Get all merchant routes failed")
        const estimateId = EstimateOrderTest(user, admin, allMerchantRoutes.zone[0], allMerchantRoutes.zone[1], positiveConfig, { feature: "Estimate Order" })
        checkRes(estimateId, "Estimate Order failed")
        OrderTest(user, estimateId, positiveConfig, { feature: "Order" })
        GetOrderTest(user, positiveConfig, { feature: "Get Order" })
        sleep(1)

        return
    }

    // Admin test
    let admin = AdminRegisterTest(config, { feature: "Admin Register" })
    checkRes(admin, "Admin Register failed")
    admin = AdminLoginTest(admin, config, { feature: "Admin Login" })
    checkRes(admin, "Admin Login failed")

    const merchantToAdd = getAllPregeneratedMerchants(client)

    const merchantFromPost = MerchantPostTest(admin, merchantToAdd.merchant.pop(), config, { feature: "Merchant Post" })
    checkRes(merchantFromPost, "Merchant Post failed")
    assignPregeneratedMerchant(client, { merchantId: merchantFromPost.merchantId, pregeneratedId: merchantFromPost.pregeneratedId })

    // slice the last element due to [].pop() method doesn't remove the last element in k6, it only nulls the last element
    merchantToAdd.merchant = merchantToAdd.merchant.slice(0, merchantToAdd.merchant.length - 1)
    const merchantFromGet = MerchantGetTest(admin, merchantToAdd.merchant, config, { feature: "Merchant Get" })
    checkRes(merchantFromGet, "Merchant Get failed")
    merchantFromGet.forEach(m => {
        assignPregeneratedMerchant(client, { merchantId: m.merchantId, pregeneratedId: m.pregeneratedId })
    });

    let merchantItems = MerchantItemGetTest(admin, merchantFromPost, 10, config, { feature: "Merchant Item Get" })
    checkRes(merchantItems, "Merchant Item Get failed")

    let merchantItem = MerchantItemPostTest(admin, merchantFromPost, config, { feature: "Merchant Item Post" })
    checkRes(merchantItem, "Merchant Item Post failed")
    merchantItems.push(merchantItem)


    // User test
    let user = UserRegisterTest(config, { feature: "User Register" })
    checkRes(user, "User Register failed")
    user = UserLoginTest(user, config, { feature: "User Login" })
    checkRes(user, "User Login failed")

    let allNearestMerchantLocations = getAllMerchantNearestLocations(client)
    checkRes(allNearestMerchantLocations, "Get all nearest merchant locations failed")
    GetNearbyMerchantTest(user, allNearestMerchantLocations.zones[0].records[0], config, { feature: "Get Nearby Merchant" })

    let allMerchantRoutes = getTwoZoneMerchantRoutes(client)
    checkRes(allMerchantRoutes, "Get all merchant routes failed")
    const estimateId = EstimateOrderTest(user, admin, allMerchantRoutes.zone[0], allMerchantRoutes.zone[1], config, { feature: "Estimate Order" })
    checkRes(estimateId, "Estimate Order failed")
    OrderTest(user, estimateId, config, { feature: "Order" })
    GetOrderTest(user, config, { feature: "Get Order" })

    client.close()
}

function checkRes(res, errMsg) {
    if (!res) {
        fail(errMsg)
    }
}