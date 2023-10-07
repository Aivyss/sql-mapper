package context

// getQueryClientKey (for internal) get an identifier of QueryClient
func getQueryClientKey(identifier string, readOnly bool) queryClientKey {
	return queryClientKey{
		identifier: identifier,
		readOnly:   readOnly,
	}
}
