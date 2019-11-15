package sl

import "github.com/Konstantin8105/errors"

// MatrixType is type of matrix
type MatrixType uint8

const (
	_    MatrixType = iota // ignore zero value of MatrixType for avoid `zero-error`
	Ssm                    // sparse symmetrical matrix
	Sltm                   // sparse lower triangular matrix
	Tm                     // triplet matrix format
)

func (m MatrixType) String() string {
	switch m {
	case Ssm:
		return "sparse symmetrical matrix"
	case Sltm:
		return "sparse lower triangular matrix"
	case Tm:
		return "triplet matrix format"
	}
	return "not defined matrix type(format)"
}

// Matrix - sparse symmetrical matrix or sparse lower triangular matrix
// in compressed-column(CCS) or triplet matrix fotmat.
//
// Example of storing symmetrical matrix or lower triangular:
//
//	General matrix:
//	[ 1 3 0 ]
//	[ 3 2 7 ]
//	[ 0 7 8 ]
//
//	Symmetrical matrix in lower triangle view:
//	[ 1 . . ]
//	[ 3 2 . ]
//	[ 0 7 8 ]
//
//	CCS view:
//	Format      = Ssm
//	Size        = 3             # amount rows and columns
//	                0 1 2 3 4   # position in `values` array
//	Values      = [ 1 3 2 7 8 ] # all non-zero values of matrix
//	RowIndexes  = [ 0 1 1 2 2 ] # row position for each `values`
//	ColPos      = [ 0 2 4 5 ]   # column positions
//		# values of 0 column : [0 ... 2)
//		# values of 1 column : [2 ... 4)
//		# values of 2 column : [4 ... 5)
//
// Example of storing symmetrical matrix in triplet format:
//
//	General matrix:
//	[ 1 3 0 ]
//	[ 3 2 7 ]
//	[ 0 7 8 ]
//
//	Symmetrical matrix in lower triangle view:
//	[ 1 . . ]
//	[ 3 2 . ]
//	[ 0 7 8 ]
//
//	Symmetrical matrix in triplet format:
//	Format      = Tm
//	Size        = 3             # amount rows and columns
//	                0 1 2 3 4   # position in `values` array
//	Values      = [ 1 3 2 7 8 ] # all non-zero values of matrix
//	RowIndexes  = [ 0 1 1 2 2 ] # row position for each `values`
//	ColPos      = [ 0 0 1 1 2 ] # column positions
//
// Note:
//
//	* all internal struct values are share for adding external features.
type Matrix struct {
	Format     MatrixType // matrix type(format)
	Size       int        // number of rows and columns
	Values     []float64  // all non-zero values of matrix
	RowIndexes []int      // row position for each `values`
	ColPos     []int      // column positions
}

// New create a new matrix in triplet format.
func New(s int) *Matrix {
	return &Matrix{
		Format:     Tm,
		Size:       s,
		Values:     make([]float64, 0, s),
		RowIndexes: make([]int, 0, s),
		ColPos:     make([]int, 0, s),
	}
}

// Put add value to matrix in triplet format.
//
// Input data:
//	r - row index
//	c - column index
//	x - value
//
// function return error if matrix or input data is not valid.
//
func (m *Matrix) Put(r, c int, x float64) error {
	// check input data
	var et errors.Tree
	if m == nil {
		et.Add("Matrix is nil")
	} else {
		if m.Format != Tm {
			et.Add(fmt.Errorf("Matrix type is not Triplet: %s", m.Format))
		} else {
			// row check
			if r < 0 {
				et.Add(fmt.Errorf("row index is negative: %d", r))
			}
			if r >= m.Size {
				et.Add(fmt.Errorf("row index is outside matrix: %d", r))
			}
			// column check
			if c < 0 {
				et.Add(fmt.Errorf("column index is negative: %d", c))
			}
			if c >= m.Size {
				et.Add(fmt.Errorf("column index is outside matrix: %d", c))
			}
			// index checking
			if r < c {
				et.Add(fmt.Errorf("row index %d is less column index %d:"+
					" not valid add in up at diagonal",
					r, c))
			}
			// value check
			if math.IsNaN(x) {
				et.Add(fmt.Errorf("value `x` is Nan value"))
			}
			if math.IsInf(x, 0) {
				et.Add(fmt.Errorf("value `x` is infinity value"))
			}
		}
	}
	if et.IsError() {
		et.Name = "function `Entry` error:"
		return et
	}

	if x != 0.0 { // ignore zero value for minimaze memory allocation
		// append new information
		m.RowIndexes = append(m.RowIndexes, r)
		m.ColPos = append(m.ColPos, c)
		m.Values = append(m.Values, x)
	}
	return nil
}

// Transform convert matrix
//
// Input data:
//	mt - matrix type to transformation
//
// function return error if matrix or input data is not valid.
//
func (m *Matrix) TransformTo(mt MatrixType) error {
	// check input data
	var et errors.Tree
	if m == nil {
		et.Add("Matrix is nil")
	} else {
		switch m {
		case Ssm, Sltm, Tm:
		default:
			et.Add("not valid type of matrix: %s", mt)
		}
	}
	if et.IsError() {
		et.Name = "function `TransformTo` error:"
		return et
	}

	// simple transformation
	if m.Format == mt {
		// matrix types(formats) are same
		return nil
	}
	if (m.Format == Ssm && mt == Sltm) || (m.Format == Sltm && mt == Ssm) {
		// matrix transformation without any changing
		m.Format = mt
		return nil
	}

	// transformation from triplet matrix format
	// TODO

	m.Format = mt
	return nil
}
