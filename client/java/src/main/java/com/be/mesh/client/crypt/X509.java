/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.crypt;


import com.be.mesh.client.cause.CryptException;
import com.be.mesh.client.mpc.MeshCode;
import org.bouncycastle.asn1.x500.X500Name;
import org.bouncycastle.asn1.x509.AlgorithmIdentifier;
import org.bouncycastle.asn1.x509.Certificate;
import org.bouncycastle.asn1.x509.SubjectPublicKeyInfo;
import org.bouncycastle.cert.X509CertificateHolder;
import org.bouncycastle.cert.X509v3CertificateBuilder;
import org.bouncycastle.crypto.params.AsymmetricKeyParameter;
import org.bouncycastle.crypto.util.PrivateKeyFactory;
import org.bouncycastle.jce.provider.BouncyCastleProvider;
import org.bouncycastle.operator.ContentSigner;
import org.bouncycastle.operator.DefaultDigestAlgorithmIdentifierFinder;
import org.bouncycastle.operator.DefaultSignatureAlgorithmIdentifierFinder;
import org.bouncycastle.operator.bc.BcRSAContentSignerBuilder;
import org.bouncycastle.operator.jcajce.JcaContentSignerBuilder;
import org.bouncycastle.pkcs.PKCS10CertificationRequest;
import org.bouncycastle.pkcs.PKCS10CertificationRequestBuilder;
import org.bouncycastle.pkcs.jcajce.JcaPKCS10CertificationRequestBuilder;

import javax.security.auth.x500.X500Principal;
import java.io.ByteArrayInputStream;
import java.io.InputStream;
import java.math.BigInteger;
import java.security.KeyPair;
import java.security.PrivateKey;
import java.security.PublicKey;
import java.security.Security;
import java.security.cert.CertificateFactory;
import java.security.cert.X509Certificate;
import java.util.Date;
import java.util.HashMap;
import java.util.Map;

/**
 * @author coyzeng@gmail.com
 */
public class X509 implements Crypt {

    static {
        if (null == Security.getProvider(BouncyCastleProvider.PROVIDER_NAME)) {
            Security.addProvider(new BouncyCastleProvider());
        }
    }

    private final Crypt ca;
    private final String signatureAlgo;
    private final String keyPrefix;

    /**
     * 生成 PKCS#10 证书请求. RSA P10 证书请求 Base64 字符串
     *
     * @param isRsaNotEcc {@code true}：使用 RSA 加密算法；{@code false}：使用 ECC（SM2）加密算法
     */
    public X509(boolean isRsaNotEcc) {
        this.ca = new RSA2();
        this.keyPrefix = isRsaNotEcc ? "" : "EC ";
        this.signatureAlgo = RSA2.SIGNATURE_ALGORITHM;
    }

    @Override
    public KeyPair generateKeyPair() {
        return this.ca.generateKeyPair();
    }

    @Override
    public Map<String, String> generateKeys(Codec codec) {
        try {
            // 使用 RSA/ECC 算法，生成密钥对（公钥、私钥）
            KeyPair keyPair = this.generateKeyPair();
            PrivateKey privateKey = keyPair.getPrivate();
            PublicKey publicKey = keyPair.getPublic();

            // 按需添加证书主题项, 有些 CSR 不需要我们在主题项中添加各字段, 而是通过额外参数提交，故这里我只简单地指定了国家码
            // 如 `C=CN, CN=ducesoft.com, E=coyzeng@ducesoft.com, OU=3303..., L=杭州, S=浙江`
            X500Principal subject = new X500Principal("C=CN");

            // 使用私钥和 SHA256WithRSA/SM3withSM2 算法创建签名者对象
            JcaContentSignerBuilder jcb = new JcaContentSignerBuilder(this.signatureAlgo);
            ContentSigner signer = jcb.setProvider(BouncyCastleProvider.PROVIDER_NAME).build(privateKey);

            // CSR
            PKCS10CertificationRequestBuilder builder = new JcaPKCS10CertificationRequestBuilder(subject, publicKey);
            PKCS10CertificationRequest csr = builder.build(signer);
            // X509证书
            X509Certificate x509 = createX509(csr, privateKey);

            // 私钥自己保存，公钥对外公布
            Map<String, String> keyStore = new HashMap<>(4);
            keyStore.put("private.key", codec.encodeKey(privateKey.getEncoded(), String.format("%sPRIVATE KEY", this.keyPrefix)));
            keyStore.put("public.key", codec.encodeKey(publicKey.getEncoded(), String.format("%sPUBLIC KEY", this.keyPrefix)));
            keyStore.put("cert.csr", codec.encodeKey(csr.getEncoded(), "CERTIFICATE REQUEST"));
            keyStore.put("x509.pem", codec.encodeKey(x509.getEncoded(), "CERTIFICATE"));
            return keyStore;
        } catch (Exception e) {
            throw new CryptException(MeshCode.CRYPT_ERROR, e);
        }
    }

    @Override
    public String signature(String data, String privateKey, Codec codec) {
        return this.ca.signature(data, privateKey, codec);
    }

    @Override
    public boolean verify(String data, String signature, String publicKey, Codec codec) {
        return this.ca.verify(data, signature, publicKey, codec);
    }

    @Override
    public String encrypt(String data, String publicKey, Codec codec) {
        return this.ca.encrypt(data, publicKey, codec);
    }

    @Override
    public String decrypt(String data, String privateKey, Codec codec) {
        return this.ca.decrypt(data, privateKey, codec);
    }

    /**
     * 签发证书.
     *
     * @param csr        证书签发请求
     * @param privateKey 私钥
     * @return x509标准证书
     * @throws Exception 异常
     */
    public X509Certificate createX509(PKCS10CertificationRequest csr, PrivateKey privateKey) throws Exception {
        AlgorithmIdentifier algoID = new DefaultSignatureAlgorithmIdentifierFinder().find(this.signatureAlgo);
        AlgorithmIdentifier digestAlgoID = new DefaultDigestAlgorithmIdentifierFinder().find(algoID);

        AsymmetricKeyParameter key = PrivateKeyFactory.createKey(privateKey.getEncoded());
        SubjectPublicKeyInfo info = csr.getSubjectPublicKeyInfo();

        Date notBefore = new Date(System.currentTimeMillis());
        Date notAfter = new Date(System.currentTimeMillis() + 30L * 365 * 24 * 60 * 60 * 1000);
        X500Name name = new X500Name("CN=trustbe.com");
        BigInteger serial = BigInteger.valueOf(System.currentTimeMillis());
        ContentSigner signer = new BcRSAContentSignerBuilder(algoID, digestAlgoID).build(key);
        X509v3CertificateBuilder builder = new X509v3CertificateBuilder(name, serial, notBefore, notAfter, csr.getSubject(), info);
        X509CertificateHolder holder = builder.build(signer);
        Certificate certificate = holder.toASN1Structure();

        CertificateFactory cf = CertificateFactory.getInstance("X.509", BouncyCastleProvider.PROVIDER_NAME);
        InputStream input = new ByteArrayInputStream(certificate.getEncoded());
        X509Certificate x509 = (X509Certificate) cf.generateCertificate(input);
        input.close();
        return x509;
    }

}
