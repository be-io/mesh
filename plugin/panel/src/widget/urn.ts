// /**
//  * Like:
//  *
//  * create.tenant.omega . 0001 . json . http2 . lx000001 . ducesoft.net
//  * -------------------   ----   ----   -----   --------   -----------
//  * name            flag   codec  proto     node       domain
//  */
//
// class URNFlag {
//     public v       :string // 2     00
//     public  proto   :string // 2     00
//     public  codec   :string // 2     00
//     public  version :string // 6     000000
//     public  zone    :string // 2     00
//     public  cluster :string // 2     00
//     public  cell    :string // 2     00
//     public  group   :string // 2     00
//     public  address :string // 12    000000000000
//     public  port    :string // 5     00080   Max port 65535
//
//     String() :string {
//     return fmt.Sprintf("%s%s%s%s%s%s%s%s%s%s",
//     this.Padding(this.v, 2),
//         this.Padding(this.proto, 2),
//         this.Padding(this.codec, 2),
//         this.PaddingChain(this.version, 2, 3),
//         this.Padding(this.zone, 2),
//         this.Padding(this.cluster, 2),
//         this.Padding(this.cell, 2),
//         this.Padding(this.group, 2),
//         this.PaddingChain(this.address, 3, 4),
//         this.Padding(this.port, 5))
// }
//
//     Padding(v: string, length :number) :string {
//     const value = strings.ReplaceAll(v, ".", "")
//     if (value.length == length) {
//     return value
// }
// if (value.length < length) {
//     return strings.Repeat("0", length-v.length) + value
// }
// return v.substring(0,length)
// }
//
//     PaddingChain(v :string, length :number, size: number) :string {
//     chain := func(frags []string) string {
//     var chars []string
//     for _, frag := range frags {
//     chars = append(chars, that.Padding(frag, length))
// }
// return strings.Join(chars, "")
// }
// frags := strings.Split(v, ".")
// min := len(frags)
// if len(frags) == size {
//     return chain(frags)
// }
// if len(frags) < size {
//     for index := 0; index < size-min; index++ {
//         frags = append(frags, "")
//     }
//     return chain(frags)
// }
// return chain(frags[0:size])
// }
// }
// class URN {
//  public domain:string
//     public node_id:string // Maybe nodeId or instId
//     public flag:URNFlag
//     public name:string
//
//     AsArray(text :string) :string[] {
//     var pairs = strings.Split(text, ".")
//     var buff = &bytes.Buffer{}
// var names []string
// for _, pair := range pairs {
//     if strings.Index(pair, "${") == 0 {
//         buff.WriteString(pair)
//         buff.WriteRune('.')
//         continue
//     }
//     if strings.Index(pair, "}") == len(pair)-1 {
//         buff.WriteString(pair)
//         names = append(names, buff.String())
//         buff.Reset()
//         continue
//     }
//     names = append(names, pair)
// }
// return names
// }
//
//     String() :string {
//     var urn []string
//     var names = that.AsArray(that.Name)
//     for x, y := 0, len(names)-1; x < y; x, y = x+1, y-1 {
//     names[x], names[y] = names[y], names[x]
// }
// urn = append(append(urn, names...), that.Flag.String(), strings.ToLower(that.NodeId))
// if "" == that.Domain {
//     urn = append(urn, strings.ToLower(MeshDomain))
// } else {
//     urn = append(urn, strings.ToLower(that.Domain))
// }
// return strings.Join(urn, ".")
// }
//
// }
//
//
// function FromURN(urn: string):URN {
//     name := &URN{Flag: &URNFlag{}}
//     if "" == urn {
//         log.Error(ctx, "Unresolved urn %s", urn)
//         return name
//     }
//     names := name.AsArray(urn)
//     if len(names) < 5 {
//         log.Error(ctx, "Unresolved urn %s", urn)
//         name.Name = urn
//         return name
//     }
//     for x, y := 0, len(names)-1; x < y; x, y = x+1, y-1 {
//         names[x], names[y] = names[y], names[x]
//     }
//     name.Domain = fmt.Sprintf("%s.%s", names[1], names[0])
//     name.NodeId = names[2]
//     name.Flag = FromURNFlag(ctx, names[3])
//     name.Name = strings.Join(names[4:], ".")
//     return name
// }
//
// function FromURNFlag( value: string) :URNFlag {
//     if "" == value {
//         log.Error(ctx, "Unresolved Flag %s", value)
//         return new(URNFlag)
//     }
//     flag := new(URNFlag)
//     flag.V = Substring(value, 0, 2)
//     flag.Proto = Substring(value, 2, 4)
//     flag.Codec = Substring(value, 4, 6)
//     flag.Version = Reduce(Substring(value, 6, 12), 2)
//     flag.Zone = Substring(value, 12, 14)
//     flag.Cluster = Substring(value, 14, 16)
//     flag.Cell = Substring(value, 16, 18)
//     flag.Group = Substring(value, 18, 20)
//     flag.Address = Reduce(Substring(value, 20, 32), 3)
//     flag.Port = Reduce(Substring(value, 32, 37), 5)
//     return flag
// }
//
// function Reduce(value: string, length :int): string {
//     var bu bytes.Buffer
//     for index := 0; index < len(value); index = index + length {
//         hasNoneZero := true
//         for offset := index; offset < index+length; offset++ {
//             if offset >= len(value) {
//                 break
//             }
//             if value[offset] != '0' {
//                 hasNoneZero = false
//             }
//             if !hasNoneZero {
//                 bu.WriteByte(value[offset])
//             }
//         }
//         if hasNoneZero {
//             bu.WriteByte('0')
//         }
//         bu.WriteByte('.')
//     }
//     if bu.Len() < 1 {
//         return strings.Repeat("0", length)
//     }
//     return bu.String()[0 : bu.Len()-1]
// }
//
// function Substring(v: string, start: number, stop: number) :string {
//     const chars:string[]= []
//     v.
//     for index, char := range v {
//         if index >= start && index < stop {
//             chars = append(chars, char)
//         }
//     }
//     return string(chars)
// }

export default {}