import express from "express";
import { createServer } from "http";
import { Server } from "socket.io";

const app = express();
const server = createServer(app);

const io = new Server(server);

io.on("connection", () => {});

server.listen(8080, () => {});
