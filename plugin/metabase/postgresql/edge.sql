/* name: CreateEdge :exec */
CREATE TABLE "edge"
(
    "id" bigserial NOT NULL PRIMARY KEY,
    "node_id"     varchar(255)  NOT NULL DEFAULT '',
    "inst_id"     varchar(255)  NOT NULL DEFAULT '',
    "inst_name"   varchar(255)  NOT NULL DEFAULT '',
    "address"     varchar(4096) NOT NULL DEFAULT '',
    "describe"    varchar(4096) NOT NULL DEFAULT '',
    "certificate" text          NOT NULL DEFAULT '',
    "status"      int           NOT NULL DEFAULT '0',
    "version"     int           NOT NULL DEFAULT '0',
    "auth_code"   varchar(4096) NOT NULL DEFAULT '',
    "extra"       varchar(4096) NOT NULL DEFAULT '',
    "expire_at"   timestamp     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "create_at"   timestamp     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "update_at"   timestamp     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "create_by"   varchar(255)  NOT NULL DEFAULT '',
    "update_by"   varchar(255)  NOT NULL DEFAULT '',
    "group"       varchar(255)  NOT NULL DEFAULT ''
);
CREATE UNIQUE INDEX if not exists "edge_uk_node_id" ON "edge" ("node_id");
CREATE INDEX if not exists "edge_idx_inst_id" ON "edge" ("inst_id");
CREATE INDEX if not exists "edge_idx_inst_name" ON "edge" ("inst_name");
COMMENT ON COLUMN "edge"."id" IS '自增ID';
COMMENT ON COLUMN "edge"."node_id" IS '节点编号';
COMMENT ON COLUMN "edge"."inst_id" IS '机构编号';
COMMENT ON COLUMN "edge"."inst_name" IS '机构名称';
COMMENT ON COLUMN "edge"."address" IS '节点地址';
COMMENT ON COLUMN "edge"."describe" IS '节点说明';
COMMENT ON COLUMN "edge"."certificate" IS '节点证书';
COMMENT ON COLUMN "edge"."status" IS '状态';
COMMENT ON COLUMN "edge"."version" IS '乐观锁版本';
COMMENT ON COLUMN "edge"."auth_code" IS '授权码';
COMMENT ON COLUMN "edge"."extra" IS '补充信息';
COMMENT ON COLUMN "edge"."expire_at" IS '过期时间';
COMMENT ON COLUMN "edge"."create_at" IS '创建时间';
COMMENT ON COLUMN "edge"."update_at" IS '更新时间';
COMMENT ON COLUMN "edge"."create_by" IS '创建人';
COMMENT ON COLUMN "edge"."update_by" IS '更新人';
COMMENT ON TABLE "edge" IS '网络节点联通表，两点之间最多一条边';

/* name: InsertEdge :execrows */
INSERT INTO edge ("node_id", "inst_id", "inst_name", "address", "describe", "certificate", "status", "version",
                          "auth_code", "extra", "expire_at", "create_at", "update_at", "create_by", "update_by",
                          "group")
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16);

/* name: DeleteEdge :exec */
DELETE FROM edge WHERE "node_id" = $1;

/* name: IndexEdge :many */
SELECT * FROM edge ORDER BY "node_id" ASC LIMIT $1 OFFSET $2;