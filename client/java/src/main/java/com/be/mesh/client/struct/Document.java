package com.be.mesh.client.struct;

import com.be.mesh.client.annotate.Index;
import lombok.Data;

import java.io.Serializable;
import java.util.HashMap;
import java.util.Map;

/**
 * @author jianyue.li
 * @date 2022/3/14 7:40 PM
 */
@Data
public class Document implements Serializable {

    @Index(0)
    private Map<String, String> metadata = new HashMap<>();

    @Index(5)
    private String content;

    @Index(10)
    private Long timestamp = System.currentTimeMillis();

}
