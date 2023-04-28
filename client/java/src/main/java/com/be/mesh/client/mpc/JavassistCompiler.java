/*
 * Copyright (c) 2000, 2023, ducesoft and/or its affiliates. All rights reserved.
 * DUCESOFT PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.mpc;

import com.be.mesh.client.annotate.*;
import com.be.mesh.client.cause.MeshException;
import com.be.mesh.client.schema.runtime.Javadoc;
import com.be.mesh.client.schema.runtime.TypeParameter;
import com.be.mesh.client.schema.runtime.TypeStruct;
import com.be.mesh.client.struct.Cause;
import com.be.mesh.client.struct.Reference;
import com.be.mesh.client.struct.Service;
import com.be.mesh.client.tool.Tool;
import com.thoughtworks.qdox.JavaProjectBuilder;
import com.thoughtworks.qdox.model.DocletTag;
import com.thoughtworks.qdox.model.JavaClass;
import com.thoughtworks.qdox.model.JavaField;
import javassist.*;
import javassist.bytecode.AnnotationsAttribute;
import javassist.bytecode.ClassFile;
import javassist.bytecode.ConstPool;
import javassist.bytecode.annotation.Annotation;
import javassist.bytecode.annotation.IntegerMemberValue;
import javassist.bytecode.annotation.StringMemberValue;
import lombok.SneakyThrows;
import lombok.extern.slf4j.Slf4j;

import java.io.InputStream;
import java.io.InputStreamReader;
import java.io.Serializable;
import java.lang.invoke.MethodHandles;
import java.lang.reflect.Method;
import java.lang.reflect.Parameter;
import java.lang.reflect.Type;
import java.util.*;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.Future;
import java.util.stream.Collectors;

/**
 * @author coyzeng@gmail.com
 */
@Slf4j
@SPI(JCompiler.JAVASSIST)
public class JavassistCompiler extends ClassPool implements JCompiler {

    private static final String PARAMETERS = "P";
    private static final String RETURNS = "R";
    private static final ThreadLocal<Class<?>> X = new ThreadLocal<>();
    private final Map<String, Class<?>> cache = new ConcurrentHashMap<>();

    public JavassistCompiler() {
        this.appendSystemPath();
    }

    @SneakyThrows
    @Override
    public Class<?> toClass(CtClass clazz) throws CannotCompileException {
        if (ClassFile.MAJOR_VERSION > ClassFile.JAVA_9) {
            Class<?> target = null != X.get() ? X.get() : JCompiler.class;
            Class<?> lookupType = MethodHandles.Lookup.class;
            Method privateLookupIn = MethodHandles.class.getDeclaredMethod("privateLookupIn", Class.class, lookupType);
            Object lookup = privateLookupIn.invoke(null, target, MethodHandles.lookup());
            return (Class<?>) lookupType.getDeclaredMethod("defineClass", byte[].class).invoke(lookup, clazz.toBytecode());
        }
        return super.toClass(clazz);
    }

