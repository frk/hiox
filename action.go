package hxio

// Action wraps a set of methods that are executed in sequence to accomplish
// the Action's objective. Each of the Action's methods is intended to represent
// a single task with responsibilities separate from, although related to, those
// of the other tasks.
//
// The names of the methods should serve as a guide to what the intended responsibility
// for each method is, however not every task will fit neatly into any of those methods
// and therefore, in such a case, it is left to the developer's judgement to decide
// which method should implement such a task.
//
// If, during execution, any of the Action's methods returns an error the execution
// will skip over to, and invoke, the Done method passing it the error value and
// leaving the methods in-between uninvoked.
type Action interface {
	BeforeValidate() error // Prepare for input validation.
	Validate() error       // Validate the input.
	AfterValidate() error  // Post processing after the input has been validated.
	BeforeExecute() error  // Prepare for main task execution.
	Execute() error        // Execute the main task.
	AfterExecute() error   // Post processing after the main task has been executed.

	// Done is the last of the Action's methods to be invoked, it is special
	// in that regard that it will be invoked irrespective of whether or not
	// one of the preceding methods returned an error.
	//
	// The in error parameter will either be nil or it will hold the error
	// value returned by one of the preceding methods. The out error paremeter
	// is used as the final return value of the ExecuteAction function and this
	// gives Done the ability to override the error returned from one of those
	// preceding methods if need be.
	Done(in error) (out error)
}

// ExecuteAction executes the given Action returning the error returned from
// the Action's Done method, if any.
func ExecuteAction(a Action) error {
	if err := execAction(a); err != nil && err != IsDone {
		return a.Done(err)
	}
	return a.Done(nil)
}

func execAction(a Action) error {
	if err := a.BeforeValidate(); err != nil {
		return err
	}
	if err := a.Validate(); err != nil {
		return err
	}
	if err := a.AfterValidate(); err != nil {
		return err
	}
	if err := a.BeforeExecute(); err != nil {
		return err
	}
	if err := a.Execute(); err != nil {
		return err
	}
	return a.AfterExecute()
}

// ActionBase is a noop helper type that can be embedded by user defined
// types that are intended to implement the Action interface but do not
// need to, nor want to, declare every single one of its methods.
type ActionBase struct{ actionbase }

// actionbase is embedded by ActionBase to artificially increase the depth level
// of the noop methods to reduce the possibility of an "ambiguous selector" issue.
type actionbase struct{}

func (actionbase) BeforeValidate() error { return nil }
func (actionbase) Validate() error       { return nil }
func (actionbase) AfterValidate() error  { return nil }
func (actionbase) BeforeExecute() error  { return nil }
func (actionbase) Execute() error        { return nil }
func (actionbase) AfterExecute() error   { return nil }
func (actionbase) Done(err error) error  { return err }
