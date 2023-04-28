package io.be.mesh.system;

import io.be.mesh.macro.SPI;
import io.be.mesh.mpc.ServiceProxy;
import io.be.mesh.prsim.DataHouse;
import io.be.mesh.struct.Document;
import io.be.mesh.struct.Page;
import io.be.mesh.struct.Paging;
import org.springframework.util.CollectionUtils;

import java.util.List;

/**
 * @author jianyue.li
 * @date 2022/3/14 7:47 PM
 */
@SPI("mesh")
public class MeshDataHouse implements DataHouse {

    private final DataHouse dataHouse = ServiceProxy.proxy(DataHouse.class);

    @Override
    public void writes(List<Document> docs) {
        if (!CollectionUtils.isEmpty(docs)) {
            dataHouse.writes(docs);
        }
    }

    @Override
    public void write(Document doc) {
        if (doc != null) {
            dataHouse.write(doc);
        }
    }

    @Override
    public Page<Object> read(Paging index) {
        return dataHouse.read(index);
    }

    @Override
    public Page<Object> indies(Paging index) {
        return dataHouse.indies(index);
    }

    @Override
    public Page<Object> tables(Paging index) {
        return dataHouse.tables(index);
    }

}
