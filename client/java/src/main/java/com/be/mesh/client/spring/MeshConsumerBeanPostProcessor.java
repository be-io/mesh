/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.spring;

import com.be.mesh.client.annotate.MPI;
import com.be.mesh.client.cause.CompatibleException;
import lombok.EqualsAndHashCode;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.BeanUtils;
import org.springframework.beans.BeansException;
import org.springframework.beans.PropertyValues;
import org.springframework.beans.TypeConverter;
import org.springframework.beans.factory.*;
import org.springframework.beans.factory.annotation.InjectionMetadata;
import org.springframework.beans.factory.annotation.Lookup;
import org.springframework.beans.factory.config.ConfigurableListableBeanFactory;
import org.springframework.beans.factory.config.DependencyDescriptor;
import org.springframework.beans.factory.config.SmartInstantiationAwareBeanPostProcessor;
import org.springframework.beans.factory.support.LookupOverride;
import org.springframework.beans.factory.support.MergedBeanDefinitionPostProcessor;
import org.springframework.beans.factory.support.RootBeanDefinition;
import org.springframework.core.BridgeMethodResolver;
import org.springframework.core.MethodParameter;
import org.springframework.core.Ordered;
import org.springframework.core.PriorityOrdered;
import org.springframework.core.annotation.AnnotationUtils;
import org.springframework.core.annotation.MergedAnnotation;
import org.springframework.core.annotation.MergedAnnotations;
import org.springframework.lang.Nullable;
import org.springframework.util.Assert;
import org.springframework.util.ClassUtils;
import org.springframework.util.ReflectionUtils;
import org.springframework.util.StringUtils;

import java.beans.PropertyDescriptor;
import java.lang.annotation.Annotation;
import java.lang.reflect.*;
import java.util.*;
import java.util.concurrent.ConcurrentHashMap;

/**
 * @author coyzeng@gmail.com
 */
