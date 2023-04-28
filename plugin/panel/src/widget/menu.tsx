import * as React from 'react';
import Paper from '@mui/material/Paper';
import MenuList from '@mui/material/MenuList';
import MenuItem from '@mui/material/MenuItem';
import ListItemText from '@mui/material/ListItemText';
import ListItemIcon from '@mui/material/ListItemIcon';
import {Box, Divider} from "@mui/material";


export class ItemOption {

    public key: string = ""
    public name: string = "";
    public icon: any = "";
    public onclick: (option: ItemOption) => void = () => {
    };
}

export default function Menu(props: { menus: ItemOption[] }) {
    return (
        <Paper sx={{width: 240, maxWidth: '100%', height: '100vh'}}>
            <MenuList>
                {
                    props.menus && props.menus.map(menu => {
                        return (
                            <Box key={menu.name}>
                                <MenuItem onClick={() => menu.onclick(menu)}>
                                    <ListItemIcon>
                                        {menu.icon}
                                    </ListItemIcon>
                                    <ListItemText>{menu.name}</ListItemText>
                                </MenuItem>
                                <Divider/>
                            </Box>
                        )
                    })
                }
            </MenuList>
        </Paper>
    );
}