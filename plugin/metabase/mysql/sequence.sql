/* name: CreateSequence :exec */
CREATE TABLE `sequence`
(
    `kind`      varchar(255) NOT NULL DEFAULT '' COMMENT '序列号类型',
    `min`       bigint(20)   NOT NULL DEFAULT '0' COMMENT '当前范围最小值',
    `max`       bigint(20)   NOT NULL DEFAULT '0' COMMENT '当前范围最大值',
    `size`      int(11)      NOT NULL DEFAULT '0' COMMENT '每次取号段大小',
    `length`    int(11)      NOT NULL DEFAULT '0' COMMENT '序列号长度不足补零',
    `status`    int(11)      NOT NULL DEFAULT '0' COMMENT '状态',
    `version`   int(11)      NOT NULL DEFAULT '0' COMMENT '乐观锁版本',
    `create_at` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_at` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`kind`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8 COMMENT ='序列号表';

/* name: InsertSequence :execrows */
INSERT INTO `sequence` (`kind`, `min`, `max`, `size`, `length`, `status`, `version`, `create_at`, `update_at`)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);

/* name: GetSequenceByKind :one */
SELECT * FROM `sequence` WHERE `kind` = ?;

/* name: SetSequenceMin :execrows */
UPDATE `sequence` SET `min` = ?, `version` = `version` + 1 WHERE `kind` = ? AND `version` = ?;

/* name: GetSequenceByKindForUpdate :one */
SELECT * FROM `sequence` WHERE `kind` = ? FOR UPDATE;