/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

const fs = require('fs');
const path = require('path');
const {inspect} = require('util');
const ts = require('typescript');
const stripComments = require('strip-comments');
const {ObjectHook, TransformPluginContext} = require("rollup");

const theFinder = new RegExp(/((?<![\(\s]\s*['"])@\w*[\w\d]\s*(?![;])[\((?=\s)])/);
const findDecorators = (fc) => theFinder.test(stripComments(fc));

class TSC {

    name = "tsc";
    tsx;
    tsconfig;
    cache = {};

    constructor(tsconfigPath, tsx) {
        this.tsx = tsx;
        this.tsconfig = this.parseTsConfig(tsconfigPath, process.cwd());
        if (this.tsconfig?.options?.sourcemap) {
            this.tsconfig.options.sourcemap = false;
            this.tsconfig.options.inlineSources = true;
            this.tsconfig.options.inlineSourceMap = true;
        }
    }

    load = (id, options) => {
        if (!id.endsWith('.ts') && !id.endsWith('.tsx')) {
            return;
        }
        if (!this.tsconfig || !this.tsconfig.options || !this.tsconfig.options.emitDecoratorMetadata) {
            return;
        }
        const code = fs.readFileSync(id, 'utf8').toString('utf-8');
        // Find the decorator and if there isn't one, return out
        const hasDecorator = findDecorators(code);
        if (!hasDecorator) {
            return;
        }
        const program = ts.transpileModule(code, {compilerOptions: this.tsconfig.options});
        this.cache[id] = program.outputText;
    }

    transform = (code, id, options) => {
        if (!this.cache[id]) {
            return code;
        }
        return this.cache[id];
    }

    parseTsConfig = (tsconfig, cwd = process.cwd()) => {
        const fileName = ts.findConfigFile(cwd, ts.sys.fileExists, tsconfig);
        // if the value was provided, but no file, fail hard
        if (tsconfig !== undefined && !fileName) {
            throw new Error(`failed to open '${fileName}'`);
        }
        let loadedConfig = {};
        let baseDir = cwd;
        if (fileName) {
            const text = ts.sys.readFile(fileName);
            if (text === undefined) {
                throw new Error(`failed to read '${fileName}'`);
            }
            const result = ts.parseConfigFileTextToJson(fileName, text);

            if (result.error !== undefined) {
                this.printDiagnostics(result.error);
                throw new Error(`failed to parse '${fileName}'`);
            }
            loadedConfig = result.config;
            baseDir = path.dirname(fileName);
        }
        const config = ts.parseJsonConfigFileContent(loadedConfig, ts.sys, baseDir);
        if (config.errors[0]) {
            this.printDiagnostics(config.errors);
        }
        return config;
    }

    printDiagnostics = (...args) => {
        console.log(inspect(args, false, 10, true));
    }
}

const tsc = (tsconfigPath = path.join(process.cwd(), './tsconfig.json'), force = false, tsx = true) => {
    return new TSC(tsconfigPath, force, tsx);
}

module.exports = tsc