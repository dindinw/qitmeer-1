// Code generated by fastssz. DO NOT EDIT.
package qitmeer_p2p_v1

import (
	"fmt"

	ssz "github.com/ferranbt/fastssz"
)

var (
	errDivideInt           = fmt.Errorf("incorrect int divide")
	errListTooBig          = fmt.Errorf("incorrect list size, too big")
	errMarshalDynamicBytes = fmt.Errorf("incorrect dynamic bytes marshalling")
	errMarshalFixedBytes   = fmt.Errorf("incorrect fixed bytes marshalling")
	errMarshalList         = fmt.Errorf("incorrect vector list")
	errMarshalVector       = fmt.Errorf("incorrect vector marshalling")
	errOffset              = fmt.Errorf("incorrect offset")
	errSize                = fmt.Errorf("incorrect size")
)

// MarshalSSZ ssz marshals the GraphState object
func (g *GraphState) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, g.SizeSSZ())
	return g.MarshalSSZTo(buf[:0])
}

// MarshalSSZTo ssz marshals the GraphState object to a target array
func (g *GraphState) MarshalSSZTo(dst []byte) ([]byte, error) {
	var err error
	offset := int(20)

	// Field (0) 'Total'
	dst = ssz.MarshalUint32(dst, g.Total)

	// Field (1) 'Layer'
	dst = ssz.MarshalUint32(dst, g.Layer)

	// Field (2) 'MainHeight'
	dst = ssz.MarshalUint32(dst, g.MainHeight)

	// Field (3) 'MainOrder'
	dst = ssz.MarshalUint32(dst, g.MainOrder)

	// Offset (4) 'Tips'
	dst = ssz.WriteOffset(dst, offset)
	offset += len(g.Tips) * 32

	// Field (4) 'Tips'
	if len(g.Tips) > 100 {
		return nil, errMarshalList
	}
	for ii := 0; ii < len(g.Tips); ii++ {
		if dst, err = g.Tips[ii].MarshalSSZTo(dst); err != nil {
			return nil, err
		}
	}

	return dst, err
}

// UnmarshalSSZ ssz unmarshals the GraphState object
func (g *GraphState) UnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size < 20 {
		return errSize
	}

	tail := buf
	var o4 uint64

	// Field (0) 'Total'
	g.Total = ssz.UnmarshallUint32(buf[0:4])

	// Field (1) 'Layer'
	g.Layer = ssz.UnmarshallUint32(buf[4:8])

	// Field (2) 'MainHeight'
	g.MainHeight = ssz.UnmarshallUint32(buf[8:12])

	// Field (3) 'MainOrder'
	g.MainOrder = ssz.UnmarshallUint32(buf[12:16])

	// Offset (4) 'Tips'
	if o4 = ssz.ReadOffset(buf[16:20]); o4 > size {
		return errOffset
	}

	// Field (4) 'Tips'
	{
		buf = tail[o4:]
		num, ok := ssz.DivideInt(len(buf), 32)
		if !ok {
			return errDivideInt
		}
		if num > 100 {
			return errListTooBig
		}
		g.Tips = make([]*Hash, num)
		for ii := 0; ii < num; ii++ {
			if g.Tips[ii] == nil {
				g.Tips[ii] = new(Hash)
			}
			if err = g.Tips[ii].UnmarshalSSZ(buf[ii*32 : (ii+1)*32]); err != nil {
				return err
			}
		}
	}
	return err
}

// SizeSSZ returns the ssz encoded size in bytes for the GraphState object
func (g *GraphState) SizeSSZ() (size int) {
	size = 20

	// Field (4) 'Tips'
	size += len(g.Tips) * 32

	return
}

// MarshalSSZ ssz marshals the ChainState object
func (c *ChainState) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, c.SizeSSZ())
	return c.MarshalSSZTo(buf[:0])
}

