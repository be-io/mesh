/*
 * Copyright (c) 2000, 2099, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */
//
// import PropTypes from 'prop-types';
// import {forwardRef} from 'react';
// import {useTheme} from '@mui/material/styles';
// import {Card, CardContent, CardHeader, Divider, Typography} from '@mui/material';
// import Highlighter from './third-party/Highlighter';
//
// const headerSX = {
//     p: 2.5,
//     '& .MuiCardHeader-action': {m: '0px auto', alignSelf: 'center'}
// };
//
// const MainCard = forwardRef(
//     (
//         {
//             border = true,
//             boxShadow,
//             children,
//             content = true,
//             contentSX = {},
//             darkTitle,
//             elevation,
//             secondary,
//             shadow,
//             sx = {},
//             title,
//             codeHighlight,
//             ...others
//         },
//         ref
//     ) => {
//         const theme = useTheme();
//         boxShadow = theme.palette.mode === 'dark' ? boxShadow || true : boxShadow;
//
//         return (
//             <Card
//                 elevation={elevation || 0}
//                 ref={ref}
//                 {...others}
//                 sx={{
//                     border: border ? '1px solid' : 'none',
//                     borderRadius: 2,
//                     borderColor: theme.palette.mode === 'dark' ? theme.palette.divider : theme.palette.grey.A800,
//                     boxShadow: boxShadow && (!border || theme.palette.mode === 'dark') ? shadow || theme.customShadows.z1 : 'inherit',
//                     ':hover': {
//                         boxShadow: boxShadow ? shadow || theme.customShadows.z1 : 'inherit'
//                     },
//                     '& pre': {
//                         m: 0,
//                         p: '16px !important',
//                         fontFamily: theme.typography.fontFamily,
//                         fontSize: '0.75rem'
//                     },
//                     ...sx
//                 }}
//             >
//                 {/* card header and action */}
//                 {!darkTitle && title && (
//                     <CardHeader sx={headerSX} titleTypographyProps={{variant: 'subtitle1'}} title={title}
//                                 action={secondary}/>
//                 )}
//                 {darkTitle && title && <CardHeader sx={headerSX} title={<Typography variant="h3">{title}</Typography>}
//                                                    action={secondary}/>}
//
//                 {/* card content */}
//                 {content && <CardContent sx={contentSX}>{children}</CardContent>}
//                 {!content && children}
//
//                 {/* card footer - clipboard & highlighter  */}
//                 {codeHighlight && (
//                     <>
//                         <Divider sx={{borderStyle: 'dashed'}}/>
//                         <Highlighter codeHighlight={codeHighlight} main>
//                             {children}
//                         </Highlighter>
//                     </>
//                 )}
//             </Card>
//         );
//     }
// );
//
// MainCard.propTypes = {
//     border: PropTypes.bool,
//     boxShadow: PropTypes.bool,
//     contentSX: PropTypes.object,
//     darkTitle: PropTypes.bool,
//     divider: PropTypes.bool,
//     elevation: PropTypes.number,
//     secondary: PropTypes.node,
//     shadow: PropTypes.string,
//     sx: PropTypes.object,
//     title: PropTypes.oneOfType([PropTypes.string, PropTypes.node]),
//     codeHighlight: PropTypes.bool,
//     content: PropTypes.bool,
//     children: PropTypes.node
// };
//
// export function MainCard;