@Slf4j
public class MeshConsumerBeanPostProcessor implements SmartInstantiationAwareBeanPostProcessor,
        MergedBeanDefinitionPostProcessor, PriorityOrdered, BeanFactoryAware {

    private final Set<Class<? extends Annotation>> autowiredAnnotationTypes = new LinkedHashSet<>(4);

    private final Set<String> lookupMethodsChecked = Collections.newSetFromMap(new ConcurrentHashMap<>(256));

    private final Map<Class<?>, Constructor<?>[]> candidateConstructorsCache = new ConcurrentHashMap<>(256);

    private final Map<String, InjectionMetadata> injectionMetadataCache = new ConcurrentHashMap<>(256);

    @Nullable
    private ConfigurableListableBeanFactory beanFactory;

    public MeshConsumerBeanPostProcessor() {
        this.autowiredAnnotationTypes.add(MPI.class);
    }


    @Override
    public int getOrder() {
        return Ordered.LOWEST_PRECEDENCE - 2;
    }

    @Override
    public void setBeanFactory(BeanFactory beanFactory) {
        if (!(beanFactory instanceof ConfigurableListableBeanFactory)) {
            throw new CompatibleException("AutowiredAnnotationBeanPostProcessor requires a ConfigurableListableBeanFactory: " + beanFactory);
        }
        this.beanFactory = (ConfigurableListableBeanFactory) beanFactory;
    }


    @Override
    public void postProcessMergedBeanDefinition(RootBeanDefinition beanDefinition, Class<?> beanType, String beanName) {
        InjectionMetadata metadata = findAutowiringMetadata(beanName, beanType, null);
        metadata.checkConfigMembers(beanDefinition);
    }

    @Override
    public void resetBeanDefinition(String beanName) {
        this.lookupMethodsChecked.remove(beanName);
        this.injectionMetadataCache.remove(beanName);
    }

    @Override
    @Nullable
    public Constructor<?>[] determineCandidateConstructors(Class<?> beanClass, final String beanName)
            throws BeanCreationException {

        // Let's check for lookup methods here...
        if (!this.lookupMethodsChecked.contains(beanName)) {
            if (AnnotationUtils.isCandidateClass(beanClass, Lookup.class)) {
                try {
                    Class<?> targetClass = beanClass;
                    do {
                        ReflectionUtils.doWithLocalMethods(targetClass, method -> {
                            Lookup lookup = method.getAnnotation(Lookup.class);
                            if (lookup != null) {
                                Assert.state(this.beanFactory != null, "No BeanFactory available");
                                LookupOverride override = new LookupOverride(method, lookup.value());
                                try {
                                    RootBeanDefinition mbd = (RootBeanDefinition)
                                            this.beanFactory.getMergedBeanDefinition(beanName);
                                    mbd.getMethodOverrides().addOverride(override);
                                } catch (NoSuchBeanDefinitionException ex) {
                                    throw new BeanCreationException(beanName,
                                            "Cannot apply @Lookup to beans without corresponding bean definition");
                                }
                            }
                        });
                        targetClass = targetClass.getSuperclass();
                    }
                    while (targetClass != null && targetClass != Object.class);

                } catch (IllegalStateException ex) {
                    throw new BeanCreationException(beanName, "Lookup method resolution failed", ex);
                }
            }
            this.lookupMethodsChecked.add(beanName);
        }

        // Quick check on the concurrent map first, with minimal locking.
        Constructor<?>[] candidateConstructors = this.candidateConstructorsCache.get(beanClass);
        if (candidateConstructors == null) {
            // Fully synchronized resolution now...
            synchronized (this.candidateConstructorsCache) {
                candidateConstructors = this.candidateConstructorsCache.get(beanClass);
                if (candidateConstructors == null) {
                    Constructor<?>[] rawCandidates;
                    try {
                        rawCandidates = beanClass.getDeclaredConstructors();
                    } catch (Throwable ex) {
                        throw new BeanCreationException(beanName,
                                "Resolution of declared constructors on bean Class [" + beanClass.getName() +
                                        "] from ClassLoader [" + beanClass.getClassLoader() + "] failed", ex);
                    }
                    List<Constructor<?>> candidates = new ArrayList<>(rawCandidates.length);
                    Constructor<?> requiredConstructor = null;
                    Constructor<?> defaultConstructor = null;
                    Constructor<?> primaryConstructor = BeanUtils.findPrimaryConstructor(beanClass);
                    int nonSyntheticConstructors = 0;
                    for (Constructor<?> candidate : rawCandidates) {
                        if (!candidate.isSynthetic()) {
                            nonSyntheticConstructors++;
                        } else if (primaryConstructor != null) {
                            continue;
                        }
                        MergedAnnotation<?> ann = findAutowiredAnnotation(candidate);
                        if (ann == null) {
                            Class<?> userClass = ClassUtils.getUserClass(beanClass);
                            if (userClass != beanClass) {
                                try {
                                    Constructor<?> superCtor =
                                            userClass.getDeclaredConstructor(candidate.getParameterTypes());
                                    ann = findAutowiredAnnotation(superCtor);
                                } catch (NoSuchMethodException ex) {
                                    // Simply proceed, no equivalent superclass constructor found...
                                }
                            }
                        }
                        if (ann != null) {
                            if (requiredConstructor != null) {
                                throw new BeanCreationException(beanName,
                                        "Invalid autowire-marked constructor: " + candidate +
                                                ". Found constructor with 'required' Autowired annotation already: " +
                                                requiredConstructor);
                            }
                            boolean required = determineRequiredStatus(ann);
                            if (required) {
                                if (!candidates.isEmpty()) {
                                    throw new BeanCreationException(beanName,
                                            "Invalid autowire-marked constructors: " + candidates +
                                                    ". Found constructor with 'required' Autowired annotation: " +
                                                    candidate);
                                }
                                requiredConstructor = candidate;
                            }
                            candidates.add(candidate);
                        } else if (candidate.getParameterCount() == 0) {
                            defaultConstructor = candidate;
                        }
                    }
                    if (!candidates.isEmpty()) {
                        // Add default constructor to list of optional constructors, as fallback.
                        if (requiredConstructor == null) {
                            if (defaultConstructor != null) {
                                candidates.add(defaultConstructor);
                            } else if (candidates.size() == 1) {
                                log.info("Inconsistent constructor declaration on bean with name '{}': single autowire-marked constructor flagged as optional - " +
                                        "this constructor is effectively required since there is no default constructor to fall back to: {}", beanName, candidates.get(0));
                            }
                        }
                        candidateConstructors = candidates.toArray(new Constructor<?>[0]);
                    } else if (rawCandidates.length == 1 && rawCandidates[0].getParameterCount() > 0) {
                        candidateConstructors = new Constructor<?>[]{rawCandidates[0]};
                    } else if (nonSyntheticConstructors == 2 && primaryConstructor != null &&
                            defaultConstructor != null && !primaryConstructor.equals(defaultConstructor)) {
                        candidateConstructors = new Constructor<?>[]{primaryConstructor, defaultConstructor};
                    } else if (nonSyntheticConstructors == 1 && primaryConstructor != null) {
                        candidateConstructors = new Constructor<?>[]{primaryConstructor};
                    } else {
                        candidateConstructors = new Constructor<?>[0];
                    }
                    this.candidateConstructorsCache.put(beanClass, candidateConstructors);
                }
            }
        }
        return (candidateConstructors.length > 0 ? candidateConstructors : null);
    }

    @Override
    public PropertyValues postProcessProperties(PropertyValues pvs, Object bean, String beanName) {
        InjectionMetadata metadata = findAutowiringMetadata(beanName, bean.getClass(), pvs);
        try {
            metadata.inject(bean, beanName, pvs);
        } catch (BeanCreationException e) {
            throw e;
        } catch (Throwable e) {
            throw new BeanCreationException(beanName, "Injection of autowired dependencies failed", e);
        }
        return pvs;
    }

    /**
     * 'Native' processing method for direct calls with an arbitrary target instance,
     * resolving all of its fields and methods which are annotated with one of the
     * configured 'autowired' annotation types.
     *
     * @param bean the target instance to process
     * @throws BeanCreationException if autowiring failed
     */
    public void processInjection(Object bean) throws BeanCreationException {
        Class<?> clazz = bean.getClass();
        InjectionMetadata metadata = findAutowiringMetadata(clazz.getName(), clazz, null);
        try {
            metadata.inject(bean, null, null);
        } catch (BeanCreationException ex) {
            throw ex;
        } catch (Throwable ex) {
            throw new BeanCreationException(
                    "Injection of autowired dependencies failed for class [" + clazz + "]", ex);
        }
    }


    private InjectionMetadata findAutowiringMetadata(String beanName, Class<?> clazz, @Nullable PropertyValues pvs) {
        // Fall back to class name as cache key, for backwards compatibility with custom callers.
        String cacheKey = (StringUtils.hasLength(beanName) ? beanName : clazz.getName());
        // Quick check on the concurrent map first, with minimal locking.
        InjectionMetadata metadata = this.injectionMetadataCache.get(cacheKey);
        if (InjectionMetadata.needsRefresh(metadata, clazz)) {
            synchronized (this.injectionMetadataCache) {
                metadata = this.injectionMetadataCache.get(cacheKey);
                if (InjectionMetadata.needsRefresh(metadata, clazz)) {
                    if (metadata != null) {
                        metadata.clear(pvs);
                    }
                    metadata = buildAutowiringMetadata(clazz);
                    this.injectionMetadataCache.put(cacheKey, metadata);
                }
            }
        }
        return metadata;
    }

    private InjectionMetadata buildAutowiringMetadata(final Class<?> clazz) {
        if (!AnnotationUtils.isCandidateClass(clazz, this.autowiredAnnotationTypes)) {
            return InjectionMetadata.EMPTY;
        }

        List<InjectionMetadata.InjectedElement> elements = new ArrayList<>();
        Class<?> targetClass = clazz;

        do {
            final List<InjectionMetadata.InjectedElement> currElements = new ArrayList<>();

            ReflectionUtils.doWithLocalFields(targetClass, field -> {
                MergedAnnotation<?> ann = findAutowiredAnnotation(field);
                if (ann != null) {
                    if (Modifier.isStatic(field.getModifiers())) {
                        log.info("Autowired annotation is not supported on static fields: {}", field);
                        return;
                    }
                    boolean required = determineRequiredStatus(ann);
                    currElements.add(new AutowiredFieldElement(field, required));
                }
            });

            ReflectionUtils.doWithLocalMethods(targetClass, method -> {
                Method bridgedMethod = BridgeMethodResolver.findBridgedMethod(method);
                if (!BridgeMethodResolver.isVisibilityBridgeMethodPair(method, bridgedMethod)) {
                    return;
                }
                MergedAnnotation<?> ann = findAutowiredAnnotation(bridgedMethod);
                if (ann != null && method.equals(ClassUtils.getMostSpecificMethod(method, clazz))) {
                    if (Modifier.isStatic(method.getModifiers())) {
                        log.info("Autowired annotation is not supported on static methods: {}", method);
                        return;
                    }
                    if (method.getParameterCount() == 0) {
                        log.info("Autowired annotation should only be used on methods with parameters: {}", method);
                    }
                    boolean required = determineRequiredStatus(ann);
                    PropertyDescriptor pd = BeanUtils.findPropertyForMethod(bridgedMethod, clazz);
                    currElements.add(new AutowiredMethodElement(method, required, pd));
                }
            });

            elements.addAll(0, currElements);
            targetClass = targetClass.getSuperclass();
        }
        while (targetClass != null && targetClass != Object.class);

        return InjectionMetadata.forElements(elements, clazz);
    }

    @Nullable
    private MergedAnnotation<?> findAutowiredAnnotation(AccessibleObject ao) {
        MergedAnnotations annotations = MergedAnnotations.from(ao);
        for (Class<? extends Annotation> type : this.autowiredAnnotationTypes) {
            MergedAnnotation<?> annotation = annotations.get(type);
            if (annotation.isPresent()) {
                return annotation;
            }
        }
        return null;
    }

    /**
     * The following (AnnotationAttributes) cast is required on JDK 9+.
     * AnnotationAttributes annotationAttributes = ann.asMap(mergedAnnotation -> new AnnotationAttributes(mergedAnnotation.getType()));
     * <p>
     * Determine if the annotated field or method requires its dependency.
     * <p>A 'required' dependency means that autowiring should fail when no beans
     * are found. Otherwise, the autowiring process will simply bypass the field
     * or method when no beans are found.
     *
     * @param ann the Autowired annotation
     * @return whether the annotation indicates that a dependency is required
     */
    @SuppressWarnings({"cast"})
    protected boolean determineRequiredStatus(MergedAnnotation<?> ann) {
        return true;
    }

    /**
     * Obtain all beans of the given type as autowire candidates.
     *
     * @param type the type of the bean
     * @return the target beans, or an empty Collection if no bean of this type is found
     * @throws BeansException if bean retrieval failed
     */
    protected <T> Map<String, T> findAutowireCandidates(Class<T> type) throws BeansException {
        if (this.beanFactory == null) {
            throw new IllegalStateException("No BeanFactory configured - " +
                    "override the getBeanOfType method or specify the 'beanFactory' property");
        }
        return BeanFactoryUtils.beansOfTypeIncludingAncestors(this.beanFactory, type);
    }

    /**
     * Register the specified bean as dependent on the autowired beans.
     */
    private void registerDependentBeans(@Nullable String beanName, Set<String> autowiredBeanNames) {
        if (beanName != null) {
            for (String autowiredBeanName : autowiredBeanNames) {
                if (this.beanFactory != null && this.beanFactory.containsBean(autowiredBeanName)) {
                    this.beanFactory.registerDependentBean(autowiredBeanName, beanName);
                }
                log.trace("MPI inject by type from bean name '{}' to bean named '{}'", beanName, autowiredBeanName);
            }
        }
    }

    /**
     * Resolve the specified cached method argument or field value.
     */
    @Nullable
    private Object resolvedCachedArgument(@Nullable String beanName, @Nullable Object cachedArgument) {
        if (cachedArgument instanceof DependencyDescriptor) {
            DependencyDescriptor descriptor = (DependencyDescriptor) cachedArgument;
            Assert.state(this.beanFactory != null, "No BeanFactory available");
            return this.beanFactory.resolveDependency(descriptor, beanName, null, null);
        } else {
            return cachedArgument;
        }
    }


    /**
     * Class representing injection information about an annotated field.
     */
    private class AutowiredFieldElement extends InjectionMetadata.InjectedElement {

        private final boolean required;

        private volatile boolean cached;

        @Nullable
        private volatile Object cachedFieldValue;

        public AutowiredFieldElement(Field field, boolean required) {
            super(field, null);
            this.required = required;
        }

        @Override
        protected void inject(Object bean, @Nullable String beanName, @Nullable PropertyValues pvs) throws Throwable {
            Field field = (Field) this.member;
            Object value;
            if (this.cached) {
                try {
                    value = resolvedCachedArgument(beanName, this.cachedFieldValue);
                } catch (NoSuchBeanDefinitionException ex) {
                    // Unexpected removal of target bean for cached argument -> re-resolve
                    value = resolveFieldValue(field, bean, beanName);
                }
            } else {
                value = resolveFieldValue(field, bean, beanName);
            }
            if (value != null) {
                ReflectionUtils.makeAccessible(field);
                field.set(bean, value);
            }
        }

        @Nullable
        private Object resolveFieldValue(Field field, Object bean, @Nullable String beanName) {
            DependencyDescriptor desc = new DependencyDescriptor(field, this.required);
            desc.setContainingClass(bean.getClass());
            Set<String> autowiredBeanNames = new LinkedHashSet<>(1);
            Assert.state(beanFactory != null, "No BeanFactory available");
            TypeConverter typeConverter = beanFactory.getTypeConverter();
            Object value;
            try {
                value = beanFactory.resolveDependency(desc, beanName, autowiredBeanNames, typeConverter);
            } catch (BeansException ex) {
                throw new UnsatisfiedDependencyException(null, beanName, new InjectionPoint(field), ex);
            }
            synchronized (this) {
                if (!this.cached) {
                    Object cachedFieldValue = null;
                    if (value != null || this.required) {
                        cachedFieldValue = desc;
                        registerDependentBeans(beanName, autowiredBeanNames);
                        if (autowiredBeanNames.size() == 1) {
                            String autowiredBeanName = autowiredBeanNames.iterator().next();
                            if (beanFactory.containsBean(autowiredBeanName) &&
                                    beanFactory.isTypeMatch(autowiredBeanName, field.getType())) {
                                cachedFieldValue = new ShortcutDependencyDescriptor(
                                        desc, autowiredBeanName, field.getType());
                            }
                        }
                    }
                    this.cachedFieldValue = cachedFieldValue;
                    this.cached = true;
                }
            }
            return value;
        }
    }


    /**
     * Class representing injection information about an annotated method.
     */
    @EqualsAndHashCode(callSuper = true)
    private class AutowiredMethodElement extends InjectionMetadata.InjectedElement {

        private final boolean required;

        private volatile boolean cached;

        @Nullable
        private volatile Object[] cachedMethodArguments;

        public AutowiredMethodElement(Method method, boolean required, @Nullable PropertyDescriptor pd) {
            super(method, pd);
            this.required = required;
        }

        @Override
        protected void inject(Object bean, @Nullable String beanName, @Nullable PropertyValues pvs) throws Throwable {
            if (checkPropertySkipping(pvs)) {
                return;
            }
            Method method = (Method) this.member;
            Object[] arguments;
            if (this.cached) {
                try {
                    arguments = resolveCachedArguments(beanName);
                } catch (NoSuchBeanDefinitionException ex) {
                    // Unexpected removal of target bean for cached argument -> re-resolve
                    arguments = resolveMethodArguments(method, bean, beanName);
                }
            } else {
                arguments = resolveMethodArguments(method, bean, beanName);
            }
            if (arguments != null) {
                try {
                    ReflectionUtils.makeAccessible(method);
                    method.invoke(bean, arguments);
                } catch (InvocationTargetException ex) {
                    throw ex.getTargetException();
                }
            }
        }

        @Nullable
        private Object[] resolveCachedArguments(@Nullable String beanName) {
            Object[] cachedMethodArguments = this.cachedMethodArguments;
            if (cachedMethodArguments == null) {
                return null;
            }
            Object[] arguments = new Object[cachedMethodArguments.length];
            for (int i = 0; i < arguments.length; i++) {
                arguments[i] = resolvedCachedArgument(beanName, cachedMethodArguments[i]);
            }
            return arguments;
        }

        @Nullable
        private Object[] resolveMethodArguments(Method method, Object bean, @Nullable String beanName) {
            int argumentCount = method.getParameterCount();
            Object[] arguments = new Object[argumentCount];
            DependencyDescriptor[] descriptors = new DependencyDescriptor[argumentCount];
            Set<String> autowiredBeans = new LinkedHashSet<>(argumentCount);
            Assert.state(beanFactory != null, "No BeanFactory available");
            TypeConverter typeConverter = beanFactory.getTypeConverter();
            for (int i = 0; i < arguments.length; i++) {
                MethodParameter methodParam = new MethodParameter(method, i);
                DependencyDescriptor currDesc = new DependencyDescriptor(methodParam, this.required);
                currDesc.setContainingClass(bean.getClass());
                descriptors[i] = currDesc;
                try {
                    Object arg = beanFactory.resolveDependency(currDesc, beanName, autowiredBeans, typeConverter);
                    if (arg == null && !this.required) {
                        arguments = null;
                        break;
                    }
                    arguments[i] = arg;
                } catch (BeansException ex) {
                    throw new UnsatisfiedDependencyException(null, beanName, new InjectionPoint(methodParam), ex);
                }
            }
            synchronized (this) {
                if (!this.cached) {
                    if (arguments != null) {
                        DependencyDescriptor[] cachedMethodArguments = Arrays.copyOf(descriptors, arguments.length);
                        registerDependentBeans(beanName, autowiredBeans);
                        if (autowiredBeans.size() == argumentCount) {
                            Iterator<String> it = autowiredBeans.iterator();
                            Class<?>[] paramTypes = method.getParameterTypes();
                            for (int i = 0; i < paramTypes.length; i++) {
                                String autowiredBeanName = it.next();
                                if (beanFactory.containsBean(autowiredBeanName) &&
                                        beanFactory.isTypeMatch(autowiredBeanName, paramTypes[i])) {
                                    cachedMethodArguments[i] = new ShortcutDependencyDescriptor(
                                            descriptors[i], autowiredBeanName, paramTypes[i]);
                                }
                            }
                        }
                        this.cachedMethodArguments = cachedMethodArguments;
                    } else {
                        this.cachedMethodArguments = null;
                    }
                    this.cached = true;
                }
            }
            return arguments;
        }
    }


    /**
     * DependencyDescriptor variant with a pre-resolved target bean name.
     */
    @EqualsAndHashCode(callSuper = true)
    private static class ShortcutDependencyDescriptor extends DependencyDescriptor {

        private static final long serialVersionUID = 1506066929177002393L;
        private final String shortcut;

        private final Class<?> requiredType;

        public ShortcutDependencyDescriptor(DependencyDescriptor original, String shortcut, Class<?> requiredType) {
            super(original);
            this.shortcut = shortcut;
            this.requiredType = requiredType;
        }

        @Override
        public Object resolveShortcut(BeanFactory beanFactory) {
            return beanFactory.getBean(this.shortcut, this.requiredType);
        }
    }
}
