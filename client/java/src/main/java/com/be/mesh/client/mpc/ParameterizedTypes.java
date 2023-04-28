/*
 * Copyright (c) 2000, 2023, ducesoft and/or its affiliates. All rights reserved.
 * DUCESOFT PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.mpc;

import com.be.mesh.client.tool.Tool;
import lombok.AllArgsConstructor;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.io.Serializable;
import java.lang.reflect.ParameterizedType;
import java.lang.reflect.Type;
import java.util.*;
import java.util.function.Function;
import java.util.stream.Collectors;
import java.util.stream.Stream;

/**
 * @author coyzeng@gmail.com
 */
@Data
@NoArgsConstructor
@AllArgsConstructor
public class ParameterizedTypes implements Serializable {

    private static final long serialVersionUID = 3876574260512549174L;
    private static final Map<String, String> symbols = new HashMap<>();

    static {
        symbols.put("byte", "B");
        symbols.put("char", "C");
        symbols.put("double", "D");
        symbols.put("float", "F");
        symbols.put("int", "I");
        symbols.put("long", "J");
        symbols.put("short", "S");
        symbols.put("boolean", "Z");
        symbols.put("void", "V");

        symbols.put("[B", "[B");
        symbols.put("[C", "[C");
        symbols.put("[D", "[D");
        symbols.put("[F", "[F");
        symbols.put("[I", "[I");
        symbols.put("[J", "[J");
        symbols.put("[S", "[S");
        symbols.put("[Z", "[Z");
        symbols.put("[V", "[V");
    }

    public ParameterizedTypes(String raw) {
        this.raw = raw;
    }

    public ParameterizedTypes(Type type) {
        if (null == type) {
            return;
        }
        if (!(type instanceof ParameterizedType)) {
            this.raw = type instanceof Class ? ((Class<?>) type).getName() : type.getTypeName();
            return;
        }
        ParameterizedType pt = (ParameterizedType) type;
        this.raw = pt.getRawType() instanceof Class ? ((Class<?>) pt.getRawType()).getName() : pt.getRawType().getTypeName();
        for (Type argument : pt.getActualTypeArguments()) {
            if (argument instanceof ParameterizedType) {
                args.add(new ParameterizedTypes(argument));
            } else if (argument instanceof Class<?>) {
                args.add(new ParameterizedTypes(((Class<?>) argument).getName()));
            } else {
                args.add(new ParameterizedTypes(argument.getTypeName()));
            }
        }
    }

    private List<ParameterizedTypes> args = new ArrayList<>();
    private String raw;

    /**
     * 基本类型终结符(这个与方法签名就比较亲近了)
     * <p>
     * boolean 的终结符不为 B 是让 byte 占用
     * <p>
     * long 的终结符不是 L 也是因为让对象终结符占用
     * <p>
     * 返回值为 void 类型的终结符是 V，对象类型终结符为 L 和 ;
     * <p>
     * 数组类型终结符 [
     * <p>
     * Ljava/util/Map<Ljava/lang/String;Ljava/util/List<Ljava/util/Map<Ljava/lang/String;Lcom/be/mesh/client/struct/Principal;>;>;>;
     */
    public String toSignature() {
        StringBuilder signature = new StringBuilder();
        this.recursionSignature(this, signature);
        return signature.toString();
    }

    public String loopSignature() {
        StringBuilder signature = new StringBuilder();
        Deque<ParameterizedTypes> queue = new ArrayDeque<>();
        queue.addLast(this);
        while (!queue.isEmpty()) {
            ParameterizedTypes ref = queue.removeLast();
            for (int index = 0; index < ref.args.size(); index++) {
                queue.addLast(ref.args.get(ref.args.size() - index - 1));
            }
            signature.append(toSymbol(ref.getRaw()));
            if (Tool.optional(ref.getArgs())) {
                signature.append(';');
            } else {
                signature.append('<');
            }
        }
        return signature.toString();
    }

    private void recursionSignature(ParameterizedTypes ref, StringBuilder signature) {
        signature.append(toSymbol(ref.getRaw()));
        if (Tool.optional(ref.getArgs())) {
            if (!isArray(ref.getRaw()) && !isNative(ref.getRaw())) {
                signature.append(';');
            }
        } else {
            signature.append('<');
        }
        for (ParameterizedTypes arg : ref.getArgs()) {
            recursionSignature(arg, signature);
        }
        if (Tool.required(ref.getArgs())) {
            signature.append('>').append(';');
        }
    }

