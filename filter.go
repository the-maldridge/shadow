package shadow

// A NumericFilter can be applied to the number fields of a type, and
// provides a way to do filtering based on numeric ranging, not-below,
// not-above etc.  The function should return true if an item matches
// the filter, or false if it does not.
type NumericFilter func(int) bool

// A StringFilter is exactly the same as a NumericFilter, but is
// applied to strings instead.
type StringFilter func(string) bool