// MarshalSSZTo ssz marshals the ChainState object to a target array
func (c *ChainState) MarshalSSZTo(dst []byte) ([]byte, error) {
	var err error
	offset := int(61)

	// Field (0) 'GenesisHash'
	if c.GenesisHash == nil {
		c.GenesisHash = new(Hash)
	}
	if dst, err = c.GenesisHash.MarshalSSZTo(dst); err != nil {
		return nil, err
	}

	// Field (1) 'ProtocolVersion'
	dst = ssz.MarshalUint32(dst, c.ProtocolVersion)

	// Field (2) 'Timestamp'
	dst = ssz.MarshalUint64(dst, c.Timestamp)

	// Field (3) 'Services'
	dst = ssz.MarshalUint64(dst, c.Services)

	// Field (4) 'DisableRelayTx'
	dst = ssz.MarshalBool(dst, c.DisableRelayTx)

	// Offset (5) 'GraphState'
	dst = ssz.WriteOffset(dst, offset)
	if c.GraphState == nil {
		c.GraphState = new(GraphState)
	}
	offset += c.GraphState.SizeSSZ()

	// Offset (6) 'UserAgent'
	dst = ssz.WriteOffset(dst, offset)
	offset += len(c.UserAgent)

	// Field (5) 'GraphState'
	if dst, err = c.GraphState.MarshalSSZTo(dst); err != nil {
		return nil, err
	}

	// Field (6) 'UserAgent'
	if len(c.UserAgent) > 256 {
		return nil, errMarshalDynamicBytes
	}
	dst = append(dst, c.UserAgent...)

	return dst, err
}

// UnmarshalSSZ ssz unmarshals the ChainState object
func (c *ChainState) UnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size < 61 {
		return errSize
	}

	tail := buf
	var o5, o6 uint64

	// Field (0) 'GenesisHash'
	if c.GenesisHash == nil {
		c.GenesisHash = new(Hash)
	}
	if err = c.GenesisHash.UnmarshalSSZ(buf[0:32]); err != nil {
		return err
	}

	// Field (1) 'ProtocolVersion'
	c.ProtocolVersion = ssz.UnmarshallUint32(buf[32:36])

	// Field (2) 'Timestamp'
	c.Timestamp = ssz.UnmarshallUint64(buf[36:44])

	// Field (3) 'Services'
	c.Services = ssz.UnmarshallUint64(buf[44:52])

	// Field (4) 'DisableRelayTx'
	c.DisableRelayTx = ssz.UnmarshalBool(buf[52:53])

	// Offset (5) 'GraphState'
	if o5 = ssz.ReadOffset(buf[53:57]); o5 > size {
		return errOffset
	}

	// Offset (6) 'UserAgent'
	if o6 = ssz.ReadOffset(buf[57:61]); o6 > size || o5 > o6 {
		return errOffset
	}

	// Field (5) 'GraphState'
	{
		buf = tail[o5:o6]
		if c.GraphState == nil {
			c.GraphState = new(GraphState)
		}
		if err = c.GraphState.UnmarshalSSZ(buf); err != nil {
			return err
		}
	}

	// Field (6) 'UserAgent'
	{
		buf = tail[o6:]
		c.UserAgent = append(c.UserAgent, buf...)
	}
	return err
}

// SizeSSZ returns the ssz encoded size in bytes for the ChainState object
func (c *ChainState) SizeSSZ() (size int) {
	size = 61

	// Field (5) 'GraphState'
	if c.GraphState == nil {
		c.GraphState = new(GraphState)
	}
	size += c.GraphState.SizeSSZ()

	// Field (6) 'UserAgent'
	size += len(c.UserAgent)

	return
}

// MarshalSSZ ssz marshals the BlockData object
func (b *BlockData) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, b.SizeSSZ())
	return b.MarshalSSZTo(buf[:0])
}

