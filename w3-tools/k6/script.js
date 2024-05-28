/* eslint-disable no-undef */
import { sleep } from 'k6';
import { config } from './config.js';
import { combine, generateBoolFromPercentage, } from './helpers/generator.js';
import exec from 'k6/execution';
import { TestLogin } from './testCases/authLogin.js';
import { TestRegister } from './testCases/authRegister.js';
import { TestNurseManagementPost } from './testCases/nurseManagementPost.js';
import { TestNurseManagementGet } from './testCases/nurseManagementGet.js';
import { TestNurseManagementPut } from './testCases/nurseManagementPut.js';
import { TestNurseManagementDelete } from './testCases/nurseManagementDelete.js';
import { TestMedicalPatientGet } from './testCases/medicalPatientGet.js';
import { TestMedicalPatientPost } from './testCases/medicalPatientPost.js';
import { TestMedicalRecordPost } from './testCases/medicalRecordPost.js';
import { TestNurseManagementAccessPost } from './testCases/nurseManagementAccessPost.js';
import { TestNurseManagementLoginPost } from './testCases/nurseManagementLoginPost.js';
import { TestUpload } from './testCases/uploadPost.js';
import { generateItUserNip, generateNurseUserNip } from './types/user.js';
import { TestMedicalRecordGet } from './testCases/medicalRecordGet.js';
import grpc from 'k6/net/grpc';

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
    if (elapsedTime < 5) return 1; // 0
    if (elapsedTime < 15) return 2; // 50
    if (elapsedTime < 35) return 3; // 100
    if (elapsedTime < 55) return 4; // 150
    if (elapsedTime < 75) return 5; // 200
    if (elapsedTime < 95) return 6; // 300
    if (elapsedTime < 115) return 7; // 600
    if (elapsedTime < 135) return 8; // 1200
    return 9; // Remaining time
}

export const options = {
    stages: stages,
};

export function setup() {
    if (config.LOAD_TEST) {
        console.log('Warmup, adding 5 users to the database.')
        client.connect('127.0.0.1:50051', {
            plaintext: true
        });
        for (let i = 0; i < 5; i++) {
            let usrIt = TestRegister(positiveConfig, generateItUserNip(), {})
            if (usrIt) {
                console.log('usrIt success:', usrIt)
                PostUsedIt(client, usrIt)
            } else {
                console.log('usrIt fail:', usrIt)
            }
        }
        client.close();
    }
}
export function teardown() {
    if (config.LOAD_TEST) {
        client.connect('127.0.0.1:50051', {
            plaintext: true
        });

        Reset(client)
        client.close();
    }
}

const positiveConfig = combine(config, {
    POSITIVE_CASE: true
})

const client = new grpc.Client();
client.load([], 'backend.proto');

function GetItNip(cli) {
    const response = cli.invoke('pb.NIPService/GetItNip', {});
    return parseInt(response.message.nip)
}
function GetNurseNip(cli) {
    const response = cli.invoke('pb.NIPService/GetNurseNip', {});
    return parseInt(response.message.nip)
}

function GetUsedNurse(cli) {
    const response = cli.invoke('pb.NIPService/GetUsedNurse', {});
    return { nip: parseInt(response.message.nip), password: response.message.password }
}

function GetUsedIt(cli) {
    const response = cli.invoke('pb.NIPService/GetUsedIt', {});
    return { nip: parseInt(response.message.nip), password: response.message.password }
}

function PostUsedIt(cli, payload) {
    cli.invoke('pb.NIPService/PostUsedIT', {
        nip: payload.nip,
        password: payload.password
    });
}

function PostUsedNurse(cli, payload) {
    cli.invoke('pb.NIPService/PostUsedNurse', {
        nip: payload.nip,
        password: payload.password
    });
}

function Reset(cli) {
    cli.invoke('pb.NIPService/ResetAll', {
    });
}

function loop(fn, times) {
    for (let index = 0; index < times; index++) {
        fn()
    }
}

