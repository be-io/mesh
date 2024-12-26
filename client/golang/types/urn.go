/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package types

import (
	"bytes"
	"context"
	"fmt"
	"github.com/opendatav/mesh/client/golang/log"
	"strings"
)

const (
	CN             = "trustbe"
	MeshDomain     = CN + ".cn"
	LocalNodeId    = "LX0000000000000"
	LocalInstId    = "JG0000000000000000"
	MeshKeywords   = "([-a-zA-Z\\d]+)"
	MeshUrnPattern = "%s\\.%s\\.((?i)JG[-a-zA-Z\\d]{2}%s[-a-zA-Z\\d]{8}|(?i)LX[-a-zA-Z\\d]{6}%s[-a-zA-Z\\d]{1}%s)\\.%s\\.(net|cn|com)"
)

type URN struct {
	Domain string
	NodeId string // Maybe nodeId or instId
	Flag   *URNFlag
	Name   string
}

type URNFlag struct {
	V       string // 2     00
	Proto   string // 2     00
	Codec   string // 2     00
	Version string // 6     000000
	Zone    string // 2     00
	Cluster string // 2     00
	Cell    string // 2     00
	Group   string // 2     00
	Address string // 12    000000000000
	Port    string // 5     00080   Max port 65535
}

func (that *URNFlag) String() string {
	return fmt.Sprintf("%s%s%s%s%s%s%s%s%s%s",
		that.Padding(that.V, 2),
		that.Padding(that.Proto, 2),
		that.Padding(that.Codec, 2),
		that.PaddingChain(that.Version, 2, 3),
		that.Padding(that.Zone, 2),
		that.Padding(that.Cluster, 2),
		that.Padding(that.Cell, 2),
		that.Padding(that.Group, 2),
		that.PaddingChain(that.Address, 3, 4),
		that.Padding(that.Port, 5))
}

func (that *URNFlag) Padding(v string, length int) string {
	value := strings.ReplaceAll(v, ".", "")
	if len(value) == length {
		return value
	}
	if len(value) < length {
		return strings.Repeat("0", length-len(v)) + value
	}
	return v[0:length]
}

func (that *URNFlag) PaddingChain(v string, length int, size int) string {
	chain := func(frags []string) string {
		var chars []string
		for _, frag := range frags {
			chars = append(chars, that.Padding(frag, length))
		}
		return strings.Join(chars, "")
	}
	frags := strings.Split(v, ".")
	min := len(frags)
	if len(frags) == size {
		return chain(frags)
	}
	if len(frags) < size {
		for index := 0; index < size-min; index++ {
			frags = append(frags, "")
		}
		return chain(frags)
	}
	return chain(frags[0:size])
}

func (that *URN) String() string {
	var urn []string
	var names = AsArray(that.Name)
	for x, y := 0, len(names)-1; x < y; x, y = x+1, y-1 {
		names[x], names[y] = names[y], names[x]
	}
	urn = append(append(urn, names...), that.Flag.String(), strings.ToLower(that.NodeId))
	if "" == that.Domain {
		urn = append(urn, strings.ToLower(MeshDomain))
	} else {
		urn = append(urn, strings.ToLower(that.Domain))
	}
	return strings.Join(urn, ".")
}

func (that *URN) IsInstitutional(rk string) bool {
	return strings.Index(strings.ToLower(rk), "jg") == 0
}

func (that *URN) MatchNode(ctx context.Context, nodeId string) bool {
	xIsInstitutional := that.IsInstitutional(that.NodeId)
	if xIsInstitutional {
		instId, err := FromInstID(that.NodeId)
		if nil != err {
			log.Error(ctx, "%s is not standard identity", that.String())
			return false
		}
		return instId.MatchNode(nodeId)
	}
	return strings.ToLower(that.NodeId) == strings.ToLower(nodeId)
}

func (that *URN) MatchInst(ctx context.Context, instId string) bool {
	xIsInstitutional := that.IsInstitutional(that.NodeId)
	institutionId, err := FromInstID(instId)
	if nil != err {
		log.Error(ctx, "%s is not standard identity", that.String())
		return false
	}
	if xIsInstitutional {
		return institutionId.Match(that.NodeId)
	}
	return institutionId.MatchNode(that.NodeId)
}

func (that *URN) MatchPDC(ctx context.Context, urn *URN) bool {
	yIsInstitutional := that.IsInstitutional(urn.NodeId)
	if yIsInstitutional {
		return that.MatchInst(ctx, urn.NodeId)
	}
	return that.MatchNode(ctx, urn.NodeId)
}

