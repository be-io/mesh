/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import {ReactElement} from "react";
import {Route,} from "react-router-dom";
import Atlas from "@/views/atlas/atlas";
import Fops from "@/views/fops/fops";
import Governance from "@/views/governance/governance";
import Keys from "@/views/keys/keys";
import Network from "@/views/network/network";
import Queue from "@/views/queue/queue";
import Store from "@/views/store/store";
import Traits from "@/views/traits/traits";
import Workflow from "@/views/workflow/workflow"


export class Citizen {
    public name: string = '';
    public icon: string = '';
    public path: string = '';
    public route?: ReactElement;
}

const Citizens: Citizen[] = [{
    name: '服务治理',
    icon: '',
    path: '/mesh/governance',
    route: <Route path="/mesh/governance/*" element={<Governance/>}/>,
}, {
    name: '联邦运维',
    icon: '',
    path: '/mesh/fops',
    route: <Route path="/mesh/fops/*" element={<Fops/>}/>,
}, {
    name: '联邦云图',
    icon: '',
    path: '/mesh/atlas',
    route: <Route path="/mesh/atlas/*" element={<Atlas/>}/>,
}, {
    name: '网络管理',
    icon: '',
    path: '/mesh/network',
    route: <Route path="/mesh/network/*" element={<Network/>}/>,
}, {
    name: '密钥管理',
    icon: '',
    path: '/mesh/keys',
    route: <Route path="/mesh/keys/*" element={<Keys/>}/>,
}, {
    name: '队列管理',
    icon: '',
    path: '/mesh/queue',
    route: <Route path="/mesh/queue/*" element={<Queue/>}/>,
}, {
    name: '存储管理',
    icon: '',
    path: '/mesh/store',
    route: <Route path="/mesh/store/*" element={<Store/>}/>,
}, {
    name: '接口管理',
    icon: '',
    path: '/mesh/traits',
    route: <Route path="/mesh/traits/*" element={<Traits/>}/>,
}, {
    name: '流程管理',
    icon: '',
    path: '/mesh/workflow',
    route: <Route path="/mesh/workflow/*" element={<Workflow/>}/>,
}];

export default Citizens;
