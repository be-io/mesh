/*
 * Copyright (c) 2000, 2099, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import {ReactElement} from "react";
import Governance from "@/views/governance/governance";
import Keys from "@/views/keys/keys";
import Network from "@/views/network/network";
import Queue from "@/views/queue/queue";
import Store from "@/views/store/store";
import {Route, Routes} from "react-router-dom";
import Home from "@/views/home/home";
import {Signin} from "@/views/home/signin";
import {Signup} from "@/views/home/signup";
import Privacy from "@/views/home/privacy";
import Terms from "@/views/home/terms";

export class Citizen {

    public name: string = '';
    public icon: string = '';
    public path: string = '';
    public route?: ReactElement;
}

export const Citizens: Citizen[] = [{
    name: '服务治理',
    icon: '',
    path: '/mesh/governance',
    route: <Governance/>,
}, {
    name: '网络管理',
    icon: '',
    path: '/mesh/network',
    route: <Network/>,
}, {
    name: '密钥管理',
    icon: '',
    path: '/mesh/keys',
    route: <Keys/>,
}, {
    name: '队列管理',
    icon: '',
    path: '/mesh/queue',
    route: <Queue/>,
}, {
    name: '存储管理',
    icon: '',
    path: '/mesh/store',
    route: <Store/>,
},];

export const Alones: Citizen[] = [{
    name: '登录',
    icon: '',
    path: '/signin',
    route: <Signin/>,
}, {
    name: '注册',
    icon: '',
    path: '/signup',
    route: <Signup/>,
}, {
    name: '隐私策略',
    icon: '',
    path: '/privacy',
    route: <Privacy/>,
}, {
    name: '服务条款',
    icon: '',
    path: '/terms',
    route: <Terms/>,
},]

export function CitizenRouter() {
    return (
        <Routes>
            <Route path="/" element={<Governance/>}/>
            <Route path="/mesh" element={<Governance/>}/>
            {
                Citizens.map(c => <Route key={c.path} path={`${c.path}/*`} element={c.route}/>)
            }
        </Routes>
    )
}

export function AloneRouter() {
    return (
        <Routes>
            <Route path="/" element={<Signin/>}/>
            {
                Alones.map(c => <Route key={c.path} path={`${c.path}/*`} element={c.route}/>)
            }
        </Routes>
    )
}
