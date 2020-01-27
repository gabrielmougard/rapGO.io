// TODO: here we put the definition of the constants
// using either environment variables or default values
export const CONVERTER_SERVER = (process.env.CONVERTER_SERVER) ?  process.env.CONVERTER_SERVER : "localhost:3001";
export const HEARTBEAT_WS_ENDPOINT = "ws://localhost:3004/heartbeat/" // in production, it  will be ws://hearbeat:3004/heartbeat/
export const HEARTBEAT_TO_PROGRESSBAR = {
    "Reaching converter server ..." : 0,
    "heartbeat1" : 10,
    "heartbeat2" : 20,
    "heartbeat3" : 30,
    "heartbeat4" : 40,
    "heartbeat5" : 50,
    "heartbeat6" : 60,
    "heartbeat7" : 70,
    "heartbeat8" : 80,
    "heartbeat9" : 90,
    "heartbeat10" : 100
}