/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import {useNavigate} from 'react-router-dom';

function Footer(props: any) {

    const navigate = useNavigate();
    const clickMenu = (menu: { name?: string; path: any; }) => {
        navigate(menu.path, {replace: true});
    };
    return (
        <div>
        </div>
    );
}

export default Footer;
