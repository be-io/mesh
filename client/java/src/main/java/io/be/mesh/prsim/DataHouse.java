package io.be.mesh.prsim;

import io.be.mesh.macro.MPI;
import io.be.mesh.macro.SPI;
import io.be.mesh.struct.Document;
import io.be.mesh.struct.Page;
import io.be.mesh.struct.Paging;

import java.util.List;

/**
 * @author jianyue.li
 */
@SPI("mesh")
public interface DataHouse {

    /**
     * 批量写入，主要使用
     */
    @MPI(value = "mesh.dh.writes", timeout = 3000, retries = 1)
    void writes(List<Document> docs);

    /**
     * 单条写入
     */
    @MPI(value = "mesh.dh.write", timeout = 3000, retries = 1)
    void write(Document doc);

    /**
     * 单条写入
     */
    @MPI("mesh.dh.read")
    Page<Object> read(Paging index);

    /**
     * Index list.
     */
    @MPI("mesh.dh.indies")
    Page<Object> indies(Paging index);

    /**
     * Talbe list.
     */
    @MPI("mesh.dh.tables")
    Page<Object> tables(Paging index);
}
