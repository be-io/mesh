/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.schema.plugin;

import com.be.mesh.client.annotate.SPI;
import com.be.mesh.client.mpc.Codec;
import com.be.mesh.client.mpc.ServiceLoader;
import com.be.mesh.client.schema.CompilePlugin;
import com.be.mesh.client.schema.context.CompileContext;
import com.be.mesh.client.schema.macro.Version;
import com.be.mesh.client.struct.Versions;
import com.be.mesh.client.tool.Tool;

import javax.tools.StandardLocation;
import java.io.*;
import java.util.HashMap;
import java.util.Map;
import java.util.stream.Collectors;

/**
 * @author coyzeng@gmail.com
 */
@SPI(value = "version", meta = Version.class)
public class VersionPlugin implements CompilePlugin {

    public static final String PATTERN = ".*?((?<!\\w[a-zA-Z])\\d+([.-]\\d+)*).*";

    @Override
    public void proceed(CompileContext context) throws Throwable {
        try {
            Codec codec = ServiceLoader.load(Codec.class).get(Codec.JSON);
            Version version = context.annotatedElement().getAnnotation(Version.class);
            try (OutputStream stream = openSchema(context, version)) {
                String v = Tool.anyone(version.version(), getVersion(context));
                Map<String, String> infos = new HashMap<>();
                infos.put(String.format("%s.commit_id", Tool.anyone(version.value(), version.name())), getCommitId(context));
                infos.put(String.format("%s.os", Tool.anyone(version.value(), version.name())), getOS());
                infos.put(String.format("%s.arch", Tool.anyone(version.value(), version.name())), System.getProperty("os.arch"));
                infos.put(String.format("%s.version", Tool.anyone(version.value(), version.name())), v);
                Versions versions = new Versions();
                versions.setVersion(v);
                versions.setInfos(infos);
                stream.write(codec.encode(versions).array());
            }
        } catch (Throwable e) {
            error(Tool.getStackTrace(e));
        }
    }

    private String getVersion(CompileContext context) throws Exception {
        String branch = getBranch(context);
        warn(String.format("Detect git branch %s. ", branch));
        if (Tool.required(branch)) {
            return branch.replaceAll(VersionPlugin.PATTERN, "$1");
        }
        String tag = getTag(context);
        warn(String.format("Detect git tag %s. ", tag));
        if (Tool.required(tag)) {
            return tag.replaceAll(VersionPlugin.PATTERN, "$1");
        }
        return "";
    }

    private String getBranch(CompileContext context) throws Exception {
        return execCMD(context, "git branch --show-current");
    }

    private String getCommitId(CompileContext context) throws Exception {
        return execCMD(context, "git rev-parse --short HEAD");
    }

    private String getTag(CompileContext context) throws Exception {
        return execCMD(context, "git describe --tags --abbrev=0");
    }

    private String execCMD(CompileContext context, String cmd) throws Exception {
        Process process = Runtime.getRuntime().exec(Tool.split(cmd, " "));
        try (InputStream stream = process.getInputStream()) {
            process.waitFor();
            BufferedReader reader = new BufferedReader(new InputStreamReader(stream));
            return reader.lines().filter(Tool::required).collect(Collectors.joining());
        } catch (InterruptedException e) {
            warn(Tool.getStackTrace(e));
            Thread.currentThread().interrupt();
            return "";
        } catch (Throwable e) {
            warn(Tool.getStackTrace(e));
            return "";
        } finally {
            process.destroy();
        }
    }

    private String getOS() {
        String os = Tool.toLowerCase(System.getProperty("os.name"));
        if ("mac os x".equals(os)) {
            return "darwin";
        }
        if (os.contains("windows") || "windows xp".equals(os) || "windows 2003".equals(os) || "windows vista".equals(os) || os.contains("windows 7") || os.contains("windows 8")) {
            return "windows";
        }
        if (os.contains("linux")) {
            return "linux";
        }
        return os;
    }

    private OutputStream openSchema(CompileContext context, Version version) throws IOException {
        String name = Tool.anyone(version.value(), version.name(), Tool.MESH_NAME.get());
        String location = String.format("META-INF/mesh/%s.version", name);
        return context.processingEnvironment().getFiler().createResource(StandardLocation.CLASS_OUTPUT, "", location).openOutputStream();
    }
}
