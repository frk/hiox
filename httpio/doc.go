// The package httpio provides utility types and interfaces to complement the httpcrud package.
//
// The main two types, RequestReader and ResponseWriter, both implement parts
// of the httpcrud.Handler interface and can therefore be utilized as composition
// blocks in a complete implementation of the interface.
//
// Both, the RequestReader and ResponseWriter, by exporting a number of interface
// type fields, allow the user to split the implementation into smaller, more specific
// subtasks to read and write the requests and responses, respectively.
package httpio
