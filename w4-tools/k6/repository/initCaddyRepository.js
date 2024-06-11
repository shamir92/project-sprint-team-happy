import grpc from "k6/net/grpc";
import { logGrpcError } from "./getterCaddyRepository.js";

/**
 * Initializes the merchant's nearest locations.
 * 
 * @param {grpc.Client} grpcClient - The gRPC client.
 * @param {Object} opts - The options for initializing the merchant's nearest locations.
 * @param {number} opts.generateCount - The count of locations to generate.
 * 
 * @returns {Object|null} - The response message or null if there was an error.
 */
export function initMerchantNearestLocations(grpcClient, opts) {
    const resp = grpcClient.invoke("pb.MerchantService/InitMerchantNearestLocations", opts);
    if (resp.status !== grpc.StatusOK) {
        logGrpcError(resp)
        return null
    }
    return resp.message
}

/**
 * 
 * @param {grpc.Client} grpcClient - The gRPC client.
 * @param {Object} opts - The options for initializing the merchant's nearest locations.
 * @param {number} opts.generateCount - The count of locations to generate.
 * 
 * @returns {null | Object}
 */
export function initPegeneratedTSPMerchants(grpcClient, opts) {
    const resp = grpcClient.invoke("pb.MerchantService/InitPegeneratedTSPMerchants", opts);
    if (resp.status !== grpc.StatusOK) {
        logGrpcError(resp)
        return null
    }
    return resp.message
}

/**
 * 
 * @param {grpc.Client} grpcClient 
 * @returns {null | Object}
 */
/**
 * Initializes zones with pregenerated merchants.
 * @param {grpc.Client} grpcClient - The gRPC client.
 * @param {Object} opts - The options for initializing zones.
 * @param {number} opts.area - The area value for initializing zones.
 * @param {number} opts.gap - The gap value for initializing zones.
 * @param {number} opts.numberOfZones - The number of zones to initialize.
 * @param {number} opts.numberOfMerchantsPerZone - The number of merchants per zone.
 * @returns {Object|null} - The response message or null if there was an error.
 */
export function initZonesWithPregeneratedMerchants(grpcClient, opts) {
    const resp = grpcClient.invoke("pb.MerchantService/InitZonesWithPregeneratedMerchants", opts);
    if (resp.status !== grpc.StatusOK) {
        logGrpcError(resp)
        return null
    }
    return resp.message
}

/**
 * 
 * @param {grpc.Client} grpcClient 
 * @returns {null | Object}
 */
export function resetAll(grpcClient) {
    const resp = grpcClient.invoke("pb.MerchantService/ResetAll", {});
    if (resp.status !== grpc.StatusOK) {
        logGrpcError(resp)
        return null
    }
    return resp.message
}