    @SuppressWarnings("unchecked")
    @Override
    public <T extends Parameters> Class<T> intype(Method method) {
        if (method.getParameters().length < 1) {
            return (Class<T>) GenericParameters.class;
        }
        return (Class<T>) cache.computeIfAbsent(getSignature(method, PARAMETERS), key -> Tool.uncheck(() -> {
            try {
                X.set(method.getDeclaringClass());
                this.insertClassPath(new ClassClassPath(this.getClass()));
                CtClass pc = this.makeClass(getTypeName(method, PARAMETERS, true));
                ConstPool cp = pc.getClassFile().getConstPool();
                pc.addInterface(this.getCtClass(Parameters.class.getName()));
                pc.addInterface(this.getCtClass(Serializable.class.getName()));
                pc.addField(CtField.make("private static final long serialVersionUID = -1L;", pc));
                CtField attachments = CtField.make("private java.util.Map attachments;", pc);
                attachments.setGenericSignature(new ParameterizedTypes(Types.MapString).toSignature());
                attachments.getFieldInfo2().addAttribute(this.makeAnnotation(cp, -1, "attachments"));
                pc.addField(attachments);
                CtMethod getAttachments = CtMethod.make("public java.util.Map getAttachments(){return this.attachments;}", pc);
                getAttachments.setGenericSignature(makeGetterSignature(Types.MapString));
                pc.addMethod(getAttachments);
                CtMethod setAttachments = CtMethod.make("public void setAttachments(java.util.Map attachments){this.attachments = attachments;}", pc);
                setAttachments.setGenericSignature(makeSetterSignature(Types.MapString));
                pc.addMethod(setAttachments);
                for (int index = 0; index < method.getParameters().length; index++) {
                    Parameter parameter = method.getParameters()[index];
                    int idx = Optional.ofNullable(parameter.getAnnotation(Index.class)).map(Index::value).orElse(index);
                    String in = Optional.ofNullable(parameter.getAnnotation(Index.class)).map(Index::name).orElse("");
                    String pn = parameter.getName();
                    CtField cf = CtField.make(String.format("private Object %s;", pn), pc);
                    cf.setGenericSignature(new ParameterizedTypes(parameter.getParameterizedType()).toSignature());
                    cf.getFieldInfo2().addAttribute(this.makeAnnotation(cp, idx, Tool.anyone(in, pn)));
                    pc.addField(cf);
                    CtMethod get = CtMethod.make(String.format("public Object get%s(){return this.%s;}", Tool.firstUpperCase(pn), pn), pc);
                    get.setGenericSignature(makeGetterSignature(parameter.getParameterizedType()));
                    pc.addMethod(get);
                    CtMethod set = CtMethod.make(String.format("public void set%s(Object %s){this.%s=%s;}", Tool.firstUpperCase(pn), pn, pn, pn), pc);
                    set.setGenericSignature(makeSetterSignature(parameter.getParameterizedType()));
                    pc.addMethod(set);
                }
                CtMethod getAttachments0 = CtMethod.make("public java.util.Map attachments(){return this.attachments;}", pc);
                getAttachments0.getMethodInfo2().addAttribute(this.makeOverride(cp));
                pc.addMethod(getAttachments0);
                CtMethod setAttachments0 = CtMethod.make("public void attachments(java.util.Map attachments){this.attachments = attachments;}", pc);
                setAttachments0.getMethodInfo2().addAttribute(this.makeOverride(cp));
                pc.addMethod(setAttachments0);
                String toArray = Arrays.stream(method.getParameters()).map(Parameter::getName).collect(Collectors.joining(", "));
                StringBuilder fromArray = new StringBuilder();
                for (int index = 0; index < method.getParameters().length; index++) {
                    Parameter p = method.getParameters()[index];
                    String type = Tool.boxType(p.getType().getCanonicalName());
                    fromArray.append(String.format("if(args.length > %d && null != args[%d]){this.%s=(%s)args[%d];}", index, index, p.getName(), type, index));
                }
                CtMethod getArguments = CtMethod.make(String.format("public Object[] arguments(){return new Object[]{this.%s};}", toArray), pc);
                getArguments.getMethodInfo2().addAttribute(this.makeOverride(cp));
                pc.addMethod(getArguments);
                CtMethod setArguments = CtMethod.make(String.format("public void arguments(Object[] args){if(null != args){%s}}", fromArray), pc);
                setArguments.getMethodInfo2().addAttribute(this.makeOverride(cp));
                pc.addMethod(setArguments);
                CtMethod getType = CtMethod.make("public Class type(){return this.getClass();}", pc);
                getType.getMethodInfo2().addAttribute(this.makeOverride(cp));
                pc.addMethod(getType);
                String toMap = Arrays.stream(method.getParameters()).map(x -> String.format("dict.put(\"%s\",this.%s)", x.getName(), x.getName())).collect(Collectors.joining(";"));
                CtMethod map = CtMethod.make(String.format("public java.util.Map map(){java.util.Map dict = new java.util.HashMap();%s;return dict;}", toMap), pc);
                map.getMethodInfo2().addAttribute(this.makeOverride(cp));
                map.setGenericSignature(makeGetterSignature(Types.MapObject));
                pc.addMethod(map);
                return pc.toClass();
            } catch (Throwable e) {
                log.error("Compile with error for {}.{}", method.getDeclaringClass().getCanonicalName(), method.getName());
                throw e;
            } finally {
                X.remove();
            }
        }));
    }

    @Override
    public <T extends Parameters> Class<T> intype(Reference reference) {
        return Tool.uncheck(() -> intype(Class.forName(reference.getNamespace()).getMethod(reference.getName())));
    }

