/* name: CreateClient :exec */
DROP TABLE IF EXISTS "oauth2_client";
CREATE TABLE "oauth2_client"
(
    "id"     varchar(255) NOT NULL PRIMARY KEY,
    "name"   varchar(255) NOT NULL,
    "secret" varchar(255) NOT NULL,
    "domain" varchar(255) NOT NULL,
    "data"   text         NOT NULL
);
COMMENT ON COLUMN "oauth2_client"."id" IS '客户端ID';
COMMENT ON COLUMN "oauth2_client"."name" IS '客户端名称';
COMMENT ON COLUMN "oauth2_client"."secret" IS '客户端密钥';
COMMENT ON COLUMN "oauth2_client"."domain" IS '客户端域名';
COMMENT ON COLUMN "oauth2_client"."data" IS '补充数据';
COMMENT ON TABLE "oauth2_client" IS 'OAuth2客户端表';

/* name: InsertClient :execrows */
INSERT INTO oauth2_client ("id", "name", "secret", "domain", "data")
VALUES ($1, $2, $3, $4, $5);

/* name: DeleteClient :exec */
DELETE FROM oauth2_client WHERE "id" = $1;

/* name: IndexClient :many */
SELECT * FROM oauth2_client ORDER BY "id" ASC LIMIT $1 OFFSET $2;