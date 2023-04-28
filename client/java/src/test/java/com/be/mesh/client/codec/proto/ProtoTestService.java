package com.be.mesh.client.codec.proto;

/**
 * @author nanfeng
 * @date 2023/2/15 1:48 PM
 */
public interface ProtoTestService {
    /**
     * 模拟proto的MPI服务
     * @return ProtobufTest.PsiExecuteRequest
     */
    ProtobufTest.PsiExecuteResponse mockProto(ProtobufTest.PsiExecuteRequest request);
}
