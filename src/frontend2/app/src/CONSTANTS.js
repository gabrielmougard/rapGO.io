// TODO: here we put the definition of the constants
// using either environment variables or default values
export const CONVERTER_SERVER = (process.env.CONVERTER_SERVER) ?  process.env.CONVERTER_SERVER : "localhost:3001";
export const HEARTBEAT_WS_ENDPOINT = "ws://localhost:3002/heartbeat/" // in production, it  will be ws://hearbeat:3002/heartbeat/
// for the two last heartbeart, the order is not defined since these heartbeats are generated from two different services concurrently...
export const NB_STEPS_HEARTBEATS = 5
export const HEARTBEAT_TRIGGER_DOWNLOAD = "Rap generated successfully !";
export const NON_CHANGING_HEARTBEAT = "Saving generated data to cloud...";