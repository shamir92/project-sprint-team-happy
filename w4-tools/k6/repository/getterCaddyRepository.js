import grpc from "k6/net/grpc";

/**
 * 
 * @param {grpc.Client} grpcClient
 * @returns {null | import("../entity/merchant").AllMerchantNearestRecord}
 */
export function getAllMerchantNearestLocations(grpcClient) {
    const resp = grpcClient.invoke("pb.MerchantService/GetAllMerchantNearestLocations", {});
    if (resp.status !== grpc.StatusOK) {
        logGrpcError(resp)
        return null
    }
    return resp.message
}

/**
 * Retrieves the nearest locations of a merchant.
 * 
 * @param {grpc.Client} grpcClient - The gRPC client.
 * 
 * @returns {import("../entity/merchant").MerchantNearestRecord|null} - The nearest locations of the merchant or null if there was an error.
 */
export function getMerchantNearestLocations(grpcClient) {
    const resp = grpcClient.invoke("pb.MerchantService/GetMerchantNearestLocations", {});
    if (resp.status !== grpc.StatusOK) {
        logGrpcError(resp);
        return null;
    }
    return resp.message;
}


/**
 * 
 * @param {grpc.Client} grpcClient 
 * @returns {null | import("../entity/merchant").AllGeneratedRoutes}
 */
export function getAllMerchantRoutes(grpcClient) {
    const resp = grpcClient.invoke("pb.MerchantService/GetAllMerchantRoutes", {});
    if (resp.status !== grpc.StatusOK) {
        logGrpcError(resp)
        return null
    }
    return resp.message
}

/**
 * Retrieves the routes for two-zone merchants.
 * 
 * @param {grpc.Client} grpcClient - The gRPC client.
 * 
 * @returns {null | import("../entity/merchant").AllGeneratedRoutes} - The routes for two-zone merchants or null if there was an error.
 */
export function getTwoZoneMerchantRoutes(grpcClient) {
    const resp = grpcClient.invoke("pb.MerchantService/GetTwoZoneMerchantRoutes", {});
    if (resp.status !== grpc.StatusOK) {
        logGrpcError(resp);
        return null;
    }
    return resp.message;
}

/**
 * Retrieves the routes for a merchant.
 * 
 * @param {grpc.Client} grpcClient - The gRPC client.
 * 
 * @returns {null | import("../entity/merchant").RouteZone} - The routes for the merchant or null if there was an error.
 */
export function getMerchantRoutes(grpcClient) {
    const resp = grpcClient.invoke("pb.MerchantService/GetMerchantRoutes", {});
    if (resp.status !== grpc.StatusOK) {
        logGrpcError(resp);
        return null;
    }
    return resp.message;
}



/**
 * 
 * @param {grpc.Client} grpcClient 
 * @returns {null | import("../entity/merchant").PregeneratedMerchant}
 */
export function getAllPregeneratedMerchants(grpcClient) {
    const resp = grpcClient.invoke("pb.MerchantService/GetAllPregeneratedMerchants", {});
    if (resp.status !== grpc.StatusOK) {
        logGrpcError(resp)
        return null
    }
    return resp.message
}

/**
 * 
 * @param {grpc.Client} grpcClient 
 * @param {import("../entity/merchant").AssignMerchant} assignMerchant
 * @returns {null | Object}
 */
export function assignPregeneratedMerchant(grpcClient, assignMerchant) {
    const resp = grpcClient.invoke("pb.MerchantService/AssignPregeneratedMerchant", assignMerchant);
    if (resp.status !== grpc.StatusOK) {
        logGrpcError(resp)
        return null
    }
    return resp.message
}


/**
 * Retrieves a pregenerated merchant.
 * 
 * @param {grpc.Client} grpcClient - The gRPC client.
 * 
 * @returns {import("../entity/merchant").Merchant|null} - The pregenerated merchant or null if there was an error.
 */
export function getPregeneratedMerchant(grpcClient) {
    const resp = grpcClient.invoke("pb.MerchantService/GetPregeneratedMerchant", {});
    if (resp.status !== grpc.StatusOK) {
        logGrpcError(resp);
        return null;
    }
    return resp.message;
}

/**
 * 
 * @param {import("k6/net/grpc").Response} resp 
 */
export function logGrpcError(resp) {
    console.error(`Error grpc: ${JSON.stringify(resp.error)}, msg: ${JSON.stringify(resp.message)}`);
}