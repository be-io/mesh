package com.be.mesh.client.system;

import com.be.mesh.client.annotate.SPI;
import com.be.mesh.client.mpc.ServiceProxy;
import com.be.mesh.client.prsim.OperateSystem;
import com.be.mesh.client.struct.OSChart;
import lombok.extern.slf4j.Slf4j;

@Slf4j
@SPI("mesh")
public class MeshOperateSystem implements OperateSystem {

    private final OperateSystem operateSystem = ServiceProxy.proxy(OperateSystem.class);

    @Override
    public void install(OSChart chart) {
        this.operateSystem.install(chart);
    }

    @Override
    public void uninstall(OSChart chart) {
        this.operateSystem.uninstall(chart);
    }

}
