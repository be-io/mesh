package com.be.mesh.client.schema.javadoc;

import com.be.mesh.client.mpc.ParameterizedTypes;
import com.be.mesh.client.schema.context.CompileContext;
import com.be.mesh.client.schema.macro.JavaDoc;
import com.be.mesh.client.tool.Once;
import com.be.mesh.client.tool.Tool;
import lombok.Getter;

import javax.lang.model.element.*;
import javax.lang.model.type.ArrayType;
import javax.lang.model.type.DeclaredType;
import javax.lang.model.type.TypeMirror;
import java.util.*;
import java.util.stream.Collectors;

public final class Parser {

    /**
     * 常见泛型参数 unused!  "T","Y","K","V","?"
     */
    @Getter
    private enum Native {

        OBJECT("java.lang.Object", "Object"),
        VOID("void", "void"),
        LONG("long", "long"),
        INT("int", "int"),
        CHAR("char", "char"),
        BOOLEAN("boolean", "boolean"),
        BYTE("byte", "byte"),
        SHORT("short", "short"),
        FLOAT("float", "float"),
        DOUBLE("double", "double"),
        UNKNOWN("unknown", "unknown"),
        ;
        private final Kind kind;

        Native(String name, String qualifiedName) {
            this.kind = ofNative(name, qualifiedName);
        }

        public Once<Kind> once() {
            return Once.with(() -> this.kind);
        }
    }

    public static List<Kind> parse(CompileContext context, JavaDoc macro, TypeElement e) {
        Parser parser = new Parser(e);
        Kind kind = parser.parse(context, macro);
        String name = Tool.required(kind.getPkg()) ? String.format("%s.%s", kind.getPkg(), kind.getName()) : kind.getName();
        List<Kind> dependencies = new ArrayList<>();
        kind.getVariables().forEach(v -> dependencies.add(v.getKind()));
        kind.getMethods().forEach(m -> {
            m.getParameters().forEach(p -> dependencies.add(p.getKind()));
            m.getReturns().forEach(r -> dependencies.add(r.getKind()));
            m.getCauses().forEach(c -> dependencies.add(c.getKind()));
        });
        dependencies.addAll(kind.getTraits());
        dependencies.addAll(kind.getSupers());
        dependencies.stream().flatMap(x -> {
            List<String> types = ParameterizedTypes.ofTree(x.getSignature()).toList();
            if (Tool.required(x.getPkg())) {
                types.add(String.format("%s.%s", x.getPkg(), x.getName()));
            } else {
                types.add(x.getName());
            }
            return types.stream();
        }).filter(x -> !Tool.equals(x, name)).distinct().forEach(x -> kind.getImports().add(x));
        return Collections.singletonList(kind);
    }

    private final Deque<TypeElement> elements = new ArrayDeque<>();
    private final Deque<Kind> kinds = new ArrayDeque<>();

    private Parser(TypeElement element) {
        this.elements.push(element);
    }

    public Kind parse(CompileContext context, JavaDoc macro) {
        return Once.with(() -> this.elements).map(Deque::pop).
                map(x -> elementType(context, macro, x, x.asType())).orElseGet(Once::empty).get();
    }

