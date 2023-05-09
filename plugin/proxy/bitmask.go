/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package proxy

var bitMasks BitMask = map[[4]byte]string{
	{0x00, 0x00, 0x00, 0x00}: "empty",
	{0x00, 0x00, 0x00, 0x01}: "close",
}

type BitMask map[[4]byte]string

func (that BitMask) index(buff []byte) string {
	if nil == buff {
		return ""
	}
	if len(buff) < 4 {
		return ""
	}
	return that[[4]byte{buff[0], buff[1], buff[2], buff[3]}]
}

func (that BitMask) IsEmpty(buff []byte) bool {
	return "empty" == that.index(buff)
}

func (that BitMask) IsClose(buff []byte) bool {
	return "close" == that.index(buff)
}

func (that BitMask) UnPacket(buff []byte) []byte {
	if len(buff) < 4 {
		return buff
	}
	return buff[4:]
}