    @SuppressWarnings("unchecked")
    @Override
    public <T extends Returns> Class<T> retype(Method method) {
        return (Class<T>) cache.computeIfAbsent(getSignature(method, RETURNS), key -> Tool.uncheck(() -> {
            try {
                X.set(method.getDeclaringClass());
                this.insertClassPath(new ClassClassPath(this.getClass()));
                CtClass pc = this.makeClass(getTypeName(method, RETURNS, true));
                ConstPool cp = pc.getClassFile().getConstPool();
                pc.addInterface(this.getCtClass(Returns.class.getName()));
                pc.addInterface(this.getCtClass(Serializable.class.getName()));
                pc.addField(CtField.make("private static final long serialVersionUID = -1L;", pc));
                CtField code = CtField.make("private String code;", pc);
                code.getFieldInfo2().addAttribute(this.makeAnnotation(cp, 0, "code"));
                pc.addField(code);
                CtMethod getCode = CtMethod.make("public String getCode(){return this.code;}", pc);
                getCode.getMethodInfo2().addAttribute(this.makeOverride(cp));
                pc.addMethod(getCode);
                CtMethod setCode = CtMethod.make("public void setCode(String code){this.code = code;}", pc);
                setCode.getMethodInfo2().addAttribute(this.makeOverride(cp));
                pc.addMethod(setCode);
                CtField message = CtField.make("private String message;", pc);
                message.getFieldInfo2().addAttribute(this.makeAnnotation(cp, 5, "message"));
                pc.addField(message);
                CtMethod getMessage = CtMethod.make("public String getMessage(){return this.message;}", pc);
                getMessage.getMethodInfo2().addAttribute(this.makeOverride(cp));
                pc.addMethod(getMessage);
                CtMethod setMessage = CtMethod.make("public void setMessage(String message){this.message = message;}", pc);
                setMessage.getMethodInfo2().addAttribute(this.makeOverride(cp));
                pc.addMethod(setMessage);
                CtField cause = CtField.make(String.format("private %s cause;", Cause.class.getCanonicalName()), pc);
                cause.getFieldInfo2().addAttribute(this.makeAnnotation(cp, 10, "cause"));
                pc.addField(cause);
                CtMethod getCause = CtMethod.make(String.format("public %s getCause(){return this.cause;}", Cause.class.getCanonicalName()), pc);
                getCause.getMethodInfo2().addAttribute(this.makeOverride(cp));
                pc.addMethod(getCause);
                CtMethod setCause = CtMethod.make(String.format("public void setCause(%s cause){this.cause = cause;}", Cause.class.getCanonicalName()), pc);
                setCause.getMethodInfo2().addAttribute(this.makeOverride(cp));
                pc.addMethod(setCause);
                CtField content = CtField.make("private Object content;", pc);
                content.getFieldInfo2().addAttribute(this.makeAnnotation(cp, 15, "content"));
                content.setGenericSignature(new ParameterizedTypes(Types.unbox(method, Future.class)).toSignature());
                pc.addField(content);
                CtMethod getContent = CtMethod.make("public Object getContent(){return this.content;}", pc);
                getContent.setGenericSignature(makeGetterSignature(Types.unbox(method, Future.class)));
                getContent.getMethodInfo2().addAttribute(this.makeOverride(cp));
                pc.addMethod(getContent);
                CtMethod setContent = CtMethod.make("public void setContent(Object content){if(null != content){this.content=content;}}", pc);
                setContent.setGenericSignature(makeSetterSignature(Types.unbox(method, Future.class)));
                setContent.getMethodInfo2().addAttribute(this.makeOverride(cp));
                pc.addMethod(setContent);
                return pc.toClass();
            } catch (Throwable e) {
                log.error("Compile with error for {}.{}", method.getDeclaringClass().getCanonicalName(), method.getName());
                throw e;
            } finally {
                X.remove();
            }
        }));
    }

    @Override
    public <T extends Returns> Class<T> retype(Service service) {
        return Tool.uncheck(() -> retype(Class.forName(service.getNamespace()).getMethod(service.getName())));
    }

