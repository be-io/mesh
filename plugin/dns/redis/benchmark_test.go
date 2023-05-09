/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package redis

import (
	"github.com/coredns/coredns/plugin/test"

	"github.com/miekg/dns"
)

var zone string = "example.com."

var benchmarkEntries = [][]string{
	{"@",
		"{\"a\":[{\"ttl\":300, \"ip\":\"2.2.2.2\"}]}",
	},
	{"x",
		"{\"a\":[{\"ttl\":300, \"ip\":\"3.3.3.3\"}]}",
	},
	{"y",
		"{\"a\":[{\"ttl\":300, \"ip\":\"4.4.4.4\"}]}",
	},
	{"z",
		"{\"a\":[{\"ttl\":300, \"ip\":\"5.5.5.5\"}]}",
	},
}

var testCasesHit = []test.Case{
	{
		Qname: "example.com.", Qtype: dns.TypeA,
		Answer: []dns.RR{
			test.A("example.com. 300 IN A 2.2.2.2"),
		},
	},
	{
		Qname: "x.example.com.", Qtype: dns.TypeA,
		Answer: []dns.RR{
			test.A("x.example.com. 300 IN A 3.3.3.3"),
		},
	},
	{
		Qname: "y.example.com.", Qtype: dns.TypeA,
		Answer: []dns.RR{
			test.A("y.example.com. 300 IN A 4.4.4.4"),
		},
	},
	{
		Qname: "z.example.com.", Qtype: dns.TypeA,
		Answer: []dns.RR{
			test.A("z.example.com. 300 IN A 5.5.5.5"),
		},
	},
}

var testCasesMiss = []test.Case{
	{
		Qname: "q.example.com.", Qtype: dns.TypeA,
		Rcode: dns.RcodeNameError,
	},
	{
		Qname: "w.example.com.", Qtype: dns.TypeA,
		Rcode: dns.RcodeNameError,
	},
	{
		Qname: "e.example.com.", Qtype: dns.TypeA,
		Rcode: dns.RcodeNameError,
	},
	{
		Qname: "r.example.com.", Qtype: dns.TypeA,
		Rcode: dns.RcodeNameError,
	},
}
