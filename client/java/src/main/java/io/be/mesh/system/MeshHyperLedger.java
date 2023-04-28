package io.be.mesh.system;

import io.be.mesh.macro.SPI;
import io.be.mesh.mpc.ServiceProxy;
import io.be.mesh.prsim.HyperLedger;
import io.be.mesh.struct.LedgerRawTxInput;
import io.be.mesh.struct.LedgerTxReceipt;

@SPI("mesh")
public class MeshHyperLedger implements HyperLedger {

    private final HyperLedger hyperLedger = ServiceProxy.proxy(HyperLedger.class);

    @Override
    public LedgerTxReceipt sendRawTx(LedgerRawTxInput param) {
        return hyperLedger.sendRawTx(param);
    }
}