// MarshalSSZTo ssz marshals the BlockData object to a target array
func (b *BlockData) MarshalSSZTo(dst []byte) ([]byte, error) {
	var err error
	offset := int(8)

	// Field (0) 'DagID'
	dst = ssz.MarshalUint32(dst, b.DagID)

	// Offset (1) 'BlockBytes'
	dst = ssz.WriteOffset(dst, offset)
	offset += len(b.BlockBytes)

	// Field (1) 'BlockBytes'
	if len(b.BlockBytes) > 1048576 {
		return nil, errMarshalDynamicBytes
	}
	dst = append(dst, b.BlockBytes...)

	return dst, err
}

// UnmarshalSSZ ssz unmarshals the BlockData object
func (b *BlockData) UnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size < 8 {
		return errSize
	}

	tail := buf
	var o1 uint64

	// Field (0) 'DagID'
	b.DagID = ssz.UnmarshallUint32(buf[0:4])

	// Offset (1) 'BlockBytes'
	if o1 = ssz.ReadOffset(buf[4:8]); o1 > size {
		return errOffset
	}

	// Field (1) 'BlockBytes'
	{
		buf = tail[o1:]
		b.BlockBytes = append(b.BlockBytes, buf...)
	}
	return err
}

// SizeSSZ returns the ssz encoded size in bytes for the BlockData object
func (b *BlockData) SizeSSZ() (size int) {
	size = 8

	// Field (1) 'BlockBytes'
	size += len(b.BlockBytes)

	return
}

// MarshalSSZ ssz marshals the SyncDAG object
func (s *SyncDAG) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, s.SizeSSZ())
	return s.MarshalSSZTo(buf[:0])
}

// MarshalSSZTo ssz marshals the SyncDAG object to a target array
func (s *SyncDAG) MarshalSSZTo(dst []byte) ([]byte, error) {
	var err error
	offset := int(8)

	// Offset (0) 'MainLocator'
	dst = ssz.WriteOffset(dst, offset)
	offset += len(s.MainLocator) * 32

	// Offset (1) 'GraphState'
	dst = ssz.WriteOffset(dst, offset)
	if s.GraphState == nil {
		s.GraphState = new(GraphState)
	}
	offset += s.GraphState.SizeSSZ()

	// Field (0) 'MainLocator'
	if len(s.MainLocator) > 32 {
		return nil, errMarshalList
	}
	for ii := 0; ii < len(s.MainLocator); ii++ {
		if dst, err = s.MainLocator[ii].MarshalSSZTo(dst); err != nil {
			return nil, err
		}
	}

	// Field (1) 'GraphState'
	if dst, err = s.GraphState.MarshalSSZTo(dst); err != nil {
		return nil, err
	}

	return dst, err
}

// UnmarshalSSZ ssz unmarshals the SyncDAG object
func (s *SyncDAG) UnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size < 8 {
		return errSize
	}

	tail := buf
	var o0, o1 uint64

	// Offset (0) 'MainLocator'
	if o0 = ssz.ReadOffset(buf[0:4]); o0 > size {
		return errOffset
	}

	// Offset (1) 'GraphState'
	if o1 = ssz.ReadOffset(buf[4:8]); o1 > size || o0 > o1 {
		return errOffset
	}

	// Field (0) 'MainLocator'
	{
		buf = tail[o0:o1]
		num, ok := ssz.DivideInt(len(buf), 32)
		if !ok {
			return errDivideInt
		}
		if num > 32 {
			return errListTooBig
		}
		s.MainLocator = make([]*Hash, num)
		for ii := 0; ii < num; ii++ {
			if s.MainLocator[ii] == nil {
				s.MainLocator[ii] = new(Hash)
			}
			if err = s.MainLocator[ii].UnmarshalSSZ(buf[ii*32 : (ii+1)*32]); err != nil {
				return err
			}
		}
	}

	// Field (1) 'GraphState'
	{
		buf = tail[o1:]
		if s.GraphState == nil {
			s.GraphState = new(GraphState)
		}
		if err = s.GraphState.UnmarshalSSZ(buf); err != nil {
			return err
		}
	}
	return err
}