    private Once<Kind> elementType(CompileContext context, JavaDoc macro, TypeElement e, TypeMirror m) {
        if (!(e.getKind().isClass() || e.getKind().isInterface())) {
            return Once.empty();
        }
        if (!e.getModifiers().contains(Modifier.PUBLIC)) {
            return Once.empty();
        }
        if (this.elements.stream().anyMatch(x -> Tool.equals(x.toString(), e.toString()))) {
            return Once.of(this.kinds.peek()).any(Kind::copy);
        }
        if (Arrays.stream(macro.ignore()).anyMatch(x -> Tool.startWith(e.toString(), x))) {
            return Once.with(() -> ofNative(e, m));
        }
        Kind kind = new Kind();
        this.elements.push(e);
        this.kinds.push(kind);
        try {
            String[] names = Tool.split(e.getQualifiedName().toString(), "\\.");
            Map<String, Map<String, String>> macros = new HashMap<>();
            Optional.ofNullable(e.getAnnotationMirrors()).orElseGet(Collections::emptyList).forEach(x -> macros.putAll(parseMacro(x, context, macro)));
            kind.setPkg(String.join(".", Arrays.asList(names).subList(0, names.length - 1)));
            kind.setName(e.getSimpleName().toString());
            kind.setSignature(ParameterizedTypes.ofTree(e.toString()).toCanonicalName(Parser::ofCanonicalName));
            kind.setImports(new ArrayList<>());
            kind.setMacros(macros);
            kind.setComments(parseComment(context, e));
            kind.setModifier(parseModifier(e.getModifiers()) | (e.getKind().isInterface() ? java.lang.reflect.Modifier.INTERFACE : 0));
            kind.setVariables(new ArrayList<>());
            kind.setMethods(new ArrayList<>());
            kind.setSupers(new ArrayList<>());
            kind.setTraits(new ArrayList<>());
            mirrorType(context, macro, e.getSuperclass()).ifPresent(cls -> {
                if (!Tool.endsWith(cls.getName(), Object.class.getCanonicalName())) {
                    kind.getSupers().add(cls);
                }
            });
            for (TypeMirror trait : e.getInterfaces()) {
                mirrorType(context, macro, trait).ifPresent(cls -> kind.getTraits().add(cls));
            }
            for (Element enclosed : e.getEnclosedElements()) {
                if (enclosed.getKind() == ElementKind.FIELD || enclosed.getKind() == ElementKind.ENUM_CONSTANT) {
                    kind.getVariables().add(parseVariable(context, macro, (VariableElement) enclosed));
                } else if (enclosed.getKind() == ElementKind.METHOD) {
                    kind.getMethods().add(parseMethod(context, macro, (ExecutableElement) enclosed));
                }
            }
            return Once.with(() -> kind);
        } finally {
            this.elements.pop();
            this.kinds.pop();
        }
    }

    /**
     * 泛型类参数 如: T K ?等
     */
    private Once<Kind> mirrorType(CompileContext context, JavaDoc macro, TypeMirror mirror) {
        if (null == mirror || Object.class.getCanonicalName().equals(mirror.toString())) {
            return Native.OBJECT.once();
        }
        switch (mirror.getKind()) {
            case BOOLEAN:
                return Native.BOOLEAN.once();
            case BYTE:
                return Native.BYTE.once();
            case SHORT:
                return Native.SHORT.once();
            case INT:
                return Native.INT.once();
            case LONG:
                return Native.LONG.once();
            case CHAR:
                return Native.CHAR.once();
            case FLOAT:
                return Native.FLOAT.once();
            case DOUBLE:
                return Native.DOUBLE.once();
            case VOID:
                return Native.VOID.once();
            case NONE:
            case TYPEVAR:
            case NULL:
                return Native.OBJECT.once();
            case ARRAY:
                TypeMirror ak = ((ArrayType) mirror).getComponentType();
                if (ak instanceof TypeElement) {
                    return this.elementType(context, macro, (TypeElement) ak, ak);
                }
                return Once.with(() -> ofNative("[]", mirror.toString()));
            case DECLARED:
                DeclaredType pk = (DeclaredType) mirror;
                if (!(pk.asElement() instanceof TypeElement)) {
                    return Once.empty();
                }
                return this.elementType(context, macro, (TypeElement) pk.asElement(), pk);
            case ERROR:
            case WILDCARD:
            case PACKAGE:
            case EXECUTABLE:
            case OTHER:
            case UNION:
            case INTERSECTION:
                // case MODULE:
                context.warn(String.format("Unsupported type %s. ", mirror));
                return Once.empty();
            default:
                return Once.empty();
        }
    }

