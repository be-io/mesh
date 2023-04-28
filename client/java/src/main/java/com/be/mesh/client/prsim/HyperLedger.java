package com.be.mesh.client.prsim;

import com.be.mesh.client.annotate.MPI;
import com.be.mesh.client.struct.LedgerRawTxInput;
import com.be.mesh.client.struct.LedgerTxReceipt;

/**
 * @author lifeng
 */
public interface HyperLedger {

    /**
     * raw tx sender
     *
     * @param param
     * @return
     */
    @MPI("mesh.ledger.tx.sendRaw")
    LedgerTxReceipt sendRawTx(LedgerRawTxInput param);

}
