import { createSocket, onConnection } from "./socket";
import { redisConnect } from "./utils/redis";
import { pgConnect } from "./utils/db";

const boot = async () => {
    await pgConnect();
    await redisConnect();

    const io = createSocket();
    io.listen(8080).of("/chat").on("connection", onConnection);
};

boot().catch((err) => {
    console.error("[ERROR] boot failed");
    process.exit(1);
});
