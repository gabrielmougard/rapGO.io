// TODO: here we put the definition of the constants
// using either environment variables or default values
export const INGESTOR_SERVER = (process.env.INGESTOR_SERVER) ?  process.env.INGESTOR_SERVER : "localhost:3003";
export const INGESTOR_BUILD_VERSION = (process.env.INGESTOR_BUILD_VERSION) ?  process.env.INGESTOR_BUILD_VERSION : "v0.1";