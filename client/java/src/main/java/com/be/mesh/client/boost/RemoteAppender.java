package com.be.mesh.client.boost;

import ch.qos.logback.classic.Level;
import ch.qos.logback.classic.PatternLayout;
import ch.qos.logback.classic.spi.ILoggingEvent;
import ch.qos.logback.core.UnsynchronizedAppenderBase;
import ch.qos.logback.core.status.OnErrorConsoleStatusListener;
import ch.qos.logback.core.status.Status;
import com.be.mesh.client.mpc.Mesh;
import com.be.mesh.client.mpc.ServiceLoader;
import com.be.mesh.client.mpc.Types;
import com.be.mesh.client.prsim.Context;
import com.be.mesh.client.prsim.DataHouse;
import com.be.mesh.client.struct.Document;
import com.be.mesh.client.tool.Mode;
import com.be.mesh.client.tool.Once;
import com.be.mesh.client.tool.Tool;
import com.lmax.disruptor.EventHandler;
import com.lmax.disruptor.EventTranslator;
import com.lmax.disruptor.ExceptionHandler;
import com.lmax.disruptor.SleepingWaitStrategy;
import com.lmax.disruptor.dsl.Disruptor;
import com.lmax.disruptor.dsl.ProducerType;
import lombok.AllArgsConstructor;
import lombok.Data;

import java.util.*;
import java.util.concurrent.Semaphore;
import java.util.concurrent.ThreadFactory;
import java.util.concurrent.atomic.AtomicInteger;

/**
 * @author jianyue.li
 */
public class RemoteAppender extends UnsynchronizedAppenderBase<ILoggingEvent> implements EventHandler<List<Document>>, ExceptionHandler<List<Document>> {

    private static final String PATTERN = "%d{yyyy-MM-dd HH:mm:ss.SSS},%X{tid}#%X{sid},%thread,%highlight(%level),%marker,%logger{20},%msg,%n%ex";
    private static final Context.Key<Boolean> DISCARD = new Context.Key<>("mesh.syslog.discard", Types.of(Boolean.class));
    private static final int BUFFER_SIZE = 10;
    private final Once<DataHouse> dataHouse = Once.with(() -> ServiceLoader.load(DataHouse.class).getDefault());
    private final Once<PatternLayout> layout = new Once<>();
    private final Once<RemoteFormatter> configuration = Once.with(RemoteFormatter::new);
    private final Once<Semaphore> limiter = Once.with(() -> new Semaphore(3));
    private final Once<Disruptor<List<Document>>> disruptor = Once.with(() -> {
        Disruptor<List<Document>> queue = new Disruptor<>(() -> new ArrayList<>(BUFFER_SIZE), 64, new RemoteThreadFactory("syslog"), ProducerType.SINGLE, new SleepingWaitStrategy());
        queue.handleEventsWith(this);
        queue.setDefaultExceptionHandler(this);
        queue.start();
        return queue;
    });

    @Override
    protected void append(ILoggingEvent event) {
        if (!Tool.MESH_MODE.get().match(Mode.RLog)) {
            return;
        }
        if (!super.isStarted()) {
            addWarn(event.getFormattedMessage());
            return;
        }
        // avoid print log recursive
        Boolean isDiscard = Mesh.context().getAttribute(DISCARD);
        if (null != isDiscard && isDiscard) {
            return;
        }
        String message = layout.map(x -> x.doLayout(event)).orElse(event.getMessage());
        Map<String, String> metadata = new HashMap<>();
        if (Tool.required(event.getMDCPropertyMap())) {
            metadata.putAll(event.getMDCPropertyMap());
        }
        if (Tool.required(Mesh.context().getAttachments())) {
            metadata.putAll(Mesh.context().getAttachments());
        }
        metadata.put("name", event.getLoggerName());
        metadata.put("level", event.getLevel().toString());
        metadata.put("host", Tool.HOST_NAME.get());
        metadata.put("ip", Tool.IP.get());
        metadata.put("mesh_runtime", String.valueOf(Tool.MESH_RUNTIME.get().toString()));
        metadata.put("mesh_trace_id", Mesh.context().getTraceId());
        metadata.put("mesh_span_id", Mesh.context().getSpanId());
        metadata.put("mesh_name", Tool.MESH_NAME.get());
        metadata.put("mesh_run_mode", String.valueOf(Mesh.context().getRunMode().getMode()));
        Document document = new Document();
        document.setMetadata(metadata);
        document.setContent(message);
        document.setTimestamp(System.currentTimeMillis() * 1000000L + System.nanoTime() % 1000000L);
        disruptor.get().publishEvent(new DocumentTranslator(document));
    }

    @AllArgsConstructor
    private static class DocumentTranslator implements EventTranslator<List<Document>> {