export default function () {
    let tags = {}

    if (config.LOAD_TEST) {
        client.connect('127.0.0.1:50051', {
            plaintext: true
        });

        if (determineStage() == 1) { // 0
            let usrIt;
            usrIt = TestRegister(config, GetItNip(client), tags)
            if (usrIt) {
                PostUsedIt(client, usrIt)
                usrIt = TestLogin(usrIt, positiveConfig, GetItNip(client), tags)

                loop(() => TestNurseManagementPost(positiveConfig, usrIt, GetNurseNip(client), tags), 8)
                loop(() => TestNurseManagementGet(positiveConfig, usrIt, tags), 10)

                if (generateBoolFromPercentage(.2)) {
                    TestNurseManagementPut(positiveConfig, usrIt, GetNurseNip(client), tags)
                    TestNurseManagementDelete(config, usrIt, tags)
                }
            }
        }
        else if (determineStage() == 2) { // 50
            let usrIt;
            let usrNurse;
            if (generateBoolFromPercentage(.2)) {
                const regUsr = TestRegister(config, GetItNip(client), tags)
                if (regUsr) {
                    PostUsedIt(client, regUsr)
                }
            }
            usrIt = TestLogin(GetUsedIt(client), positiveConfig, GetItNip(client), tags)
            if (usrIt) {
                if (generateBoolFromPercentage(.5)) {
                    TestNurseManagementPost(positiveConfig, usrIt, GetNurseNip(client), tags)
                }
                loop(() => TestNurseManagementGet(positiveConfig, usrIt, tags), 15)

                if (generateBoolFromPercentage(.2)) {
                    TestNurseManagementPut(positiveConfig, usrIt, GetNurseNip(client), tags)
                    TestNurseManagementDelete(positiveConfig, usrIt, tags)
                }

                const rawNurse = TestNurseManagementPost(positiveConfig, usrIt, GetNurseNip(client), tags)
                if (generateBoolFromPercentage(.6)) {
                    if (rawNurse) {
                        const accessNurse = TestNurseManagementAccessPost(config, usrIt, rawNurse, tags)
                        if (accessNurse) {
                            usrNurse = TestNurseManagementLoginPost(positiveConfig, accessNurse, tags)
                            if (usrNurse) {
                                TestUpload(positiveConfig, usrIt, usrNurse, tags)
                                PostUsedNurse(client, usrNurse)
                                if (generateBoolFromPercentage(.8)) {
                                    loop(() => TestMedicalPatientPost(positiveConfig, usrIt, usrNurse, tags), 5)
                                    TestMedicalPatientGet(positiveConfig, usrIt, usrNurse, tags)
                                }
                                if (generateBoolFromPercentage(.5)) {
                                    loop(() => TestMedicalRecordPost(positiveConfig, usrIt, usrNurse, tags), 5)
                                    TestMedicalRecordGet(positiveConfig, usrIt, usrNurse, tags)
                                }
                            }
                        }
                    }
                }
            }
        }
        else if (determineStage() == 3) { // 100
            let usrIt;
            let usrNurse;
            if (generateBoolFromPercentage(.1)) {
                const regUsr = TestRegister(config, GetItNip(client), tags)
                if (regUsr) {
                    PostUsedIt(client, regUsr)
                }
            }
            usrIt = TestLogin(GetUsedIt(client), positiveConfig, GetItNip(client), tags)

            if (generateBoolFromPercentage(.3)) {
                TestNurseManagementPost(config, usrIt, GetNurseNip(client), tags)
            }
            loop(() => TestNurseManagementGet(positiveConfig, usrIt, tags), 20)

            if (generateBoolFromPercentage(.1)) {
                TestNurseManagementPut(config, usrIt, GetNurseNip(client), tags)
                TestNurseManagementDelete(config, usrIt, tags)
            }

            if (generateBoolFromPercentage(.2)) {
                const rawNurse = TestNurseManagementPost(positiveConfig, usrIt, GetNurseNip(client), tags)
                if (generateBoolFromPercentage(.6)) {
                    if (rawNurse) {
                        const accessNurse = TestNurseManagementAccessPost(config, usrIt, rawNurse, tags)
                        if (accessNurse) {
                            usrNurse = TestNurseManagementLoginPost(positiveConfig, accessNurse, tags)
                            if (usrNurse) {
                                TestUpload(positiveConfig, usrIt, usrNurse, tags)
                                PostUsedNurse(client, usrNurse)
                                if (generateBoolFromPercentage(.8)) {
                                    loop(() => TestMedicalPatientPost(positiveConfig, usrIt, usrNurse, tags), 5)
                                    TestMedicalPatientGet(positiveConfig, usrIt, usrNurse, tags)
                                }
                                if (generateBoolFromPercentage(.5)) {
                                    loop(() => TestMedicalRecordPost(positiveConfig, usrIt, usrNurse, tags), 5)
                                    TestMedicalRecordGet(positiveConfig, usrIt, usrNurse, tags)
                                }
                            }
                        }
                    }
                }
            } else {
                usrNurse = TestNurseManagementLoginPost(positiveConfig, GetUsedNurse(client), tags)
                if (usrNurse) {
                    TestUpload(config, usrIt, usrNurse, tags)
                    if (generateBoolFromPercentage(.8)) {
                        loop(() => TestMedicalPatientPost(positiveConfig, usrIt, usrNurse, tags), 5)
                        TestMedicalPatientGet(positiveConfig, usrIt, usrNurse, tags)
                    }
                    if (generateBoolFromPercentage(.5)) {
                        loop(() => TestMedicalRecordPost(config, usrIt, usrNurse, tags), 5)
                        TestMedicalRecordGet(config, usrIt, usrNurse, tags)
                    }
                }
            }
        }
        else if (determineStage() == 4) { // 150
            let usrIt;
            let usrNurse;
            if (generateBoolFromPercentage(.01)) {
                const regUsr = TestRegister(config, GetItNip(client), tags)
                if (regUsr) {
                    PostUsedIt(client, regUsr)
                }
            }
            usrIt = TestLogin(GetUsedIt(client), positiveConfig, GetItNip(client), tags)

            if (generateBoolFromPercentage(.2)) {
                TestNurseManagementPost(config, usrIt, GetNurseNip(client), tags)
            }
            loop(() => TestNurseManagementGet(positiveConfig, usrIt, tags), 20)

            if (generateBoolFromPercentage(.01)) {
                TestNurseManagementPut(config, usrIt, GetNurseNip(client), tags)
                TestNurseManagementDelete(positiveConfig, usrIt, tags)
            }

            if (generateBoolFromPercentage(.1)) {
                const rawNurse = TestNurseManagementPost(positiveConfig, usrIt, GetNurseNip(client), tags)
                if (generateBoolFromPercentage(.6)) {
                    if (rawNurse) {
                        const accessNurse = TestNurseManagementAccessPost(config, usrIt, rawNurse, tags)
                        if (accessNurse) {
                            usrNurse = TestNurseManagementLoginPost(positiveConfig, accessNurse, tags)
                            if (usrNurse) {
                                TestUpload(positiveConfig, usrIt, usrNurse, tags)
                                PostUsedNurse(client, usrNurse)
                                if (generateBoolFromPercentage(.8)) {
                                    loop(() => TestMedicalPatientPost(positiveConfig, usrIt, usrNurse, tags), 5)
                                    TestMedicalPatientGet(positiveConfig, usrIt, usrNurse, tags)
                                }
                                if (generateBoolFromPercentage(.5)) {
                                    loop(() => TestMedicalRecordPost(positiveConfig, usrIt, usrNurse, tags), 5)
                                    TestMedicalRecordGet(positiveConfig, usrIt, usrNurse, tags)
                                }
                            }
                        }
                    }
                }
            } else {
                usrNurse = TestNurseManagementLoginPost(positiveConfig, GetUsedNurse(client), tags)
                if (usrNurse) {
                    TestUpload(positiveConfig, usrIt, usrNurse, tags)
                    if (generateBoolFromPercentage(.8)) {
                        TestMedicalPatientPost(positiveConfig, usrIt, usrNurse, tags)
                        TestMedicalPatientGet(positiveConfig, usrIt, usrNurse, tags)
                    }
                    if (generateBoolFromPercentage(.5)) {
                        TestMedicalRecordPost(config, usrIt, usrNurse, tags)
                        TestMedicalRecordGet(positiveConfig, usrIt, usrNurse, tags)
                    }
                }
            }
        }
        else if (determineStage() == 5) { // 200
            let usrIt;
            let usrNurse;
            if (generateBoolFromPercentage(.01)) {
                const regUsr = TestRegister(config, GetItNip(client), tags)
                if (regUsr) {
                    PostUsedIt(client, regUsr)
                }
            }
            usrIt = TestLogin(GetUsedIt(client), config, GetItNip(client), tags)

            if (generateBoolFromPercentage(.2)) {
                TestNurseManagementPost(positiveConfig, usrIt, GetNurseNip(client), tags)
            }
            loop(() => TestNurseManagementGet(positiveConfig, usrIt, tags), 20)

            if (generateBoolFromPercentage(.01)) {
                TestNurseManagementPut(config, usrIt, GetNurseNip(client), tags)
                TestNurseManagementDelete(positiveConfig, usrIt, tags)
            }

            if (generateBoolFromPercentage(.01)) {
                const rawNurse = TestNurseManagementPost(positiveConfig, usrIt, GetNurseNip(client), tags)
                if (generateBoolFromPercentage(.6)) {
                    if (rawNurse) {
                        const accessNurse = TestNurseManagementAccessPost(config, usrIt, rawNurse, tags)
                        if (accessNurse) {
                            usrNurse = TestNurseManagementLoginPost(positiveConfig, accessNurse, tags)
                            if (usrNurse) {
                                TestUpload(positiveConfig, usrIt, usrNurse, tags)
                                PostUsedNurse(client, usrNurse)
                                if (generateBoolFromPercentage(.8)) {
                                    loop(() => TestMedicalPatientPost(positiveConfig, usrIt, usrNurse, tags), 5)
                                    TestMedicalPatientGet(positiveConfig, usrIt, usrNurse, tags)
                                }
                                if (generateBoolFromPercentage(.5)) {
                                    loop(() => TestMedicalRecordPost(positiveConfig, usrIt, usrNurse, tags), 5)
                                    TestMedicalRecordGet(positiveConfig, usrIt, usrNurse, tags)
                                }
                            }
                        }
                    }
                }
            } else {
                usrNurse = TestNurseManagementLoginPost(positiveConfig, GetUsedNurse(client), tags)
                if (usrNurse) {
                    TestUpload(positiveConfig, usrIt, usrNurse, tags)
                    if (generateBoolFromPercentage(.9)) {
                        TestMedicalPatientPost(config, usrIt, usrNurse, tags)
                        loop(() => TestMedicalPatientGet(positiveConfig, usrIt, usrNurse, tags), 10)
                    }
                    if (generateBoolFromPercentage(.9)) {
                        TestMedicalRecordPost(config, usrIt, usrNurse, tags)
                        loop(() => TestMedicalRecordGet(positiveConfig, usrIt, usrNurse, tags), 10)
                    }
                }
            }
        }
        else if (determineStage() == 6) { // 300
            let usrIt;
            let usrNurse;
            if (generateBoolFromPercentage(.001)) {
                const regUsr = TestRegister(positiveConfig, GetItNip(client), tags)
                if (regUsr) {
                    PostUsedIt(client, regUsr)
                }
            }
            usrIt = TestLogin(GetUsedIt(client), positiveConfig, GetItNip(client), tags)

            if (generateBoolFromPercentage(.2)) {
                TestNurseManagementPost(positiveConfig, usrIt, GetNurseNip(client), tags)
            }
            loop(() => TestNurseManagementGet(positiveConfig, usrIt, tags), 20)

            if (generateBoolFromPercentage(.01)) {
                TestNurseManagementPut(positiveConfig, usrIt, GetNurseNip(client), tags)
                TestNurseManagementDelete(positiveConfig, usrIt, tags)
            }

            if (generateBoolFromPercentage(.01)) {
                const rawNurse = TestNurseManagementPost(positiveConfig, usrIt, GetNurseNip(client), tags)
                if (generateBoolFromPercentage(.6)) {
                    if (rawNurse) {
                        const accessNurse = TestNurseManagementAccessPost(config, usrIt, rawNurse, tags)
                        if (accessNurse) {
                            usrNurse = TestNurseManagementLoginPost(positiveConfig, accessNurse, tags)
                            if (usrNurse) {
                                TestUpload(positiveConfig, usrIt, usrNurse, tags)
                                PostUsedNurse(client, usrNurse)
                                if (generateBoolFromPercentage(.8)) {
                                    loop(() => TestMedicalPatientPost(positiveConfig, usrIt, usrNurse, tags), 5)
                                    TestMedicalPatientGet(positiveConfig, usrIt, usrNurse, tags)
                                }
                                if (generateBoolFromPercentage(.5)) {
                                    loop(() => TestMedicalRecordPost(positiveConfig, usrIt, usrNurse, tags), 5)
                                    TestMedicalRecordGet(positiveConfig, usrIt, usrNurse, tags)
                                }
                            }
                        }
                    }
                }
            } else {
                usrNurse = TestNurseManagementLoginPost(positiveConfig, GetUsedNurse(client), tags)
                if (usrNurse) {
                    TestUpload(positiveConfig, usrIt, usrNurse, tags)
                    if (generateBoolFromPercentage(.9)) {
                        TestMedicalPatientPost(positiveConfig, usrIt, usrNurse, tags)
                        loop(() => TestMedicalPatientGet(positiveConfig, usrIt, usrNurse, tags), 10)
                    }
                    if (generateBoolFromPercentage(.9)) {
                        TestMedicalRecordPost(positiveConfig, usrIt, usrNurse, tags)
                        loop(() => TestMedicalRecordGet(config, usrIt, usrNurse, tags), 10)
                    }
                }
            }
        }
        else if (determineStage() == 7) { // 600
            let usrIt;
            let usrNurse;
            usrIt = TestLogin(GetUsedIt(client), positiveConfig, GetItNip(client), tags)

            if (generateBoolFromPercentage(.01)) {
                TestNurseManagementPost(config, usrIt, GetNurseNip(client), tags)
            }
            loop(() => TestNurseManagementGet(positiveConfig, usrIt, tags), 20)

            if (generateBoolFromPercentage(.01)) {
                TestNurseManagementPut(config, usrIt, GetNurseNip(client), tags)
                TestNurseManagementDelete(config, usrIt, tags)
            }

            usrNurse = TestNurseManagementLoginPost(positiveConfig, GetUsedNurse(client), tags)
            if (usrNurse) {
                TestUpload(positiveConfig, usrIt, usrNurse, tags)
                if (generateBoolFromPercentage(.9)) {
                    TestMedicalPatientPost(positiveConfig, usrIt, usrNurse, tags)
                    loop(() => TestMedicalPatientGet(positiveConfig, usrIt, usrNurse, tags), 30)
                }
                if (generateBoolFromPercentage(.9)) {
                    TestMedicalRecordPost(config, usrIt, usrNurse, tags)
                    loop(() => TestMedicalRecordGet(positiveConfig, usrIt, usrNurse, tags), 30)
                }
            }
        }
        else if (determineStage() == 8) { // 1200
            let usrIt;
            let usrNurse;
            usrIt = TestLogin(GetUsedIt(client), positiveConfig, GetItNip(client), tags)

            if (generateBoolFromPercentage(.01)) {
                TestNurseManagementPost(positiveConfig, usrIt, GetNurseNip(client), tags)
            }
            loop(() => TestNurseManagementGet(positiveConfig, usrIt, tags), 30)

            if (generateBoolFromPercentage(.01)) {
                TestNurseManagementPut(config, usrIt, GetNurseNip(client), tags)
                TestNurseManagementDelete(positiveConfig, usrIt, tags)
            }

            usrNurse = TestNurseManagementLoginPost(positiveConfig, GetUsedNurse(client), tags)
            if (usrNurse) {
                TestUpload(positiveConfig, usrIt, usrNurse, tags)
                if (generateBoolFromPercentage(.9)) {
                    TestMedicalPatientPost(config, usrIt, usrNurse, tags)
                    loop(() => TestMedicalPatientGet(positiveConfig, usrIt, usrNurse, tags), 50)
                }
                if (generateBoolFromPercentage(.9)) {
                    TestMedicalRecordPost(config, usrIt, usrNurse, tags)
                    loop(() => TestMedicalRecordGet(positiveConfig, usrIt, usrNurse, tags), 50)
                }
            }
        }
        else {
            let usrIt;
            let usrNurse;
            usrIt = TestLogin(GetUsedIt(client), positiveConfig, GetItNip(client), tags)

            if (generateBoolFromPercentage(.01)) {
                TestNurseManagementPost(positiveConfig, usrIt, GetNurseNip(client), tags)
            }
            loop(() => TestNurseManagementGet(positiveConfig, usrIt, tags), 20)

            if (generateBoolFromPercentage(.01)) {
                TestNurseManagementPut(positiveConfig, usrIt, GetNurseNip(client), tags)
                TestNurseManagementDelete(positiveConfig, usrIt, tags)
            }

            usrNurse = TestNurseManagementLoginPost(positiveConfig, GetUsedNurse(client), tags)
            if (usrNurse) {
                TestUpload(positiveConfig, usrIt, usrNurse, tags)
                if (generateBoolFromPercentage(.9)) {
                    TestMedicalPatientPost(positiveConfig, usrIt, usrNurse, tags)
                    loop(() => TestMedicalPatientGet(positiveConfig, usrIt, usrNurse, tags), 80)
                }
                if (generateBoolFromPercentage(.9)) {
                    TestMedicalRecordPost(positiveConfig, usrIt, usrNurse, tags)
                    loop(() => TestMedicalRecordGet(positiveConfig, usrIt, usrNurse, tags), 80)
                }
            }
        }
        client.close();
    } else {
        let usr
        for (let index = 0; index < 5; index++) {
            usr = TestRegister(config, generateItUserNip(), tags)
        }
        if (usr) {
            TestLogin(usr, config, generateItUserNip(), tags)
            TestNurseManagementPost(config, usr, generateNurseUserNip(), tags)
            TestNurseManagementGet(config, usr, tags)
            TestNurseManagementPut(config, usr, generateNurseUserNip(), tags)
            TestNurseManagementDelete(config, usr, tags)
            const rawNurse = TestNurseManagementPost(positiveConfig, usr, generateNurseUserNip(), tags)
            const accessNurse = TestNurseManagementAccessPost(config, usr, rawNurse, tags)
            const nurse = TestNurseManagementLoginPost(config, accessNurse, tags)
            if (nurse) {
                TestMedicalPatientPost(config, usr, nurse, tags)
                TestMedicalPatientGet(config, usr, nurse, tags)
                TestMedicalRecordPost(config, usr, nurse, tags)
                TestMedicalRecordGet(config, usr, nurse, tags)
                TestUpload(config, usr, nurse, tags)
            }
        }
    }

    sleep(1)
}
