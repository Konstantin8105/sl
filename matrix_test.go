package sl_test

import (
	"fmt"
	"github.com/Konstantin8105/sl"
	"os"
)

func Example() {
	m := sl.New(3)
	//	RowIndexes  = [ 2 0 1 1 1 2 ]
	//	ColPos      = [ 1 0 1 1 0 2 ]
	//	Values      = [ 7 1 1 1 3 8 ]
	for _, err := range []error{
		m.Put(2, 1, 7.0),
		m.Put(0, 0, 1.0),
		m.Put(1, 1, 1.0),
		m.Put(1, 1, 1.0),
		m.Put(1, 0, 3.0),
		m.Put(2, 2, 8.0),
	} {
		if err != nil {
			panic(err)
		}
	}

	fmt.Fprintf(os.Stdout, "%s\n", m)

	fmt.Fprintf(os.Stdout, "Transform to Ssm:\n")
	if err := m.TransformTo(sl.Ssm); err != nil {
		panic(err)
	}

	fmt.Fprintf(os.Stdout, "%s\n", m)

	// Output:
	// Type       : triplet matrix format
	// Size       : 3
	// Values     : [7 1 1 1 3 8]
	// RowIndexes : [2 0 1 1 1 2]
	// ColPos     : [1 0 1 1 0 2]
	//
	// Transform to Ssm:
	// Type       : sparse symmetrical matrix
	// Size       : 3
	// Values     : [1 3 2 7 8]
	// RowIndexes : [0 1 1 2 2]
	// ColPos     : [0 2 4 5]
}
