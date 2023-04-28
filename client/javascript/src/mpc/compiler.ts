/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import {Cause, Class, dict, Dict, index, MethodInspector, Parameter, Return, spi} from "@/macro";

@spi("mesh")
export abstract class Compiler {

    /**
     * Get the parameter class type.
     */
    abstract intype(method: MethodInspector): Class<Parameter>;

    /**
     * Get the return class type.
     */
    abstract retype(method: MethodInspector): Class<Return>;

    /**
     * Compile source code.
     */
    abstract compile(code: string, name: string): any;
}

export function returnType<T extends (...args: any) => any>(fn: T) {
    return {} as T;
}

@spi("mesh")
export class TSCompiler extends Compiler {

    private static readonly parameters: Dict<string, Class<Parameter>> = new dict();
    private static readonly returns: Dict<string, Class<Return>> = new dict();

    compile(code: string, name: string): any {

    }

    intype(method: MethodInspector): Class<Parameter> {
        return TSCompiler.parameters.computeIfy(method.getSignature(), k => {
            class dyn implements Parameter {

                @index(-1, 'attachments', Map)
                private attachments: Map<string, string> = new Map()

                getAttachments(): Map<string, string> {
                    return this.attachments;
                }

                map(): Map<string, any> {
                    return Object.assign(new Map(), this);
                }

                setAttachments(attachments: Map<string, string>): void {
                    (attachments || {}).forEach((v, k) => this.attachments.set(k, v));
                }

                getArguments(): any[] {
                    return method.getParameters().map((arg) => Object.getOwnPropertyDescriptor(this, arg.name)?.value);
                }

                setArguments(args: any[]): void {
                    method.getParameters().forEach((arg, idx) => {
                        if (args && args.length > arg.index) {
                            Reflect.set(this, arg.name, args[arg.index]);
                        }
                    });
                }
            }

            return dyn;
        })
    }

    retype(method: MethodInspector): Class<Return> {
        return TSCompiler.returns.computeIfy(method.getSignature(), k => {
            class dyn implements Return {

                @index(0)
                private code: string = "";
                @index(5)
                private message: string = "";
                @index(10, 'content', method.getReturnType())
                private content: any;
                @index(15, 'cause', Cause)
                private cause: Cause | undefined;

                getCause(): Cause | undefined {
                    return this.cause;
                }

                getCode(): string {
                    return this.code;
                }

                getContent(): any {
                    return this.content;
                }

                getMessage(): string {
                    return this.message;
                }

                setCause(cause: Cause): void {
                    this.cause = cause;
                }

                setCode(code: string): void {
                    this.code = code;
                }

                setContent(content: any): void {
                    this.content = content;
                }

                setMessage(message: string): void {
                    this.message = message;
                }

            }

            return dyn;
        })
    }

}