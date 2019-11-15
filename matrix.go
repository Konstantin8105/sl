package sl

import (
	"fmt"
	"github.com/Konstantin8105/errors"
	"math"
	"sort"
)

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

func (m Matrix) String() string {
	s := fmt.Sprintf("Type       : %s\n", m.Format)
	s += fmt.Sprintf("Size       : %d\n", m.Size)
	s += fmt.Sprintf("Values     : %v\n", m.Values)
	s += fmt.Sprintf("RowIndexes : %v\n", m.RowIndexes)
	s += fmt.Sprintf("ColPos     : %v\n", m.ColPos)
	return s
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
		et.Add(fmt.Errorf("Matrix is nil"))
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
		et.Add(fmt.Errorf("Matrix is nil"))
	} else {
		switch mt {
		case Ssm, Sltm, Tm:
		default:
			et.Add(fmt.Errorf("not valid type of matrix: %s", mt))
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

	// transformation from triplet matrix format.

	// Example of transformation
	//
	// from:
	//
	//	Format      = Tm
	//	Size        = 3
	//	Values      = [ 7 1 1 1 3 8 ]
	//	RowIndexes  = [ 2 0 1 1 1 2 ]
	//	ColPos      = [ 1 0 1 1 0 2 ]
	//
	// Note:
	//	* Matrix[1,1] in triplet format have 2 values and in result will be summ
	//	* Acceptable not sorted positions in triplet format
	//
	// to:
	//
	//	Format      = `mt`
	//	Size        = 3
	//	Values      = [ 1 3 1 7 1 8 ]
	//	RowIndexes  = [ 0 1 1 2 1 2 ]
	//	ColPos      = [ 0 2 5 6 ]
	//
	// Note:
	//	* compressed on same matrix position

	// sorting by ColPos to:
	//	Format      = Tm
	//	Size        = 3
	//	Values      = [ 1 3 1 7 1 8 ]
	//	RowIndexes  = [ 0 1 1 2 1 2 ]
	//	ColPos      = [ 0 0 1 1 1 2 ]
	sort.Slice(m.ColPos, func(i, j int) bool {
		if m.ColPos[i] < m.ColPos[j] {
			// swap
			m.Values[i], m.Values[j], m.RowIndexes[i], m.RowIndexes[j] =
				m.Values[j], m.Values[i], m.RowIndexes[j], m.RowIndexes[i]
			return true
		}
		return false
	})

	compressCol := func(colpos []int) []int {
		// compress ColPos
		// from:
		//	ColPos = [ 0 0 1 1 1 2 ]
		// to:
		//	           0 1 2 3   # position
		//	cp     = [ 0 2 3 1 ]
		cp := make([]int, m.Size+1)
		for i := range colpos {
			cp[colpos[i]+1]++
		}
		// cum summ
		// from:
		//	           0 1 2 3  # position
		//	cp     = [ 0 2 3 1 ]
		//	            |=====|
		// to:
		//	cp     = [ 0 2 5 6 ]
		for i := range cp {
			if i == 0 {
				continue
			}
			cp[i] += cp[i-1]
		}
		return cp
	}
	cp := compressCol(m.ColPos)

	// soring by RowIndexes to:
	//
	//	cp     = [ 0 2 5 6 ]
	//
	//	Format      = Tm
	//	Size        = 3
	//	Values      = [ 1 3 1 1 7 8 ]
	//	RowIndexes  = [ 0 1 1 1 2 2 ]
	//	ColPos      = [ 0 0 1 1 1 2 ]
	for k := 1; k < len(cp); k++ {
		sort.Slice(m.RowIndexes[cp[k-1]:cp[k]], func(i, j int) bool {
			if m.RowIndexes[cp[k-1]+i] < m.RowIndexes[cp[k-1]+j] {
				// swap
				m.Values[cp[k-1]+i], m.Values[cp[k-1]+j] =
					m.Values[cp[k-1]+j], m.Values[cp[k-1]+i]
				return true
			}
			return false
		})
	}

	// summary with same row and columns:
	//
	//	Format      = Tm
	//	Size        = 3
	//	Values      = [ 1 3 0 2 7 8 ]
	//	RowIndexes  = [ 0 1 1 1 2 2 ]
	//	ColPos      = [ 0 0 1 1 1 2 ]
	for i := range m.Values {
		if i == 0 {
			continue
		}
		if (m.ColPos[i] != m.ColPos[i-1]) ||
			(m.RowIndexes[i] != m.RowIndexes[i-1]) {
			continue
		}
		m.Values[i] += m.Values[i-1]
		m.Values[i-1] = 0.0
	}

	// calculate non-zero values
	var amountNonzeroValues int
	for i := range m.Values {
		if m.Values[i] == 0.0 {
			continue
		}
		amountNonzeroValues++
	}

	// reallocate memory
	var (
		v     = make([]float64, amountNonzeroValues)
		r     = make([]int, amountNonzeroValues)
		c     = make([]int, amountNonzeroValues)
		count int
	)
	for i := range m.Values {
		if m.Values[i] == 0.0 {
			continue
		}
		v[count] = m.Values[i]
		r[count] = m.RowIndexes[i]
		c[count] = m.ColPos[i]
		count++
	}
	m.Values, v = v, m.Values
	m.RowIndexes, r = r, m.RowIndexes
	m.ColPos = compressCol(c)

	// free memory for reuse memory
	// Free(v)
	// Free(r)
	// Free(c)

	m.Format = mt
	return nil
}

// TODO: empty rows
