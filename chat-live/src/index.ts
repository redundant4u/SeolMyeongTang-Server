import express from "express";
import { createServer } from "http";
import { createSocket, onConnection } from "./socket";

const app = express();
const server = createServer(app);

const io = createSocket(server);
io.of("/chat").on("connection", onConnection);

server.listen(8080, () => {});
