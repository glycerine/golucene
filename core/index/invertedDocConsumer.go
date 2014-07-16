package index

import (
	"github.com/balzaczyy/golucene/core/index/model"
	"github.com/balzaczyy/golucene/core/util"
)

// type InvertedDocConsumer interface {
// 	// Abort (called after hitting abort error)
// 	abort()
// 	addField(*DocInverterPerField, *model.FieldInfo) InvertedDocConsumerPerField
// 	// Flush a new segment
// 	flush(map[string]InvertedDocConsumerPerField, *model.SegmentWriteState) error
// 	startDocument()
// 	finishDocument() error
// }

/*
This class is passed each token produced by the analyzer on each
field during indexing, and it stores these tokens in a hash table,
and allocates separate byte streams per token. Consumers of this
class, eg FreqProxTermsWriter and TermVectorsConsumer, write their
own byte streams under each term.
*/
type TermsHash interface {
	startDocument()
	finishDocument() error
	abort()
	setTermBytePool(*util.ByteBlockPool)
}

type TermsHashImpl struct {
	nextTermsHash TermsHash

	intPool      *util.IntBlockPool
	bytePool     *util.ByteBlockPool
	termBytePool *util.ByteBlockPool
	bytesUsed    util.Counter

	docState *docState

	trackAllocations bool
}

func newTermsHash(docWriter *DocumentsWriterPerThread,
	trackAllocations bool, nextTermsHash TermsHash) *TermsHashImpl {

	ans := &TermsHashImpl{
		docState:         docWriter.docState,
		trackAllocations: trackAllocations,
		nextTermsHash:    nextTermsHash,
		intPool:          util.NewIntBlockPool(docWriter.intBlockAllocator),
		bytePool:         util.NewByteBlockPool(docWriter.byteBlockAllocator),
	}
	if trackAllocations {
		ans.bytesUsed = docWriter._bytesUsed
	} else {
		ans.bytesUsed = util.NewCounter()
	}
	if nextTermsHash != nil {
		ans.termBytePool = ans.bytePool
		nextTermsHash.setTermBytePool(ans.bytePool)
	}
	return ans
}

func (hash *TermsHashImpl) setTermBytePool(p *util.ByteBlockPool) {
	hash.termBytePool = p
}

func (hash *TermsHashImpl) abort() {
	defer func() {
		if hash.nextTermsHash != nil {
			hash.nextTermsHash.abort()
		}
	}()
	hash.reset()
}

/* Clear all state */
func (hash *TermsHashImpl) reset() {
	// we don't reuse so we drop everything and don't fill with 0
	hash.intPool.Reset(false, false)
	hash.bytePool.Reset(false, false)
}

func (hash *TermsHashImpl) flush(fieldsToFlush map[string]*TermsHashPerField,
	state *model.SegmentWriteState) error {
	panic("not implemented yet")
	// childFields := make(map[string]TermsHashConsumerPerField)
	// var nextChildFieldFields map[string]InvertedDocConsumerPerField
	// if hash.nextTermsHash != nil {
	// 	nextChildFieldFields = make(map[string]InvertedDocConsumerPerField)
	// }

	// for k, v := range fieldsToFlush {
	// 	perField := v.(*TermsHashPerField)
	// 	childFields[k] = perField.consumer
	// 	if hash.nextTermsHash != nil {
	// 		nextChildFieldFields[k] = perField.nextPerField
	// 	}
	// }

	// err := hash.consumer.flush(childFields, state)
	// if err == nil && hash.nextTermsHash != nil {
	// 	err = hash.nextTermsHash.flush(nextChildFieldFields, state)
	// }
	// return err
}

// func (h *TermsHash) addField(docInverterPerField *DocInverterPerField,
// 	fieldInfo *model.FieldInfo) InvertedDocConsumerPerField {
// 	return newTermsHashPerField(docInverterPerField, h, h.nextTermsHash, fieldInfo)
// }

func (h *TermsHashImpl) finishDocument() error {
	if h.nextTermsHash != nil {
		return h.nextTermsHash.finishDocument()
	}
	return nil
}

func (h *TermsHashImpl) startDocument() {
	if h.nextTermsHash != nil {
		h.nextTermsHash.startDocument()
	}
}
