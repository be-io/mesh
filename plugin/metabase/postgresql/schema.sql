DROP TABLE IF EXISTS "edge";
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

-- ---------------------------------------------------------------------------------------------------------------------
DROP TABLE IF EXISTS "sequence";
CREATE TABLE "sequence"
(
    "kind"      varchar(255) NOT NULL PRIMARY KEY,
    "min"       bigint       NOT NULL DEFAULT '0',
    "max"       bigint       NOT NULL DEFAULT '0',
    "size"      int          NOT NULL DEFAULT '0',
    "length"    int          NOT NULL DEFAULT '0',
    "status"    int          NOT NULL DEFAULT '0',
    "version"   int          NOT NULL DEFAULT '0',
    "create_at" timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "update_at" timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP
);
COMMENT ON COLUMN "sequence"."kind" IS '序列号类型';
COMMENT ON COLUMN "sequence"."min" IS '当前范围最小值';
COMMENT ON COLUMN "sequence"."max" IS '当前范围最大值';
COMMENT ON COLUMN "sequence"."size" IS '每次取号段大小';
COMMENT ON COLUMN "sequence"."length" IS '序列号长度不足补零';
COMMENT ON COLUMN "sequence"."status" IS '状态';
COMMENT ON COLUMN "sequence"."version" IS '乐观锁版本';
COMMENT ON COLUMN "sequence"."create_at" IS '创建时间';
COMMENT ON COLUMN "sequence"."update_at" IS '更新时间';
COMMENT ON TABLE "sequence" IS '序列号表';

-- ---------------------------------------------------------------------------------------------------------------------
DROP TABLE IF EXISTS "cooperator";
CREATE TABLE "cooperator"
(
    "id" bigserial NOT NULL PRIMARY KEY,
    "created_by"   varchar(255) NOT NULL DEFAULT '',
    "created_time" timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_by"   varchar(255) NOT NULL DEFAULT '',
    "updated_time" timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "is_deleted"   SMALLINT     NOT NULL DEFAULT '0',
    "is_disabled"  SMALLINT     NOT NULL DEFAULT '0',
    "tenant_id"    varchar(255) NOT NULL DEFAULT '',
    "version"      bigint       NOT NULL DEFAULT '0',
    "code"         varchar(255) NOT NULL DEFAULT '',
    "description"  varchar(255) NOT NULL DEFAULT '',
    "is_local"     SMALLINT     NOT NULL DEFAULT '',
    "name"         varchar(255) NOT NULL DEFAULT '',
    "node_id"      varchar(255) NOT NULL DEFAULT '',
    "group"        varchar(255) NOT NULL DEFAULT ''
);
CREATE INDEX if not exists "cooperator_idx_code" ON "cooperator" ("idx_code");
CREATE INDEX if not exists "cooperator_idx_node_id" ON "cooperator" ("node_id");
CREATE INDEX if not exists "cooperator_idx_tenant_id" ON "cooperator" ("tenant_id");
COMMENT ON COLUMN "cooperator"."id" IS '序列号';
COMMENT ON COLUMN "cooperator"."created_by" IS '创建者';
COMMENT ON COLUMN "cooperator"."created_time" IS '创建时间';
COMMENT ON COLUMN "cooperator"."updated_by" IS '更新者';
COMMENT ON COLUMN "cooperator"."updated_time" IS '更新时间';
COMMENT ON COLUMN "cooperator"."is_deleted" IS '是否删除';
COMMENT ON COLUMN "cooperator"."is_disabled" IS '是否禁用';
COMMENT ON COLUMN "cooperator"."tenant_id" IS '租户ID';
COMMENT ON COLUMN "cooperator"."version" IS '版本';
COMMENT ON COLUMN "cooperator"."code" IS '机构号';
COMMENT ON COLUMN "cooperator"."description" IS '机构说明';
COMMENT ON COLUMN "cooperator"."is_local" IS '是否本地机构';
COMMENT ON COLUMN "cooperator"."name" IS '机构名';
COMMENT ON COLUMN "cooperator"."node_id" IS '节点编号';
COMMENT ON COLUMN "cooperator"."group" IS '联盟中心节点机构id-多个用逗号分割';
COMMENT ON TABLE "cooperator" IS '机构表';

-- ---------------------------------------------------------------------------------------------------------------------
DROP TABLE IF EXISTS "kv";
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

-- ---------------------------------------------------------------------------------------------------------------------
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

-- ---------------------------------------------------------------------------------------------------------------------
DROP TABLE IF EXISTS "oauth2_token";
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

