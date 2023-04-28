package com.be.mesh.client.system;

import com.be.mesh.client.annotate.SPI;
import com.be.mesh.client.mpc.ServiceProxy;
import com.be.mesh.client.prsim.HyperLedger;
import com.be.mesh.client.struct.LedgerRawTxInput;
import com.be.mesh.client.struct.LedgerTxReceipt;

@SPI("mesh")
public class MeshHyperLedger implements HyperLedger {

    private final HyperLedger hyperLedger = ServiceProxy.proxy(HyperLedger.class);

    @Override
    public LedgerTxReceipt sendRawTx(LedgerRawTxInput param) {
        return hyperLedger.sendRawTx(param);
    }
}