    /**
     * return String.format("T%s", raw.replace(".", "/"));
     * <p>
     * Get signature of primitive type.
     * <pre>
     *   B	    byte
     *   C	    char
     *   D	    double
     *   F	    float
     *   I	    int
     *   J	    long
     *   S	    short
     *   Z	    boolean
     *   T      Generic
     * </pre>
     */
    public String toSymbol(String raw) {
        if (Tool.required(symbols.get(raw))) {
            return symbols.get(raw);
        }
        if (isArray(raw)) {
            return raw.replace(".", "/");
        }
        try {
            Class.forName(raw);
            return String.format("L%s", raw.replace(".", "/"));
        } catch (ClassNotFoundException e) {
            return String.format("L%s", "java/lang/Object");
        }
    }

    private boolean isNative(String raw) {
        if (Tool.required(symbols.get(raw))) {
            return true;
        }
        if (isArray(raw)) {
            return isNative(raw.substring(1));
        }
        return false;
    }

    private boolean isArray(String raw) {
        return raw.startsWith("[");
    }

    public Stream<String> toGenerics() {
        Set<String> types = new HashSet<>();
        Deque<ParameterizedTypes> queue = new ArrayDeque<>();
        queue.push(this);
        while (!queue.isEmpty()) {
            ParameterizedTypes type = queue.pop();
            types.add(type.getRaw());
            type.getArgs().forEach(queue::push);
        }
        return types.stream().distinct().filter(x -> x.startsWith("T"));
    }

    public boolean isParameterized() {
        return null != this.args && !args.isEmpty();
    }

    public String toCanonicalName() {
        return this.toCanonicalName(x -> x);
    }

    public String toCanonicalName(Function<String, String> former) {
        StringBuilder names = new StringBuilder();
        names.append(former.apply(this.raw));
        if (Tool.required(this.args)) {
            names.append('<');
            names.append(this.args.stream().map(x -> x.toCanonicalName(former)).collect(Collectors.joining(",")));
            names.append('>');
        }
        return names.toString();
    }

    public List<String> toList() {
        List<String> as = new ArrayList<>();
        Deque<ParameterizedTypes> queue = new ArrayDeque<>();
        queue.push(this);
        while (null != queue.peek()) {
            ParameterizedTypes types = queue.pop();
            if (Tool.optional(types.getRaw())) {
                continue;
            }
            as.add(types.getRaw());
            if (Tool.required(types.getArgs())) {
                for (ParameterizedTypes arg : types.getArgs()) {
                    queue.push(arg);
                }
            }
        }
        return as;
    }

    public static String getRawName(Type type) {
        if (type instanceof Class) {
            return ((Class<?>) type).getName();
        }
        if (type instanceof ParameterizedType) {
            return getRawName(((ParameterizedType) type).getRawType());
        }
        return type.getTypeName();
    }

    public static ParameterizedTypes ofTree(String signature) {
        if (!Tool.contains(signature, "<") || Tool.optional(signature)) {
            ParameterizedTypes types = new ParameterizedTypes();
            types.setRaw(signature);
            return types;
        }
        int offset = signature.indexOf("<");
        List<ParameterizedTypes> childrens = new ArrayList<>();
        String[] vs = new String[]{Tool.substring(signature, 0, offset), Tool.substring(signature, offset + 1, -1)};
        String cs = Tool.substring(vs[1], 0, vs[1].length());
        StringBuilder chars = new StringBuilder();
        int l = 0, r = 0;
        for (int index = 0; index < cs.length(); index++) {
            if ('<' == cs.charAt(index)) {
                l++;
            }
            if ('>' == cs.charAt(index)) {
                r++;
            }
            if (',' != cs.charAt(index) || chars.length() > 0) {
                chars.append(cs.charAt(index));
            }
            if (l == r && 0 != l) {
                int lof = signature.indexOf("<");
                String substring = chars.toString();
                String[] cps = new String[]{Tool.substring(substring, 0, lof), Tool.substring(substring, lof + 1, -1)};
                String[] cpp = Tool.split(cps[0], ",");
                if (cpp.length > 1) {
                    for (int idx = 0; idx < cpp.length - 1; idx++) {
                        ParameterizedTypes types = new ParameterizedTypes();
                        types.setRaw(cpp[idx]);
                        childrens.add(types);
                    }
                    ParameterizedTypes types = new ParameterizedTypes();
                    types.setRaw(String.format("%s<%s", cpp[cpp.length - 1], cps[1]));
                    childrens.add(types);
                } else {
                    childrens.add(ofTree(substring));
                }
                l = 0;
                r = 0;
                chars.delete(0, chars.length());
                continue;
            }
            if (index == cs.length() - 1 && chars.length() > 0) {
                for (String expr : Tool.split(chars.toString(), ",")) {
                    ParameterizedTypes types = new ParameterizedTypes();
                    types.setRaw(expr);
                    childrens.add(types);
                }
            }
        }
        ParameterizedTypes types = new ParameterizedTypes();
        types.setRaw(vs[0]);
        types.setArgs(childrens);
        return types;
    }
}