    @Override
    public List<TypeStruct> documents() {
        JavaProjectBuilder project = new JavaProjectBuilder();
        Eden eden = ServiceLoader.load(Eden.class).getDefault();
        eden.inferTypes().forEach(type -> {
            try (InputStream input = this.getClass().getResourceAsStream("/" + type.getName().replace(".", "/") + ".java")) {
                if (null == input) {
                    return;
                }
                try (InputStreamReader reader = new InputStreamReader(input)) {
                    project.addSource(reader);
                }
            } catch (Exception e) {
                log.error("", e);
                log.error("Service document can not load in {}", type.getName());
                throw new MeshException("Document can not found.");
            }
        });
        return eden.inferTypes().stream().flatMap(x -> Arrays.stream(x.getDeclaredMethods())).map(me -> {
            TypeStruct struct = new TypeStruct();
            String name = me.getDeclaringClass().getName().replace(".", "/") + ".java";
            JavaClass javaClass = project.getClassByName(name);
            Javadoc javaDoc = new Javadoc(javaClass);
            String methodComment = javaDoc.getMethodComment(me.getName());
            List<TypeParameter> parameters = Arrays.stream(me.getParameters()).map(parameter -> {
                JavaClass pc = project.getClassByName(parameter.getType().getName().replace(".", "/") + ".java");
                String paramDoc = javaDoc.getMethodParameterComment(me.getName(), parameter.getName());
                TypeParameter ps = new TypeParameter();
                ps.setAliasName(parameter.getName());
                ps.setComment(paramDoc);
                ps.setFullName(parameter.getType().getName());
                ps.setName(parameter.getName());
                ps.setRequired(true);
                ps.setArray(pc.isArray());
                if (!pc.isPrimitive()) {
                    ps.setParameters(parseField(pc.getFields(), 0));
                }
                return ps;
            }).collect(Collectors.toList());

            Class<?> rtType = me.getReturnType();
            JavaClass rtTypeClass = project.getClassByName(rtType.getName().replace(".", "/") + ".java");
            TypeParameter rtParameter = new TypeParameter();
            rtParameter.setAliasName(rtType.getSimpleName());
            rtParameter.setComment(rtTypeClass.getComment());
            rtParameter.setFullName(rtType.getName());
            rtParameter.setName(rtType.getSimpleName());
            rtParameter.setRequired(true);

            MPS mps = me.getAnnotation(MPS.class);
            if (null != mps) {
                struct.setCommand(Tool.anyone(mps.name(), me.getDeclaringClass().getName() + "." + me.getName()));
                struct.setVersion(mps.version());
            }
            Binding binding = me.getAnnotation(Binding.class);
            if (null != binding) {
                struct.setCommand(Tool.anyone(binding.topic(), me.getDeclaringClass().getName() + "." + me.getName()));
                struct.setVersion(binding.version());
            }
            Bindings bindings = me.getAnnotation(Bindings.class);
            if (null != bindings) {
                struct.setCommand(Tool.anyone(bindings.value()[0].topic(), me.getDeclaringClass().getName() + "." + me.getName()));
                struct.setVersion(bindings.value()[0].topic());
            }
            struct.setAliasName(me.getDeclaringClass().getSimpleName());
            struct.setComment(methodComment);
            struct.setFullName(me.getDeclaringClass().getName());
            struct.setMethod(me.getName());
            struct.setClassComment(javaDoc.getComment());
            struct.setAuthor(Optional.ofNullable(javaClass.getTagByName("author")).map(DocletTag::getName).orElse(""));
            struct.setInput(parameters);
            struct.setOutput(Collections.singletonList(rtParameter));
            return struct;
        }).collect(Collectors.toList());
    }

    @SuppressWarnings("unchecked")
    @Override
    public <T> Class<? extends T> compile(Class<T> interfaces, String implement) {
        return (Class<T>) cache.computeIfAbsent(implement, key -> Tool.uncheck(() -> {
            try {
                X.set(interfaces);
                this.insertClassPath(new ClassClassPath(this.getClass()));
                CtClass pc = this.makeClass(String.format("%s%s", interfaces.getCanonicalName(), Math.abs(implement.hashCode())));
                pc.addInterface(this.getCtClass(interfaces.getName()));
                pc.addMethod(CtMethod.make(implement, pc));
                return pc.toClass();
            } catch (Throwable e) {
                log.error("Compile with error for {}.{}", interfaces.getCanonicalName(), implement);
                throw e;
            } finally {
                X.remove();
            }
        }));
    }

    private List<TypeParameter> parseField(List<JavaField> fields, int revers) {
        return fields.stream().map(field -> {
            JavaClass fc = field.getType();
            TypeParameter fs = new TypeParameter();
            fs.setAliasName(field.getName());
            fs.setComment(field.getComment());
            fs.setFullName(fc.getFullyQualifiedName());
            fs.setName(field.getName());
            fs.setRequired(true);
            fs.setArray(fc.isArray());
            if (!fc.isPrimitive() && revers < 4) {
                fs.setParameters(parseField(fc.getFields(), revers + 1));
            }
            return fs;
        }).collect(Collectors.toList());
    }

    private AnnotationsAttribute makeAnnotation(ConstPool cp, int idx, String name) {
        AnnotationsAttribute attribute = new AnnotationsAttribute(cp, AnnotationsAttribute.visibleTag);
        Annotation annotation = new Annotation(Index.class.getName(), cp);
        annotation.addMemberValue("value", new IntegerMemberValue(cp, idx));
        annotation.addMemberValue("name", new StringMemberValue(name, cp));
        attribute.addAnnotation(annotation);
        return attribute;
    }

    private AnnotationsAttribute makeOverride(ConstPool cp) {
        AnnotationsAttribute attribute = new AnnotationsAttribute(cp, AnnotationsAttribute.visibleTag);
        Annotation annotation = new Annotation(Override.class.getName(), cp);
        attribute.addAnnotation(annotation);
        return attribute;
    }

    private String makeGetterSignature(Type type) throws Exception {
        return String.format("()%s", new ParameterizedTypes(type).toSignature());
    }

    private String makeSetterSignature(Type type) throws Exception {
        return String.format("(%s)V", new ParameterizedTypes(type).toSignature());
    }

}
