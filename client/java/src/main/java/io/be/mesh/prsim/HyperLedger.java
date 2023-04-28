package io.be.mesh.prsim;

import io.be.mesh.macro.MPI;
import io.be.mesh.struct.LedgerRawTxInput;
import io.be.mesh.struct.LedgerTxReceipt;

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