    private Method parseMethod(CompileContext ctx, JavaDoc macro, ExecutableElement e) {
        Method method = new Method();
        method.setName(e.getSimpleName().toString());
        method.setMacros(new HashMap<>());
        method.setComments(parseComment(ctx, e));
        method.setModifier(parseModifier(e.getModifiers()));
        method.setParameters(e.getParameters().stream().map(v -> parseParameter(ctx, macro, v)).collect(Collectors.toList()));
        method.setReturns(new ArrayList<>());
        method.setCauses(new ArrayList<>());
        // parse return type
        mirrorType(ctx, macro, e.getReturnType()).ifPresent(kind -> {
            Return r = new Return();
            r.setName(kind.getName());
            r.setComments(new ArrayList<>());
            r.setKind(kind);
            method.getReturns().add(r);
        });
        for (TypeMirror m : e.getThrownTypes()) {
            mirrorType(ctx, macro, m).ifPresent(kind -> {
                Throw thr = new Throw();
                thr.setKind(kind);
                thr.setName(kind.getName());
                method.getCauses().add(thr);
            });
        }

        for (AnnotationMirror am : e.getAnnotationMirrors()) {
            method.getMacros().putAll(parseMacro(am, ctx, macro));
        }
        return method;
    }

    private List<String> parseComment(CompileContext ctx, Element e) {
        String explain = Optional.ofNullable(ctx.utilities().getDocComment(e)).orElse("");
        List<String> comments = new ArrayList<>();
        for (String comment : explain.split("\n")) {
            comments.add(comment.trim());
        }
        return comments;
    }

    private Map<String, Map<String, String>> parseMacro(AnnotationMirror an, CompileContext ctx, JavaDoc macro) {
        Map<String, Map<String, String>> macros = new HashMap<>();
        String annoTypeName = an.getAnnotationType().toString();
        mirrorType(ctx, macro, an.getAnnotationType()).ifPresent(kind -> {
            Map<String, String> metadata = new HashMap<>();
            ctx.utilities().getElementValuesWithDefaults(an).forEach((k, v) -> {
                String val = v.getValue().toString();
                metadata.put(k.getSimpleName().toString(), parseMacroMetadata(val, annoTypeName));
            });
            an.getElementValues().forEach((k, v) -> {
                String val = v.getValue().toString();
                metadata.put(k.getSimpleName().toString(), parseMacroMetadata(val, annoTypeName));
            });
            macros.put(kind.getName(), metadata);
        });
        return macros;
    }

    private String parseMacroMetadata(String v, String annoTypeName) {
        StringBuilder be = new StringBuilder();
        if (v.contains(",")) {
            be.append("[");
            Arrays.stream(v.split(",")).collect(Collectors.toList()).forEach(e -> {
                be.append("\"").append(formatMacroMetadata(e, annoTypeName)).append("\",");
            });
            be.setLength(be.length() - 1);
            be.append("]");
            return be.toString();
        } else {
            if (!(v.startsWith("\"") && v.endsWith("\""))) {
                v = "[\"" + formatMacroMetadata(v, annoTypeName) + "\"]";
            }
            return v;
        }
    }

    private String formatMacroMetadata(String v, String anTypeName) {
        if (anTypeName.equals("java.lang.annotation.Target")) {
            return "java.lang.annotation.ElementType." + v;
        } else if (anTypeName.equals("java.lang.annotation.Retention")) {
            return "java.lang.annotation.RetentionPolicy." + v;
        }
        return v;
    }

    private Parameter parseParameter(CompileContext ctx, JavaDoc macro, VariableElement v) {
        Parameter parameter = new Parameter();
        parameter.setName(v.getSimpleName().toString());
        parameter.setMacros(new HashMap<>());
        parameter.setComments(new ArrayList<>());
        parameter.setKind(mirrorType(ctx, macro, v.asType()).get());
        parameter.setValue("");
        if (null != v.getAnnotationMirrors()) {
            for (AnnotationMirror an : v.getAnnotationMirrors()) {
                parameter.getMacros().putAll(parseMacro(an, ctx, macro));
            }
        }
        return parameter;
    }

