/*
 * Copyright (c) 2000, 2099, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import {
    Box,
    Button,
    Grid,
    IconButton,
    ImageList,
    ImageListItem,
    ImageListItemBar,
    Paper,
    Tooltip,
    Typography
} from "@mui/material";
import CircleIcon from "@mui/icons-material/Circle";
import * as React from "react";
import {OverridableStringUnion} from "@mui/types";
import {SvgIconPropsColorOverrides} from "@mui/material/SvgIcon/SvgIcon";

export class SetOption {

    public name: string = '';
    public icon: any = '';
    public describe: string = '';
    public replicas: number = 0;
    public status: string = 'success';
    public color?: OverridableStringUnion<
        | 'inherit'
        | 'action'
        | 'disabled'
        | 'primary'
        | 'secondary'
        | 'error'
        | 'info'
        | 'success'
        | 'warning',
        SvgIconPropsColorOverrides
    >;

}

export default function Block(props: { option: SetOption }) {
    return (
        <Grid item xs={3} sm={3} md={3}>
            <Box bgcolor="background.default" sx={{margin: 'auto 20px auto auto'}}>
                <Paper elevation={3}>
                    <Box>
                        <ImageList cols={1}>
                            <ImageListItem>
                                <Box>
                                    <img
                                        src={`${props.option.icon}?w=248&fit=crop&auto=format`}
                                        srcSet={`${props.option.icon}?w=248&fit=crop&auto=format&dpr=2 2x`}
                                        alt={props.option.name}
                                        loading="lazy"
                                    />
                                </Box>
                                <ImageListItemBar
                                    title={props.option.name}
                                    subtitle=""
                                    actionIcon={
                                        <IconButton
                                            sx={{color: 'rgba(255, 255, 255, 0.54)'}}
                                            aria-label={`info about ${props.option.name}`}
                                        >
                                        </IconButton>
                                    }
                                />
                            </ImageListItem>
                        </ImageList>
                    </Box>
                    <Box sx={{margin: 'auto 10px auto 10px'}}>
                        <Box sx={{
                            display: 'flex',
                            justifyContent: 'flex-start',
                            alignItems: 'center',
                            marginTop: '20px',
                            flexDirection: 'row'
                        }}>
                            <Typography align="left" variant="caption" color="text.secondary" noWrap={false}
                                        sx={{lineHeight: 'auto'}}>
                                状态&nbsp;&nbsp;
                            </Typography>
                            <CircleIcon fontSize="small"
                                        color={props.option.color}/>
                            <Typography variant="caption" color="text.secondary" noWrap={false}>
                                &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;实例&nbsp;&nbsp;
                            </Typography>
                            <Box>{props.option.replicas}</Box>
                        </Box>
                        <Tooltip title={props.option.describe} placement="bottom-start">
                            <Typography align="left" variant="caption" color="text.secondary" noWrap={true}
                                        paragraph={true}>
                                {props.option.describe}
                            </Typography>
                        </Tooltip>
                    </Box>
                    <Box>
                        <Button size="small">去运维</Button>
                        <Button size="small">看监控</Button>
                    </Box>
                </Paper>
            </Box>
        </Grid>
    )
}