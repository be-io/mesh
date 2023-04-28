/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package com.be.mesh.client.schema.compiler;


import com.be.mesh.client.cause.CompatibleException;
import lombok.Getter;
import lombok.extern.slf4j.Slf4j;

import javax.tools.JavaCompiler;
import javax.tools.*;
import javax.tools.JavaFileObject.Kind;
import java.io.*;
import java.net.*;
import java.util.*;
import java.util.jar.JarEntry;
import java.util.jar.JarFile;

/**
 * JdkCompiler. (SPI, Singleton, ThreadSafe)
 */
@Slf4j
public class JdkCompiler extends AbstractCompiler {

    private static final String DEFAULT_JAVA_VERSION = "1.8";
    private static final DynamicClassLoader CLASS_LOADER = new DynamicClassLoader();

    @Override
    public Class<?> doCompile(String name, String sourceCode) throws Throwable {
        try {
            return CLASS_LOADER.findClass(name);
        } catch (ClassNotFoundException e) {
            //
        }
        List<String> options = new ArrayList<>();
        options.add("-source");
        options.add(DEFAULT_JAVA_VERSION);
        options.add("-target");
        options.add(DEFAULT_JAVA_VERSION);
        options.add("-proc:none");
        options.add("-classpath");
        options.add(System.getProperty("java.class.path"));
        options.add("-encoding");
        options.add("UTF-8");
        options.add("-XDuseUnsharedTable");

        DiagnosticCollector<JavaFileObject> collector = new DiagnosticCollector<>();
        ClassLoader classloader = this.getClass().getClassLoader();
        List<URL> entries = new ArrayList<>();
        if (classloader instanceof URLClassLoader) {
            URLClassLoader ucl = (URLClassLoader) classloader;
            entries.addAll(Arrays.asList(ucl.getURLs()));
        }
        JavaCompiler compiler = ToolProvider.getSystemJavaCompiler();
        JavaByteObject byteObject = new JavaByteObject(name);
        try (StandardJavaFileManager standard = compiler.getStandardFileManager(collector, null, null); JavaFileManager manager = new DynamicJavaFileManager(standard, entries, byteObject)) {
            JavaStringObject source = new JavaStringObject(name, sourceCode);
            JavaCompiler.CompilationTask task = compiler.getTask(null, manager, collector, options, null, Collections.singletonList(source));
            Boolean result = task.call();
            if (result == null || !result) {
                throw new CompatibleException("Compilation failed. class: %s, diagnostics:%s ", name, getCompileErrorMessage(collector));
            }
        }
        return CLASS_LOADER.defineClassIfAbsent(name, byteObject);
    }

    private String getCompileErrorMessage(DiagnosticCollector<JavaFileObject> diagnostics) {
        StringBuilder stack = new StringBuilder();
        for (Diagnostic<?> diagnostic : diagnostics.getDiagnostics()) {
            stack.append(diagnostic.toString());
        }
        return stack.toString();
    }

    private static final class DynamicJavaFileManager extends ForwardingJavaFileManager<JavaFileManager> {

        private final JavaByteObject byteObject;
        private final List<URL> entryURLs;

        public DynamicJavaFileManager(JavaFileManager fileManager, List<URL> entries, JavaByteObject byteObject) {
            super(fileManager);
            this.entryURLs = entries;
            this.byteObject = byteObject;
        }

        @Override
        public JavaFileObject getJavaFileForOutput(Location location, String className, Kind kind, FileObject sibling) throws IOException {
            return byteObject;
        }

        @Override
        public String inferBinaryName(Location loc, JavaFileObject file) {
            if (file instanceof ClassByteObject) {
                return file.getName();
            }
            return super.inferBinaryName(loc, file);
        }

        @Override
        public Iterable<JavaFileObject> list(Location location, String packageName, Set<Kind> kinds, boolean recurse) throws IOException {
            Iterable<JavaFileObject> supers = super.list(location, packageName, kinds, recurse);
            List<JavaFileObject> files = new ArrayList<>();
            // com/be/metabase/repository/KeyPairRepository.class
            String path = packageName.replace(".", "/");
            for (URL url : this.entryURLs) {
                try {
                    URLConnection connection = url.openConnection();
                    if (connection instanceof JarURLConnection) {
                        try (JarFile jarFile = ((JarURLConnection) connection).getJarFile()) {
                            Enumeration<JarEntry> entries = jarFile.entries();
                            while (entries.hasMoreElements()) {
                                JarEntry entry = entries.nextElement();
                                if (!entry.isDirectory() && entry.getName().startsWith(path)) {
                                    files.add(new ClassByteObject(entry.getName(), com.be.mesh.client.tool.Tool.readBytes(jarFile.getInputStream(entry))));
                                }
                            }
                        }
                    }
                } catch (Exception e) {
                    log.error("Cant read classes from jar {}, {}", url.getFile(), e.getMessage());
                }
            }
            for (JavaFileObject file : supers) {
                files.add(file);
            }
            return files;
        }
    }

    public static final class JavaStringObject extends SimpleJavaFileObject {
        private final String source;

        public JavaStringObject(String name, String source) {
            super(URI.create("string:///" + name.replace(".", "/") + Kind.SOURCE.extension), Kind.SOURCE);
            this.source = source;
        }

        @Override
        public CharSequence getCharContent(boolean ignoreEncodingErrors) throws IOException {
            return source;
        }
    }

    public static final class JavaByteObject extends SimpleJavaFileObject {

        private final ByteArrayOutputStream bytecodes;

        public JavaByteObject(String name) {
            super(URI.create("bytes:///" + name.replace(".", "/")), Kind.CLASS);
            bytecodes = new ByteArrayOutputStream();
        }

        @Override
        public OutputStream openOutputStream() throws IOException {
            return bytecodes;
        }

        public byte[] getBytes() {
            return bytecodes.toByteArray();
        }
    }

    public static final class ClassByteObject extends SimpleJavaFileObject {

        private final byte[] bytecodes;

        public ClassByteObject(String name, byte[] bytecodes) {
            super(URI.create("bytes:///" + name.replace(".", "/").replace("/class", Kind.CLASS.extension)), Kind.CLASS);
            this.bytecodes = bytecodes;
        }

        @Override
        public InputStream openInputStream() throws IOException {
            return new ByteArrayInputStream(bytecodes);
        }

        @Override
        public String getName() {
            String name = super.getName().replace("/", ".");
            if (name.startsWith(".") && name.endsWith(Kind.CLASS.extension)) {
                return name.substring(1, name.length() - 6);
            }
            if (name.startsWith(".")) {
                return name.substring(1);
            }
            if (name.endsWith(Kind.CLASS.extension)) {
                return name.substring(0, name.length() - 6);
            }
            return name;
        }
    }

    private static final class DynamicClassLoader extends ClassLoader {

        @Getter
        private final Map<String, JavaFileObject> bytecodes = new HashMap<>();
        private final Map<String, Class<?>> classes = new HashMap<>();

        DynamicClassLoader() {
            super(Thread.currentThread().getContextClassLoader());
        }

        @Override
        protected Class<?> findClass(String name) throws ClassNotFoundException {
            if (null != classes.get(name)) {
                return classes.get(name);
            }
            return super.findClass(name);
        }

        public Class<?> defineClassIfAbsent(String name, JavaByteObject byteObject) {
            try {
                return findClass(name);
            } catch (ClassNotFoundException e) {
                //
            }
            Class<?> type = defineClass(name, byteObject.getBytes(), 0, byteObject.getBytes().length);
            classes.put(name, type);
            bytecodes.put(name, byteObject);
            return type;
        }

    }

}