    private Variable parseVariable(CompileContext ctx, JavaDoc macro, VariableElement e) {
        Variable variable = new Variable();
        variable.setName(e.getSimpleName().toString());
        variable.setMacros(new HashMap<>());
        variable.setComments(parseComment(ctx, e));
        variable.setModifier(parseModifier(e.getModifiers()));
        variable.setKind(mirrorType(ctx, macro, e.asType()).get());
        variable.setValue(Optional.ofNullable(e.getConstantValue()).map(Object::toString).orElse(""));
        return variable;
    }

    public int parseModifier(Set<Modifier> modifiers) {
        return modifiers.stream().map(x -> {
            switch (x) {
                case PUBLIC:
                    return java.lang.reflect.Modifier.PUBLIC;
                case PROTECTED:
                    return java.lang.reflect.Modifier.PROTECTED;
                case PRIVATE:
                    return java.lang.reflect.Modifier.PRIVATE;
                case ABSTRACT:
                    return java.lang.reflect.Modifier.ABSTRACT;
                case DEFAULT:
                case STATIC:
                    return java.lang.reflect.Modifier.STATIC;
                case FINAL:
                    return java.lang.reflect.Modifier.FINAL;
                case TRANSIENT:
                    return java.lang.reflect.Modifier.TRANSIENT;
                case VOLATILE:
                    return java.lang.reflect.Modifier.VOLATILE;
                case SYNCHRONIZED:
                    return java.lang.reflect.Modifier.SYNCHRONIZED;
                case NATIVE:
                    return java.lang.reflect.Modifier.NATIVE;
                case STRICTFP:
                default:
                    return 0;
            }
        }).reduce((x, y) -> x | y).orElse(0);
    }

    public static Kind ofNative(String name, String qn) {
        Kind kind = new Kind();
        kind.setPkg("");
        kind.setName(name);
        kind.setImports(new ArrayList<>());
        kind.setMacros(new HashMap<>());
        kind.setComments(new ArrayList<>());
        kind.setModifier(0);
        kind.setVariables(new ArrayList<>());
        kind.setMethods(new ArrayList<>());
        kind.setSupers(new ArrayList<>());
        kind.setTraits(new ArrayList<>());
        kind.setSignature(ParameterizedTypes.ofTree(ofCanonicalArray(qn)).toCanonicalName(Parser::ofCanonicalName));
        return kind;
    }

    public static Kind ofNative(TypeElement e, TypeMirror m) {
        Kind kind = new Kind();
        kind.setPkg(ofCanonicalPkg(e.getQualifiedName().toString()));
        kind.setName(e.getSimpleName().toString());
        kind.setImports(new ArrayList<>());
        kind.setMacros(new HashMap<>());
        kind.setComments(new ArrayList<>());
        kind.setModifier(0);
        kind.setVariables(new ArrayList<>());
        kind.setMethods(new ArrayList<>());
        kind.setSupers(new ArrayList<>());
        kind.setTraits(new ArrayList<>());
        kind.setSignature(ParameterizedTypes.ofTree(m.toString()).toCanonicalName(Parser::ofCanonicalName));
        return kind;
    }

    private static String ofCanonicalName(String name) {
        if (Tool.startWith(name, "java")) {
            return name.substring(name.lastIndexOf(".") + 1);
        }
        return name;
    }

    private static String ofCanonicalPkg(String qn) {
        if (Tool.startWith(qn, "java")) {
            return "";
        }
        return qn.lastIndexOf(".") < 0 ? qn : qn.substring(0, qn.lastIndexOf("."));
    }

    private static String ofCanonicalArray(String signature) {
        if (!Tool.endsWith(signature, "[]")) {
            return signature;
        }
        String ss = Tool.substring(signature, 0, signature.length() - 2);
        return String.format("List<%s>", ofCanonicalArray(ss));
    }

}
