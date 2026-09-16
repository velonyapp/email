package api

func mapError(err error) error {
	if err == nil {
		return nil
	}

	switch {

	default:
		return err // let this be i want to debug
		// return kerrors.InternalServer(
		// 	"",
		// 	"internal server error",
		// )
	}
}
