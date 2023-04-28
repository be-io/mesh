package com.be.mesh.client.struct;

import lombok.Data;

import java.util.HashMap;
import java.util.Map;

/**
 * 存证输入
 *
 * @author lifeng
 */
@Data
public class LedgerRawTxInput {

    /**
     * 签发者：1号流通应用等
     */
    private String issueType;
    /**
     * 申请机构
     */
    private String applicant;
    /**
     * 申请机构id
     */
    private String applicantInst;
    /**
     * 发布机构
     */
    private String publisher;
    /**
     * 发布机构id
     */
    private String publisherInst;
    /**
     * 存证号
     */
    private String ledgerNo;
    /**
     * 存证名称
     */
    private String ledgerDesc;
    /**
     * 存证状态
     */
    private String issueStatus;
    /**
     * 存证日期
     */
    private String issueDate;
    /**
     * 存证内容 格式自定义
     */
    private Map<String, Object> content = new HashMap<>();
    /**
     * 存证可选命令
     * 保留为空
     */
    private Map<String, String> options = new HashMap<>();
}
