package com.be.mesh.client.struct;

import lombok.Data;

import java.io.Serializable;
import java.util.List;
import java.util.Map;

@Data
public class OSChart implements Serializable {

    private String releaseName;

    private String repoName;

    private String chartName;

    private String k8sNamespace;

    private Map<String, Object> chartArgs;

    private List<Byte> cfs;

    private String chartType;

    private String uniqueName;

    private List<Byte> kubeConfig;

}
