/* name: CreateToken :exec */
CREATE TABLE "oauth2_token"
(
    "code"      varchar(255) NOT NULL PRIMARY KEY,
    "access"    varchar(255) NOT NULL,
    "refresh"   varchar(255) NOT NULL,
    "data"      text         NOT NULL,
    "create_at" timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "expire_at" timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX if not exists "oauth2_token_udx_access" ON "oauth2_token" ("access");
CREATE UNIQUE INDEX if not exists "oauth2_token_udx_refresh" ON "oauth2_token" ("refresh");
CREATE INDEX if not exists "oauth2_token_idx_expire_at" ON "oauth2_token" ("expire_at");
COMMENT ON COLUMN "oauth2_token"."code" IS '授权码';
COMMENT ON COLUMN "oauth2_token"."access" IS '准入TOKEN';
COMMENT ON COLUMN "oauth2_token"."refresh" IS '刷新TOKEN';
COMMENT ON COLUMN "oauth2_token"."data" IS '补充数据';
COMMENT ON COLUMN "oauth2_token"."create_at" IS '创建时间';
COMMENT ON COLUMN "oauth2_token"."expire_at" IS '过期时间';
COMMENT ON TABLE "oauth2_token" IS 'OAuth2令牌表';

/* name: InsertToken :execrows */
INSERT INTO oauth2_token ("create_at", "expire_at", "code", "access", "refresh", "data")
VALUES ($1, $2, $3, $4, $5, $6);

/* name: DeleteToken :exec */
DELETE FROM oauth2_token WHERE "code" = $1;


/* name: IndexToken :many */
SELECT * FROM oauth2_token ORDER BY "code" ASC LIMIT $1 OFFSET $2;