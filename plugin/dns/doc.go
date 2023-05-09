/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package dns

// Doc
//
// A记录： 将域名指向一个IPv4地址（例如：100.100.100.100），需要增加A记录
// CNAME记录： 如果将域名指向一个域名，实现与被指向域名相同的访问效果，需要增加CNAME记录。这个域名一般是主机服务商提供的一个域名
// MX记录： 建立电子邮箱服务，将指向邮件服务器地址，需要设置MX记录。建立邮箱时，一般会根据邮箱服务商提供的MX记录填写此记录
// NS记录： 域名解析服务器记录，如果要将子域名指定某个域名服务器来解析，需要设置NS记录
// TXT记录： 可任意填写，可为空。一般做一些验证记录时会使用此项，如：做SPF（反垃圾邮件）记录
// AAAA记录： 将主机名（或域名）指向一个IPv6地址（例如：ff03:0:0:0:0:0:0:c1），需要添加AAAA记录
// SRV记录： 添加服务记录服务器服务记录时会添加此项，SRV记录了哪台计算机提供了哪个服务。格式为：服务的名字.协议的类型（例如：_example-server._tcp）。
// SOA记录： SOA叫做起始授权机构记录，NS用于标识多台域名解析服务器，SOA记录用于在众多NS记录中那一台是主服务器
// PTR记录： PTR记录是A记录的逆向记录，又称做IP反查记录或指针记录，负责将IP反向解析为域名
// 显性URL转发记录： 将域名指向一个http(s)协议地址，访问域名时，自动跳转至目标地址。例如：将ducesoft.net显性转发到x.ducesoft.net后，访问ducesoft.net时，地址栏显示的地址为：x.ducesoft.net。
// 隐性UR转发记录L： 将域名指向一个http(s)协议地址，访问域名时，自动跳转至目标地址，隐性转发会隐藏真实的目标地址。例如：将ducesoft.net显性转发到x.ducesoft.net后，访问ducesoft.net时，地址栏显示的地址仍然是：ducesoft.net。
type Doc struct {
}
