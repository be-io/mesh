/*
 * Copyright (c) 2000, 2099, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

import {createTheme} from '@mui/material';

// A custom theme for this app
const theme = createTheme({
  palette: {
    mode: 'light',
    primary: {
      main: 'rgb(50, 129, 246)',
    },
    background: {
      default: 'rgb(255,255,255)',
      paper: 'rgb(255,255,255)',
    },
  },
  typography: {
    button: {
      textTransform: 'none',
    },
  },
});

export default theme;
