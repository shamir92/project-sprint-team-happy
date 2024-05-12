/* eslint-disable no-undef */
import { sleep } from 'k6';
import { config } from './config.js';
import { clone, generateRandomNumber } from './helpers/generator.js';
import { TestLogin } from './testCases/authLogin.js';
import { TestRegister } from './testCases/authRegister.js';
import { TestCustomerCheckout, TestCustomerCheckoutHistory, TestCustomerGet, TestCustomerRegister } from './testCases/productCheckout.js';
import { TestCustomerGetProduct } from './testCases/productCustomer.js';
import { TestProductManagementDelete, TestProductManagementGet, TestProductManagementPost, TestProductManagementPut } from './testCases/productManagement.js';
import exec from 'k6/execution';

const stages = []

if (config.LOAD_TEST) {
    stages.push(
        { target: 50, iterations: 1, duration: "5s" },
        { target: 100, iterations: 1, duration: "10s" },
        { target: 150, iterations: 1, duration: "20s" },
        { target: 200, iterations: 1, duration: "20s" },
        { target: 250, iterations: 1, duration: "20s" },
        { target: 300, iterations: 1, duration: "20s" },
        { target: 600, iterations: 1, duration: "20s" },
        { target: 1200, iterations: 1, duration: "20s" },
    );
} else {
    stages.push({
        target: 1,
        iterations: 1
    });
}

function determineStage() {
    let elapsedTime = (exec.instance.currentTestRunDuration / 1000).toFixed(0);
    if (elapsedTime < 5) return 1; // First 5 seconds
    if (elapsedTime < 15) return 2; // Next 10 seconds
    if (elapsedTime < 35) return 3; // Next 20 seconds
    if (elapsedTime < 55) return 4; // Next 20 secondsd
    return 5; // Remaining time
}

export const options = {
    stages,
    // summaryTrendStats: ['avg', 'min', 'med', 'max', 'p(95)', 'p(99)'],
};

const positiveCaseConfig = Object.assign(clone(config), {
    POSITIVE_CASE: true
})

function calculatePercentage(percentage, __VU) {
    return (__VU - 1) % Math.ceil(__VU / Math.round(__VU * percentage)) === 0;
}

const users = []
function getRandomUser() {
    const i = generateRandomNumber(0, users.length - 1)
    return users[i]
}

