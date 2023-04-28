/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
package com.be.mesh.client.mpc;

import lombok.AllArgsConstructor;
import lombok.Getter;

import java.util.Arrays;

/**
 * @author coyzeng@gmail.com
 */
@Getter
@AllArgsConstructor
public enum Runmode {

    // 正常模式
    ROUTINE(1),
    // 评测模式
    PERFORM(2),
    // 高防模式
    DEFENSE(4),
    // 调试模式,
    DEBUG(8),
    // 压测模式
    LOAD_TEST(16),
    // Mock模式
    MOCK(32),
    ;

    private final int mode;

    public boolean isDebug() {
        return matches(DEBUG, this.getMode());
    }

    public boolean isLoadTest() {
        return matches(LOAD_TEST, this.getMode());
    }

    public boolean isRoutine() {
        return matches(ROUTINE, this.getMode());
    }

    public boolean isPerform() {
        return matches(PERFORM, this.getMode());
    }

    public boolean isDefense() {
        return matches(DEFENSE, this.getMode());
    }

    public boolean isMock() {
        return matches(MOCK, this.getMode());
    }

    public static Runmode from(int code) {
        return Arrays.stream(Runmode.values()).filter(x -> (x.mode & code) == x.mode).findFirst().orElse(ROUTINE);
    }

    public static boolean matches(Runmode runmode, int code) {
        return (runmode.mode & code) == runmode.mode;
    }

    public static boolean isDebug(int code) {
        return matches(DEBUG, code);
    }

    public static boolean isLoadTest(int code) {
        return matches(LOAD_TEST, code);
    }

    public static boolean isRoutine(int code) {
        return matches(ROUTINE, code);
    }

    public static boolean isPerform(int code) {
        return matches(PERFORM, code);
    }

    public static boolean isDefense(int code) {
        return matches(DEFENSE, code);
    }

    public static boolean isMock(int code) {
        return matches(MOCK, code);
    }

}
