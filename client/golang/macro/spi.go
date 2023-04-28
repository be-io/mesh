/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package macro

import "fmt"

func init() {

}

func init() {
	var _ fmt.Stringer = new(Att)
	var _ SPI = new(SPIAnnotation)
}

type SPI interface {
	// Att
	// SPI provider attribute
	Att() *Att
}

type Att struct {
	// Name spi name
	Name string
	// Pattern to customized
	Pattern string
	// Priority
	Priority int
	// Prototype
	Prototype bool
	// Alias
	Alias []string
	// Constructor
	Constructor func() SPI
}

func (that *Att) String() string {
	return fmt.Sprintf("%p", that)
}

type SPIAnnotation struct {
	Attribute *Att
}

func (that *SPIAnnotation) Att() *Att {
	return that.Attribute
}

type Initializer interface {

	// Init the spi provider
	Init() error
}
