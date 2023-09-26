/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.boost;

import com.be.mesh.client.annotate.SPI;
import com.be.mesh.client.mpc.ServiceLoader;
import com.be.mesh.client.mpc.*;
import com.be.mesh.client.prsim.Cache;
import com.be.mesh.client.prsim.Context;
import com.be.mesh.client.struct.CacheEntity;
import com.be.mesh.client.struct.Entity;
import com.be.mesh.client.tool.Mode;
import com.be.mesh.client.tool.Once;
import com.be.mesh.client.tool.Tool;
import jakarta.servlet.Filter;
import jakarta.servlet.*;
import jakarta.servlet.http.*;
import lombok.AllArgsConstructor;
import lombok.Setter;
import lombok.extern.slf4j.Slf4j;
import org.slf4j.MDC;

import java.io.*;
import java.time.Duration;
import java.util.*;
import java.util.function.Supplier;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

/**
 * @author coyzeng@gmail.com
 */
@Slf4j
@Setter
@SPI(value = "boost", priority = Integer.MIN_VALUE)
public class ServletBoostFilter implements Filter {

    public static final String FORM_CONTENT_TYPE = "application/x-www-form-urlencoded";
    public static final String STREAM_CONTENT_TYPE = "application/octet-stream";
    public static final String MULTI_PART_CONTENT_TYPE = "multipart/form-data";
    public static final String DEFAULT_CHARSET = "UTF-8";
    public static final String SESSION_PREFIX = "http.session";
    public static final String SESSION_ID = "MSESSIONID";

    private static final String IP_ADDRESS_REGEX = "^((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$";
    private static final Pattern IP_ADDRESS_PATTERN = Pattern.compile(IP_ADDRESS_REGEX);
    private static final Cache CACHE = ServiceLoader.load(Cache.class).getDefault();
    private boolean skipStream = false;
    private String removePath = "";

    private String getRemovePath() {
        return Tool.anyone(removePath, Tool.getProperty("", "MESH_SESSION_REMOVE_PATH", "mesh.session.remove.path"));
    }

    @Override
    public void doFilter(ServletRequest request, ServletResponse response, FilterChain chain) throws IOException, ServletException {
        Mesh.contextSafeUncheck(() -> {
            try {
                if (!(request instanceof HttpServletRequest) || !(response instanceof HttpServletResponse)) {
                    chain.doFilter(request, response);
                    return;
                }
                HttpServletRequest hr = (HttpServletRequest) request;
                MpcContext context = new MpcContext();
                context.setTraceId(Tool.anyone(hr.getHeader(Context.Metadata.MESH_TRACE_ID.getKey()), Mesh.context().getTraceId()));
                context.setSpanId(Tool.anyone(hr.getHeader(Context.Metadata.MESH_SPAN_ID.getKey()), Mesh.context().getSpanId()));
                String rm = hr.getHeader(Context.Metadata.MESH_RUN_MODE.getKey());
                if (Tool.required(rm) && Tool.isNumeric(rm)) {
                    context.setRunMode(Runmode.from(Integer.parseInt(rm)));
                }
                Mesh.context().rewriteContext(context);
                MDC.put("tid", Mesh.context().getTraceId());
                this.doHTTPFilter((HttpServletRequest) request, (HttpServletResponse) response, chain);
            } finally {
                MDC.clear();
            }
        });
    }

    public void doHTTPFilter(HttpServletRequest request, HttpServletResponse response, FilterChain chain) throws IOException, ServletException {
        Cookie sc = createCookieIfAbsent(request);
        if (Tool.contains(this.getRemovePath(), request.getRequestURI())) {
            sc.setMaxAge(0);
            response.addCookie(sc);
        } else {
            response.addCookie(sc);
        }
        RereadServletRequest newReq = new RereadServletRequest(request, skipStream, sc);
        chain.doFilter(newReq, response);
        try {
            if (Tool.equals(this.getRemovePath(), request.getRequestURI())) {
                newReq.getSession().invalidate();
            }
        } catch (Exception e) {
            log.error("Expire session with unexpect errr.", e);
        }
    }

    private Cookie createCookieIfAbsent(HttpServletRequest request) {
        int maxAge = Tool.getProperty(7 * 24 * 60 * 60, "MESH_SESSION_MAX_AGE", "mesh.session.max.age");
        if (Tool.required(request.getCookies())) {
            for (Cookie cookie : request.getCookies()) {
                if (Tool.equals(SESSION_ID, cookie.getName())) {
                    if (cookie.getMaxAge() < 1) {
                        cookie.setDomain(getDomain(request.getServerName()));
                        cookie.setPath("/");
                        cookie.setMaxAge(maxAge);
                    }
                    return cookie;
                }
            }
        }
        Cookie cookie = new Cookie(SESSION_ID, Tool.newTraceId());
        cookie.setDomain(getDomain(request.getServerName()));
        cookie.setPath("/");
        cookie.setMaxAge(maxAge);
        return cookie;
    }

    private String getDomain(String serverName) {
        String domain = serverName;
        if (!IP_ADDRESS_PATTERN.matcher(serverName).matches()) {
            domain = String.format(".%s", serverName);
        }
        return domain;
    }

    private static String getSet(Cookie cookie) {
        return String.format("%s.%s", SESSION_PREFIX, cookie.getValue());
    }

    public static class RereadServletRequest extends HttpServletRequestWrapper {

        private final Once<byte[]> buffer = Once.empty();
        private final Once<HttpSession> session = Once.empty();
        private final boolean skipStream;
        private final Cookie sc;