        private final Document document;

        @Override
        public void translateTo(List<Document> event, long sequence) {
            Document doc = new Document();
            doc.setMetadata(document.getMetadata());
            doc.setContent(document.getContent());
            doc.setTimestamp(document.getTimestamp());
            event.add(doc);
        }
    }

    @Override
    public void start() {
        if (null != getStatusManager() && getStatusManager().getCopyOfStatusListenerList().isEmpty()) {
            StatusPrinter statusListener = new StatusPrinter(Status.INFO);
            statusListener.setContext(getContext());
            statusListener.start();
            getStatusManager().add(statusListener);
        }
        PatternLayout patternLayout = new PatternLayout();
        patternLayout.setContext(super.context);
        patternLayout.setPattern(this.configuration.map(RemoteFormatter::getPattern).orElse(PATTERN));
        patternLayout.start();
        this.layout.put(patternLayout);
        super.start();
        addInfo("Remote syslog has been started. ");
    }

    @Override
    public void stop() {
        if (!super.isStarted()) {
            return;
        }
        int size = disruptor.get().getRingBuffer().getBufferSize();
        for (int cursor = 0; cursor < size; cursor++) {
            if (disruptor.get().get(cursor).isEmpty()) {
                continue;
            }
            try {
                this.onEvent(disruptor.get().get(cursor), cursor, size - 1 == cursor);
            } catch (Exception e) {
                addError("Remote syslog shutdown clean with error. ", e);
            }
        }
        disruptor.get().shutdown();
        super.stop();
        addInfo("Remote syslog has been stopped. ");
    }

    public void setFormatter(RemoteFormatter formatter) {
        Optional.ofNullable(formatter).ifPresent(configuration::put);
        limiter.put(new Semaphore(configuration.get().getLimit()));
    }

    @Override
    public void onEvent(List<Document> event, long sequence, boolean endOfBatch) throws Exception {
        if (Mode.NOLOG.match(Tool.MESH_MODE.get())) {
            return;
        }
        Mesh.contextSafeUncheck(() -> {
            Mesh.context().setAttribute(DISCARD, true);
            if (event.size() > BUFFER_SIZE - 1) {
                try {
                    dataHouse.get().writes(event);
                } catch (Exception e) {
                    this.handleEventException(e, sequence, event);
                } finally {
                    event.clear();
                }
                return;
            }
            if (!limiter.get().tryAcquire()) {
                return;
            }
            try {
                dataHouse.get().writes(event);
            } catch (Exception e) {
                this.handleEventException(e, sequence, event);
            } finally {
                limiter.get().release();
                event.clear();
            }
        });
    }

    @Override
    public void handleEventException(Throwable e, long sequence, List<Document> event) {
        try {
            for (Document document : event) {
                Level level = Level.toLevel(document.getMetadata().get("level"), Level.INFO);
                if (Level.ERROR == level) {
                    addError(document.getContent());
                } else if (Level.WARN == level) {
                    addWarn(document.getContent());
                } else {
                    addInfo(document.getContent());
                }
            }
            addError(String.format("Remote syslog flush %d with error. ", sequence), e);
        } catch (Exception c) {
            addError("Unexpected cause. ", c);
        } finally {
            event.clear();
        }
    }

    @Override
    public void handleOnStartException(Throwable e) {
        addError("Remote syslog cant start. ", e);
    }

    @Override
    public void handleOnShutdownException(Throwable e) {
        addError("Remote syslog cant shutdown. ", e);
    }

    private static class StatusPrinter extends OnErrorConsoleStatusListener {
        private final int minLevel;

        public StatusPrinter(int minLevel) {
            super.setRetrospective(0L);
            this.minLevel = minLevel;
        }

        @Override
        public void addStatusEvent(Status status) {
            if (status.getEffectiveLevel() >= minLevel) super.addStatusEvent(status);
        }
    }

    @Data
    public static class RemoteFormatter {
        /**
         * Logback pattern to use for log record's message
         */
        private String pattern = PATTERN;
        /**
         * policy while message queue full, default wait until queue is available
         * wait,discard,log
         */
        private String policy = "discard";
        /**
         * message queue size
         */
        private int size = 64;
        /**
         * message batch write limits.
         */
        private int limit = 3;
    }

    private static class RemoteThreadFactory implements ThreadFactory {
        private final String prefix;
        private final AtomicInteger counter = new AtomicInteger(0);

        public RemoteThreadFactory(String prefix) {
            this.prefix = prefix;
        }

        @Override
        public Thread newThread(Runnable rn) {
            Thread thread = new Thread(rn, prefix + "-" + counter.getAndIncrement());
            thread.setDaemon(true);
            return thread;
        }
    }


}
