package sl

// MatrixType is type of matrix
type MatrixType uint8

const (
	_    MatrixType = iota // ignore zero value of MatrixType for avoid `zero-error`
	Ssm                    // sparse symmetrical matrix
	Sltm                   // sparse lower triangular matrix
	Tm                     // triplet matrix format
)

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
//	* all internal struct values are share for adding external features.
type Matrix struct {
	Format     MatrixType // matrix type
	Size       int        // number of rows and columns
	Values     []float64  // all non-zero values of matrix
	RowIndexes []int      // row position for each `values`
	ColPos     []int      // column positions
}