// SizeSSZ returns the ssz encoded size in bytes for the SyncDAG object
func (s *SyncDAG) SizeSSZ() (size int) {
	size = 8

	// Field (0) 'MainLocator'
	size += len(s.MainLocator) * 32

	// Field (1) 'GraphState'
	if s.GraphState == nil {
		s.GraphState = new(GraphState)
	}
	size += s.GraphState.SizeSSZ()

	return
}

// MarshalSSZ ssz marshals the SubDAG object
func (s *SubDAG) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, s.SizeSSZ())
	return s.MarshalSSZTo(buf[:0])
}

// MarshalSSZTo ssz marshals the SubDAG object to a target array
func (s *SubDAG) MarshalSSZTo(dst []byte) ([]byte, error) {
	var err error
	offset := int(40)

	// Field (0) 'SyncPoint'
	if s.SyncPoint == nil {
		s.SyncPoint = new(Hash)
	}
	if dst, err = s.SyncPoint.MarshalSSZTo(dst); err != nil {
		return nil, err
	}

	// Offset (1) 'GraphState'
	dst = ssz.WriteOffset(dst, offset)
	if s.GraphState == nil {
		s.GraphState = new(GraphState)
	}
	offset += s.GraphState.SizeSSZ()

	// Offset (2) 'Blocks'
	dst = ssz.WriteOffset(dst, offset)
	for ii := 0; ii < len(s.Blocks); ii++ {
		offset += 4
		offset += s.Blocks[ii].SizeSSZ()
	}

	// Field (1) 'GraphState'
	if dst, err = s.GraphState.MarshalSSZTo(dst); err != nil {
		return nil, err
	}

	// Field (2) 'Blocks'
	if len(s.Blocks) > 500 {
		return nil, errMarshalList
	}
	{
		offset = 4 * len(s.Blocks)
		for ii := 0; ii < len(s.Blocks); ii++ {
			dst = ssz.WriteOffset(dst, offset)
			offset += s.Blocks[ii].SizeSSZ()
		}
	}
	for ii := 0; ii < len(s.Blocks); ii++ {
		if dst, err = s.Blocks[ii].MarshalSSZTo(dst); err != nil {
			return nil, err
		}
	}

	return dst, err
}

// UnmarshalSSZ ssz unmarshals the SubDAG object
func (s *SubDAG) UnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size < 40 {
		return errSize
	}

	tail := buf
	var o1, o2 uint64

	// Field (0) 'SyncPoint'
	if s.SyncPoint == nil {
		s.SyncPoint = new(Hash)
	}
	if err = s.SyncPoint.UnmarshalSSZ(buf[0:32]); err != nil {
		return err
	}

	// Offset (1) 'GraphState'
	if o1 = ssz.ReadOffset(buf[32:36]); o1 > size {
		return errOffset
	}

	// Offset (2) 'Blocks'
	if o2 = ssz.ReadOffset(buf[36:40]); o2 > size || o1 > o2 {
		return errOffset
	}

	// Field (1) 'GraphState'
	{
		buf = tail[o1:o2]
		if s.GraphState == nil {
			s.GraphState = new(GraphState)
		}
		if err = s.GraphState.UnmarshalSSZ(buf); err != nil {
			return err
		}
	}

	// Field (2) 'Blocks'
	{
		buf = tail[o2:]
		num, err := ssz.DecodeDynamicLength(buf, 500)
		if err != nil {
			return err
		}
		s.Blocks = make([]*BlockData, num)
		err = ssz.UnmarshalDynamic(buf, num, func(indx int, buf []byte) (err error) {
			if s.Blocks[indx] == nil {
				s.Blocks[indx] = new(BlockData)
			}
			if err = s.Blocks[indx].UnmarshalSSZ(buf); err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return err
		}
	}
	return err
}