        public RereadServletRequest(HttpServletRequest request, boolean skipStream, Cookie sc) {
            super(request);
            this.skipStream = skipStream;
            this.sc = sc;
        }

        @Override
        public ServletInputStream getInputStream() throws IOException {
            if (skipStream) {
                String contentType = getContentType();
                if (Tool.required(contentType) && (contentType.contains(STREAM_CONTENT_TYPE) || contentType.contains(MULTI_PART_CONTENT_TYPE))) {
                    return getRequest().getInputStream();
                }
            }
            byte[] bytes = buffer.get(() -> Tool.uncheck(() -> {
                try (InputStream input = getRequest().getInputStream()) {
                    return Tool.readBytes(input);
                }
            }));
            return new RereadServletInputStream(new ByteArrayInputStream(bytes));
        }

        @Override
        public BufferedReader getReader() throws IOException {
            return new BufferedReader(new InputStreamReader(getInputStream(), getCharacterEncoding()));
        }

        @Override
        public String getCharacterEncoding() {
            return Tool.anyone(super.getCharacterEncoding(), DEFAULT_CHARSET);
        }

        @Override
        public HttpSession getSession(boolean create) {
            if (Tool.MESH_MODE.get().match(Mode.RCache)) {
                return super.getSession(create);
            }
            Supplier<HttpSession> hsf = () -> {
                HttpSession hs = super.getSession(create);
                if (null != hs) {
                    return hs;
                }
                return super.getSession(true);
            };
            if (create) {
                session.put(new ClusterHttpSession(sc, hsf.get()));
            }
            return session.get(() -> new ClusterHttpSession(sc, hsf.get()));
        }

        @Override
        public HttpSession getSession() {
            if (Tool.MESH_MODE.get().match(Mode.RCache)) {
                return super.getSession();
            }
            return session.get(() -> new ClusterHttpSession(sc, super.getSession()));
        }
    }

    @AllArgsConstructor
    public static final class RereadServletInputStream extends ServletInputStream {

        private final ByteArrayInputStream buffer;

        @Override
        public boolean isFinished() {
            return buffer.available() < 1;
        }

        @Override
        public boolean isReady() {
            return true;
        }

        @Override
        public void setReadListener(ReadListener readListener) {
            //
        }

        @Override
        public int read() throws IOException {
            return buffer.read();
        }

    }

    @AllArgsConstructor
    public static final class ClusterHttpSession implements HttpSession {

        private final Cookie sc;
        private final HttpSession session;

        @Override
        public long getCreationTime() {
            return this.session.getCreationTime();
        }

        @Override
        public String getId() {
            return this.sc.getValue();
        }

        @Override
        public long getLastAccessedTime() {
            return this.session.getLastAccessedTime();
        }

        @Override
        public ServletContext getServletContext() {
            return this.session.getServletContext();
        }

        @Override
        public void setMaxInactiveInterval(int interval) {
            this.session.setMaxInactiveInterval(interval);
        }

        @Override
        public int getMaxInactiveInterval() {
            return this.sc.getMaxAge();
        }

        @Override
        public Object getAttribute(String name) {
            return Optional.ofNullable(CACHE.hget(getSet(this.sc), name)).map(CacheEntity::getEntity).map(Entity::getBuffer).map(x -> Tool.uncheck(() -> Tool.deserialize(x))).orElse(null);
        }

        @Override
        public Enumeration<String> getAttributeNames() {
            List<String> names = Optional.ofNullable(CACHE.hkeys(getSet(this.sc))).orElseGet(Collections::emptyList);
            Iterator<String> iter = names.stream().iterator();
            return new Enumeration<String>() {
                @Override
                public boolean hasMoreElements() {
                    return iter.hasNext();
                }

                @Override
                public String nextElement() {
                    return iter.next();
                }
            };
        }

        @Override
        public void setAttribute(String name, Object value) {
            if (null == value) {
                return;
            }
            byte[] buffer = Tool.uncheck(() -> Tool.serialize(value));
            Entity entity = new Entity();
            entity.setCodec(Codec.JSON);
            entity.setSchema("");
            entity.setBuffer(buffer);
            CacheEntity cacheEntity = new CacheEntity();
            cacheEntity.setVersion("1.0.0");
            cacheEntity.setTimestamp(System.currentTimeMillis());
            cacheEntity.setDuration(Duration.ofSeconds(sc.getMaxAge()).toMillis());
            cacheEntity.setUnique(name);
            cacheEntity.setEntity(entity);
            CACHE.hset(getSet(this.sc), cacheEntity);
        }

        @Override
        public void removeAttribute(String name) {
            CACHE.hdel(getSet(this.sc), name);
        }

        @Override
        public void invalidate() {
            CACHE.remove(getSet(this.sc));
        }

        @Override
        public boolean isNew() {
            return true;
        }

//        @Override
//        public HttpSessionContext getSessionContext() {
//            return this.session.getSessionContext();
//        }
//
//        @Override
//        public Object getValue(String name) {
//            return this.getAttribute(name);
//        }
//
//        @Override
//        public String[] getValueNames() {
//            return Collections.list(this.getAttributeNames()).toArray(new String[]{});
//        }
//
//        @Override
//        public void putValue(String name, Object value) {
//            this.setAttribute(name, value);
//        }
//
//        @Override
//        public void removeValue(String name) {
//            this.removeAttribute(name);
//        }

    }

}