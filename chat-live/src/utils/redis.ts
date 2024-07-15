import { createClient } from "redis";

const redisConnect = async () => {
    Promise.all([pubClient.connect(), subClient.connect()]).then(() => console.log("redis connect"));
};

const pubClient = createClient({ url: "redis://localhost:6379" });
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
