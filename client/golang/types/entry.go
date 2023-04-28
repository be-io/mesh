package types

type Entry struct {
	Key      string  `index:"0" json:"key" yaml:"key" xml:"key" comment:""`
	Value    *Entity `index:"5" json:"value" yaml:"value" xml:"value" comment:""`
	UpdateAt Time    `index:"10" json:"update_at" yaml:"update_at" xml:"update_at" comment:""`
}
