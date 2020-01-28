// TODO: here we put the definition of the constants
// using either environment variables or default values
export const CONVERTER_SERVER = (process.env.CONVERTER_SERVER) ?  process.env.CONVERTER_SERVER : "localhost:3001";
export const HEARTBEAT_WS_ENDPOINT = "ws://localhost:3002/heartbeat/" // in production, it  will be ws://hearbeat:3002/heartbeat/
// for the two last heartbeart, the order is not defined since these heartbeats are generated from two different services concurrently...
export const HEARTBEAT_TO_PROGRESSBAR = {
    "Reaching converter server ..." : 0,
    "Saving raw data to cloud..." : 17,
    "Internal error. The file couldn't be uploaded to the storage." : 17,
    "Starting core processing..." : 34,
    "Fetching metadata..." : 40,
    "Metadata fetched successfully !" : 70,
    "Rap generated successfully !" : 100,
    "Saving generated data to cloud..." : 100
}