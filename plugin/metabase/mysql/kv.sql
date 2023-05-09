/* name: CreateKV :exec */
CREATE TABLE `kv`
(
    `key`       varchar(255) NOT NULL COMMENT '配置KEY',
    `value`     text         NOT NULL COMMENT '配置内容',
    `create_at` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_at` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
    `create_by` varchar(255) NOT NULL DEFAULT '' COMMENT '创建人',
    `update_by` varchar(255) NOT NULL DEFAULT '' COMMENT '更新人',
    PRIMARY KEY (`key`)
) DEFAULT CHARSET = utf8 COMMENT ='KV配置表';

/* name: InsertKV :execrows */
INSERT INTO `kv` (`key`, `value`, `create_at`, `update_at`, `create_by`, `update_by`)
VALUES (?, ?, ?, ?, ?, ?);

/* name: DeleteKV :execrows */
DELETE FROM `kv` WHERE `key` = ?;

/* name: IndexKV :many */
SELECT * FROM `kv` ORDER BY `key` ASC LIMIT ?, ?;