package com.be.mesh.client.prsim;

import com.be.mesh.client.annotate.MPI;
import com.be.mesh.client.annotate.SPI;
import com.be.mesh.client.struct.OSChart;

@SPI("mesh")
public interface OperateSystem {

    /**
     * install
     *
     * @param chart
     */
    @MPI("mesh.os.install")
    void install(OSChart chart);

    /**
     * uninstall
     *
     * @param chart
     */
    @MPI("mesh.os.uninstall")
    void uninstall(OSChart chart);

}
