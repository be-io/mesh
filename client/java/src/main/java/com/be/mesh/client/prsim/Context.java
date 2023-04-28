/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.prsim;

import com.be.mesh.client.mpc.Mesh;
import com.be.mesh.client.mpc.Runmode;
import com.be.mesh.client.mpc.Types;
import com.be.mesh.client.struct.Location;
import com.be.mesh.client.struct.Principal;
import com.be.mesh.client.tool.Tool;
import lombok.AllArgsConstructor;
import lombok.Getter;

import java.io.Serializable;
import java.util.Deque;
import java.util.Map;
import java.util.Optional;
import java.util.function.Consumer;
import java.util.function.Function;
import java.util.function.Supplier;

/**
 * https://www.rfc-editor.org/rfc/rfc7540#section-8.1.2
 *
 * @author coyzeng@gmail.com
 */
public interface Context extends Serializable {

    @Getter
    @AllArgsConstructor
    enum Metadata {
        MESH_SPAN_ID("mesh-span-id"),
        MESH_TIMESTAMP("mesh-timestamp"),
        MESH_RUN_MODE("mesh-run-mode"),
        MESH_CONSUMER("mesh-consumer"),
        MESH_PROVIDER("mesh-provider"),
        MESH_URN("mesh-urn"),
        MESH_INCOMING_HOST("mesh-incoming-host"),
        MESH_OUTGOING_HOST("mesh-outgoing-host"),
        MESH_INCOMING_PROXY("mesh-incoming-proxy"),
        MESH_OUTGOING_PROXY("mesh-outgoing-proxy"),
        MESH_SUBSET("mesh-subset"),
        MESH_TIMEOUT("mesh-timeout"),
        // PTP
        MESH_VERSION("x-ptp-version"),
        MESH_TECH_PROVIDER_CODE("x-ptp-tech-provider-code"),
        MESH_TRACE_ID("x-ptp-trace-id"),
        MESH_TOKEN("x-ptp-token"),
        MESH_URI("x-ptp-uri"),
        MESH_FROM_NODE_ID("x-ptp-from-node-id"),
        MESH_FROM_INST_ID("x-ptp-from-inst-id"),
        MESH_TARGET_NODE_ID("x-ptp-target-node-id"),
        MESH_TARGET_INST_ID("x-ptp-target-inst-id"),
        MESH_SESSION_ID("x-ptp-session-id"),
        MESH_TOPIC("x-ptp-topic"),
        MESH_TIMEOUT("x-ptp-timeout"),
        ;
        private final String key;

        public final String get() {
            return this.get(Mesh.context().getAttachments());
        }

        @SafeVarargs
        public final String get(Map<String, String> attachments, Supplier<String>... defaults) {
            if (Tool.required(attachments)) {
                String v = attachments.get(this.getKey());
                if (Tool.required(v)) {
                    return v;
                }
                v = attachments.get(this.getKey().replace("-", "_"));
                if (Tool.required(v)) {
                    return v;
                }
                for (Map.Entry<String, String> kv : attachments.entrySet()) {
                    if (null != kv.getValue() &&
                            (Tool.equals(this.getKey(), Tool.toLowerCase(kv.getKey())) || Tool.equals(this.getKey(), Tool.toLowerCase(kv.getKey()).replace("_", "-")))) {
                        return kv.getValue();
                    }
                }
            }
            if (Tool.optional(defaults)) {
                return "";
            }
            return defaults[0].get();
        }

        public void set(Map<String, String> attachments, String value) {
            if (null != attachments && null != value) {
                attachments.put(this.getKey(), value);
            }
        }
    }

    /**
     * Get the request trace id.
     *
     * @return The request trace id.
     */
    String getTraceId();

    /**
     * Get the request span id.
     *
     * @return The request span id.
     */
    String getSpanId();

    /**
     * Get the request create time.
     *
     * @return The request create time.
     */
    long getTimestamp();

    /**
     * Get the request run mode. {@link Runmode}
     *
     * @return The request run mode.
     */
    Runmode getRunMode();

    /**
     * Mesh resource uniform name. Like: create.tenant.omega.json.http2.000.lx000001.trustbe.cn
     *
     * @return Uniform name.
     */
    String getUrn();

    /**
     * Get the consumer network principal.
     *
     * @return Consumer network principal.
     */
    Location getConsumer();

    /**
     * Get the provider network principal.
     *
     * @return Provider network principal.
     */
    Location getProvider();

    /**
     * Dispatch attachments.
     *
     * @return Dispatch attachments.
     */
    Map<String, String> getAttachments();

    /**
     * Get the mpc broadcast network principals.
     *
     * @return Broadcast principals.
     */
    Deque<Principal> getPrincipals();

    /**
     * Get the context attributes. The attributes don't be serialized.
     *
     * @return attributes
     */
    Map<String, Object> getAttributes();

    /**
     * Like getAttachments, but attribute wont be transfer in invoke chain.
     *
     * @param key attribute key
     * @param <T> attribute type
     * @return attribute value
     */
    <T> T getAttribute(Key<T> key);

    /**
     * Like putAttachments, but attribute won't be transfer in invoke chain.
     *
     * @param key   attribute key
     * @param value attribute
     * @param <T>   attribute type
     */
    <T> void setAttribute(Key<T> key, T value);

    /**
     * Compute attribute if absent, is a context attribute cache.
     *
     * @param key      attribute key
     * @param supplier attribute value supplier
     * @param <T>      attribute type
     * @return attribute value
     */
    <T> T computeAttribute(Key<T> key, Supplier<T> supplier);

    /**
     * Rewrite the urn.
     *
     * @param urn urn
     */
    void rewriteUrn(String urn);

    /**
     * Rewrite the context by another context.
     *
     * @param context another context
     */
    void rewriteContext(Context context);

    /**
     * Open a new context.
     */
    Context resume();

    /**
     * Context key.
     *
     * @param <T>
     */
    @Getter
    @AllArgsConstructor
    final class Key<T> implements Serializable {

        private static final long serialVersionUID = -1427493963117316519L;
        // Attribute name, must be unique
        private final String name;
        private final Types<T> type;

        // Get attribute from context
        @SuppressWarnings("unchecked")
        public Optional<T> getIfPresent() {
            if (null != Mesh.context().getAttribute(this)) {
                return Optional.ofNullable(Mesh.context().getAttribute(this));
            }
            String value = Mesh.context().getAttachments().get(this.name);
            if (String.class == this.type.getRawType() && Tool.required(value)) {
                return Optional.ofNullable((T) value);
            }
            return Optional.empty();
        }

        public <V> Optional<V> map(Function<T, V> fn) {
            return getIfPresent().map(fn);
        }

        public void ifPresent(Consumer<T> consumer) {
            getIfPresent().ifPresent(consumer);
        }

        public T orElse(T value) {
            return Tool.anyone(getIfPresent().orElse(value), value);
        }

        public boolean isPresent() {
            return Tool.required(orElse(null));
        }

    }
}
