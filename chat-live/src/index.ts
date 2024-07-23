import { createSocket, onConnection } from "./socket";
import { redisConnect } from "./utils/redis";
import { pgConnect } from "./utils/db";
import env from "./utils/env";

const boot = async () => {
    await pgConnect();
    await redisConnect();

    const io = createSocket().path(env.CHAT_WS_PATH);
    io.listen(env.PORT).of("/chat").on("connection", onConnection);
};

boot().catch((err) => {
    console.error("[ERROR] boot failed");
    process.exit(1);
});
