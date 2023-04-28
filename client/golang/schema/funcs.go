/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package schema

import (
	"fmt"
	"strings"
	"unicode"
)

func Underline(name string) string {
	var chars []uint8
	for index := 0; index < len(name); index++ {
		if unicode.IsUpper(rune(name[index])) {
			if 0 != index && (index < 1 || unicode.IsLower(rune(name[index-1]))) {
				chars = append(chars, '_')
			}
			chars = append(chars, byte(unicode.ToLower(rune(name[index]))))
			continue
		}
		chars = append(chars, name[index])
	}
	return string(chars)
}

func Hump(name string) string {
	var chars []uint8
	for index := 0; index < len(name); index++ {
		if name[index] == '_' {
			continue
		}
		if index > 0 && name[index-1] == '_' {
			chars = append(chars, byte(unicode.ToUpper(rune(name[index]))))
			continue
		}
		chars = append(chars, name[index])
	}
	return string(chars)
}

func IsPublic(modifiers int) bool    { return (modifiers & 0x00000001) != 0 }
func IsPrivate(modifiers int) bool   { return (modifiers & 0x00000002) != 0 }
func IsProtected(modifiers int) bool { return (modifiers & 0x00000004) != 0 }
func IsStatic(modifiers int) bool    { return (modifiers & 0x00000008) != 0 }
func IsFinal(modifiers int) bool     { return (modifiers & 0x00000010) != 0 }
func IsVolatile(modifiers int) bool  { return (modifiers & 0x00000040) != 0 }
func IsTransient(modifiers int) bool { return (modifiers & 0x00000080) != 0 }
func IsAbstract(modifiers int) bool  { return (modifiers & 0x00000400) != 0 }
func IsInterface(modifiers int) bool { return (modifiers & 0x00000200) != 0 }

func IsSetMethod(methodName string, parameters []*Parameter) bool {
	if parameters == nil || len(parameters) != 1 {
		return false
	}
	return strings.HasPrefix(strings.ToLower(methodName), "set")
}

func GetSetField(methodName string, parameters []*Parameter) string {
	fieldName := parameters[0].Name
	return "{ this." + strings.ToLower(methodName[3:4]) + methodName[4:] + "=" + fieldName + ";}"
}

func IsGetMethod(methodName string, parameters []*Parameter) bool {
	if parameters != nil && len(parameters) > 0 {
		return false
	}
	lowerMethodName := strings.ToLower(methodName)
	return strings.HasPrefix(lowerMethodName, "get") || strings.HasPrefix(lowerMethodName, "is")
}

func GetGetField(methodName string) string {
	if strings.HasPrefix(strings.ToLower(methodName), "get") {
		return "{ return this." + strings.ToLower(methodName[3:4]) + methodName[4:] + "; }\n"
	} else if strings.HasPrefix(strings.ToLower(methodName), "is") {
		return "{ return this." + strings.ToLower(methodName[2:3]) + methodName[3:] + "; }\n"
	}
	return methodName
}

func GetMacroMetadata(name string, macro map[string]string) []string {
	return strings.Split(strings.TrimSuffix(strings.TrimPrefix(macro[name], "["), "]"), ",")
}

func GenericTree(signature string) *Tree {
	if !strings.Contains(signature, "<") || "" == signature {
		return &Tree{Root: signature}
	}
	var childrens []*Tree
	vs := strings.SplitN(signature, "<", 2)
	cs := strings.TrimSuffix(vs[1], ">")
	var chars []uint8
	var l = 0
	var r = 0
	for index := 0; index < len(cs); index++ {
		if '<' == cs[index] {
			l++
		}
		if '>' == cs[index] {
			r++
		}
		if ',' != cs[index] || len(chars) > 0 {
			chars = append(chars, cs[index])
		}
		if l == r && 0 != l {
			cps := strings.SplitN(string(chars), "<", 2)
			cpp := strings.Split(cps[0], ",")
			if len(cpp) > 1 {
				for idx := 0; idx < len(cpp)-1; idx++ {
					childrens = append(childrens, &Tree{Root: cpp[idx]})
				}
				childrens = append(childrens, GenericTree(fmt.Sprintf("%s<%s", cpp[len(cpp)-1], cps[1])))
			} else {
				childrens = append(childrens, GenericTree(string(chars)))
			}
			l = 0
			r = 0
			chars = nil
			continue
		}
		if index == len(cs)-1 && len(chars) > 0 {
			for _, expr := range strings.Split(string(chars), ",") {
				childrens = append(childrens, &Tree{Root: expr})
			}
		}
	}
	return &Tree{Root: vs[0], Childrens: childrens}
}

func FormatTree(tree *Tree, format func(root string, children []string) string) string {
	var childrens []string
	for _, ch := range tree.Childrens {
		childrens = append(childrens, FormatTree(ch, format))
	}
	return format(tree.Root, childrens)
}

func Wrap() string {
	return "\n"
}

func Ident(p string, n int, s string) string {
	return fmt.Sprintf("%s%s%s", p, strings.Repeat("    ", n), s)
}

var Funcs = map[string]any{
	"IsPublic":         IsPublic,
	"IsPrivate":        IsPrivate,
	"IsProtected":      IsProtected,
	"IsStatic":         IsStatic,
	"IsFinal":          IsFinal,
	"IsVolatile":       IsVolatile,
	"IsTransient":      IsTransient,
	"IsAbstract":       IsAbstract,
	"IsSetMethod":      IsSetMethod,
	"GetSetField":      GetSetField,
	"IsGetMethod":      IsGetMethod,
	"GetGetField":      GetGetField,
	"GetMacroMetadata": GetMacroMetadata,
	"Underline":        Underline,
	"Hump":             Hump,
	"GenericTree":      GenericTree,
	"FormatTree":       FormatTree,
	"Ident":            Ident,
	"Wrap":             Wrap,
}
