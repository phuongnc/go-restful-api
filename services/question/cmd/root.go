package cmd

// RootExec func
func RootExec() {
	runtime := newRuntime()
	runtime.serve()
}