// SizeSSZ returns the ssz encoded size in bytes for the SubDAG object
func (s *SubDAG) SizeSSZ() (size int) {
	size = 40

	// Field (1) 'GraphState'
	if s.GraphState == nil {
		s.GraphState = new(GraphState)
	}
	size += s.GraphState.SizeSSZ()

	// Field (2) 'Blocks'
	for ii := 0; ii < len(s.Blocks); ii++ {
		size += 4
		size += s.Blocks[ii].SizeSSZ()
	}

	return
}

// MarshalSSZ ssz marshals the ErrorResponse object
func (e *ErrorResponse) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, e.SizeSSZ())
	return e.MarshalSSZTo(buf[:0])
}

// MarshalSSZTo ssz marshals the ErrorResponse object to a target array
func (e *ErrorResponse) MarshalSSZTo(dst []byte) ([]byte, error) {
	var err error
	offset := int(4)

	// Offset (0) 'Message'
	dst = ssz.WriteOffset(dst, offset)
	offset += len(e.Message)

	// Field (0) 'Message'
	if len(e.Message) > 256 {
		return nil, errMarshalDynamicBytes
	}
	dst = append(dst, e.Message...)

	return dst, err
}

// UnmarshalSSZ ssz unmarshals the ErrorResponse object
func (e *ErrorResponse) UnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size < 4 {
		return errSize
	}

	tail := buf
	var o0 uint64

	// Offset (0) 'Message'
	if o0 = ssz.ReadOffset(buf[0:4]); o0 > size {
		return errOffset
	}

	// Field (0) 'Message'
	{
		buf = tail[o0:]
		e.Message = append(e.Message, buf...)
	}
	return err
}

// SizeSSZ returns the ssz encoded size in bytes for the ErrorResponse object
func (e *ErrorResponse) SizeSSZ() (size int) {
	size = 4

	// Field (0) 'Message'
	size += len(e.Message)

	return
}

// MarshalSSZ ssz marshals the MetaData object
func (m *MetaData) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, m.SizeSSZ())
	return m.MarshalSSZTo(buf[:0])
}

// MarshalSSZTo ssz marshals the MetaData object to a target array
func (m *MetaData) MarshalSSZTo(dst []byte) ([]byte, error) {
	var err error

	// Field (0) 'SeqNumber'
	dst = ssz.MarshalUint64(dst, m.SeqNumber)

	// Field (1) 'Subnets'
	if dst, err = ssz.MarshalFixedBytes(dst, m.Subnets, 8); err != nil {
		return nil, errMarshalFixedBytes
	}

	return dst, err
}

// UnmarshalSSZ ssz unmarshals the MetaData object
func (m *MetaData) UnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size != 16 {
		return errSize
	}

	// Field (0) 'SeqNumber'
	m.SeqNumber = ssz.UnmarshallUint64(buf[0:8])

	// Field (1) 'Subnets'
	m.Subnets = append(m.Subnets, buf[8:16]...)

	return err
}

// SizeSSZ returns the ssz encoded size in bytes for the MetaData object
func (m *MetaData) SizeSSZ() (size int) {
	size = 16
	return
}

// MarshalSSZ ssz marshals the Hash object
func (h *Hash) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, h.SizeSSZ())
	return h.MarshalSSZTo(buf[:0])
}

// MarshalSSZTo ssz marshals the Hash object to a target array
func (h *Hash) MarshalSSZTo(dst []byte) ([]byte, error) {
	var err error

	// Field (0) 'Hash'
	if dst, err = ssz.MarshalFixedBytes(dst, h.Hash, 32); err != nil {
		return nil, errMarshalFixedBytes
	}

	return dst, err
}

// UnmarshalSSZ ssz unmarshals the Hash object
func (h *Hash) UnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size != 32 {
		return errSize
	}

	// Field (0) 'Hash'
	h.Hash = append(h.Hash, buf[0:32]...)

	return err
}

// SizeSSZ returns the ssz encoded size in bytes for the Hash object
func (h *Hash) SizeSSZ() (size int) {
	size = 32
	return
}