func (that *URN) Match(ctx context.Context, urn *URN) bool {
	if !that.MatchPDC(ctx, urn) {
		return false
	}
	return that.Name == urn.Name
}

func (that *URN) ID(ctx context.Context) string {
	return NodeSEQ(that.NodeId)
}

func AsArray(text string) []string {
	var pairs = strings.Split(text, ".")
	var buff = &bytes.Buffer{}
	var names []string
	for _, pair := range pairs {
		if strings.Index(pair, "${") == 0 {
			buff.WriteString(pair)
			buff.WriteRune('.')
			continue
		}
		if strings.Index(pair, "}") == len(pair)-1 {
			buff.WriteString(pair)
			names = append(names, buff.String())
			buff.Reset()
			continue
		}
		names = append(names, pair)
	}
	return names
}

func FromURN(ctx context.Context, urn string) *URN {
	name := &URN{Flag: &URNFlag{}}
	if "" == urn {
		log.Debug(ctx, "Unresolved urn %s", urn)
		return name
	}
	names := AsArray(urn)
	if len(names) < 5 {
		log.Debug(ctx, "Unresolved urn %s", urn)
		name.Name = urn
		return name
	}
	for x, y := 0, len(names)-1; x < y; x, y = x+1, y-1 {
		names[x], names[y] = names[y], names[x]
	}
	name.Domain = fmt.Sprintf("%s.%s", names[1], names[0])
	name.NodeId = names[2]
	name.Flag = FromURNFlag(ctx, names[3])
	name.Name = strings.Join(names[4:], ".")
	return name
}

func FromURNFlag(ctx context.Context, value string) *URNFlag {
	if "" == value {
		log.Error(ctx, "Unresolved Flag %s", value)
		return new(URNFlag)
	}
	flag := new(URNFlag)
	flag.V = Substring(value, 0, 2)
	flag.Proto = Substring(value, 2, 4)
	flag.Codec = Substring(value, 4, 6)
	flag.Version = Reduce(Substring(value, 6, 12), 2)
	flag.Zone = Substring(value, 12, 14)
	flag.Cluster = Substring(value, 14, 16)
	flag.Cell = Substring(value, 16, 18)
	flag.Group = Substring(value, 18, 20)
	flag.Address = Reduce(Substring(value, 20, 32), 3)
	flag.Port = Reduce(Substring(value, 32, 37), 5)
	return flag
}

func Reduce(value string, length int) string {
	var bu bytes.Buffer
	for index := 0; index < len(value); index = index + length {
		hasNoneZero := true
		for offset := index; offset < index+length; offset++ {
			if offset >= len(value) {
				break
			}
			if value[offset] != '0' {
				hasNoneZero = false
			}
			if !hasNoneZero {
				bu.WriteByte(value[offset])
			}
		}
		if hasNoneZero {
			bu.WriteByte('0')
		}
		bu.WriteByte('.')
	}
	if bu.Len() < 1 {
		return strings.Repeat("0", length)
	}
	return bu.String()[0 : bu.Len()-1]
}

func Substring(v string, start int, stop int) string {
	var chars []rune
	for index, char := range v {
		if index >= start && index < stop {
			chars = append(chars, char)
		}
	}
	return string(chars)
}

func LocURN(ctx context.Context, name string) string {
	return AnyURN(ctx, name, LocalNodeId).String()
}

func AnyURN(ctx context.Context, name string, nodeId string) *URN {
	return &URN{
		Domain: MeshDomain,
		NodeId: nodeId,
		Flag:   FromURNFlag(ctx, "0001000000000000000000000000000000000"),
		Name:   name,
	}
}

func URNMatcher(name string, seq string, domain string, nodeIds ...string) string {
	return URNFlagMatcher(name, MeshKeywords, seq, domain, nodeIds...)
}

func URNFlagMatcher(name string, flag string, seq string, domain string, nodeIds ...string) string {
	return fmt.Sprintf(MeshUrnPattern, fmt.Sprintf("%s\\.%s", MeshKeywords, name), flag, seq, seq, fmt.Sprintf("|(?i)%s", strings.Join(nodeIds, "|(?i)")), domain)
}

func URNPatternMatcher(pattern string, seq string, domain string, nodeIds ...string) string {
	return fmt.Sprintf(MeshUrnPattern, pattern, MeshKeywords, seq, seq, fmt.Sprintf("|(?i)%s", strings.Join(nodeIds, "|(?i)")), domain)
}

// MatchURNDomain
// Warning: dont change if not clearly.
func MatchURNDomain(domain string) bool {
	names := strings.Split(domain, ".")
	if len(names) < 3 {
		return false
	}
	id := names[len(names)-3]
	return len(id) == 15 || len(id) == 18
}