export default function () {
    let tags = {}

    if (config.LOAD_TEST) {
        if (determineStage() == 1) {
            let user = TestRegister(positiveCaseConfig, tags)
            users.push(user)
            for (let i = 0; i < 10; i++) {
                TestProductManagementPost(getRandomUser(), positiveCaseConfig, tags)
                TestProductManagementGet(getRandomUser(), positiveCaseConfig, tags)
            }
            if (calculatePercentage(0.5, __VU)) {
                TestProductManagementPut(getRandomUser(), positiveCaseConfig, tags)
            }
            if (calculatePercentage(0.1, __VU)) {
                TestProductManagementDelete(getRandomUser(), positiveCaseConfig, tags)
            }
            for (let i = 0; i < 20; i++) {
                TestCustomerGetProduct(getRandomUser(), positiveCaseConfig, tags)
            }

            if (calculatePercentage(0.5, __VU)) {
                TestCustomerRegister(getRandomUser(), positiveCaseConfig, tags)
                TestCustomerCheckout(getRandomUser(), positiveCaseConfig, tags)
            }

        } else if (determineStage() == 2) {
            let user = TestRegister(positiveCaseConfig, tags)
            users.push(user)
            TestLogin(getRandomUser(), positiveCaseConfig, tags)
            for (let i = 0; i < 10; i++) {
                TestProductManagementPost(getRandomUser(), positiveCaseConfig, tags)
                TestProductManagementGet(getRandomUser(), positiveCaseConfig, tags)
            }
            if (calculatePercentage(0.5, __VU)) {
                TestProductManagementPut(getRandomUser(), positiveCaseConfig, tags)
            }
            if (calculatePercentage(0.1, __VU)) {
                TestProductManagementDelete(getRandomUser(), positiveCaseConfig, tags)
            }
            for (let i = 0; i < 20; i++) {
                TestCustomerGetProduct(getRandomUser(), positiveCaseConfig, tags)
            }

            if (calculatePercentage(0.5, __VU)) {
                TestCustomerRegister(getRandomUser(), positiveCaseConfig, tags)
                TestCustomerCheckout(getRandomUser(), positiveCaseConfig, tags)
            }
        } else if (determineStage() == 3) {
            let user = TestRegister(positiveCaseConfig, tags)
            users.push(user)
            TestLogin(getRandomUser(), positiveCaseConfig, tags)
            for (let i = 0; i < 10; i++) {
                if (calculatePercentage(0.2, __VU)) {
                    TestProductManagementPost(getRandomUser(), positiveCaseConfig, tags)
                }
                TestProductManagementGet(getRandomUser(), positiveCaseConfig, tags)
            }
            if (calculatePercentage(0.8, __VU)) {
                TestProductManagementPut(getRandomUser(), positiveCaseConfig, tags)
            }
            if (calculatePercentage(0.1, __VU)) {
                TestProductManagementDelete(getRandomUser(), positiveCaseConfig, tags)
            }
            for (let i = 0; i < 30; i++) {
                if (calculatePercentage(0.5, __VU)) {
                    TestCustomerGetProduct(getRandomUser(), positiveCaseConfig, tags)
                } else {
                    TestCustomerGetProduct(getRandomUser(), config, tags)
                }
            }

            let productCheckout;
            if (calculatePercentage(0.6, __VU)) {
                TestCustomerRegister(getRandomUser(), config, tags)
                productCheckout = TestCustomerCheckout(getRandomUser(), config, tags)
                if (calculatePercentage(0.1, __VU)) {
                    TestCustomerCheckoutHistory(user, productCheckout[generateRandomNumber(0, productCheckout.length - 1)],
                        positiveCaseConfig, tags)
                }
            }


        } else if (determineStage() == 4) {
            let user = TestRegister(positiveCaseConfig, tags)
            users.push(user)
            TestLogin(getRandomUser(), positiveCaseConfig, tags)
            for (let i = 0; i < 5; i++) {
                TestProductManagementGet(getRandomUser(), positiveCaseConfig, tags)
            }
            if (calculatePercentage(0.8, __VU)) {
                TestProductManagementPut(getRandomUser(), positiveCaseConfig, tags)
            }
            if (calculatePercentage(0.1, __VU)) {
                TestProductManagementDelete(getRandomUser(), positiveCaseConfig, tags)
            }
            for (let i = 0; i < 60; i++) {
                if (calculatePercentage(0.5, __VU)) {
                    TestCustomerGetProduct(getRandomUser(), positiveCaseConfig, tags)
                } else {
                    TestCustomerGetProduct(getRandomUser(), config, tags)
                }
            }

            let productCheckout;
            if (calculatePercentage(0.6, __VU)) {
                TestCustomerRegister(getRandomUser(), config, tags)
                productCheckout = TestCustomerCheckout(getRandomUser(), config, tags)
                if (calculatePercentage(0.1, __VU)) {
                    TestCustomerCheckoutHistory(user, productCheckout[generateRandomNumber(0, productCheckout.length - 1)],
                        positiveCaseConfig, tags)
                }
            }
        } else if (determineStage() == 5) {
            let user = TestRegister(positiveCaseConfig, tags)
            users.push(user)
            TestLogin(getRandomUser(), positiveCaseConfig, tags)
            for (let i = 0; i < 5; i++) {
                TestProductManagementGet(getRandomUser(), positiveCaseConfig, tags)
            }
            if (calculatePercentage(0.8, __VU)) {
                TestProductManagementPut(getRandomUser(), positiveCaseConfig, tags)
            }
            if (calculatePercentage(0.1, __VU)) {
                TestProductManagementDelete(getRandomUser(), positiveCaseConfig, tags)
            }
            for (let i = 0; i < 120; i++) {
                if (calculatePercentage(0.5, __VU)) {
                    TestCustomerGetProduct(getRandomUser(), positiveCaseConfig, tags)
                } else {
                    TestCustomerGetProduct(getRandomUser(), config, tags)
                }
            }

            let productCheckout;
            if (calculatePercentage(0.6, __VU)) {
                TestCustomerRegister(getRandomUser(), config, tags)
                productCheckout = TestCustomerCheckout(getRandomUser(), config, tags)
                if (calculatePercentage(0.1, __VU)) {
                    TestCustomerCheckoutHistory(user, productCheckout[generateRandomNumber(0, productCheckout.length - 1)],
                        positiveCaseConfig, tags)
                }
            }
        }

    } else {
        let user = TestRegister(config, tags)
        if (user) {
            user = TestLogin(user, config, tags)
            TestProductManagementPost(user, config, tags)
            TestProductManagementGet(user, config, tags)
            TestProductManagementPut(user, config, tags)
            TestProductManagementDelete(user, config, tags)

            TestCustomerGetProduct(user, config, tags)
            TestCustomerRegister(user, config, tags)
            TestCustomerGet(user, config, tags)
            const productCheckout = TestCustomerCheckout(user, config, tags)
            TestCustomerCheckoutHistory(user, productCheckout, config, tags)
        }
    }

    sleep(1)
}
