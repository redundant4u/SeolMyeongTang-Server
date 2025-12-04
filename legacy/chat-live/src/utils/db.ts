import { Pool, QueryResultRow } from "pg";
import env from "./env";

const pool = new Pool({
    host: env.DB_HOST,
    user: env.DB_USER,
    password: env.DB_PASSWORD,
    database: env.DB_NAME,
    port: env.DB_PORT,
    max: 10,
});

const pgConnect = async () => {
    pool.connect()
        .then(() => console.log("db connect"))
        .catch((err) => {
            console.error(`[ERROR] ${err}`);
            process.exit(1);
        });
};

/* eslint-disable  @typescript-eslint/no-explicit-any */
const query = async <T extends QueryResultRow>(q: string, values?: any[]): Promise<T[]> =>
    pool.query<T>(q, values).then((result) => result.rows);

export { pgConnect, query };
