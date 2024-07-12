import { createClient } from "redis";

const redisConnect = async () => {
    await Promise.all([pubClient.connect(), subClient.connect()]);
};

const pubClient = createClient({ url: "redis://localhost:6379" });
const subClient = pubClient.duplicate();

pubClient.on("error", (err) => {
    console.log(`redis connection failed: ${err}`);
    process.exit(1);
});
subClient.on("error", (err) => {
    console.log(`redis connection failed: ${err}`);
    process.exit(1);
});

export { redisConnect, pubClient, subClient };
