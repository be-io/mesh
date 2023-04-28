package com.be.mesh.client.schema.javadoc;

import com.be.mesh.client.annotate.Index;
import lombok.Data;

import java.io.Serializable;
import java.util.*;

@Data
public class Kind implements Serializable {

    private static final long serialVersionUID = -4704612630848058134L;
    @Index(value = 0, name = "0")
    private String pkg;
    @Index(value = 1, name = "1")
    private String name;
    @Index(value = 2, name = "2")
    private List<String> imports;
    @Index(value = 3, name = "3")
    private Map<String, Map<String, String>> macros;
    @Index(value = 4, name = "4")
    private List<String> comments;
    @Index(value = 5, name = "5")
    private int modifier; // {@link java.lang.reflect.Modifier}
    @Index(value = 6, name = "6")
    private List<Variable> variables;
    @Index(value = 7, name = "7")
    private List<Method> methods;
    @Index(value = 8, name = "8")
    private List<Kind> supers;
    @Index(value = 9, name = "9")
    private List<Kind> traits;
    @Index(value = 10, name = "10")
    private String signature;

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (!(o instanceof Kind)) return false;
        return Objects.equals(getPkg(), ((Kind) o).getPkg()) && Objects.equals(getName(), ((Kind) o).getName());
    }

    @Override
    public int hashCode() {
        return Objects.hash(getPkg(), getName());
    }

    public Kind copy() {
        Kind kind = new Kind();
        kind.setPkg(this.getPkg());
        kind.setName(this.getName());
        kind.setImports(new ArrayList<>(this.getImports()));
        kind.setMacros(new HashMap<>(this.getMacros()));
        kind.setComments(new ArrayList<>(this.getComments()));
        kind.setModifier(this.getModifier());
        kind.setVariables(new ArrayList<>());
        kind.setMethods(new ArrayList<>());
        kind.setSupers(new ArrayList<>());
        kind.setTraits(new ArrayList<>());
        kind.setSignature(this.getSignature());
        return kind;
    }
}
