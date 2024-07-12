import { createServer } from "http";
import { createSocket, onConnection } from "./socket";
import { redisConnect } from "./redis";

const boot = async () => {
    await redisConnect();

    const io = createSocket();
    io.listen(8080).of("/chat").on("connection", onConnection);
};

boot().catch((err) => {
    console.log(err);
});
