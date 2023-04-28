package com.be.mesh.client.schema.javadoc;

import com.be.mesh.client.annotate.Index;
import lombok.Data;

import java.io.Serializable;
import java.util.List;
import java.util.Map;

@Data
public class Parameter implements Serializable {

    private static final long serialVersionUID = -4712741489649117836L;
    @Index(value = 0, name = "0")
    private String name;
    @Index(value = 1, name = "1")
    private Map<String, Map<String, String>> macros;
    @Index(value = 2, name = "2")
    private List<String> comments;
    @Index(value = 3, name = "3")
    private Kind kind;
    @Index(value = 4, name = "4")
    private String value;
}
