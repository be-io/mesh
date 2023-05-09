CREATE TABLE `edge`
(
    `id`          bigint(20)    NOT NULL AUTO_INCREMENT COMMENT '自增ID',
    `node_id`     varchar(255)  NOT NULL DEFAULT '' COMMENT '节点编号',
    `inst_id`     varchar(255)  NOT NULL DEFAULT '' COMMENT '机构编号',
    `inst_name`   varchar(255)  NOT NULL DEFAULT '' COMMENT '机构名称',
    `address`     varchar(4096) NOT NULL DEFAULT '' COMMENT '节点地址',
    `describe`    varchar(4096) NOT NULL DEFAULT '' COMMENT '节点说明',
    `certificate` text          NOT NULL COMMENT '节点证书',
    `status`      int(11)       NOT NULL DEFAULT '0' COMMENT '状态',
    `version`     int(11)       NOT NULL DEFAULT '0' COMMENT '乐观锁版本',
    `auth_code`   varchar(4096) NOT NULL DEFAULT '' COMMENT '授权码',
    `extra`       varchar(4096) NOT NULL DEFAULT '' COMMENT '补充信息',
    `expire_at`   datetime      NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '过期时间',
    `create_at`   datetime      NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_at`   datetime      NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
    `create_by`   varchar(255)  NOT NULL DEFAULT '' COMMENT '创建人',
    `update_by`   varchar(255)  NOT NULL DEFAULT '' COMMENT '更新人',
    `group`       varchar(255)  NOT NULL DEFAULT '' COMMENT '联盟中心节点机构id-多个用逗号分割',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_node_id` (`node_id`),
    KEY `idx_inst_id` (`inst_id`),
    KEY `idx_inst_name` (`inst_name`)
) ENGINE = InnoDB
  AUTO_INCREMENT = 10
  DEFAULT CHARSET = utf8 COMMENT ='网络节点联通表，两点之间最多一条边';

-- ---------------------------------------------------------------------------------------------------------------------
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

-- ---------------------------------------------------------------------------------------------------------------------
CREATE TABLE `cooperator`
(
    `id`           bigint(20)   NOT NULL AUTO_INCREMENT,
    `created_by`   varchar(255) NOT NULL DEFAULT '' COMMENT '创建者',
    `created_time` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_by`   varchar(255) NOT NULL DEFAULT '' COMMENT '更新者',
    `updated_time` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '更新时间',
    `is_deleted`   tinyint(4)   NOT NULL DEFAULT '0' COMMENT '是否删除',
    `is_disabled`  tinyint(4)   NOT NULL DEFAULT '0' COMMENT '是否禁用',
    `tenant_id`    varchar(255) NOT NULL DEFAULT '' COMMENT '租户ID',
    `version`      int(11)      NOT NULL DEFAULT '0' COMMENT '版本',
    `code`         varchar(255) NOT NULL DEFAULT '' COMMENT '机构号',
    `description`  varchar(255) NOT NULL DEFAULT '' COMMENT '机构说明',
    `is_local`     int(11)      NOT NULL DEFAULT 0 COMMENT '是否本地机构',
    `name`         varchar(255) NOT NULL DEFAULT '' COMMENT '机构名',
    `node_id`      varchar(255) NOT NULL DEFAULT '' COMMENT '节点编号',
    `group`        varchar(255) NOT NULL DEFAULT '' COMMENT '联盟中心节点机构id-多个用逗号分割',
    PRIMARY KEY (`id`),
    KEY `idx_code` (`code`),
    KEY `idx_node_id` (`node_id`),
    KEY `idx_tenant_id` (`tenant_id`)
) ENGINE = InnoDB
  AUTO_INCREMENT = 26
  DEFAULT CHARSET = utf8 COMMENT '机构表';

-- ---------------------------------------------------------------------------------------------------------------------
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

-- ---------------------------------------------------------------------------------------------------------------------
CREATE TABLE oauth2_client
(
    `id`     VARCHAR(255) NOT NULL COMMENT '客户端ID',
    `name`   VARCHAR(255) NOT NULL COMMENT '客户端名称',
    `secret` VARCHAR(255) NOT NULL COMMENT '客户端密钥',
    `domain` VARCHAR(255) NOT NULL COMMENT '客户端域名',
    `data`   TEXT         NOT NULL COMMENT '补充数据',
    PRIMARY KEY (`id`)
) DEFAULT CHARSET = utf8 COMMENT ='OAuth2客户端表';

-- ---------------------------------------------------------------------------------------------------------------------
CREATE TABLE oauth2_token
(
    `code`      VARCHAR(255) NOT NULL COMMENT '授权码',
    `access`    VARCHAR(255) NOT NULL COMMENT '准入TOKEN',
    `refresh`   VARCHAR(255) NOT NULL COMMENT '刷新TOKEN',
    `data`      TEXT         NOT NULL COMMENT '补充数据',
    `create_at` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `expire_at` datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '过期时间',
    PRIMARY KEY (`code`),
    UNIQUE KEY `udx_access` (`access`),
    UNIQUE KEY `udx_refresh` (`refresh`),
    KEY `idx_expire_at` (`expire_at`)
) DEFAULT CHARSET = utf8 COMMENT ='OAuth2令牌表';