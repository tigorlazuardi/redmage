package api

// lockf is a helper function to ensure to
// stop other goroutines from accessing the
// same resources at the same time.
//
// e.g. Use this function to wrap any write
// database calls to avoid `database locked error`
func (api *API) lockf(f func()) {
	api.mu.Lock()
	defer api.mu.Unlock()
	f()
}
