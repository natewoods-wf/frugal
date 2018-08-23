package music

import (
	"fmt"

	"git.apache.org/thrift.git/lib/go/thrift"
)

// possible alternatives to dealing with maps

// Skeleton is a thin wireframe for a document.
//
// It allows us to scroll through 100k page documents without needing to load all content.
type Skeleton struct {
	GUID           string                 `thrift:"GUID,1" db:"GUID" json:"GUID"`
	GridStructures map[string]interface{} `thrift:"GridStructures,2" db:"GridStructures" json:"GridStructures"`
	Pages          []interface{}          `thrift:"Pages,3" db:"Pages" json:"Pages"`
}

func NewSkeleton() *Skeleton {
	return &Skeleton{}
}

func (p *Skeleton) Read(iprot thrift.TProtocol) error { return nil }

type SkeletonRoot struct {
	OutlineMaxStamp         int64                `thrift:"OutlineMaxStamp,1" db:"OutlineMaxStamp" json:"OutlineMaxStamp"`
	SectionToCustomSkeleton map[string]*Skeleton `thrift:"SectionToCustomSkeleton,2" db:"SectionToCustomSkeleton" json:"SectionToCustomSkeleton"`
	MainSkeleton            *Skeleton            `thrift:"MainSkeleton,3" db:"MainSkeleton" json:"MainSkeleton"`
}

func (p *SkeletonRoot) ReadField2(iprot thrift.TProtocol) error {
	_, _, size, err := iprot.ReadMapBegin()
	if err != nil {
		return thrift.PrependError("error reading map begin: ", err)
	}
	p.SectionToCustomSkeleton = make(map[string]*Skeleton, size)
	for i := 0; i < size; i++ {
		var elem11 string
		if v, err := iprot.ReadString(); err != nil {
			return thrift.PrependError("error reading field 0: ", err)
		} else {
			elem11 = v
		}
		elem12 := NewSkeleton()
		if err := elem12.Read(iprot); err != nil {
			return thrift.PrependError(fmt.Sprintf("%T error reading struct: ", elem12), err)
		}
		(p.SectionToCustomSkeleton)[elem11] = elem12
	}
	if err := iprot.ReadMapEnd(); err != nil {
		return thrift.PrependError("error reading map end: ", err)
	}
	return nil
}

func ReadMap(iprot thrift.TProtocol, constructor func(size int), readPair func() error) error {
	_, _, size, err := iprot.ReadMapBegin()
	if err != nil {
		return err
	}
	constructor(size)
	for i := 0; i < size; i++ {
		if err := readPair(); err != nil {
			return err
		}
	}
	return iprot.ReadMapEnd()
}

func (p *SkeletonRoot) ReadField2Option2(iprot thrift.TProtocol) error {
	return ReadMap(iprot, func(size int) {
		p.SectionToCustomSkeleton = make(map[string]*Skeleton, size)
	}, func() (err error) {
		var (
			key   string
			value = NewSkeleton()
		)
		key, err = iprot.ReadString()
		if err == nil {
			err = value.Read(iprot)
		}
		if err == nil {
			p.SectionToCustomSkeleton[key] = value
		}
		return err
	})
}

func ReadMap2(iprot thrift.TProtocol,
	construct func(size int),
	readKey, readValue func() (interface{}, error),
	assign func(key, value interface{})) error {
	_, _, size, err := iprot.ReadMapBegin()
	if err != nil {
		return err
	}
	construct(size)
	for i := 0; i < size; i++ {
		if key, err := readKey(); err != nil {
			return err
		} else if val, err := readValue(); err != nil {
			return err
		} else {
			assign(key, val)
		}
	}
	return iprot.ReadMapEnd()
}

func (p *SkeletonRoot) ReadField2Option3(iprot thrift.TProtocol) error {
	return ReadMap2(iprot, func(size int) {
		p.SectionToCustomSkeleton = make(map[string]*Skeleton, size)
	}, func() (interface{}, error) {
		return iprot.ReadString()
	}, func() (interface{}, error) {
		val := NewSkeleton()
		return val, val.Read(iprot)
	}, func(key, value interface{}) {
		p.SectionToCustomSkeleton[key.(string)] = value.(*Skeleton)
	})
}
