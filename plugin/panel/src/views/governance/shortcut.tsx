/*
 * Copyright (c) 2000, 2099, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import {Terminal} from 'xterm';
import {SearchAddon} from 'xterm-addon-search';
import {FitAddon} from 'xterm-addon-fit';
import {CanvasAddon} from 'xterm-addon-canvas';
import {Box} from "@mui/material";
import {createRef, useEffect} from "react";
import "xterm/css/xterm.css";
import service from "@/services/service";
import {Codec, context, Mesh, ServiceLoader, Versions} from '@mesh/mesh';


export default function Shortcut() {

    const ref = createRef<HTMLDivElement>();
    const version = new Versions();
    const commands: string[] = [];
    const codec = ServiceLoader.load(Codec).getDefault();
    const terminal = new Terminal({
        windowsMode: ['Windows', 'Win16', 'Win32', 'WinCE'].indexOf(navigator.platform) >= 0,
        fontFamily: '"Cascadia Code", Menlo, monospace',
        cursorBlink: true,
        screenReaderMode: false,
        allowProposedApi: true,
        rows: 30,
        theme: {
            foreground: '#F8F8F8',
            background: '#2D2E2C',
            selectionBackground: '#5DA5D533',
            black: '#1E1E1D',
            brightBlack: '#262625',
            red: '#CE5C5C',
            brightRed: '#FF7272',
            green: '#5BCC5B',
            brightGreen: '#72FF72',
            yellow: '#CCCC5B',
            brightYellow: '#FFFF72',
            blue: '#5D5DD3',
            brightBlue: '#7279FF',
            magenta: '#BC5ED1',
            brightMagenta: '#E572FF',
            cyan: '#5DA5D5',
            brightCyan: '#72F0FF',
            white: '#F8F8F8',
            brightWhite: '#FFFFFF'
        },
    })
    terminal.registerLinkProvider({
        provideLinks(bufferLineNumber, callback) {
            callback([
                {
                    text: 'MESH',
                    range: {start: {x: 0, y: 0}, end: {x: 1, y: 1}},
                    activate() {
                        window.open('https://github.com/mesh/mesh', '_blank');
                    }
                },
            ]);
        }
    });
    terminal.onData(e => {
        switch (e) {
            case '\u0003': // Ctrl+C
                terminal.write('^C');
                terminal.write('\r\n$ ');
                break;
            case '\r': // Enter
                onCommandInput(commands.join(''));
                commands.splice(0, commands.length);
                terminal.write('\r\n$ ');
                break;
            case '\u007F': // Backspace (DEL)
                // Do not delete the prompt
                if (terminal.buffer.active.cursorX > 2) {
                    terminal.write('\b \b');
                    if (commands.length > 0) {
                        commands.splice(commands.length - 1, 1);
                    }
                }
                break;
            default: // Print all other characters
                if (e >= String.fromCharCode(0x20) && e <= String.fromCharCode(0x7E) || e >= '\u00a0') {
                    terminal.write(e);
                    commands.push(e);
                }
        }
    });
    const fitAddon = new FitAddon();
    const searchAddon = new SearchAddon();
    const canvasAddon = new CanvasAddon();

    const initTerm = (e: HTMLDivElement, v: Versions) => {
        // Cancel wheel events from scrolling the page if the terminal has scroll back
        e.addEventListener('wheel', e => {
            if (terminal.buffer.active.baseY > 0) {
                e.preventDefault();
            }
        });
        terminal.open(e);
        terminal.focus();
        terminal.loadAddon(canvasAddon);
        terminal.loadAddon(fitAddon);
        fitAddon.fit();
        terminal.loadAddon(searchAddon);
        searchAddon.findNext('');
        terminal.writeln(logo(v));
        terminal.write("$ ");
    }

    const logo = (v: Versions): string => {
        return [
            ``,
            `	 __  __           _     `,
            `	|  \\/  | ___  ___| |__  `,
            `	| |\\/| |/ _ \\/ __| '_ \\      A \x1b[32mlightweight\x1b[0m, \x1b[33;1mdistributed\x1b[0m, \x1b[35;1mrelational\x1b[0m`,
            `	| |  | |  __/\\__ \\ | | |     network architecture for MPC`,
            `	|_|  |_|\\___||___/_| |_|`,
            ``,
            `	(v\x1b[31;1m${v?.version || '1.0.0.0'}\x1b[0m, build \x1b[36m${Object.getOwnPropertyDescriptor(v?.infos || {}, 'mesh.commit_id')?.value || '59a06ccd'}\x1b[0m)`,
            ``,
            ``,
        ].join('\r\n')
    }

    const onCommandInput = (command: string): void => {
        switch (command) {
            case "":
                break;
            case "clear":
                terminal.reset();
                terminal.write(logo(version));
                break
            default:
                const ctx = context();
                ctx.setAttribute(Mesh.UNAME, "mesh.dot.exe");
                service.endpoint.fuzzy(ctx, codec.encode(command)).then(v => {
                    const std = codec.decode(v, Object) as string;
                    std.split('\n').filter(v => "" != v).forEach(z => {
                        if (z.startsWith('{') && z.endsWith('}')) {
                            const dict = codec.decodeString(z, Map);
                            terminal.write(`\u001b[34m${dict.get('timestamp')}.000\u001b[0m \u001b[32m${formatLevel(dict.get('level'))} \u001b[0m${dict.get('msg')}`);
                            terminal.write('\r\n$ ')
                            return;
                        }
                        terminal.write(z.trim().replace('\u001b[90m', '\u001b[34m'));
                        terminal.write('\r\n$ ')
                    });
                })
        }
    }

    const formatLevel = (level: string): string => {
        switch (level) {
            case 'trace':
                return "TRC";
            case "debug":
                return "DBG";
            case "info":
                return "INF";
            case "warn":
                return "WRN";
            case "error":
                return "ERR";
            case "fatal":
                return "FTL";
            case "panic":
                return "PNC";
            default:
                return "???";
        }
    }

    useEffect(() => {
        const element = ref.current;
        if (!element) {
            return;
        }
        service.network.version(context(), '').then(v => {
            initTerm(element, v);
            Object.assign(version, v);
        }).catch(e => {
            initTerm(element, new Versions());
        });
        return () => {
            canvasAddon.dispose();
            fitAddon.dispose();
            searchAddon.dispose();
            terminal.dispose();
        };
    }, []);


    return (
        <Box>
            <div ref={ref}></div>
        </Box>
    )
}