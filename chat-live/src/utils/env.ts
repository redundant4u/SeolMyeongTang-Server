import { z } from "zod";
import * as dotenv from "dotenv";

dotenv.config();

const envSchema = z.object({
    NODE_ENV: z.enum(["production", "development"]),
    PORT: z.string().transform((val) => parseInt(val, 10)),
    CHAT_WS_PATH: z.string(),
    DB_HOST: z.string(),
    DB_USER: z.string(),
    DB_PASSWORD: z.string(),
    DB_NAME: z.string(),
    DB_PORT: z.string().transform((val) => parseInt(val, 10)),
    REDIS_HOST: z.string(),
    REDIS_PORT: z.string().transform((val) => parseInt(val, 10)),
});

const env = envSchema.safeParse(process.env);
if (!env.success) {
    console.error("[ERROR] env schema is wrong");
    process.exit(1);
}

export default env.data;
