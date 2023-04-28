/* name: CreateKV :exec */
CREATE TABLE "kv"
(
    "key"       varchar(255) NOT NULL PRIMARY KEY,
    "value"     text         NOT NULL,
    "create_at" timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "update_at" timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "create_by" varchar(255) NOT NULL DEFAULT '',
    "update_by" varchar(255) NOT NULL DEFAULT ''
);
COMMENT ON COLUMN "kv"."key" IS '配置KEY';
COMMENT ON COLUMN "kv"."value" IS '配置内容';
COMMENT ON COLUMN "kv"."create_at" IS '创建时间';
COMMENT ON COLUMN "kv"."update_at" IS '更新时间';
COMMENT ON COLUMN "kv"."create_by" IS '创建人';
COMMENT ON COLUMN "kv"."update_by" IS '更新人';
COMMENT ON TABLE "kv" IS 'KV配置表';

/* name: InsertKV :execrows */
INSERT INTO "kv" ("key", "value", "create_at", "update_at", "create_by", "update_by")
VALUES ($1, $2, $3, $4, $5, $6);

/* name: DeleteKV :execrows */
DELETE FROM "kv" WHERE "key" = $1;

/* name: IndexKV :many */
SELECT * FROM "kv" ORDER BY "key" ASC LIMIT $1 OFFSET $2;