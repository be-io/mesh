package com.be.mesh.client.struct;

import lombok.Data;

/**
 * @author lifeng
 */
@Data
public class LedgerTxReceipt {

    private String transactionHash;

    private String transactionIndex;

    private String blockHash;

    private String blockNumber;

    private String gasUsed;

    private String root;

    private int status;

    private byte[] contractAddress;

    private String from;

    private String to;

    private String input;

    private String output;

    private String logsBloom;
}
