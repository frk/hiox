package httpcrud

import (
	"errors"
	"testing"

	"github.com/frk/compare"
)

func TestExecuteAction(t *testing.T) {
	var aerr = errors.New("action error")
	var berr = errors.New("action error (b)")

	tests := []struct {
		a    Action
		want Action
		err  error
	}{{
		a: &fakeact{},
		want: &fakeact{
			beforeValidate: true,
			validate:       true,
			afterValidate:  true,
			beforeExecute:  true,
			execute:        true,
			afterExecute:   true,
			done:           true,
		},
	}, {
		a: &fakeact{
			validate: aerr,
			done:     aerr,
		},
		want: &fakeact{
			beforeValidate: true,
			validate:       true,
			donein:         aerr,
			done:           true,
		},
		err: aerr,
	}, {
		a: &fakeact{
			beforeExecute: aerr,
			done:          nil,
		},
		want: &fakeact{
			beforeValidate: true,
			validate:       true,
			afterValidate:  true,
			beforeExecute:  true,
			donein:         aerr,
			done:           true,
		},
	}, {
		a: &fakeact{
			execute: aerr,
			done:    berr,
		},
		want: &fakeact{
			beforeValidate: true,
			validate:       true,
			afterValidate:  true,
			beforeExecute:  true,
			execute:        true,
			donein:         aerr,
			done:           true,
		},
		err: berr,
	}}

	for _, tt := range tests {
		err := ExecuteAction(tt.a)
		if e := compare.Compare(err, tt.err); e != nil {
			t.Error(e)
		}
		if e := compare.Compare(tt.a, tt.want); e != nil {
			t.Error(e)
		}
	}
}

type fakeact struct {
	beforeValidate interface{}
	validate       interface{}
	afterValidate  interface{}
	beforeExecute  interface{}
	execute        interface{}
	afterExecute   interface{}
	done           interface{}
	donein         error
}

func (f *fakeact) BeforeValidate() (err error) {
	if f.beforeValidate != nil {
		err = f.beforeValidate.(error)
	}
	f.beforeValidate = true
	return err
}

func (f *fakeact) Validate() (err error) {
	if f.validate != nil {
		err = f.validate.(error)
	}
	f.validate = true
	return err
}

func (f *fakeact) AfterValidate() (err error) {
	if f.afterValidate != nil {
		err = f.afterValidate.(error)
	}
	f.afterValidate = true
	return err
}

func (f *fakeact) BeforeExecute() (err error) {
	if f.beforeExecute != nil {
		err = f.beforeExecute.(error)
	}
	f.beforeExecute = true
	return err
}

func (f *fakeact) Execute() (err error) {
	if f.execute != nil {
		err = f.execute.(error)
	}
	f.execute = true
	return err
}

func (f *fakeact) AfterExecute() (err error) {
	if f.afterExecute != nil {
		err = f.afterExecute.(error)
	}
	f.afterExecute = true
	return err
}

func (f *fakeact) Done(in error) (err error) {
	f.donein = in
	if f.done != nil {
		err = f.done.(error)
	}
	f.done = true
	return err
}
