/* name: CreateSequence :exec */
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

/* name: InsertSequence :execrows */
INSERT INTO sequence ("kind", "min", "max", "size", "length", "status", "version", "create_at", "update_at")
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);

/* name: GetSequenceByKind :one */
SELECT *
FROM sequence
WHERE "kind" = $1;

/* name: SetSequenceMin :execrows */
UPDATE sequence
SET "min"     = $1,
    "version" = "version" + 1
WHERE "kind" = $2
  AND "version" = $3;

/* name: GetSequenceByKindForUpdate :one */
SELECT *
FROM sequence
WHERE "kind" = $1 FOR UPDATE;