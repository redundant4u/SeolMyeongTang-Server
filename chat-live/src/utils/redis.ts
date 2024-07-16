import { createClient } from "redis";
import env from "./env";

const redisConnect = async () => {
    Promise.all([pubClient.connect(), subClient.connect()]).then(() => console.log("redis connect"));
};

const pubClient = createClient({ url: `redis://${env.REDIS_HOST}:${env.REDIS_PORT}` });
const subClient = pubClient.duplicate();

pubClient.on("error", (err) => {
    console.error(`[ERROR] ${err}`);
    process.exit(1);
});
subClient.on("error", (err) => {
    console.error(`[ERROR] ${err}`);
    process.exit(1);
});

export { redisConnect, pubClient, subClient };
