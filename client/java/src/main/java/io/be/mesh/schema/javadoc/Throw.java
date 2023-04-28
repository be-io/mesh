package io.be.mesh.schema.javadoc;

import io.be.mesh.macro.Index;
import lombok.Data;

import java.io.Serializable;
import java.util.List;

@Data
public class Throw implements Serializable {

    private static final long serialVersionUID = -1690111767209524586L;
    @Index(value = 0, name = "0")
    private String name;
    @Index(value = 1, name = "1")
    private List<String> comments;
    @Index(value = 2, name = "2")
    private Kind kind;
}
