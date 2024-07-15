import { Pool, QueryResult, QueryResultRow } from "pg";

const pool = new Pool({
    host: "localhost",
    user: "wjstjf",
    password: "wjstjf",
    database: "wjstjf",
    port: 5432,
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

const query = async <T extends QueryResultRow>(q: string, values?: any[]): Promise<T[]> =>
    pool.query<T>(q, values).then((result) => result.rows);

export { pgConnect, query };
